package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrate membuat table-table yang dibutuhkan bila belum ada.
func Migrate() {
	// Table untuk broadcast notifications (digunakan oleh ops dashboard)
	_, err := db().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS notifications (
			id          VARCHAR(36) PRIMARY KEY,
			title       TEXT NOT NULL,
			message     TEXT NOT NULL,
			type        VARCHAR(50) NOT NULL DEFAULT 'info',
			target_role VARCHAR(50) NOT NULL DEFAULT 'all',
			is_active   BOOLEAN NOT NULL DEFAULT TRUE,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by  VARCHAR(255),
			updated_by  VARCHAR(255)
		);
		CREATE INDEX IF NOT EXISTS idx_notifications_active ON notifications(is_active, target_role);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notifications: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notifications ready")
	}

	// Table untuk notification templates (mapping event_type → template konten)
	_, err = db().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS notification_templates (
			id          VARCHAR(36) PRIMARY KEY,
			name        TEXT NOT NULL,
			event_type  VARCHAR(100) NOT NULL,
			channel     VARCHAR(50) NOT NULL DEFAULT 'email',
			subject     TEXT,
			content     TEXT NOT NULL,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by  VARCHAR(255),
			updated_by  VARCHAR(255),
			UNIQUE (event_type, channel)
		);
		CREATE INDEX IF NOT EXISTS idx_notif_tpl_channel ON notification_templates(channel);
		CREATE INDEX IF NOT EXISTS idx_notif_tpl_event   ON notification_templates(event_type);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notification_templates: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notification_templates ready")
	}

	// Table untuk registry event_type yang pernah di-dispatch oleh service manapun.
	// Diisi otomatis saat POST /api/notifications/send dipanggil.
	_, err = db().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS notification_event_types (
			event_type  VARCHAR(100) PRIMARY KEY,
			description TEXT,
			registered_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notification_event_types: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notification_event_types ready")
	}

	// Table untuk log setiap pengiriman notifikasi (sent | failed | pending).
	// Dipakai untuk monitoring dan retry mechanism.
	_, err = db().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS notification_logs (
			id            VARCHAR(36) PRIMARY KEY,
			event_type    VARCHAR(100) NOT NULL,
			channel       VARCHAR(50) NOT NULL,
			to_address    TEXT NOT NULL,
			subject       TEXT,
			content       TEXT,
			status        VARCHAR(20) NOT NULL DEFAULT 'pending',
			retry_count   INT NOT NULL DEFAULT 0,
			max_retries   INT NOT NULL DEFAULT 3,
			next_retry_at TIMESTAMPTZ,
			sent_at       TIMESTAMPTZ,
			error_msg     TEXT,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_notif_logs_status ON notification_logs(status);
		CREATE INDEX IF NOT EXISTS idx_notif_logs_retry  ON notification_logs(status, next_retry_at);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notification_logs: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notification_logs ready")
	}
}

// Seed adalah no-op. Template diisi manual via API /api/notification-templates.
// Event types di-seed dengan event yang sudah pasti, yaitu dari register dan kyc service.
func Seed() {
	knownTypes := []string{"otp_verification", "kyc_approved", "kyc_rejected"}
	for _, et := range knownTypes {
		db().Exec(context.Background(),
			`INSERT INTO notification_event_types (event_type) VALUES ($1) ON CONFLICT DO NOTHING`, et)
	}
}

// db adalah shorthand internal agar tidak perlu akses global DB langsung di migrate.
func db() *pgxpool.Pool {
	return DB
}
