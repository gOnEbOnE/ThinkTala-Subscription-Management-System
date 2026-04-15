package orders

import (
	"context"
	"errors"
	"fmt"
)

// ==========================================
// SERVICE INTERFACE
// ==========================================

type Service interface {
	CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO) (*Order, error)

	// Worker Job Processors
	ProcessCreateOrderJob(ctx context.Context, payload interface{}) (interface{}, error)
}

// ==========================================
// SERVICE IMPLEMENTATION
// ==========================================

type orderService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &orderService{repo: repo}
}

// PBI-45: CreateOrder — validasi + buat pesanan baru
func (s *orderService) CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO) (*Order, error) {
	// Validasi input wajib (PBI-45 AC: request body wajib package_id & payment_method)
	if dto.PackageID == "" {
		return nil, errors.New("package_id wajib diisi")
	}
	if dto.PaymentMethod == "" {
		return nil, errors.New("payment_method wajib diisi")
	}
	if userID == "" {
		return nil, errors.New("user tidak teridentifikasi, silakan login ulang")
	}

	// Validasi: paket harus valid dan berstatus ACTIVE (PBI-45 AC)
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

	// Buat order — total_price otomatis diambil dari data paket (PBI-45 AC)
	return s.repo.CreateOrder(ctx, userID, dto, pkg)
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
