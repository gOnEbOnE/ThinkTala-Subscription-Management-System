package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds database configuration compatible with pgxpool
type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string // disable, require, verify-full
	TimeZone string

	// Pool Configuration
	MaxConns        int32         // Total maksimal koneksi (pengganti MaxOpenConns)
	MinConns        int32         // Koneksi standby (pengganti MaxIdleConns)
	MaxConnLifetime time.Duration // Berapa lama koneksi hidup sebelum direcycle
	MaxConnIdleTime time.Duration // Berapa lama koneksi nganggur sebelum diputus

	// TAMBAHAN: Field ini wajib ada karena dipanggil di main.go
	LogLevel string
}

// DBWrapper wraps pgxpool.Pool to provide helper methods
type DBWrapper struct {
	Pool *pgxpool.Pool
}

// Connect initializes a pgx database connection pool
func Connect(cfg Config) (*DBWrapper, error) {
	// 1. Build Connection String (DSN)
	// Format: postgres://user:password@host:port/dbname?sslmode=disable&pool_max_conns=...
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// 2. Parse Config dari DSN untuk mendapatkan default pgx config
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	// 3. Apply Tuning Pool (Penting buat High Traffic!)
	pgxConfig.MaxConns = cfg.MaxConns
	pgxConfig.MinConns = cfg.MinConns
	pgxConfig.MaxConnLifetime = cfg.MaxConnLifetime
	pgxConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	// Optional: Custom Logger bisa dipasang disini jika perlu debug query nanti
	// pgxConfig.ConnConfig.Tracer = ...

	// 4. Create Pool
	// Gunakan context background dengan timeout untuk inisialisasi awal
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// 5. Ping Test (Memastikan koneksi benar-benar hidup)
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database failed: %w", err)
	}

	return &DBWrapper{Pool: pool}, nil
}

// Close safely closes the connection pool
func (d *DBWrapper) Close() {
	// SAFETY CHECK: Kalau pool nil (bypass mode), jangan lakukan apa-apa
	if d == nil || d.Pool == nil {
		return
	}
	d.Pool.Close()
}

// Stats returns connection pool statistics formatted for health check
func (d *DBWrapper) Stats() map[string]any {
	// SAFETY CHECK: Kalau DB mati, return status disabled
	if d == nil || d.Pool == nil {
		return map[string]any{"status": "disabled (bypass mode)"}
	}

	stats := d.Pool.Stat()

	return map[string]any{
		"max_conns":           stats.MaxConns(),
		"total_conns":         stats.TotalConns(),    // Total koneksi (Idle + Active)
		"acquired_conns":      stats.AcquiredConns(), // Koneksi yang sedang dipakai worker (Active)
		"idle_conns":          stats.IdleConns(),     // Koneksi nganggur (siap pakai)
		"new_conns_count":     stats.NewConnsCount(), // Total koneksi yang pernah dibuat sejak start
		"max_lifetime_closed": stats.MaxLifetimeDestroyCount(),
		"max_idle_closed":     stats.MaxIdleDestroyCount(),
	}
}

// Helper untuk query row tunggal (Shortcut)
func (d *DBWrapper) QueryRow(ctx context.Context, sql string, args ...any) interface{ Scan(dest ...any) error } {
	// SAFETY CHECK: Handle jika DB mati agar tidak panic saat Scan()
	if d == nil || d.Pool == nil {
		return &dummyScanner{err: errors.New("database connection is disabled")}
	}
	return d.Pool.QueryRow(ctx, sql, args...)
}

// Helper untuk eksekusi tanpa return rows (INSERT/UPDATE/DELETE)
func (d *DBWrapper) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	// SAFETY CHECK: Return error jika DB mati
	if d == nil || d.Pool == nil {
		return 0, errors.New("database connection is disabled")
	}
	tag, err := d.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// =================================================================
// INTERNAL HELPER (Untuk Bypass Mode)
// =================================================================

// dummyScanner digunakan saat DB mati agar method Scan(...) tidak panic
type dummyScanner struct {
	err error
}

func (ds *dummyScanner) Scan(dest ...any) error {
	return ds.err
}
