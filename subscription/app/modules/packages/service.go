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

	// Worker Job Processors
	ProcessCreatePackageJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessGetAdminPackagesJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessGetCatalogJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessUpdatePackageJob(ctx context.Context, payload interface{}) (interface{}, error)
	ProcessDeletePackageJob(ctx context.Context, payload interface{}) (interface{}, error)
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
	// Validasi business logic
	if payload.Name == "" {
		return nil, errors.New("nama paket tidak boleh kosong")
	}
	if payload.Price <= 0 {
		return nil, errors.New("harga harus lebih besar dari 0")
	}
	if payload.Duration <= 0 {
		return nil, errors.New("durasi harus lebih besar dari 0 bulan")
	}
	if payload.Quota <= 0 {
		return nil, errors.New("kuota harus lebih besar dari 0")
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
		return nil, errors.New("harga harus lebih besar dari 0")
	}
	if payload.Duration <= 0 {
		return nil, errors.New("durasi harus lebih besar dari 0 bulan")
	}
	if payload.Quota <= 0 {
		return nil, errors.New("kuota harus lebih besar dari 0")
	}

	// Cek apakah data eksis sebelum update
	existing, err := s.repo.GetPackageByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error validasi paket: %v", err)
	}
	if existing == nil {
		return nil, errors.New("paket tidak ditemukan atau sudah dihapus")
	}

	return s.repo.UpdatePackage(ctx, id, payload)
}

// PBI-35: Delete Package (Soft Delete)
func (s *packageService) DeletePackage(ctx context.Context, id string) error {
	// Pengecekan dilakukan langsung di Repo affected rows
	err := s.repo.DeletePackage(ctx, id)
	if err != nil {
		if err.Error() == "paket tidak ditemukan" {
			return errors.New("paket tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}
