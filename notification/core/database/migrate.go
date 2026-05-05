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

	// Seed default template for password_reset if not yet created via API
	db().Exec(context.Background(), `
		INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		VALUES (
			gen_random_uuid()::text,
			'Reset Kata Sandi',
			'password_reset',
			'email',
			'Reset Kata Sandi ThinkNalyze',
			'<!DOCTYPE html>
<html lang="id">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#f4f6f9;font-family:''Segoe UI'',Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background:#f4f6f9;padding:40px 0;">
    <tr><td align="center">
      <table width="520" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);">
        <tr><td style="background:#1a1c2e;padding:28px 36px;text-align:center;">
          <span style="color:#ffffff;font-size:22px;font-weight:700;">ThinkNalyze</span>
        </td></tr>
        <tr><td style="padding:36px 36px 28px;">
          <p style="margin:0 0 16px;font-size:16px;color:#1a1c2e;font-weight:600;">Halo, {{full_name}}</p>
          <p style="margin:0 0 20px;font-size:14px;color:#4a5568;line-height:1.6;">
            Kami menerima permintaan untuk mereset kata sandi akun ThinkNalyze Anda.
            Klik tombol di bawah untuk membuat kata sandi baru. Tautan ini berlaku selama <strong>15 menit</strong>.
          </p>
          <table cellpadding="0" cellspacing="0" width="100%"><tr><td align="center" style="padding:8px 0 24px;">
            <a href="{{reset_url}}" style="display:inline-block;background:#4e73df;color:#ffffff;text-decoration:none;font-weight:600;font-size:15px;padding:14px 36px;border-radius:8px;">
              Reset Kata Sandi
            </a>
          </td></tr></table>
          <p style="margin:0 0 8px;font-size:13px;color:#718096;">Atau salin tautan berikut ke browser Anda:</p>
          <p style="margin:0 0 24px;font-size:12px;color:#4e73df;word-break:break-all;">{{reset_url}}</p>
          <hr style="border:none;border-top:1px solid #e8ecf0;margin:0 0 20px;">
          <p style="margin:0;font-size:12px;color:#a0aec0;line-height:1.6;">
            Jika Anda tidak meminta reset kata sandi, abaikan email ini — akun Anda tetap aman.
          </p>
        </td></tr>
        <tr><td style="background:#f7f9fc;padding:16px 36px;text-align:center;">
          <p style="margin:0;font-size:11px;color:#a0aec0;">&copy; 2026 ThinkNalyze. All rights reserved.</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>',
			'system'
		)
		ON CONFLICT (event_type, channel) DO NOTHING
	`)
	log.Println("[NOTIFICATION] password_reset template seeded")
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
