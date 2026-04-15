package database

import (
	"context"
	"log"
)

func MigrateAndSeed(db interface{}) {
	ctx := context.Background()

	// 1. Perbaikan Casting:
	wrapper, ok := db.(*DBWrapper)
	if !ok {
		log.Fatalf("[FATAL] Input MigrateAndSeed bukan *database.DBWrapper")
	}

	// Safety check: jika DB disabled (Pool nil), skip migration
	if wrapper == nil || wrapper.Pool == nil {
		log.Println("[WARN] Database pool is nil, skipping migration & seeding")
		return
	}

	pool := wrapper.Pool

	migrationSQL := `
	CREATE SCHEMA IF NOT EXISTS subscription;

	CREATE TABLE IF NOT EXISTS subscription.packages (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price DECIMAL(15,2) NOT NULL,
		price_yearly DECIMAL(15,2) NOT NULL DEFAULT 0,
		duration INT NOT NULL,
		quota INT NOT NULL,
		status VARCHAR(50) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'DELETED')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Add price_yearly column if table already exists without it
	DO $$ BEGIN
		ALTER TABLE subscription.packages ADD COLUMN IF NOT EXISTS price_yearly DECIMAL(15,2) NOT NULL DEFAULT 0;
	EXCEPTION WHEN others THEN NULL;
	END $$;

	CREATE TABLE IF NOT EXISTS subscription.orders (
		id UUID PRIMARY KEY,
		invoice_number VARCHAR(100) UNIQUE NOT NULL,
		user_id UUID NOT NULL,
		package_id UUID,
		payment_method VARCHAR(100) NOT NULL DEFAULT '',
		total_price DECIMAL(15,2) NOT NULL DEFAULT 0,
		status VARCHAR(50) DEFAULT 'PENDING_PAYMENT' CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED', 'EXPIRED')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Safe migration untuk tabel yang sudah ada
	DO $$ BEGIN
		ALTER TABLE subscription.orders ADD COLUMN IF NOT EXISTS payment_method VARCHAR(100) NOT NULL DEFAULT '';
	EXCEPTION WHEN others THEN NULL;
	END $$;
	`

	log.Println("Menjalankan Migrasi PostgreSQL (Subscription)...")
	_, err := pool.Exec(ctx, migrationSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}

	// 2. Query Seeder (Insert Default Packages Jika Kosong)
	seederSQL := `
	INSERT INTO subscription.packages (id, name, price, price_yearly, duration, quota, status, created_at)
	VALUES 
	    ('a1b2c3d4-e5f6-7890-1234-567890abcdef', 'Paket Basic', 150000.00, 1500000.00, 1, 100, 'ACTIVE', CURRENT_TIMESTAMP),
	    ('b2c3d4e5-f678-9012-3456-7890abcdef12', 'Paket Pro', 500000.00, 5000000.00, 6, 1000, 'ACTIVE', CURRENT_TIMESTAMP),
	    ('c3d4e5f6-7890-1234-5678-90abcdef1234', 'Paket Enterprise', 1200000.00, 12000000.00, 12, 5000, 'ACTIVE', CURRENT_TIMESTAMP),
	    ('d4e5f6a7-8901-2345-6789-0abcdef12345', 'Paket Starter', 49000.00, 490000.00, 1, 50, 'INACTIVE', CURRENT_TIMESTAMP),
	    ('e5f6a7b8-9012-3456-7890-abcdef123456', 'Paket Trial', 0.00, 0.00, 1, 10, 'INACTIVE', CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO NOTHING;
	`

	log.Println("Menjalankan Seeder (Subscription)...")
	_, err = pool.Exec(ctx, seederSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan seeder: %v", err)
	}

	// 3. Seed dummy orders using real CLIENT user IDs
	ordersSeederSQL := `
	-- Wipe ALL old orders and re-seed with correct client user_ids
	DELETE FROM subscription.orders WHERE invoice_number LIKE 'INV-2026-000%';

	INSERT INTO subscription.orders (id, invoice_number, user_id, package_id, total_price, status, created_at) VALUES
	    (gen_random_uuid(), 'INV-2026-00001', '019cd5b3-d8f9-7d8d-810b-13fc26d137fa', 'a1b2c3d4-e5f6-7890-1234-567890abcdef', 150000.00,  'PAID',      '2026-01-05 10:00:00'),
	    (gen_random_uuid(), 'INV-2026-00002', '019cd65b-7a61-7047-a97b-a615cc4eb1fa', 'b2c3d4e5-f678-9012-3456-7890abcdef12', 500000.00,  'PAID',      '2026-01-12 14:30:00'),
	    (gen_random_uuid(), 'INV-2026-00003', '019cd582-9cba-7749-8924-c96b95643828', 'c3d4e5f6-7890-1234-5678-90abcdef1234', 1200000.00, 'PENDING',   '2026-02-01 09:15:00'),
	    (gen_random_uuid(), 'INV-2026-00004', '019cd5b3-d8f9-7d8d-810b-13fc26d137fa', 'b2c3d4e5-f678-9012-3456-7890abcdef12', 500000.00,  'CANCELLED', '2026-02-10 11:00:00'),
	    (gen_random_uuid(), 'INV-2026-00005', '019cd613-fe86-70be-b51e-b3bef12557a4', 'a1b2c3d4-e5f6-7890-1234-567890abcdef', 150000.00,  'PAID',      '2026-02-18 16:45:00'),
	    (gen_random_uuid(), 'INV-2026-00006', '019cd582-9cba-7749-8924-c96b95643828', 'c3d4e5f6-7890-1234-5678-90abcdef1234', 1200000.00, 'PENDING',   '2026-03-01 08:00:00'),
	    (gen_random_uuid(), 'INV-2026-00007', '019cd65b-7a61-7047-a97b-a615cc4eb1fa', 'a1b2c3d4-e5f6-7890-1234-567890abcdef', 150000.00,  'EXPIRED',   '2026-03-05 13:20:00'),
	    (gen_random_uuid(), 'INV-2026-00008', '019cd613-fe86-70be-b51e-b3bef12557a4', 'b2c3d4e5-f678-9012-3456-7890abcdef12', 500000.00,  'PAID',      '2026-03-08 17:00:00');
	`

	log.Println("Menjalankan Seeder Orders (Dummy)...")
	_, err = pool.Exec(ctx, ordersSeederSQL)
	if err != nil {
		log.Printf("[WARN] Gagal seed dummy orders: %v", err)
	}

	log.Println("Migrasi dan Seeding selesai! (Subscription)")
}
