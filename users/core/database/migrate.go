package database

import (
	"context"
	"log"
)

func MigrateAndSeed(db interface{}) {
	ctx := context.Background()

	wrapper, ok := db.(*DBWrapper)
	if !ok {
		log.Fatalf("[FATAL] Input MigrateAndSeed bukan *database.DBWrapper")
	}

	pool := wrapper.Pool

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

	-- Notifications table (managed by Operasional)
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

	-- Subscription packages table (managed by Operasional)
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

	-- KYC documents table (reviewed by Compliance)
	CREATE TABLE IF NOT EXISTS kyc_documents (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		document_type VARCHAR(50) NOT NULL CHECK (document_type IN ('ktp', 'passport', 'sim', 'selfie')),
		file_url VARCHAR(500) NOT NULL,
		status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
		reject_reason TEXT,
		reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
		reviewed_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_kyc_user ON kyc_documents(user_id);
	CREATE INDEX IF NOT EXISTS idx_kyc_status ON kyc_documents(status);
	CREATE INDEX IF NOT EXISTS idx_notif_active ON notifications(is_active);
	CREATE INDEX IF NOT EXISTS idx_subs_active ON subscription_packages(is_active);
	`

	log.Println("Menjalankan Migrasi PostgreSQL...")
	_, err := pool.Exec(ctx, migrationSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}

	seederSQL := `
	-- Seed Level
	INSERT INTO levels (id, name, code) 
	VALUES (1, 'Super Admin', 'SUPERADMIN') 
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO levels (id, name, code) 
	VALUES (2, 'Client', 'CLIENT') 
	ON CONFLICT (id) DO NOTHING;

	-- Seed Group
	INSERT INTO groups (id, name, code) 
	VALUES ('1e98c63f-5474-4506-826c-ded22b59b3db', 'Owner', 'Owner') 
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO groups (id, name, code) 
	VALUES ('2e98c63f-5474-4506-826c-ded22b59b3dc', 'Member', 'Member') 
	ON CONFLICT (id) DO NOTHING;

	-- Seed Role
	INSERT INTO roles (id, name, code, group_id) 
	VALUES ('cf47ce1c-1455-4a20-bafe-c2b7c2ab9993', 'CEO', 'CEO', '1e98c63f-5474-4506-826c-ded22b59b3db') 
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO roles (id, name, code, group_id) 
	VALUES ('df47ce1c-1455-4a20-bafe-c2b7c2ab9994', 'Client', 'CLIENT', '2e98c63f-5474-4506-826c-ded22b59b3dc') 
	ON CONFLICT (id) DO NOTHING;

	-- Seed User
	INSERT INTO users (id, name, email, password, group_id, level_id, role_id, status, created_at, created_by) 
	VALUES (
		'10ef7bff-4c69-4b56-aec8-ef7427601952', 
		'Muhammad Abror', 
		'abrorcapital@gmail.com',
		'$2a$12$l4hsF/yaFp07frhIrqfm7uJNfkOfuNXsZWEqv4YgFhs3/eYbPAqoS', 
		'1e98c63f-5474-4506-826c-ded22b59b3db', 
		1, 
		'cf47ce1c-1455-4a20-bafe-c2b7c2ab9993', 
		'active',
		CURRENT_TIMESTAMP,
		'10ef7bff-4c69-4b56-aec8-ef7427601952'
	) 
	ON CONFLICT (email) DO NOTHING;
	`

	log.Println("Menjalankan Seeder...")
	_, err = pool.Exec(ctx, seederSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan seeder: %v", err)
	}

	log.Println("Migrasi dan Seeding selesai!")
}
