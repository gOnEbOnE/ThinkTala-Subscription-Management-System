package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tickets/core/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenDB() (*sql.DB, error) {
	host := config.EnvOrDefault("read_db_host", config.EnvOrDefault("DB_HOST", "localhost"))
	port := config.EnvOrDefault("read_db_port", config.EnvOrDefault("DB_PORT", "5432"))
	user := config.EnvOrDefault("read_db_user", config.EnvOrDefault("DB_USER", "postgres"))
	pass := config.EnvOrDefault("read_db_pass", config.EnvOrDefault("DB_PASS", "postgres"))
	name := config.EnvOrDefault("read_db_name", config.EnvOrDefault("DB_NAME", "postgres"))
	sslMode := config.EnvOrDefault("read_db_ssl_mode", "disable")
	timeZone := config.EnvOrDefault("read_db_timezone", "Asia/Jakarta")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, pass, name, sslMode, timeZone)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
