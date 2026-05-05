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
			description TEXT NOT NULL DEFAULT '',
			message     TEXT NOT NULL DEFAULT '',
			type        VARCHAR(50) NOT NULL DEFAULT 'info',
			target_role VARCHAR(50) NOT NULL DEFAULT 'all',
			cta_url     TEXT,
			image_url   TEXT,
			expiry_date TIMESTAMPTZ,
			is_active   BOOLEAN NOT NULL DEFAULT TRUE,
			is_pinned   BOOLEAN NOT NULL DEFAULT FALSE,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by  VARCHAR(255),
			updated_by  VARCHAR(255)
		);

		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS description TEXT;
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS cta_url TEXT;
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS image_url TEXT;
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS expiry_date TIMESTAMPTZ;
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS is_pinned BOOLEAN NOT NULL DEFAULT FALSE;

		UPDATE notifications
		SET description = COALESCE(NULLIF(description, ''), message)
		WHERE description IS NULL OR description = '';

		UPDATE notifications
		SET message = COALESCE(NULLIF(message, ''), description, '')
		WHERE message IS NULL OR message = '';

		CREATE INDEX IF NOT EXISTS idx_notifications_active ON notifications(is_active, target_role);
		CREATE INDEX IF NOT EXISTS idx_notifications_pinned ON notifications(is_pinned, created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
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
// Event types di-seed dengan event yang sudah pasti, yaitu dari register, kyc, dan pembayaran.
func Seed() {
	knownTypes := []string{"otp_verification", "kyc_approved", "kyc_rejected", "payment_verified", "payment_rejected", "password_reset"}
	for _, et := range knownTypes {
		db().Exec(context.Background(),
			`INSERT INTO notification_event_types (event_type) VALUES ($1) ON CONFLICT DO NOTHING`, et)
	}

	// Seed default template for password_reset using parameterized query (avoids gen_random_uuid dependency)
	resetTemplateHTML := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;background:#f4f6fb;margin:0;padding:40px 20px;">
  <div style="max-width:480px;margin:0 auto;background:#fff;border-radius:16px;padding:36px;box-shadow:0 4px 24px rgba(78,115,223,0.10);border:1px solid #eaecf0;">
    <div style="margin-bottom:24px;">
      <span style="font-size:1.2rem;font-weight:700;color:#1a1c23;">Think<span style="color:#4e73df;">Nalyze</span></span>
    </div>
    <h2 style="font-size:1.2rem;font-weight:700;color:#1a1c23;margin:0 0 8px;">Reset Kata Sandi</h2>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 8px;">Halo, {{full_name}}.</p>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 24px;">
      Kami menerima permintaan reset kata sandi untuk akun Anda.
      Klik tombol di bawah untuk melanjutkan.
      Tautan ini berlaku selama <strong>15 menit</strong> dan hanya dapat digunakan <strong>satu kali</strong>.
    </p>
    <a href="{{reset_url}}" style="display:block;text-align:center;background:linear-gradient(135deg,#4e73df,#6f42c1);color:#fff;text-decoration:none;padding:14px 24px;border-radius:10px;font-weight:600;font-size:0.95rem;margin-bottom:20px;">
      Reset Kata Sandi Saya
    </a>
    <p style="color:#9ca3af;font-size:0.78rem;line-height:1.6;margin:0;border-top:1px solid #f0f0f0;padding-top:16px;">
      Jika Anda tidak meminta reset ini, abaikan email ini. Kata sandi Anda tidak akan berubah.
    </p>
  </div>
</body>
</html>`

	_, seedErr := db().Exec(context.Background(),
		`INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (event_type, channel) DO NOTHING`,
		uuid.New().String(),
		"Reset Kata Sandi",
		"password_reset",
		"email",
		"Reset Kata Sandi ThinkNalyze",
		resetTemplateHTML,
		"system",
	)
	if seedErr != nil {
		log.Printf("[NOTIFICATION] WARN: failed to seed password_reset template: %v", seedErr)
	} else {
		log.Println("[NOTIFICATION] password_reset template seed OK (or already exists)")
	}
	// Remove obsolete event types that were inserted by mistake
	obsolete := []string{"account_deactivated", "account_reactivated", "user_created"}
	for _, et := range obsolete {
		db().Exec(context.Background(),
			`DELETE FROM notification_event_types WHERE event_type = $1`, et)
	}
}

// db adalah shorthand internal agar tidak perlu akses global DB langsung di migrate.
func db() *pgxpool.Pool {
	return DB
}
