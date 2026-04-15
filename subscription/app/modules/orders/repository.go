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
	CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO, pkg *PackageInfo) (*Order, error)
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

// GetPackageByID mengambil data paket untuk keperluan validasi dan pricing
func (r *orderRepo) GetPackageByID(ctx context.Context, packageID string) (*PackageInfo, error) {
	query := `
		SELECT id, name, price, status
		FROM subscription.packages
		WHERE id = $1 AND status != 'DELETED'
		LIMIT 1
	`
	var p PackageInfo
	err := r.db.Pool.QueryRow(ctx, query, packageID).Scan(&p.ID, &p.Name, &p.Price, &p.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("gagal mengambil data paket: %w", err)
	}
	return &p, nil
}

// CreateOrder menyimpan order baru ke database
func (r *orderRepo) CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO, pkg *PackageInfo) (*Order, error) {
	id := uuid.New().String()
	invoice := generateInvoiceNumber()

	query := `
		INSERT INTO subscription.orders
			(id, invoice_number, user_id, package_id, payment_method, total_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'PENDING_PAYMENT', NOW(), NOW())
		RETURNING id, invoice_number, package_id, payment_method, total_price, status, created_at
	`

	var o Order
	err := r.db.Pool.QueryRow(ctx, query, id, invoice, userID, dto.PackageID, dto.PaymentMethod, pkg.Price).
		Scan(&o.ID, &o.InvoiceNumber, &o.PackageID, &o.PaymentMethod, &o.TotalPrice, &o.Status, &o.CreatedAt)
	if err != nil {
		// Handle duplicate invoice (sangat jarang, retry sekali)
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "invoice_number") {
			invoice = generateInvoiceNumber()
			err = r.db.Pool.QueryRow(ctx, query, uuid.New().String(), invoice, userID, dto.PackageID, dto.PaymentMethod, pkg.Price).
				Scan(&o.ID, &o.InvoiceNumber, &o.PackageID, &o.PaymentMethod, &o.TotalPrice, &o.Status, &o.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("gagal membuat order setelah retry: %w", err)
			}
			return &o, nil
		}
		return nil, fmt.Errorf("gagal membuat order: %w", err)
	}
	return &o, nil
}

// generateInvoiceNumber membuat nomor invoice unik format: INV-{YEAR}-{7 digit random}
func generateInvoiceNumber() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("INV-%d-%07d", time.Now().Year(), rng.Intn(9999999))
}
