package database

import (
	"context"
	"log"

	"github.com/google/uuid"
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
			updated_by  VARCHAR(255)
		);
		CREATE INDEX IF NOT EXISTS idx_notif_tpl_channel ON notification_templates(channel);
		CREATE INDEX IF NOT EXISTS idx_notif_tpl_event   ON notification_templates(event_type);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notification_templates: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notification_templates ready")
	}
}

// Seed mengisi dummy data awal bila table masih kosong.
func Seed() {
	var count int
	err := db().QueryRow(context.Background(), "SELECT COUNT(*) FROM notification_templates").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	templates := []struct {
		name, eventType, channel, subject, content string
	}{
		{
			"Welcome Email", "user_register", "email",
			"Selamat datang di ThinkTala!",
			"Halo {{name}},\n\nTerima kasih telah mendaftar di platform kami. Email Anda: {{email}}\n\nBest regards,\nThinkTala Team",
		},
		{
			"OTP Verification Email", "otp_verification", "email",
			"Kode Verifikasi OTP Anda",
			"Halo {{name}},\n\nKode OTP Anda: {{otp}}\nKode ini berlaku selama 10 menit.\n\nJangan bagikan kode ini kepada siapapun.",
		},
		{
			"KYC Approved Notification", "user_kyc_approved", "email",
			"Verifikasi KYC Berhasil",
			"Selamat {{name}},\n\nVerifikasi identitas Anda telah disetujui. Anda sekarang dapat mengakses semua fitur platform.",
		},
		{
			"KYC Rejected Notification", "user_kyc_rejected", "email",
			"Verifikasi KYC Ditolak",
			"Halo {{name}},\n\nMohon maaf, verifikasi identitas Anda ditolak karena: {{reason}}.\nSilakan upload ulang dokumen yang valid.",
		},
		{
			"Payment Success Notification", "payment_success", "email",
			"Pembayaran Berhasil",
			"Halo {{name}},\n\nPembayaran sebesar {{amount}} telah berhasil diproses. Terima kasih!",
		},
	}

	for _, t := range templates {
		var subj *string
		if t.subject != "" {
			subj = &t.subject
		}
		_, err := db().Exec(context.Background(), `
			INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6, 'system', 'system')
		`, uuid.New().String(), t.name, t.eventType, t.channel, subj, t.content)
		if err != nil {
			log.Printf("[WARN] seed template '%s': %v", t.name, err)
		}
	}
	log.Println("[NOTIFICATION] Seeded 7 dummy templates")
}

// db adalah shorthand internal agar tidak perlu akses global DB langsung di migrate.
func db() *pgxpool.Pool {
	return DB
}
