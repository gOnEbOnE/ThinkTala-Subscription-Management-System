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
			view_count  INTEGER NOT NULL DEFAULT 0,
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
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS view_count INTEGER NOT NULL DEFAULT 0;

		UPDATE notifications
		SET description = COALESCE(NULLIF(description, ''), message)
		WHERE description IS NULL OR description = '';

		UPDATE notifications
		SET message = COALESCE(NULLIF(message, ''), description, '')
		WHERE message IS NULL OR message = '';

		ALTER TABLE notifications
			DROP CONSTRAINT IF EXISTS notifications_type_check;
		ALTER TABLE notifications
			ADD CONSTRAINT notifications_type_check
			CHECK (type IN ('system', 'promo', 'warning', 'info', 'analysis', 'education', 'event'));

		CREATE INDEX IF NOT EXISTS idx_notifications_active ON notifications(is_active, target_role);
		CREATE INDEX IF NOT EXISTS idx_notifications_pinned ON notifications(is_pinned, created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notifications: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notifications ready")
	}

	// Tambahkan kolom telegram_chat_id ke tabel users (karena tabel users ada di database yang sama)
	_, err = db().Exec(context.Background(), `
		ALTER TABLE users ADD COLUMN IF NOT EXISTS telegram_chat_id VARCHAR(50) DEFAULT NULL;
		CREATE INDEX IF NOT EXISTS idx_users_telegram_chat_id ON users (telegram_chat_id) WHERE telegram_chat_id IS NOT NULL;
	`)
	if err != nil {
		log.Printf("[WARN] migrate users (telegram_chat_id): %v", err)
	} else {
		log.Println("[NOTIFICATION] Users table updated with telegram_chat_id")
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
			user_id       VARCHAR(36),
			user_name     TEXT,
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
			provider_response TEXT,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS user_id VARCHAR(36);
		ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS user_name TEXT;
		ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS provider_response TEXT;
		CREATE INDEX IF NOT EXISTS idx_notif_logs_status ON notification_logs(status);
		CREATE INDEX IF NOT EXISTS idx_notif_logs_retry  ON notification_logs(status, next_retry_at);
		CREATE INDEX IF NOT EXISTS idx_notif_logs_sent_at ON notification_logs(sent_at);
		CREATE INDEX IF NOT EXISTS idx_notif_logs_user_id ON notification_logs(user_id);
		CREATE INDEX IF NOT EXISTS idx_notif_logs_channel ON notification_logs(channel);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notification_logs: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notification_logs ready")
	}

	// Table untuk tracking read status di notification drawer.
	_, err = db().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS notification_reads (
			user_id     VARCHAR(36) NOT NULL,
			source_type VARCHAR(20) NOT NULL,
			source_id   VARCHAR(36) NOT NULL,
			read_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, source_type, source_id)
		);
		CREATE INDEX IF NOT EXISTS idx_notif_reads_user ON notification_reads(user_id);
		CREATE INDEX IF NOT EXISTS idx_notif_reads_source ON notification_reads(source_type, source_id);
	`)
	if err != nil {
		log.Printf("[WARN] migrate notification_reads: %v", err)
	} else {
		log.Println("[NOTIFICATION] Table notification_reads ready")
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

	seedNotifications := []struct {
		id          string
		title       string
		description string
		typeValue   string
		targetRole  string
		ctaURL      *string
		imageURL    *string
		isPinned    bool
	}{
		{
			id:          "11111111-1111-1111-1111-111111111111",
			title:       "Maintenance Window Tonight",
			description: "Platform will be briefly unavailable at 23:00-23:20 WIB. Trading and dashboard will resume normally.",
			typeValue:   "system",
			targetRole:  "client",
			ctaURL:      strPtr("https://www.thinktala.com/status"),
			imageURL:    nil,
			isPinned:    false,
		},
		{
			id:    "22222222-2222-2222-2222-222222222222",
			title: "Risk Management Playbook for New Traders",
			description: "<h2>Why it matters</h2><p>Markets can move fast. A simple risk routine keeps your decisions consistent and protects capital when volatility spikes.</p><h3>Core checklist</h3><ul><li>Set a max loss per trade (1-2% of capital).</li><li>Use a stop loss before entering the trade.</li><li>Size the position based on stop distance.</li><li>Take partial profits at predefined levels.</li></ul><p><strong>Pro tip:</strong> Write your rule set and review it weekly. Small adjustments compound quickly.</p><h3>Example workflow</h3><p>Plan the entry, define the stop, calculate size, and log the outcome. Repeat this sequence to build a stable edge.</p>",
			typeValue:  "education",
			targetRole: "client",
			ctaURL:     strPtr("https://www.thinktala.com/learn/risk-playbook"),
			imageURL:   strPtr("/assets/img/ecosystem.png"),
			isPinned:   false,
		},
	}

	for _, item := range seedNotifications {
		_, seedErr := db().Exec(context.Background(), `
			INSERT INTO notifications
				(id, title, description, message, type, target_role, cta_url, image_url, is_active, is_pinned, created_at, updated_at, created_by, updated_by)
			VALUES
				($1, $2, $3, $3, $4, $5, $6, $7, TRUE, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $9, $9)
			ON CONFLICT (id) DO NOTHING
		`, item.id, item.title, item.description, item.typeValue, item.targetRole, item.ctaURL, item.imageURL, item.isPinned, "seed")
		if seedErr != nil {
			log.Printf("[NOTIFICATION] WARN: failed to seed notification %s: %v", item.id, seedErr)
		}
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

	// Seed template for otp_verification
	otpHTML := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;background:#f4f6fb;margin:0;padding:40px 20px;">
  <div style="max-width:520px;margin:0 auto;background:#fff;border-radius:16px;padding:36px;box-shadow:0 4px 24px rgba(78,115,223,0.10);border:1px solid #eaecf0;">
    <div style="margin-bottom:24px;">
      <span style="font-size:1.2rem;font-weight:700;color:#1a1c23;">Think<span style="color:#4e73df;">Nalyze</span></span>
    </div>
    <h2 style="font-size:1.2rem;font-weight:700;color:#1a1c23;margin:0 0 8px;">Kode Verifikasi Anda</h2>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 20px;">Halo, {{user_name}}. Gunakan kode berikut untuk menyelesaikan verifikasi akun Anda.</p>
    <div style="background:#f0f4ff;border:1px solid #c7d2fe;border-radius:10px;padding:20px;text-align:center;margin-bottom:20px;">
      <p style="font-size:2rem;font-weight:800;letter-spacing:0.5rem;color:#4e73df;margin:0;">{{otp_code}}</p>
    </div>
    <p style="color:#6c757d;font-size:0.85rem;line-height:1.6;margin:0 0 8px;">Kode ini berlaku selama <strong>10 menit</strong>. Jangan bagikan kepada siapapun.</p>
    <p style="color:#9ca3af;font-size:0.78rem;line-height:1.6;margin:0;border-top:1px solid #f0f0f0;padding-top:16px;">
      Jika Anda tidak meminta kode ini, abaikan email ini.
    </p>
  </div>
</body>
</html>`

	_, seedErr = db().Exec(context.Background(),
		`INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (event_type, channel) DO NOTHING`,
		uuid.New().String(),
		"Kode OTP Verifikasi",
		"otp_verification",
		"email",
		"Kode Verifikasi Akun ThinkNalyze Anda",
		otpHTML,
		"system",
	)
	if seedErr != nil {
		log.Printf("[NOTIFICATION] WARN: failed to seed otp_verification template: %v", seedErr)
	} else {
		log.Println("[NOTIFICATION] otp_verification template seed OK (or already exists)")
	}

	// Seed template for kyc_approved
	kycApprovedHTML := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;background:#f4f6fb;margin:0;padding:40px 20px;">
  <div style="max-width:520px;margin:0 auto;background:#fff;border-radius:16px;padding:36px;box-shadow:0 4px 24px rgba(78,115,223,0.10);border:1px solid #eaecf0;">
    <div style="margin-bottom:24px;">
      <span style="font-size:1.2rem;font-weight:700;color:#1a1c23;">Think<span style="color:#4e73df;">Nalyze</span></span>
    </div>
    <h2 style="font-size:1.2rem;font-weight:700;color:#1a1c23;margin:0 0 8px;">Verifikasi KYC Berhasil ✓</h2>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 8px;">Halo, {{user_name}}.</p>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 20px;">
      Selamat! Verifikasi identitas (KYC) Anda telah disetujui. Anda sekarang dapat menikmati layanan ThinkNalyze secara penuh.
    </p>
    <div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:10px;padding:16px;margin-bottom:20px;">
      <p style="margin:0;color:#166534;font-size:0.9rem;font-weight:600;">Status: KYC Disetujui</p>
    </div>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 16px;">Langkah selanjutnya: pilih paket investasi yang sesuai dengan kebutuhan Anda.</p>
    <p style="color:#9ca3af;font-size:0.78rem;line-height:1.6;margin:0;border-top:1px solid #f0f0f0;padding-top:16px;">
      Terima kasih telah bergabung dengan ThinkNalyze.
    </p>
  </div>
</body>
</html>`

	_, seedErr = db().Exec(context.Background(),
		`INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (event_type, channel) DO NOTHING`,
		uuid.New().String(),
		"KYC Disetujui",
		"kyc_approved",
		"email",
		"Verifikasi KYC Anda Telah Disetujui - ThinkNalyze",
		kycApprovedHTML,
		"system",
	)
	if seedErr != nil {
		log.Printf("[NOTIFICATION] WARN: failed to seed kyc_approved template: %v", seedErr)
	} else {
		log.Println("[NOTIFICATION] kyc_approved template seed OK (or already exists)")
	}

	// Seed template for kyc_rejected
	kycRejectedHTML := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;background:#f4f6fb;margin:0;padding:40px 20px;">
  <div style="max-width:520px;margin:0 auto;background:#fff;border-radius:16px;padding:36px;box-shadow:0 4px 24px rgba(78,115,223,0.10);border:1px solid #eaecf0;">
    <div style="margin-bottom:24px;">
      <span style="font-size:1.2rem;font-weight:700;color:#1a1c23;">Think<span style="color:#4e73df;">Nalyze</span></span>
    </div>
    <h2 style="font-size:1.2rem;font-weight:700;color:#1a1c23;margin:0 0 8px;">Verifikasi KYC Ditolak</h2>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 8px;">Halo, {{user_name}}.</p>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 20px;">
      Kami mohon maaf, verifikasi identitas (KYC) Anda belum dapat kami setujui saat ini.
    </p>
    <div style="background:#fef2f2;border:1px solid #fecaca;border-radius:10px;padding:16px;margin-bottom:20px;">
      <p style="margin:0 0 4px;color:#991b1b;font-size:0.9rem;font-weight:600;">Alasan Penolakan:</p>
      <p style="margin:0;color:#7f1d1d;font-size:0.9rem;">{{reject_reason}}</p>
    </div>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 16px;">
      Silakan perbaiki data dan dokumen Anda, lalu ajukan ulang verifikasi KYC melalui portal ThinkNalyze.
    </p>
    <p style="color:#9ca3af;font-size:0.78rem;line-height:1.6;margin:0;border-top:1px solid #f0f0f0;padding-top:16px;">
      Jika Anda membutuhkan bantuan, hubungi kami melalui menu Support Ticket.
    </p>
  </div>
</body>
</html>`

	_, seedErr = db().Exec(context.Background(),
		`INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (event_type, channel) DO NOTHING`,
		uuid.New().String(),
		"KYC Ditolak",
		"kyc_rejected",
		"email",
		"Verifikasi KYC Ditolak - ThinkNalyze",
		kycRejectedHTML,
		"system",
	)
	if seedErr != nil {
		log.Printf("[NOTIFICATION] WARN: failed to seed kyc_rejected template: %v", seedErr)
	} else {
		log.Println("[NOTIFICATION] kyc_rejected template seed OK (or already exists)")
	}

	// Seed template for payment_verified
	paymentVerifiedHTML := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;background:#f4f6fb;margin:0;padding:40px 20px;">
  <div style="max-width:520px;margin:0 auto;background:#fff;border-radius:16px;padding:36px;box-shadow:0 4px 24px rgba(78,115,223,0.10);border:1px solid #eaecf0;">
    <div style="margin-bottom:24px;">
      <span style="font-size:1.2rem;font-weight:700;color:#1a1c23;">Think<span style="color:#4e73df;">Nalyze</span></span>
    </div>
    <h2 style="font-size:1.2rem;font-weight:700;color:#1a1c23;margin:0 0 8px;">Pembayaran Berhasil Diverifikasi ✓</h2>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 8px;">Halo, {{client_name}}.</p>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 20px;">
      Pembayaran Anda untuk paket <strong>{{package_name}}</strong> dengan nomor invoice <strong>{{invoice_number}}</strong> telah berhasil diverifikasi.
      Langganan Anda kini aktif selama <strong>{{duration_months}} bulan</strong>.
    </p>
    <div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:10px;padding:16px;margin-bottom:20px;">
      <p style="margin:0;color:#166534;font-size:0.9rem;font-weight:600;">Status: Pembayaran Terverifikasi</p>
    </div>
    <p style="color:#9ca3af;font-size:0.78rem;line-height:1.6;margin:0;border-top:1px solid #f0f0f0;padding-top:16px;">
      Terima kasih telah mempercayakan investasi Anda kepada ThinkNalyze.
    </p>
  </div>
</body>
</html>`

	_, seedErr = db().Exec(context.Background(),
		`INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (event_type, channel) DO NOTHING`,
		uuid.New().String(),
		"Pembayaran Terverifikasi",
		"payment_verified",
		"email",
		"Pembayaran Anda Telah Diverifikasi - ThinkNalyze",
		paymentVerifiedHTML,
		"system",
	)
	if seedErr != nil {
		log.Printf("[NOTIFICATION] WARN: failed to seed payment_verified template: %v", seedErr)
	} else {
		log.Println("[NOTIFICATION] payment_verified template seed OK (or already exists)")
	}

	// Seed template for payment_rejected
	paymentRejectedHTML := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;background:#f4f6fb;margin:0;padding:40px 20px;">
  <div style="max-width:520px;margin:0 auto;background:#fff;border-radius:16px;padding:36px;box-shadow:0 4px 24px rgba(78,115,223,0.10);border:1px solid #eaecf0;">
    <div style="margin-bottom:24px;">
      <span style="font-size:1.2rem;font-weight:700;color:#1a1c23;">Think<span style="color:#4e73df;">Nalyze</span></span>
    </div>
    <h2 style="font-size:1.2rem;font-weight:700;color:#1a1c23;margin:0 0 8px;">Pembayaran Ditolak</h2>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 8px;">Halo, {{client_name}}.</p>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 20px;">
      Pembayaran Anda untuk paket <strong>{{package_name}}</strong> dengan nomor invoice <strong>{{invoice_number}}</strong> tidak dapat diverifikasi.
    </p>
    <div style="background:#fef2f2;border:1px solid #fecaca;border-radius:10px;padding:16px;margin-bottom:20px;">
      <p style="margin:0 0 4px;color:#991b1b;font-size:0.9rem;font-weight:600;">Alasan Penolakan:</p>
      <p style="margin:0;color:#7f1d1d;font-size:0.9rem;">{{verification_note}}</p>
    </div>
    <p style="color:#6c757d;font-size:0.9rem;line-height:1.6;margin:0 0 16px;">
      Silakan unggah ulang bukti transfer yang benar atau hubungi tim support kami untuk bantuan lebih lanjut.
    </p>
    <p style="color:#9ca3af;font-size:0.78rem;line-height:1.6;margin:0;border-top:1px solid #f0f0f0;padding-top:16px;">
      Jika Anda merasa ini adalah kesalahan, silakan hubungi kami melalui menu Support Ticket.
    </p>
  </div>
</body>
</html>`

	_, seedErr = db().Exec(context.Background(),
		`INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (event_type, channel) DO NOTHING`,
		uuid.New().String(),
		"Pembayaran Ditolak",
		"payment_rejected",
		"email",
		"Pembayaran Ditolak - ThinkNalyze",
		paymentRejectedHTML,
		"system",
	)
	if seedErr != nil {
		log.Printf("[NOTIFICATION] WARN: failed to seed payment_rejected template: %v", seedErr)
	} else {
		log.Println("[NOTIFICATION] payment_rejected template seed OK (or already exists)")
	}

	// Remove obsolete event types that were inserted by mistake
	obsolete := []string{"account_deactivated", "account_reactivated", "user_created"}
	for _, et := range obsolete {
		db().Exec(context.Background(),
			`DELETE FROM notification_event_types WHERE event_type = $1`, et)
	}
}

func strPtr(value string) *string {
	return &value
}

// db adalah shorthand internal agar tidak perlu akses global DB langsung di migrate.
func db() *pgxpool.Pool {
	return DB
}
