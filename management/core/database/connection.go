package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func NewPoolFromEnv(ctx context.Context) (*pgxpool.Pool, error) {
	dbUser := GetEnv("DB_USER", GetEnv("read_db_user", "postgres"))
	dbPassword := GetEnv("DB_PASSWORD", GetEnv("read_db_pass", "postgres"))
	dbHost := GetEnv("DB_HOST", GetEnv("read_db_host", "localhost"))
	dbPort := GetEnv("DB_PORT", GetEnv("read_db_port", "5432"))
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

func GetServicePort(primaryKey, fallback string) string {
	if v := GetEnv(primaryKey, ""); v != "" {
		return v
	}
	if v := GetEnv("PORT", ""); v != "" {
		return v
	}
	return fallback
}
