package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func postgresOptional() bool {
	return strings.EqualFold(strings.TrimSpace(GetEnv("POSTGRES_OPTIONAL", GetEnv("postgres_optional", ""))), "true")
}

// NewPoolFromEnv opens a PostgreSQL pool using DB_* / read_db_* env vars (same as docker-compose defaults when unset).
func NewPoolFromEnv(ctx context.Context) (*pgxpool.Pool, error) {
	dbUser := GetEnv("DB_USER", GetEnv("read_db_user", "postgres"))
	dbPassword := GetEnv("DB_PASSWORD", GetEnv("read_db_pass", "postgres"))
	dbHost := GetEnv("DB_HOST", GetEnv("read_db_host", "localhost"))
	// Host port 5433 matches docker-compose.yml (postgres:15 mapped to host 5433).
	dbPort := GetEnv("DB_PORT", GetEnv("read_db_port", "5433"))
	dbName := GetEnv("DB_NAME", GetEnv("read_db_name", "thinknalyze"))
	dbSSLMode := GetEnv("DB_SSLMODE", GetEnv("read_db_ssl_mode", "disable"))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

// NewPoolFromEnvAllowNil returns (nil, nil) when connection fails and POSTGRES_OPTIONAL=true; otherwise same errors as NewPoolFromEnv.
func NewPoolFromEnvAllowNil(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := NewPoolFromEnv(ctx)
	if err != nil && postgresOptional() {
		log.Printf("[management/database] PostgreSQL unavailable (postgres_optional=true): %v", err)
		return nil, nil
	}
	return pool, err
}

func GetServicePort(primaryKey, fallback string) string {
	if v := GetEnv(primaryKey, ""); v != "" {
		return v
	}
	if v := GetEnv("PORT", ""); v != "" {
		return v
	}
	return fallback
}
