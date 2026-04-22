package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func firstEnv(keys ...string) string {
	for _, key := range keys {
		v := strings.TrimSpace(os.Getenv(key))
		v = strings.Trim(v, "\"")
		if v != "" {
			return v
		}
	}
	return ""
}

func OpenDB() (*sql.DB, error) {
	host := firstEnv("DB_HOST", "read_db_host", "PGHOST")
	if host == "" {
		host = "localhost"
	}
	port := firstEnv("DB_PORT", "read_db_port", "PGPORT")
	if port == "" {
		port = "5432"
	}
	user := firstEnv("DB_USER", "read_db_user", "PGUSER")
	if user == "" {
		user = "postgres"
	}
	pass := firstEnv("DB_PASSWORD", "DB_PASS", "read_db_pass", "PGPASSWORD")
	if pass == "" {
		pass = "postgres"
	}
	name := firstEnv("DB_NAME", "read_db_name", "PGDATABASE")
	if name == "" {
		name = "postgres"
	}
	sslMode := firstEnv("DB_SSLMODE", "read_db_ssl_mode", "PGSSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}
	timeZone := firstEnv("APP_TIMEZONE", "read_db_timezone")
	if timeZone == "" {
		timeZone = "Asia/Jakarta"
	}

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
