package packages

import (
	"context"
	"errors"
	"fmt"
)

// ==========================================
// 1. SERVICE INTERFACE
// ==========================================

type Service interface {
	CreatePackage(ctx context.Context, payload CreatePackageDTO) (*Package, error)
	GetAdminPackages(ctx context.Context, status, minPrice, maxPrice string) ([]Package, error)
	GetCatalogPackages(ctx context.Context) ([]Package, error)
	UpdatePackage(ctx context.Context, id string, payload UpdatePackageDTO) (*Package, error)
	DeletePackage(ctx context.Context, id string) error
	TogglePackageStatus(ctx context.Context, id string) (*Package, error)

	// Worker Job Processors
	ProcessCreatePackageJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessGetAdminPackagesJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessGetCatalogJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessUpdatePackageJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessDeletePackageJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessTogglePackageStatusJob(ctx context.Context, payload interface{}) (interface{}, error)
}

// ==========================================
// 2. SERVICE IMPLEMENTATION
// ==========================================

type packageService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &packageService{repo: repo}
}

// PBI-32: Create Package
func (s *packageService) CreatePackage(ctx context.Context, payload CreatePackageDTO) (*Package, error) {
	if payload.Name == "" {
		return nil, errors.New("nama paket tidak boleh kosong")
	}
	if payload.Price <= 0 {
		return nil, errors.New("harga dasar harus lebih besar dari 0")
	}
	if payload.Quota <= 0 {
		return nil, errors.New("kuota harus lebih besar dari 0")
	}
	if len(payload.PricingTiers) == 0 {
		return nil, errors.New("minimal harus ada 1 pilihan durasi harga")
	}
	for _, t := range payload.PricingTiers {
		if t.DurationMonths <= 0 {
			return nil, errors.New("durasi harga harus lebih besar dari 0 bulan")
		}
		if t.Price < 0 {
			return nil, errors.New("harga tidak boleh negatif")
		}
	}

	existing, err := s.repo.GetPackageByName(ctx, payload.Name)
	if err != nil {
		return nil, fmt.Errorf("error validasi nama paket: %v", err)
	}
	if existing != nil {
		return nil, errors.New("nama paket sudah digunakan, gunakan nama lain")
	}

	return s.repo.CreatePackage(ctx, payload)
}

// PBI-33, PBI-36: Get Packages for Admin with Filtering
func (s *packageService) GetAdminPackages(ctx context.Context, status, minPrice, maxPrice string) ([]Package, error) {
	return s.repo.GetPackages(ctx, status, minPrice, maxPrice)
}

// PBI-37: Get Catalog Packages for Client (ACTIVE only)
func (s *packageService) GetCatalogPackages(ctx context.Context) ([]Package, error) {
	return s.repo.GetPackages(ctx, "ACTIVE", "", "")
}

// PBI-34: Update Package
func (s *packageService) UpdatePackage(ctx context.Context, id string, payload UpdatePackageDTO) (*Package, error) {
	if payload.Name == "" {
		return nil, errors.New("nama paket tidak boleh kosong")
	}
	if payload.Price <= 0 {
		return nil, errors.New("harga dasar harus lebih besar dari 0")
	}
	if payload.Quota <= 0 {
		return nil, errors.New("kuota harus lebih besar dari 0")
	}
	if len(payload.PricingTiers) == 0 {
		return nil, errors.New("minimal harus ada 1 pilihan durasi harga")
	}
	for _, t := range payload.PricingTiers {
		if t.DurationMonths <= 0 {
			return nil, errors.New("durasi harga harus lebih besar dari 0 bulan")
		}
		if t.Price < 0 {
			return nil, errors.New("harga tidak boleh negatif")
		}
	}

	existing, err := s.repo.GetPackageByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error validasi paket: %v", err)
	}
	if existing == nil {
		return nil, errors.New("paket tidak ditemukan atau sudah dihapus")
	}

	duplicateName, err := s.repo.GetPackageByName(ctx, payload.Name)
	if err != nil {
		return nil, fmt.Errorf("error validasi nama paket: %v", err)
	}
	if duplicateName != nil && duplicateName.ID != id {
		return nil, errors.New("nama paket sudah digunakan, gunakan nama lain")
	}

	return s.repo.UpdatePackage(ctx, id, payload)
}

// Toggle Package Status: ACTIVE <-> INACTIVE
func (s *packageService) TogglePackageStatus(ctx context.Context, id string) (*Package, error) {
	existing, err := s.repo.GetPackageByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error validasi paket: %v", err)
	}
	if existing == nil {
		return nil, errors.New("paket tidak ditemukan atau sudah dihapus")
	}

	newStatus := "INACTIVE"
	if existing.Status == "INACTIVE" {
		newStatus = "ACTIVE"
	}

	return s.repo.TogglePackageStatus(ctx, id, newStatus)
}

// PBI-35: Delete Package (Soft Delete)
func (s *packageService) DeletePackage(ctx context.Context, id string) error {
	existing, err := s.repo.GetPackageByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error validasi paket: %v", err)
	}
	if existing == nil {
		return errors.New("paket tidak ditemukan atau sudah dihapus")
	}
	if existing.Status == "ACTIVE" {
		return errors.New("tidak dapat menghapus paket yang sedang aktif")
	}

	// PBI-35: Cek apakah ada pelanggan aktif (PAID/PENDING) yang menggunakan paket ini
	subscriberCount, err := s.repo.CountActiveSubscribers(ctx, id)
	if err != nil {
		return fmt.Errorf("error mengecek pelanggan: %v", err)
	}
	if subscriberCount > 0 {
		return errors.New("tidak dapat menghapus paket yang masih memiliki pelanggan aktif")
	}

	// Pengecekan dilakukan langsung di Repo affected rows
	err = s.repo.DeletePackage(ctx, id)
	if err != nil {
		if err.Error() == "paket tidak ditemukan" {
			return errors.New("paket tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}
