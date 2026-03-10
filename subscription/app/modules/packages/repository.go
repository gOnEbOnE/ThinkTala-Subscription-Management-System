package packages

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

// ==========================================
// REPOSITORY INTERFACE
// ==========================================

type Repository interface {
	CreatePackage(ctx context.Context, data CreatePackageDTO) (*Package, error)
	GetPackages(ctx context.Context, status, minPrice, maxPrice string) ([]Package, error)
	GetPackageByID(ctx context.Context, id string) (*Package, error)
	UpdatePackage(ctx context.Context, id string, data UpdatePackageDTO) (*Package, error)
	DeletePackage(ctx context.Context, id string) error
	TogglePackageStatus(ctx context.Context, id string, newStatus string) (*Package, error)
}

// ==========================================
// REPOSITORY IMPLEMENTATION
// ==========================================

type packageRepo struct {
	db *database.DBWrapper
}

func NewRepository(db *database.DBWrapper) Repository {
	return &packageRepo{db: db}
}

// CreatePackage membuat data paket baru di database
func (r *packageRepo) CreatePackage(ctx context.Context, data CreatePackageDTO) (*Package, error) {
	id := uuid.New().String()

	status := data.Status
	if status == "" {
		status = "ACTIVE"
	}

	query := `
        INSERT INTO subscription.packages (id, name, price, price_yearly, duration, quota, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
        RETURNING id, name, price, price_yearly, duration, quota, status, created_at, updated_at
    `

	var p Package
	err := r.db.Pool.QueryRow(ctx, query, id, data.Name, data.Price, data.PriceYearly, data.Duration, data.Quota, status).Scan(
		&p.ID, &p.Name, &p.Price, &p.PriceYearly, &p.Duration, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal membuat paket: %w", err)
	}

	return &p, nil
}

// GetPackages mengambil data paket dengan filter opsional (status, minPrice, maxPrice)
func (r *packageRepo) GetPackages(ctx context.Context, status, minPrice, maxPrice string) ([]Package, error) {
	var packages []Package

	// Base query
	query := `
        SELECT id, name, price, price_yearly, duration, quota, status, created_at, updated_at 
        FROM subscription.packages 
        WHERE 1=1
    `

	// Dynamic Filters
	var args []interface{}
	var whereClauses []string
	argID := 1

	// Filter by Status (Biasanya 'ACTIVE' untuk client, atau 'Semua' untuk admin)
	if status != "" && status != "Semua" && status != "All" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argID))
		args = append(args, strings.ToUpper(status))
		argID++
	} else if status == "" {
		// Default jika tidak ada filter dari admin, jangan tampilkan yang DELETED
		whereClauses = append(whereClauses, fmt.Sprintf("status != $%d", argID))
		args = append(args, "DELETED")
		argID++
	}

	// Filter by Min Price
	if minPrice != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("price >= $%d", argID))
		args = append(args, minPrice)
		argID++
	}

	// Filter by Max Price
	if maxPrice != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argID))
		args = append(args, maxPrice)
		argID++
	}

	// Append WHERE clauses if any
	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Order result
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		// Kembalikan empty array jika tabel tidak ada (belum di migrate misalnya)
		if strings.Contains(err.Error(), "relation \"subscription.packages\" does not exist") {
			return []Package{}, nil
		}
		return nil, fmt.Errorf("gagal mengambil data paket: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Package
		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.PriceYearly, &p.Duration, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data paket: %w", err)
		}
		packages = append(packages, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Jika kosong, pastikan kembalikan array kosong (bukan nil)
	if packages == nil {
		return []Package{}, nil
	}

	return packages, nil
}

// GetPackageByID mencari paket berdasarkan ID spesifik
func (r *packageRepo) GetPackageByID(ctx context.Context, id string) (*Package, error) {
	query := `
        SELECT id, name, price, price_yearly, duration, quota, status, created_at, updated_at 
        FROM subscription.packages 
        WHERE id = $1 AND status != 'DELETED' 
        LIMIT 1
    `

	var p Package
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.PriceYearly, &p.Duration, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not Found
		}
		return nil, fmt.Errorf("gagal mencari paket: %w", err)
	}

	return &p, nil
}

// UpdatePackage meng-update data paket
func (r *packageRepo) UpdatePackage(ctx context.Context, id string, data UpdatePackageDTO) (*Package, error) {
	query := `
        UPDATE subscription.packages 
        SET 
            name = $1, 
            price = $2, 
            price_yearly = $3,
            duration = $4, 
            quota = $5,
            status = COALESCE(NULLIF($6, ''), status),
            updated_at = NOW() 
        WHERE id = $7 AND status != 'DELETED'
        RETURNING id, name, price, price_yearly, duration, quota, status, created_at, updated_at
    `

	var p Package
	err := r.db.Pool.QueryRow(ctx, query, data.Name, data.Price, data.PriceYearly, data.Duration, data.Quota, data.Status, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.PriceYearly, &p.Duration, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // ID not found or already deleted
		}
		return nil, fmt.Errorf("gagal mengupdate paket: %w", err)
	}

	return &p, nil
}

// TogglePackageStatus mengubah status paket (ACTIVE <-> INACTIVE)
func (r *packageRepo) TogglePackageStatus(ctx context.Context, id string, newStatus string) (*Package, error) {
	query := `
        UPDATE subscription.packages
        SET status = $1, updated_at = NOW()
        WHERE id = $2 AND status != 'DELETED'
        RETURNING id, name, price, price_yearly, duration, quota, status, created_at, updated_at
    `
	var p Package
	err := r.db.Pool.QueryRow(ctx, query, newStatus, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.PriceYearly, &p.Duration, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mengubah status paket: %w", err)
	}
	return &p, nil
}

// DeletePackage melakukan *soft delete* pada data paket
func (r *packageRepo) DeletePackage(ctx context.Context, id string) error {
	query := `
        UPDATE subscription.packages 
        SET 
            status = 'DELETED',
            updated_at = NOW()
        WHERE id = $1 AND status != 'DELETED'
    `

	cmdTag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus paket: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("paket tidak ditemukan")
	}

	return nil
}
