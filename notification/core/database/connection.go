package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the shared database connection pool for the notification service.
var DB *pgxpool.Pool

// Init membuka koneksi ke PostgreSQL menggunakan env variables.
func Init() error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		GetEnv("DB_USER", "postgres"),
		GetEnv("DB_PASSWORD", "postgres"),
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "5432"),
		GetEnv("DB_NAME", "thinknalyze"),
		GetEnv("DB_SSLMODE", "disable"),
	)

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("gagal koneksi DB: %w", err)
	}

	log.Println("[NOTIFICATION] Connected to database")
	return nil
}

// GetEnv mengambil nilai env variable, jika kosong kembalikan fallback.
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
