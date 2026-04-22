package database

import (
	"context"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func hashPwd(plain string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password: %v", err)
	}
	return string(h)
}

func MigrateAndSeed(db interface{}) {
	ctx := context.Background()

	wrapper, ok := db.(*DBWrapper)
	if !ok {
		log.Fatalf("[FATAL] Input MigrateAndSeed bukan *database.DBWrapper")
	}

	pool := wrapper.Pool
	if pool == nil {
		log.Printf("[WARN] Database pool is nil, skipping migration & seeding")
		return
	}

	// ===== MIGRASI (tanpa parameter, bisa batch) =====
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS levels (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		code VARCHAR(50) UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS groups (
		id UUID PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		code VARCHAR(50) UNIQUE NOT NULL,
		parent UUID REFERENCES groups(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS roles (
		id UUID PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		code VARCHAR(50) UNIQUE NOT NULL,
		group_id UUID REFERENCES groups(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		phone VARCHAR(20) NULL,
		birthdate DATE NULL,
		photo VARCHAR(255) NULL,
		group_id UUID REFERENCES groups(id) ON DELETE SET NULL,
		level_id INT REFERENCES levels(id) ON DELETE SET NULL,
		role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
		status VARCHAR(20) DEFAULT 'inactive' CHECK (status IN ('banned', 'active', 'inactive')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by UUID REFERENCES users(id) ON DELETE SET NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_by UUID REFERENCES users(id) ON DELETE SET NULL
	);

	ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20) NULL;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS birthdate DATE NULL;

	CREATE TABLE IF NOT EXISTS otp_codes (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		email VARCHAR(255) NOT NULL,
		code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_otp_email ON otp_codes(email);
	CREATE INDEX IF NOT EXISTS idx_otp_user ON otp_codes(user_id);

	CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		type VARCHAR(50) DEFAULT 'info' CHECK (type IN ('info', 'warning', 'promo', 'system')),
		target_role VARCHAR(50) DEFAULT 'all',
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by UUID REFERENCES users(id) ON DELETE SET NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_by UUID REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS subscription_packages (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		code VARCHAR(50) UNIQUE NOT NULL,
		description TEXT,
		price DECIMAL(15,2) NOT NULL DEFAULT 0,
		duration_days INT NOT NULL DEFAULT 30,
		features JSONB DEFAULT '[]',
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by UUID REFERENCES users(id) ON DELETE SET NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_by UUID REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS kyc_submissions (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		full_name VARCHAR(255) NOT NULL,
		nik VARCHAR(16) UNIQUE NOT NULL,
		address TEXT NOT NULL,
		birthdate DATE NOT NULL,
		phone VARCHAR(20) NOT NULL,
		ktp_image TEXT NOT NULL,
		status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
		reject_reason TEXT,
		reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
		reviewed_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_kyc_nik ON kyc_submissions(nik);
	CREATE INDEX IF NOT EXISTS idx_kyc_sub_user ON kyc_submissions(user_id);
	CREATE INDEX IF NOT EXISTS idx_kyc_sub_status ON kyc_submissions(status);
	CREATE INDEX IF NOT EXISTS idx_notif_active ON notifications(is_active);
	CREATE INDEX IF NOT EXISTS idx_subs_active ON subscription_packages(is_active);

	-- Upgrade ktp_image to TEXT if it was previously VARCHAR(500)
	ALTER TABLE kyc_submissions ALTER COLUMN ktp_image TYPE TEXT;

	-- Add rejected_fields column if not exists (TEXT[] for storing field names)
	ALTER TABLE kyc_submissions ADD COLUMN IF NOT EXISTS rejected_fields TEXT[];
	`

	log.Println("Menjalankan Migrasi PostgreSQL...")
	_, err := pool.Exec(ctx, migrationSQL)
	if err != nil {
		// Ignore duplicate object errors from concurrent migrations (account & users share DB)
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "already exists") {
			log.Printf("[WARN] Migrasi race condition (safe to ignore): %v", err)
		} else {
			log.Fatalf("Gagal menjalankan migrasi: %v", err)
		}
	}

	// ===== SEEDER STATIC (tanpa parameter $1, bisa batch) =====
	seederStaticSQL := `
	INSERT INTO levels (id, name, code) VALUES (1, 'Super Admin', 'SUPERADMIN') ON CONFLICT (id) DO NOTHING;
	INSERT INTO levels (id, name, code) VALUES (2, 'Client', 'CLIENT') ON CONFLICT (id) DO NOTHING;
	INSERT INTO levels (id, name, code) VALUES (3, 'Staff', 'STAFF') ON CONFLICT (id) DO NOTHING;

	INSERT INTO groups (id, name, code) VALUES ('1e98c63f-5474-4506-826c-ded22b59b3db', 'Owner', 'Owner') ON CONFLICT (id) DO NOTHING;
	INSERT INTO groups (id, name, code) VALUES ('2e98c63f-5474-4506-826c-ded22b59b3dc', 'Member', 'Member') ON CONFLICT (id) DO NOTHING;
	INSERT INTO groups (id, name, code) VALUES ('3e98c63f-5474-4506-826c-ded22b59b3dd', 'Internal', 'Internal') ON CONFLICT (id) DO NOTHING;

	INSERT INTO roles (id, name, code, group_id) VALUES ('cf47ce1c-1455-4a20-bafe-c2b7c2ab9993', 'CEO', 'CEO', '1e98c63f-5474-4506-826c-ded22b59b3db') ON CONFLICT (id) DO NOTHING;
	INSERT INTO roles (id, name, code, group_id) VALUES ('df47ce1c-1455-4a20-bafe-c2b7c2ab9994', 'Client', 'CLIENT', '2e98c63f-5474-4506-826c-ded22b59b3dc') ON CONFLICT (id) DO NOTHING;
	INSERT INTO roles (id, name, code, group_id) VALUES ('af47ce1c-1455-4a20-bafe-c2b7c2ab9995', 'Operasional', 'OPERASIONAL', '3e98c63f-5474-4506-826c-ded22b59b3dd') ON CONFLICT (id) DO NOTHING;
	INSERT INTO roles (id, name, code, group_id) VALUES ('bf47ce1c-1455-4a20-bafe-c2b7c2ab9996', 'Compliance', 'COMPLIANCE', '3e98c63f-5474-4506-826c-ded22b59b3dd') ON CONFLICT (id) DO NOTHING;
	INSERT INTO roles (id, name, code, group_id) VALUES ('cf47ce1c-1455-4a20-bafe-c2b7c2ab9997', 'Management', 'MANAGEMENT', '3e98c63f-5474-4506-826c-ded22b59b3dd') ON CONFLICT (id) DO NOTHING;
	INSERT INTO roles (id, name, code, group_id) VALUES ('ef47ce1c-1455-4a20-bafe-c2b7c2ab9997', 'Customer Support', 'ADMIN_SUPPORT', '3e98c63f-5474-4506-826c-ded22b59b3dd') ON CONFLICT (id) DO NOTHING;

	INSERT INTO users (id, name, email, password, group_id, level_id, role_id, status, created_at, created_by) 
	VALUES (
		'10ef7bff-4c69-4b56-aec8-ef7427601952', 'Muhammad Abror', 'abrorcapital@gmail.com',
		'$2a$12$l4hsF/yaFp07frhIrqfm7uJNfkOfuNXsZWEqv4YgFhs3/eYbPAqoS',
		'1e98c63f-5474-4506-826c-ded22b59b3db', 1, 'cf47ce1c-1455-4a20-bafe-c2b7c2ab9993',
		'active', CURRENT_TIMESTAMP, '10ef7bff-4c69-4b56-aec8-ef7427601952'
	) ON CONFLICT (email) DO NOTHING;

	INSERT INTO subscription_packages (id, name, code, description, price, duration_days, features, is_active) VALUES 
		('a1a1a1a1-1111-1111-1111-111111111111', 'Free Plan', 'FREE', 'Basic access with limited features', 0, 9999, '["Market Insight (Delayed)", "Ask Nizza (5/day)"]', TRUE)
	ON CONFLICT (code) DO NOTHING;

	INSERT INTO subscription_packages (id, name, code, description, price, duration_days, features, is_active) VALUES 
		('b2b2b2b2-2222-2222-2222-222222222222', 'Pro Plan', 'PRO', 'Professional trader access', 299000, 30, '["Real-time Market Insight", "Ask Nizza Unlimited", "Deep Scanner", "Algo Studio"]', TRUE)
	ON CONFLICT (code) DO NOTHING;

	INSERT INTO subscription_packages (id, name, code, description, price, duration_days, features, is_active) VALUES 
		('c3c3c3c3-3333-3333-3333-333333333333', 'Elite Plan', 'ELITE', 'Full institutional-grade access', 799000, 30, '["All Pro Features", "Auto Trade", "Priority Support", "Custom Alerts", "API Access"]', TRUE)
	ON CONFLICT (code) DO NOTHING;
	`

	log.Println("Menjalankan Seeder (static data)...")
	_, err = pool.Exec(ctx, seederStaticSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan seeder static: %v", err)
	}

	// ===== SEEDER DENGAN PARAMETER (satu per satu, karena pgx tidak support multiple prepared statements) =====
	log.Println("Menjalankan Seeder (user accounts)...")

	// Super Admin baru (password: Super123)
	superAdminPwd := hashPwd("Super123")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, name, email, password, phone, group_id, level_id, role_id, status, created_at, created_by) 
		 VALUES (
			'11ef7bff-4c69-4b56-aec8-ef7427601960', 'Super Admin', 'superadmin@thinktala.com',
			$1, '080000000000',
			'1e98c63f-5474-4506-826c-ded22b59b3db', 1, 'cf47ce1c-1455-4a20-bafe-c2b7c2ab9993',
			'active', CURRENT_TIMESTAMP, '10ef7bff-4c69-4b56-aec8-ef7427601952'
		) ON CONFLICT (email) DO UPDATE SET password = $1`,
		superAdminPwd)
	if err != nil {
		log.Fatalf("Gagal seed Super Admin: %v", err)
	}
	log.Println("  ✓ Super Admin (superadmin@thinktala.com / Super123)")

	// Operasional (password: Operas123)
	opsPwd := hashPwd("Operas123")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, name, email, password, phone, group_id, level_id, role_id, status, created_at, created_by) 
		 VALUES (
			'20ef7bff-4c69-4b56-aec8-ef7427601953', 'Staff Operasional', 'ops@thinktala.com',
			$1, '081234567890',
			'3e98c63f-5474-4506-826c-ded22b59b3dd', 3, 'af47ce1c-1455-4a20-bafe-c2b7c2ab9995',
			'active', CURRENT_TIMESTAMP, '10ef7bff-4c69-4b56-aec8-ef7427601952'
		) ON CONFLICT (email) DO UPDATE SET password = $1`,
		opsPwd)
	if err != nil {
		log.Fatalf("Gagal seed Operasional: %v", err)
	}
	log.Println("  ✓ Operasional (ops@thinktala.com / Operas123)")

	// Compliance (password: Comply123)
	compliancePwd := hashPwd("Comply123")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, name, email, password, phone, group_id, level_id, role_id, status, created_at, created_by) 
		 VALUES (
			'30ef7bff-4c69-4b56-aec8-ef7427601954', 'Staff Compliance', 'compliance@thinktala.com',
			$1, '081234567891',
			'3e98c63f-5474-4506-826c-ded22b59b3dd', 3, 'bf47ce1c-1455-4a20-bafe-c2b7c2ab9996',
			'active', CURRENT_TIMESTAMP, '10ef7bff-4c69-4b56-aec8-ef7427601952'
		) ON CONFLICT (email) DO UPDATE SET password = $1`,
		compliancePwd)
	if err != nil {
		log.Fatalf("Gagal seed Compliance: %v", err)
	}
	log.Println("  ✓ Compliance (compliance@thinktala.com / Comply123)")

	// Management (password: Manage123)
	managementPwd := hashPwd("Manage123")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, name, email, password, phone, group_id, level_id, role_id, status, created_at, created_by)
		 VALUES (
			'40ef7bff-4c69-4b56-aec8-ef7427601955', 'Manager Analytics', 'management@thinktala.com',
			$1, '081234567892',
			'3e98c63f-5474-4506-826c-ded22b59b3dd', 3, 'cf47ce1c-1455-4a20-bafe-c2b7c2ab9997',
			'active', CURRENT_TIMESTAMP, '10ef7bff-4c69-4b56-aec8-ef7427601952'
		) ON CONFLICT (email) DO UPDATE SET password = $1`,
		managementPwd)
	if err != nil {
		log.Fatalf("Gagal seed Management: %v", err)
	}
	log.Println("  ✓ Management (management@thinktala.com / Manage123)")

	// Customer Support (password: Support123)
	supportPwd := hashPwd("Support123")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, name, email, password, phone, group_id, level_id, role_id, status, created_at, created_by)
		 VALUES (
			'45ef7bff-4c69-4b56-aec8-ef7427601955', 'Customer Support', 'support@thinktala.com',
			$1, '081234567893',
			'3e98c63f-5474-4506-826c-ded22b59b3dd', 3, 'ef47ce1c-1455-4a20-bafe-c2b7c2ab9997',
			'active', CURRENT_TIMESTAMP, '10ef7bff-4c69-4b56-aec8-ef7427601952'
		) ON CONFLICT (email) DO UPDATE SET password = $1`,
		supportPwd)
	if err != nil {
		log.Fatalf("Gagal seed Customer Support: %v", err)
	}
	log.Println("  ✓ Customer Support (support@thinktala.com / Support123)")

	log.Println("Migrasi dan Seeding selesai!")
}
