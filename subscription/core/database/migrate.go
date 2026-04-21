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

	if wrapper == nil || wrapper.Pool == nil {
		log.Println("[WARN] Database pool is nil, skipping migration & seeding")
		return
	}

	pool := wrapper.Pool

	migrationSQL := `
	CREATE SCHEMA IF NOT EXISTS subscription;

	-- ============================================================
	-- PACKAGES: tier fitur (quota, dll). Durasi TIDAK lagi di sini.
	-- ============================================================
	CREATE TABLE IF NOT EXISTS subscription.packages (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price DECIMAL(15,2) NOT NULL,             -- harga dasar / bulan (referensi)
		quota INT NOT NULL,
		status VARCHAR(50) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'DELETED')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Safe migration: hapus kolom lama yang tidak dipakai lagi (best-effort)
	DO $$ BEGIN
		ALTER TABLE subscription.packages DROP COLUMN IF EXISTS duration;
		ALTER TABLE subscription.packages DROP COLUMN IF EXISTS price_yearly;
	EXCEPTION WHEN others THEN NULL;
	END $$;

	-- ============================================================
	-- PACKAGE PRICING: harga per durasi per paket
	-- ============================================================
	CREATE TABLE IF NOT EXISTS subscription.package_pricing (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		package_id UUID NOT NULL REFERENCES subscription.packages(id) ON DELETE CASCADE,
		duration_months INT NOT NULL CHECK (duration_months > 0),
		price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
		label VARCHAR(100) DEFAULT '',       -- e.g. "Hemat 20%", "Special Price"
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (package_id, duration_months)
	);

	-- ============================================================
	-- ORDERS: tambah duration_months
	-- ============================================================
	CREATE TABLE IF NOT EXISTS subscription.orders (
		id UUID PRIMARY KEY,
		invoice_number VARCHAR(100) UNIQUE NOT NULL,
		user_id UUID NOT NULL,
		package_id UUID,
		duration_months INT NOT NULL DEFAULT 1,
		payment_method VARCHAR(100) NOT NULL DEFAULT '',
		payment_proof BYTEA,
		payment_proof_filename VARCHAR(255) DEFAULT '',
		payment_proof_content_type VARCHAR(100) DEFAULT '',
		payment_proof_uploaded_at TIMESTAMP,
		verification_note TEXT DEFAULT '',
		total_price DECIMAL(15,2) NOT NULL DEFAULT 0,
		status VARCHAR(50) DEFAULT 'PENDING_PAYMENT' CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED', 'EXPIRED')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Safe migration untuk tabel orders yang sudah ada
	DO $$ BEGIN
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS payment_method VARCHAR(100) NOT NULL DEFAULT '';
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS duration_months INT NOT NULL DEFAULT 1;
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS payment_proof BYTEA;
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS payment_proof_filename VARCHAR(255) DEFAULT '';
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS payment_proof_content_type VARCHAR(100) DEFAULT '';
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS payment_proof_uploaded_at TIMESTAMP;
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS verification_note TEXT DEFAULT '';
	EXCEPTION WHEN others THEN NULL;
	END $$;

	-- ============================================================
	-- SUBSCRIPTIONS: aktivasi langganan setelah pembayaran terverifikasi
	-- ============================================================
	CREATE TABLE IF NOT EXISTS subscription.subscriptions (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL UNIQUE REFERENCES subscription.orders(id) ON DELETE CASCADE,
		user_id UUID NOT NULL,
		package_id UUID NOT NULL REFERENCES subscription.packages(id),
		start_date DATE NOT NULL,
		end_date DATE NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'EXPIRED', 'CANCELLED')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_subscriptions_user_status
		ON subscription.subscriptions(user_id, status);

	CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date
		ON subscription.subscriptions(end_date);
	`

	log.Println("Menjalankan Migrasi PostgreSQL (Subscription)...")
	_, err := pool.Exec(ctx, migrationSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}

	// Seed default packages
	seederSQL := `
	INSERT INTO subscription.packages (id, name, price, quota, status, created_at)
	VALUES
	    ('a1b2c3d4-e5f6-7890-1234-567890abcdef', 'Paket Basic',      150000.00, 100,  'ACTIVE',   CURRENT_TIMESTAMP),
	    ('b2c3d4e5-f678-9012-3456-7890abcdef12', 'Paket Pro',        500000.00, 1000, 'ACTIVE',   CURRENT_TIMESTAMP),
	    ('c3d4e5f6-7890-1234-5678-90abcdef1234', 'Paket Enterprise', 1200000.00,5000, 'ACTIVE',   CURRENT_TIMESTAMP),
	    ('d4e5f6a7-8901-2345-6789-0abcdef12345', 'Paket Starter',    49000.00,  50,   'INACTIVE', CURRENT_TIMESTAMP),
	    ('e5f6a7b8-9012-3456-7890-abcdef123456', 'Paket Trial',      0.00,      10,   'INACTIVE', CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO NOTHING;
	`

	log.Println("Menjalankan Seeder (Subscription)...")
	_, err = pool.Exec(ctx, seederSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan seeder: %v", err)
	}

	// Seed pricing tiers
	pricingSeederSQL := `
	INSERT INTO subscription.package_pricing (package_id, duration_months, price, label) VALUES
	    -- Paket Basic
	    ('a1b2c3d4-e5f6-7890-1234-567890abcdef',  1,   150000.00, ''),
	    ('a1b2c3d4-e5f6-7890-1234-567890abcdef',  3,   405000.00, 'Hemat 10%'),
	    ('a1b2c3d4-e5f6-7890-1234-567890abcdef',  6,   720000.00, 'Hemat 20%'),
	    ('a1b2c3d4-e5f6-7890-1234-567890abcdef', 12,  1260000.00, 'Hemat 30%'),
	    -- Paket Pro
	    ('b2c3d4e5-f678-9012-3456-7890abcdef12',  1,   500000.00, ''),
	    ('b2c3d4e5-f678-9012-3456-7890abcdef12',  3,  1350000.00, 'Hemat 10%'),
	    ('b2c3d4e5-f678-9012-3456-7890abcdef12',  6,  2400000.00, 'Hemat 20%'),
	    ('b2c3d4e5-f678-9012-3456-7890abcdef12', 12,  4200000.00, 'Hemat 30%'),
	    -- Paket Enterprise
	    ('c3d4e5f6-7890-1234-5678-90abcdef1234',  1,  1200000.00, ''),
	    ('c3d4e5f6-7890-1234-5678-90abcdef1234',  3,  3240000.00, 'Hemat 10%'),
	    ('c3d4e5f6-7890-1234-5678-90abcdef1234',  6,  5760000.00, 'Hemat 20%'),
	    ('c3d4e5f6-7890-1234-5678-90abcdef1234', 12, 10080000.00, 'Hemat 30%')
	ON CONFLICT (package_id, duration_months) DO NOTHING;
	`

	log.Println("Menjalankan Seeder Pricing Tiers...")
	_, err = pool.Exec(ctx, pricingSeederSQL)
	if err != nil {
		log.Printf("[WARN] Gagal seed pricing tiers: %v", err)
	}

	cleanupDummyOrdersSQL := `
	DELETE FROM subscription.orders
	WHERE invoice_number IN (
		'INV-2026-00001',
		'INV-2026-00002',
		'INV-2026-00003',
		'INV-2026-00004',
		'INV-2026-00005'
	);
	`

	log.Println("Membersihkan data dummy order subscription...")
	if _, err = pool.Exec(ctx, cleanupDummyOrdersSQL); err != nil {
		log.Printf("[WARN] Gagal membersihkan dummy orders: %v", err)
	}

	log.Println("Migrasi dan Seeding selesai! (Subscription)")
}
