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
	GetPackageByName(ctx context.Context, name string) (*Package, error)
	UpdatePackage(ctx context.Context, id string, data UpdatePackageDTO) (*Package, error)
	DeletePackage(ctx context.Context, id string) error
	TogglePackageStatus(ctx context.Context, id string, newStatus string) (*Package, error)
	CountActiveSubscribers(ctx context.Context, packageID string) (int, error)
	// Pricing helpers (dipakai oleh orders module)
	GetPricingTier(ctx context.Context, packageID string, durationMonths int) (*PricingTier, error)
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

// ==========================================
// INTERNAL HELPERS
// ==========================================

// loadPricingTiers mengambil semua tier harga untuk sebuah paket
func (r *packageRepo) loadPricingTiers(ctx context.Context, packageID string) ([]PricingTier, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, duration_months, price, label
		 FROM subscription.package_pricing
		 WHERE package_id = $1
		 ORDER BY duration_months ASC`,
		packageID,
	)
	if err != nil {
		return []PricingTier{}, nil // tabel mungkin belum ada
	}
	defer rows.Close()

	var tiers []PricingTier
	for rows.Next() {
		var t PricingTier
		if err := rows.Scan(&t.ID, &t.DurationMonths, &t.Price, &t.Label); err != nil {
			continue
		}
		tiers = append(tiers, t)
	}
	if tiers == nil {
		return []PricingTier{}, nil
	}
	return tiers, nil
}

// savePricingTiers menghapus tier lama lalu insert yang baru (replace-all)
func (r *packageRepo) savePricingTiers(ctx context.Context, packageID string, tiers []PricingTierDTO) error {
	// Delete existing
	_, err := r.db.Pool.Exec(ctx,
		`DELETE FROM subscription.package_pricing WHERE package_id = $1`, packageID)
	if err != nil {
		return fmt.Errorf("gagal hapus pricing tiers lama: %w", err)
	}

	// Insert new
	for _, t := range tiers {
		_, err := r.db.Pool.Exec(ctx,
			`INSERT INTO subscription.package_pricing (id, package_id, duration_months, price, label)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (package_id, duration_months) DO UPDATE SET price = $4, label = $5`,
			uuid.New().String(), packageID, t.DurationMonths, t.Price, t.Label,
		)
		if err != nil {
			return fmt.Errorf("gagal insert pricing tier %d bulan: %w", t.DurationMonths, err)
		}
	}
	return nil
}

// ==========================================
// CRUD
// ==========================================

// CreatePackage membuat paket baru beserta pricing tiers-nya
func (r *packageRepo) CreatePackage(ctx context.Context, data CreatePackageDTO) (*Package, error) {
	id := uuid.New().String()
	status := data.Status
	if status == "" {
		status = "ACTIVE"
	}

	var p Package
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO subscription.packages (id, name, price, quota, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		 RETURNING id, name, price, quota, status, created_at, updated_at`,
		id, data.Name, data.Price, data.Quota, status,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat paket: %w", err)
	}

	// Simpan pricing tiers
	if len(data.PricingTiers) > 0 {
		if err := r.savePricingTiers(ctx, p.ID, data.PricingTiers); err != nil {
			return nil, err
		}
	}

	p.PricingTiers, _ = r.loadPricingTiers(ctx, p.ID)
	return &p, nil
}

// GetPackages mengambil semua paket dengan pricing tiers
func (r *packageRepo) GetPackages(ctx context.Context, status, minPrice, maxPrice string) ([]Package, error) {
	query := `SELECT id, name, price, quota, status, created_at, updated_at FROM subscription.packages WHERE 1=1`
	var args []interface{}
	var where []string
	argN := 1

	if status != "" && status != "Semua" && status != "All" {
		where = append(where, fmt.Sprintf("status = $%d", argN))
		args = append(args, strings.ToUpper(status))
		argN++
	} else if status == "" {
		where = append(where, fmt.Sprintf("status != $%d", argN))
		args = append(args, "DELETED")
		argN++
	}
	if minPrice != "" {
		where = append(where, fmt.Sprintf("price >= $%d", argN))
		args = append(args, minPrice)
		argN++
	}
	if maxPrice != "" {
		where = append(where, fmt.Sprintf("price <= $%d", argN))
		args = append(args, maxPrice)
		argN++
	}
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "relation \"subscription.packages\" does not exist") {
			return []Package{}, nil
		}
		return nil, fmt.Errorf("gagal mengambil data paket: %w", err)
	}
	defer rows.Close()

	var pkgs []Package
	for rows.Next() {
		var p Package
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal scan data paket: %w", err)
		}
		p.PricingTiers, _ = r.loadPricingTiers(ctx, p.ID)
		pkgs = append(pkgs, p)
	}
	if pkgs == nil {
		return []Package{}, nil
	}
	return pkgs, nil
}

// GetPackageByID mengambil satu paket beserta pricing tiers
func (r *packageRepo) GetPackageByID(ctx context.Context, id string) (*Package, error) {
	var p Package
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, name, price, quota, status, created_at, updated_at
		 FROM subscription.packages WHERE id = $1 AND status != 'DELETED' LIMIT 1`, id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mencari paket: %w", err)
	}
	p.PricingTiers, _ = r.loadPricingTiers(ctx, p.ID)
	return &p, nil
}

// GetPackageByName mencari paket berdasarkan nama
func (r *packageRepo) GetPackageByName(ctx context.Context, name string) (*Package, error) {
	var p Package
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, name, price, quota, status, created_at, updated_at
		 FROM subscription.packages WHERE LOWER(name) = LOWER($1) AND status != 'DELETED' LIMIT 1`, name,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mencari paket by name: %w", err)
	}
	p.PricingTiers, _ = r.loadPricingTiers(ctx, p.ID)
	return &p, nil
}

// UpdatePackage update paket + replace pricing tiers
func (r *packageRepo) UpdatePackage(ctx context.Context, id string, data UpdatePackageDTO) (*Package, error) {
	var p Package
	err := r.db.Pool.QueryRow(ctx,
		`UPDATE subscription.packages
		 SET name = $1, price = $2, quota = $3,
		     status = COALESCE(NULLIF($4, ''), status), updated_at = NOW()
		 WHERE id = $5 AND status != 'DELETED'
		 RETURNING id, name, price, quota, status, created_at, updated_at`,
		data.Name, data.Price, data.Quota, data.Status, id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mengupdate paket: %w", err)
	}

	if len(data.PricingTiers) > 0 {
		if err := r.savePricingTiers(ctx, p.ID, data.PricingTiers); err != nil {
			return nil, err
		}
	}

	p.PricingTiers, _ = r.loadPricingTiers(ctx, p.ID)
	return &p, nil
}

// TogglePackageStatus ubah status ACTIVE <-> INACTIVE
func (r *packageRepo) TogglePackageStatus(ctx context.Context, id string, newStatus string) (*Package, error) {
	var p Package
	err := r.db.Pool.QueryRow(ctx,
		`UPDATE subscription.packages SET status = $1, updated_at = NOW()
		 WHERE id = $2 AND status != 'DELETED'
		 RETURNING id, name, price, quota, status, created_at, updated_at`,
		newStatus, id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mengubah status paket: %w", err)
	}
	p.PricingTiers, _ = r.loadPricingTiers(ctx, p.ID)
	return &p, nil
}

// CountActiveSubscribers hitung jumlah order aktif per paket
func (r *packageRepo) CountActiveSubscribers(ctx context.Context, packageID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM subscription.orders WHERE package_id = $1 AND status IN ('PAID','PENDING_PAYMENT')`,
		packageID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("gagal menghitung pelanggan aktif: %w", err)
	}
	return count, nil
}

// DeletePackage soft delete
func (r *packageRepo) DeletePackage(ctx context.Context, id string) error {
	cmd, err := r.db.Pool.Exec(ctx,
		`UPDATE subscription.packages SET status = 'DELETED', updated_at = NOW()
		 WHERE id = $1 AND status != 'DELETED'`, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus paket: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("paket tidak ditemukan")
	}
	return nil
}

// GetPricingTier mengambil satu tier harga berdasarkan package_id + durasi (dipakai orders)
func (r *packageRepo) GetPricingTier(ctx context.Context, packageID string, durationMonths int) (*PricingTier, error) {
	var t PricingTier
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, duration_months, price, label
		 FROM subscription.package_pricing
		 WHERE package_id = $1 AND duration_months = $2 LIMIT 1`,
		packageID, durationMonths,
	).Scan(&t.ID, &t.DurationMonths, &t.Price, &t.Label)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mengambil pricing tier: %w", err)
	}
	return &t, nil
}
