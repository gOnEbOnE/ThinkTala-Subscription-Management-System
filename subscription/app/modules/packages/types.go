package packages

import "time"

// ==========================================
// 1. DATA TRANSFER OBJECTS (DTO)
// Request & Response Payload Structure
// ==========================================

// CreatePackageDTO format request body untuk create
type CreatePackageDTO struct {
	Name     string  `json:"name" validate:"required,min=3"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Duration int     `json:"duration" validate:"required,gt=0"`
	Quota    int     `json:"quota" validate:"required,gt=0"`
}

// UpdatePackageDTO format request body untuk update
type UpdatePackageDTO struct {
	Name     string  `json:"name" validate:"required,min=3"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Duration int     `json:"duration" validate:"required,gt=0"`
	Quota    int     `json:"quota" validate:"required,gt=0"`
}

// ==========================================
// 2. MODEL STRUCTS
// Database Row Mapping Structure
// ==========================================

type Package struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Duration  int       `json:"duration"` // in months
	Quota     int       `json:"quota"`
	Status    string    `json:"status"` // ACTIVE, INACTIVE, DELETED
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
