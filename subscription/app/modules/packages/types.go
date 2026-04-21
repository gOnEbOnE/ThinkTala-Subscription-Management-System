package packages

import "time"

// ==========================================
// 1. PRICING TIER
// ==========================================

// PricingTierDTO - input saat create/update paket
type PricingTierDTO struct {
	DurationMonths int     `json:"duration_months"` // 1, 3, 6, 12
	Price          float64 `json:"price"`           // total harga untuk durasi ini
	Label          string  `json:"label"`           // "Hemat 20%", kosong OK
}

// PricingTier - data dari DB
type PricingTier struct {
	ID             string  `json:"id"`
	DurationMonths int     `json:"duration_months"`
	Price          float64 `json:"price"`
	Label          string  `json:"label"`
}

// ==========================================
// 2. DATA TRANSFER OBJECTS (DTO)
// ==========================================

// CreatePackageDTO - request body untuk create paket
type CreatePackageDTO struct {
	Name         string           `json:"name" validate:"required,min=3"`
	Price        float64          `json:"price" validate:"required,gt=0"` // harga dasar / referensi (bulan 1)
	Quota        int              `json:"quota" validate:"required,gt=0"`
	Status       string           `json:"status"`
	PricingTiers []PricingTierDTO `json:"pricing_tiers"` // minimal harus ada 1
}

// UpdatePackageDTO - request body untuk update paket
type UpdatePackageDTO struct {
	Name         string           `json:"name" validate:"required,min=3"`
	Price        float64          `json:"price" validate:"required,gt=0"`
	Quota        int              `json:"quota" validate:"required,gt=0"`
	Status       string           `json:"status"`
	PricingTiers []PricingTierDTO `json:"pricing_tiers"`
}

// ==========================================
// 3. MODEL STRUCT
// ==========================================

// Package - mapping dari tabel subscription.packages (+ pricing tiers)
type Package struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Price        float64       `json:"price"` // harga dasar / bulan (referensi)
	Quota        int           `json:"quota"`
	Status       string        `json:"status"`        // ACTIVE, INACTIVE, DELETED
	PricingTiers []PricingTier `json:"pricing_tiers"` // semua tier harga
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
