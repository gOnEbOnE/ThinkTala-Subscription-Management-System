package orders

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

// ==========================================
// REPOSITORY INTERFACE
// ==========================================

type Repository interface {
	GetPackageByID(ctx context.Context, packageID string) (*PackageInfo, error)
	GetLatestKYCStatusByUser(ctx context.Context, userID string) (string, error)
	GetPricingTier(ctx context.Context, packageID string, durationMonths int) (float64, error)
	CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO, totalPrice float64) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*OrderRecord, error)
	ListOrdersByUser(ctx context.Context, userID string) ([]ClientOrderListItem, error)
	ListOrdersForAdmin(ctx context.Context) ([]AdminOrderListItem, error)
	UpdateOrderStatus(ctx context.Context, orderID, newStatus, verificationNote string) error
	SavePaymentProof(ctx context.Context, orderID string, file PaymentProofFile) (*UploadPaymentProofResult, error)
	GetPaymentProof(ctx context.Context, orderID string) (*PaymentProofFile, error)
	CreateSubscriptionFromOrder(ctx context.Context, orderID string) (*ActivationResult, error)
	ListActiveSubscriptionsByUser(ctx context.Context, userID string) ([]SubscriptionStatus, error)
	GetActiveSubscriptionByUser(ctx context.Context, userID string) (*SubscriptionStatus, error)
}

// ==========================================
// REPOSITORY IMPLEMENTATION
// ==========================================

type orderRepo struct {
	db *database.DBWrapper
}

func NewRepository(db *database.DBWrapper) Repository {
	return &orderRepo{db: db}
}

// GetPackageByID validasi paket aktif
func (r *orderRepo) GetPackageByID(ctx context.Context, packageID string) (*PackageInfo, error) {
	var p PackageInfo
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, name, status, price FROM subscription.packages
		 WHERE id = $1 AND status != 'DELETED' LIMIT 1`, packageID,
	).Scan(&p.ID, &p.Name, &p.Status, &p.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mengambil data paket: %w", err)
	}
	return &p, nil
}

// GetLatestKYCStatusByUser mengambil status KYC terbaru milik user.
// Return empty string jika user belum pernah submit KYC.
func (r *orderRepo) GetLatestKYCStatusByUser(ctx context.Context, userID string) (string, error) {
	var status string
	err := r.db.Pool.QueryRow(ctx,
		`SELECT status
		 FROM kyc_submissions
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT 1`,
		userID,
	).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("gagal mengambil status KYC: %w", err)
	}
	return strings.ToLower(strings.TrimSpace(status)), nil
}

// GetPricingTier mengambil harga total dari package_pricing berdasarkan package_id + duration
func (r *orderRepo) GetPricingTier(ctx context.Context, packageID string, durationMonths int) (float64, error) {
	var price float64
	err := r.db.Pool.QueryRow(ctx,
		`SELECT price FROM subscription.package_pricing
		 WHERE package_id = $1 AND duration_months = $2 LIMIT 1`,
		packageID, durationMonths,
	).Scan(&price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil // tier tidak tersedia
		}
		return 0, fmt.Errorf("gagal mengambil pricing tier: %w", err)
	}
	return price, nil
}

// CreateOrder insert order baru ke database
func (r *orderRepo) CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO, totalPrice float64) (*Order, error) {
	id := uuid.New().String()
	invoice := generateInvoiceNumber()
	durationMonths := dto.DurationMonths
	if durationMonths <= 0 {
		durationMonths = 1
	}

	var o Order
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO subscription.orders
		   (id, invoice_number, user_id, package_id, duration_months, payment_method, total_price, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING_PAYMENT', NOW(), NOW())
		 RETURNING id, invoice_number, package_id, duration_months, payment_method, total_price, status,
		          FALSE AS has_payment_proof, NULL::timestamp AS payment_proof_uploaded_at, created_at`,
		id, invoice, userID, dto.PackageID, durationMonths, dto.PaymentMethod, totalPrice,
	).Scan(&o.ID, &o.InvoiceNumber, &o.PackageID, &o.DurationMonths, &o.PaymentMethod, &o.TotalPrice, &o.Status, &o.HasPaymentProof, &o.PaymentProofUploadedAt, &o.CreatedAt)

	if err != nil {
		// Retry sekali jika duplicate invoice
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "invoice_number") {
			invoice = generateInvoiceNumber()
			err = r.db.Pool.QueryRow(ctx,
				`INSERT INTO subscription.orders
				   (id, invoice_number, user_id, package_id, duration_months, payment_method, total_price, status, created_at, updated_at)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING_PAYMENT', NOW(), NOW())
				 RETURNING id, invoice_number, package_id, duration_months, payment_method, total_price, status,
				          FALSE AS has_payment_proof, NULL::timestamp AS payment_proof_uploaded_at, created_at`,
				uuid.New().String(), invoice, userID, dto.PackageID, durationMonths, dto.PaymentMethod, totalPrice,
			).Scan(&o.ID, &o.InvoiceNumber, &o.PackageID, &o.DurationMonths, &o.PaymentMethod, &o.TotalPrice, &o.Status, &o.HasPaymentProof, &o.PaymentProofUploadedAt, &o.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("gagal membuat order setelah retry: %w", err)
			}
			return &o, nil
		}
		return nil, fmt.Errorf("gagal membuat order: %w", err)
	}
	return &o, nil
}

func (r *orderRepo) getOrderByIDWithFallback(ctx context.Context, orderID string, withUsers bool) (*OrderRecord, error) {
	baseWithUsers := `SELECT
		o.id,
		o.invoice_number,
		o.user_id,
		o.package_id,
		COALESCE(p.name, '-') AS package_name,
		COALESCE(u.name, 'Unknown') AS client_name,
		COALESCE(u.email, '-') AS client_email,
		o.duration_months,
		o.total_price,
		o.payment_method,
		o.status,
		COALESCE(o.verification_note, '') AS verification_note,
		(o.payment_proof IS NOT NULL AND octet_length(o.payment_proof) > 0) AS has_payment_proof,
		o.payment_proof_uploaded_at,
		COALESCE(o.payment_proof_filename, '') AS payment_proof_filename,
		COALESCE(o.payment_proof_content_type, '') AS payment_proof_content_type,
		o.created_at
	FROM subscription.orders o
	LEFT JOIN subscription.packages p ON p.id = o.package_id
	LEFT JOIN users u ON u.id = o.user_id
	WHERE o.id = $1
	LIMIT 1`

	baseWithoutUsers := `SELECT
		o.id,
		o.invoice_number,
		o.user_id,
		o.package_id,
		COALESCE(p.name, '-') AS package_name,
		'Unknown' AS client_name,
		'-' AS client_email,
		o.duration_months,
		o.total_price,
		o.payment_method,
		o.status,
		COALESCE(o.verification_note, '') AS verification_note,
		(o.payment_proof IS NOT NULL AND octet_length(o.payment_proof) > 0) AS has_payment_proof,
		o.payment_proof_uploaded_at,
		COALESCE(o.payment_proof_filename, '') AS payment_proof_filename,
		COALESCE(o.payment_proof_content_type, '') AS payment_proof_content_type,
		o.created_at
	FROM subscription.orders o
	LEFT JOIN subscription.packages p ON p.id = o.package_id
	WHERE o.id = $1
	LIMIT 1`

	query := baseWithoutUsers
	if withUsers {
		query = baseWithUsers
	}

	var rec OrderRecord
	err := r.db.Pool.QueryRow(ctx, query, orderID).Scan(
		&rec.OrderID,
		&rec.InvoiceNumber,
		&rec.UserID,
		&rec.PackageID,
		&rec.PackageName,
		&rec.ClientName,
		&rec.ClientEmail,
		&rec.DurationMonths,
		&rec.TotalPrice,
		&rec.PaymentMethod,
		&rec.Status,
		&rec.VerificationNote,
		&rec.HasPaymentProof,
		&rec.PaymentProofUploadedAt,
		&rec.PaymentProofFileName,
		&rec.PaymentProofContentType,
		&rec.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rec, nil
}

// GetOrderByID mengambil satu order lengkap untuk kebutuhan client/admin
func (r *orderRepo) GetOrderByID(ctx context.Context, orderID string) (*OrderRecord, error) {
	rec, err := r.getOrderByIDWithFallback(ctx, orderID, true)
	if err != nil {
		if strings.Contains(err.Error(), `relation "users" does not exist`) {
			fallback, fbErr := r.getOrderByIDWithFallback(ctx, orderID, false)
			if fbErr != nil {
				return nil, fmt.Errorf("gagal mengambil detail order: %w", fbErr)
			}
			return fallback, nil
		}
		return nil, fmt.Errorf("gagal mengambil detail order: %w", err)
	}
	return rec, nil
}

// ListOrdersByUser mengambil riwayat pesanan milik user client
func (r *orderRepo) ListOrdersByUser(ctx context.Context, userID string) ([]ClientOrderListItem, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT
			o.id,
			o.invoice_number,
			COALESCE(p.name, '-') AS package_name,
			o.total_price,
			o.payment_method,
			o.status,
			(o.payment_proof IS NOT NULL AND octet_length(o.payment_proof) > 0) AS has_payment_proof,
			o.payment_proof_uploaded_at,
			o.created_at
		 FROM subscription.orders o
		 LEFT JOIN subscription.packages p ON p.id = o.package_id
		 WHERE o.user_id = $1
		 ORDER BY o.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil riwayat order client: %w", err)
	}
	defer rows.Close()

	list := make([]ClientOrderListItem, 0)
	for rows.Next() {
		var item ClientOrderListItem
		if err := rows.Scan(
			&item.OrderID,
			&item.InvoiceNumber,
			&item.PackageName,
			&item.TotalPrice,
			&item.PaymentMethod,
			&item.Status,
			&item.HasPaymentProof,
			&item.PaymentProofUploadedAt,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("gagal scan riwayat order client: %w", err)
		}
		list = append(list, item)
	}

	if list == nil {
		return []ClientOrderListItem{}, nil
	}
	return list, nil
}

// ListOrdersForAdmin mengambil daftar order untuk operasional
func (r *orderRepo) ListOrdersForAdmin(ctx context.Context) ([]AdminOrderListItem, error) {
	queryWithUsers := `SELECT
		o.id,
		o.invoice_number,
		COALESCE(u.name, 'Unknown') AS client_name,
		COALESCE(p.name, '-') AS package_name,
		o.total_price,
		o.payment_method,
		o.status,
		(o.payment_proof IS NOT NULL AND octet_length(o.payment_proof) > 0) AS has_payment_proof,
		o.payment_proof_uploaded_at,
		o.created_at
	FROM subscription.orders o
	LEFT JOIN users u ON u.id = o.user_id
	LEFT JOIN subscription.packages p ON p.id = o.package_id
	ORDER BY o.created_at DESC`

	queryWithoutUsers := `SELECT
		o.id,
		o.invoice_number,
		'Unknown' AS client_name,
		COALESCE(p.name, '-') AS package_name,
		o.total_price,
		o.payment_method,
		o.status,
		(o.payment_proof IS NOT NULL AND octet_length(o.payment_proof) > 0) AS has_payment_proof,
		o.payment_proof_uploaded_at,
		o.created_at
	FROM subscription.orders o
	LEFT JOIN subscription.packages p ON p.id = o.package_id
	ORDER BY o.created_at DESC`

	rows, err := r.db.Pool.Query(ctx, queryWithUsers)
	if err != nil && strings.Contains(err.Error(), `relation "users" does not exist`) {
		rows, err = r.db.Pool.Query(ctx, queryWithoutUsers)
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar order admin: %w", err)
	}
	defer rows.Close()

	list := make([]AdminOrderListItem, 0)
	for rows.Next() {
		var item AdminOrderListItem
		if err := rows.Scan(
			&item.OrderID,
			&item.InvoiceNumber,
			&item.ClientName,
			&item.PackageName,
			&item.TotalPrice,
			&item.PaymentMethod,
			&item.Status,
			&item.HasPaymentProof,
			&item.PaymentProofUploadedAt,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("gagal scan daftar order admin: %w", err)
		}
		list = append(list, item)
	}

	if list == nil {
		return []AdminOrderListItem{}, nil
	}
	return list, nil
}

// UpdateOrderStatus mengubah status order
func (r *orderRepo) UpdateOrderStatus(ctx context.Context, orderID, newStatus, verificationNote string) error {
	cmd, err := r.db.Pool.Exec(ctx,
		`UPDATE subscription.orders
		 SET status = $2,
		     verification_note = $3,
		     updated_at = NOW()
		 WHERE id = $1`,
		orderID, newStatus, strings.TrimSpace(verificationNote),
	)
	if err != nil {
		return fmt.Errorf("gagal mengubah status order: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("pesanan tidak ditemukan")
	}
	return nil
}

// SavePaymentProof menyimpan bukti transfer ke orders
func (r *orderRepo) SavePaymentProof(ctx context.Context, orderID string, file PaymentProofFile) (*UploadPaymentProofResult, error) {
	var uploadedAt time.Time
	err := r.db.Pool.QueryRow(ctx,
		`UPDATE subscription.orders
		 SET payment_proof = $2,
		     payment_proof_filename = $3,
		     payment_proof_content_type = $4,
		     payment_proof_uploaded_at = NOW(),
		     updated_at = NOW()
		 WHERE id = $1
		 RETURNING payment_proof_uploaded_at`,
		orderID, file.Data, file.FileName, file.ContentType,
	).Scan(&uploadedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal menyimpan bukti transfer: %w", err)
	}

	return &UploadPaymentProofResult{
		OrderID:                orderID,
		HasPaymentProof:        true,
		PaymentProofUploadedAt: &uploadedAt,
		Message:                "Bukti transfer berhasil diunggah",
	}, nil
}

// GetPaymentProof mengambil file bukti transfer berdasarkan order
func (r *orderRepo) GetPaymentProof(ctx context.Context, orderID string) (*PaymentProofFile, error) {
	var file PaymentProofFile
	err := r.db.Pool.QueryRow(ctx,
		`SELECT payment_proof, COALESCE(payment_proof_filename, ''), COALESCE(payment_proof_content_type, ''), payment_proof_uploaded_at
		 FROM subscription.orders
		 WHERE id = $1
		   AND payment_proof IS NOT NULL
		   AND octet_length(payment_proof) > 0
		 LIMIT 1`,
		orderID,
	).Scan(&file.Data, &file.FileName, &file.ContentType, &file.UploadedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mengambil bukti transfer: %w", err)
	}

	return &file, nil
}

// CreateSubscriptionFromOrder mengaktifkan subscription dari order yang sudah PAID
func (r *orderRepo) CreateSubscriptionFromOrder(ctx context.Context, orderID string) (*ActivationResult, error) {
	order, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("pesanan tidak ditemukan")
	}
	if order.Status != "PAID" {
		return nil, fmt.Errorf("aktivasi hanya bisa diproses untuk order berstatus PAID")
	}

	var existingID string
	err = r.db.Pool.QueryRow(ctx,
		`SELECT id FROM subscription.subscriptions WHERE order_id = $1 LIMIT 1`,
		orderID,
	).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("subscription untuk order ini sudah aktif")
	}
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("gagal mengecek subscription existing: %w", err)
	}

	now := time.Now().In(time.FixedZone("WIB", 7*3600))
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	duration := order.DurationMonths
	if duration <= 0 {
		duration = 1
	}
	endDate := startDate.AddDate(0, duration, 0)

	res := &ActivationResult{}
	err = r.db.Pool.QueryRow(ctx,
		`INSERT INTO subscription.subscriptions
			(id, order_id, user_id, package_id, start_date, end_date, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, 'ACTIVE', NOW(), NOW())
		 RETURNING id, start_date, end_date, status`,
		uuid.New().String(), order.OrderID, order.UserID, order.PackageID, startDate, endDate,
	).Scan(&res.SubscriptionID, &res.StartDate, &res.EndDate, &res.Status)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat subscription: %w", err)
	}

	return res, nil
}

// ListActiveSubscriptionsByUser mengambil semua subscription aktif milik client.
func (r *orderRepo) ListActiveSubscriptionsByUser(ctx context.Context, userID string) ([]SubscriptionStatus, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT
			s.id,
			s.package_id,
			COALESCE(p.name, '-') AS package_name,
			s.start_date,
			s.end_date,
			s.status
		 FROM subscription.subscriptions s
		 LEFT JOIN subscription.packages p ON p.id = s.package_id
		 WHERE s.user_id = $1
		   AND s.status = 'ACTIVE'
		   AND s.end_date >= CURRENT_DATE
		 ORDER BY s.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar subscription aktif: %w", err)
	}
	defer rows.Close()

	list := make([]SubscriptionStatus, 0)
	for rows.Next() {
		var item SubscriptionStatus
		if err := rows.Scan(
			&item.SubscriptionID,
			&item.PackageID,
			&item.PackageName,
			&item.StartDate,
			&item.EndDate,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("gagal scan subscription aktif: %w", err)
		}
		list = append(list, item)
	}

	if list == nil {
		return []SubscriptionStatus{}, nil
	}
	return list, nil
}

// GetActiveSubscriptionByUser mempertahankan perilaku lama: ambil 1 subscription aktif pertama.
func (r *orderRepo) GetActiveSubscriptionByUser(ctx context.Context, userID string) (*SubscriptionStatus, error) {
	list, err := r.ListActiveSubscriptionsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func generateInvoiceNumber() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("INV-%d-%07d", time.Now().Year(), rng.Intn(9999999))
}
