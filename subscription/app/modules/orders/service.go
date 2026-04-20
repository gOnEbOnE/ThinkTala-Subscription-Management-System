package orders

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ==========================================
// SERVICE INTERFACE
// ==========================================

type Service interface {
	CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO) (*Order, error)
	ListOrdersForClient(ctx context.Context, userID string) ([]ClientOrderListItem, error)
	GetOrderDetailForClient(ctx context.Context, userID, orderID string) (*ClientOrderDetail, error)
	UploadPaymentProof(ctx context.Context, userID, orderID string, file PaymentProofFile) (*UploadPaymentProofResult, error)
	GetPaymentProofForClient(ctx context.Context, userID, orderID string) (*PaymentProofFile, error)
	ListOrdersForAdmin(ctx context.Context) ([]AdminOrderListItem, error)
	GetOrderDetailForAdmin(ctx context.Context, orderID string) (*AdminOrderDetail, error)
	GetPaymentProofForAdmin(ctx context.Context, orderID string) (*PaymentProofFile, error)
	VerifyOrder(ctx context.Context, orderID, action, rejectReason string) (*VerifyResult, error)
	ActivateOrderSystem(ctx context.Context, orderID string) (*ActivationResult, error)
	GetActiveSubscriptions(ctx context.Context, userID string) ([]SubscriptionStatus, error)
	GetActiveSubscription(ctx context.Context, userID string) (*SubscriptionStatus, error)
	ProcessCreateOrderJob(ctx context.Context, payload interface{}) (interface{}, error)
}

var (
	ErrOrderNotFound       = errors.New("order_not_found")
	ErrOrderForbidden      = errors.New("order_forbidden")
	ErrPaymentProofMissing = errors.New("payment_proof_missing")
)

// ==========================================
// SERVICE IMPLEMENTATION
// ==========================================

type orderService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &orderService{repo: repo}
}

// CreateOrder — validasi + tentukan harga dari package_pricing + buat order
func (s *orderService) CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO) (*Order, error) {
	// Validasi input wajib
	if dto.PackageID == "" {
		return nil, errors.New("package_id wajib diisi")
	}
	dto.PaymentMethod = strings.TrimSpace(dto.PaymentMethod)
	if dto.PaymentMethod == "" {
		return nil, errors.New("payment_method wajib diisi")
	}
	if strings.ToUpper(dto.PaymentMethod) != "TRANSFER BANK" {
		return nil, errors.New("metode pembayaran yang didukung saat ini hanya Transfer Bank")
	}
	dto.PaymentMethod = "Transfer Bank"
	if dto.DurationMonths <= 0 {
		dto.DurationMonths = 1
	}
	if userID == "" {
		return nil, errors.New("user tidak teridentifikasi, silakan login ulang")
	}

	// Guard bisnis: user wajib lulus KYC sebelum bisa membuat pesanan.
	kycStatus, err := s.repo.GetLatestKYCStatusByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	switch kycStatus {
	case "approved":
		// Lanjut proses pembuatan order.
	case "pending":
		return nil, errors.New("verifikasi KYC Anda masih dalam proses review, pembelian paket belum dapat dilakukan")
	case "rejected":
		return nil, errors.New("verifikasi KYC Anda ditolak, silakan perbaiki data KYC terlebih dahulu")
	case "":
		return nil, errors.New("Anda belum melakukan verifikasi KYC dan tidak bisa melakukan pembelian paket")
	default:
		return nil, errors.New("status KYC Anda belum memenuhi syarat untuk pembelian paket")
	}

	// Validasi paket aktif
	pkg, err := s.repo.GetPackageByID(ctx, dto.PackageID)
	if err != nil {
		return nil, fmt.Errorf("gagal memvalidasi paket: %w", err)
	}
	if pkg == nil {
		return nil, errors.New("paket tidak ditemukan")
	}
	if pkg.Status != "ACTIVE" {
		return nil, errors.New("paket tidak tersedia untuk pembelian saat ini")
	}

	// Durasi 1 bulan harus selalu pakai harga dasar bulanan.
	if dto.DurationMonths == 1 {
		return s.repo.CreateOrder(ctx, userID, dto, pkg.Price)
	}

	// Ambil harga dari pricing tier yang sesuai
	// - Jika ada tier khusus (misal 12 bln = Rp 4.200.000) → pakai itu
	// - Jika tidak ada → fallback: harga_dasar × durasi_bulan
	price, err := s.repo.GetPricingTier(ctx, dto.PackageID, dto.DurationMonths)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil harga: %w", err)
	}
	if price == 0 {
		// Fallback: kalkulasi otomatis dari harga dasar
		price = pkg.Price * float64(dto.DurationMonths)
	}

	return s.repo.CreateOrder(ctx, userID, dto, price)
}

func (s *orderService) ListOrdersForClient(ctx context.Context, userID string) ([]ClientOrderListItem, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user tidak teridentifikasi, silakan login ulang")
	}
	return s.repo.ListOrdersByUser(ctx, userID)
}

func (s *orderService) GetOrderDetailForClient(ctx context.Context, userID, orderID string) (*ClientOrderDetail, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user tidak teridentifikasi, silakan login ulang")
	}
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("id pesanan wajib diisi")
	}

	rec, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrOrderNotFound
	}
	if rec.UserID != userID {
		return nil, ErrOrderForbidden
	}

	return &ClientOrderDetail{
		OrderID:                rec.OrderID,
		InvoiceNumber:          rec.InvoiceNumber,
		PackageName:            rec.PackageName,
		TotalPrice:             rec.TotalPrice,
		PaymentMethod:          rec.PaymentMethod,
		Status:                 rec.Status,
		VerificationNote:       rec.VerificationNote,
		HasPaymentProof:        rec.HasPaymentProof,
		PaymentProofUploadedAt: rec.PaymentProofUploadedAt,
		PaymentProofURL: func() string {
			if !rec.HasPaymentProof {
				return ""
			}
			return "/api/orders/" + rec.OrderID + "/payment-proof"
		}(),
		CreatedAt: rec.CreatedAt,
	}, nil
}

func (s *orderService) UploadPaymentProof(ctx context.Context, userID, orderID string, file PaymentProofFile) (*UploadPaymentProofResult, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user tidak teridentifikasi, silakan login ulang")
	}
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("id pesanan wajib diisi")
	}
	if len(file.Data) == 0 {
		return nil, errors.New("file bukti transfer wajib diunggah")
	}
	if len(file.Data) > 5*1024*1024 {
		return nil, errors.New("ukuran file bukti transfer maksimal 5MB")
	}

	file.ContentType = strings.ToLower(strings.TrimSpace(file.ContentType))
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedContentTypes[file.ContentType] {
		return nil, errors.New("format file tidak didukung, gunakan JPG, PNG, atau WEBP")
	}

	rec, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrOrderNotFound
	}
	if rec.UserID != userID {
		return nil, ErrOrderForbidden
	}
	if rec.Status != "PENDING_PAYMENT" {
		return nil, errors.New("bukti transfer hanya dapat diunggah saat status pesanan PENDING_PAYMENT")
	}

	res, err := s.repo.SavePaymentProof(ctx, orderID, file)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return res, nil
}

func (s *orderService) GetPaymentProofForClient(ctx context.Context, userID, orderID string) (*PaymentProofFile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user tidak teridentifikasi, silakan login ulang")
	}
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("id pesanan wajib diisi")
	}

	rec, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrOrderNotFound
	}
	if rec.UserID != userID {
		return nil, ErrOrderForbidden
	}
	if !rec.HasPaymentProof {
		return nil, ErrPaymentProofMissing
	}

	proof, err := s.repo.GetPaymentProof(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if proof == nil {
		return nil, ErrPaymentProofMissing
	}

	return proof, nil
}

func (s *orderService) ListOrdersForAdmin(ctx context.Context) ([]AdminOrderListItem, error) {
	return s.repo.ListOrdersForAdmin(ctx)
}

func (s *orderService) GetOrderDetailForAdmin(ctx context.Context, orderID string) (*AdminOrderDetail, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("id pesanan wajib diisi")
	}

	rec, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrOrderNotFound
	}

	return &AdminOrderDetail{
		OrderID:                rec.OrderID,
		InvoiceNumber:          rec.InvoiceNumber,
		UserID:                 rec.UserID,
		ClientName:             rec.ClientName,
		ClientEmail:            rec.ClientEmail,
		PackageID:              rec.PackageID,
		PackageName:            rec.PackageName,
		DurationMonths:         rec.DurationMonths,
		TotalPrice:             rec.TotalPrice,
		PaymentMethod:          rec.PaymentMethod,
		Status:                 rec.Status,
		VerificationNote:       rec.VerificationNote,
		HasPaymentProof:        rec.HasPaymentProof,
		PaymentProofUploadedAt: rec.PaymentProofUploadedAt,
		PaymentProofURL: func() string {
			if !rec.HasPaymentProof {
				return ""
			}
			return "/api/admin/orders/" + rec.OrderID + "/payment-proof"
		}(),
		CreatedAt: rec.CreatedAt,
	}, nil
}

func (s *orderService) GetPaymentProofForAdmin(ctx context.Context, orderID string) (*PaymentProofFile, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("id pesanan wajib diisi")
	}

	rec, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrOrderNotFound
	}
	if !rec.HasPaymentProof {
		return nil, ErrPaymentProofMissing
	}

	proof, err := s.repo.GetPaymentProof(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if proof == nil {
		return nil, ErrPaymentProofMissing
	}

	return proof, nil
}

func (s *orderService) VerifyOrder(ctx context.Context, orderID, action, rejectReason string) (*VerifyResult, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("id pesanan wajib diisi")
	}

	action = strings.ToUpper(strings.TrimSpace(action))
	if action != "APPROVE" && action != "REJECT" {
		return nil, errors.New("action tidak valid, gunakan APPROVE atau REJECT")
	}

	rejectReason = strings.TrimSpace(rejectReason)
	if action == "REJECT" && rejectReason == "" {
		return nil, errors.New("alasan reject wajib diisi")
	}

	rec, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrOrderNotFound
	}
	if rec.Status != "PENDING_PAYMENT" {
		return nil, errors.New("aksi verifikasi hanya dapat diproses jika status pesanan PENDING_PAYMENT")
	}
	if !rec.HasPaymentProof {
		return nil, errors.New("bukti transfer belum diunggah oleh client")
	}

	newStatus := "CANCELLED"
	if action == "APPROVE" {
		newStatus = "PAID"
	}

	verificationNote := ""
	if action == "REJECT" {
		verificationNote = rejectReason
	}

	if err := s.repo.UpdateOrderStatus(ctx, orderID, newStatus, verificationNote); err != nil {
		return nil, err
	}

	if action == "APPROVE" {
		if _, err := s.repo.CreateSubscriptionFromOrder(ctx, orderID); err != nil {
			// Best-effort rollback supaya status order konsisten jika auto-activate gagal.
			_ = s.repo.UpdateOrderStatus(ctx, orderID, "PENDING_PAYMENT", "")
			return nil, fmt.Errorf("gagal aktivasi subscription otomatis: %w", err)
		}
	}

	return &VerifyResult{
		Message:          "Status pembayaran berhasil diperbarui",
		OrderID:          orderID,
		NewStatus:        newStatus,
		VerificationNote: verificationNote,
	}, nil
}

func (s *orderService) ActivateOrderSystem(ctx context.Context, orderID string) (*ActivationResult, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("id pesanan wajib diisi")
	}

	activation, err := s.repo.CreateSubscriptionFromOrder(ctx, orderID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	activation.Message = "Subscription berhasil diaktifkan"
	return activation, nil
}

func (s *orderService) GetActiveSubscriptions(ctx context.Context, userID string) ([]SubscriptionStatus, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user tidak teridentifikasi, silakan login ulang")
	}
	return s.repo.ListActiveSubscriptionsByUser(ctx, userID)
}

func (s *orderService) GetActiveSubscription(ctx context.Context, userID string) (*SubscriptionStatus, error) {
	list, err := s.GetActiveSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// ProcessCreateOrderJob — ZaFramework concurrency worker processor
func (s *orderService) ProcessCreateOrderJob(ctx context.Context, payload interface{}) (interface{}, error) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid payload type for CreateOrderJob")
	}
	userID, _ := data["user_id"].(string)
	dto, ok := data["dto"].(CreateOrderDTO)
	if !ok {
		return nil, fmt.Errorf("invalid dto type for CreateOrderJob")
	}
	return s.CreateOrder(ctx, userID, dto)
}
