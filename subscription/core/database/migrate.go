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
		total_price DECIMAL(15,2) NOT NULL DEFAULT 0,
		status VARCHAR(50) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PAID', 'CANCELLED', 'EXPIRED')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
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

	// 3. Seed dummy orders (for demo/testing)
	ordersSeederSQL := `
	INSERT INTO subscription.orders (id, invoice_number, user_id, package_id, total_price, status, created_at)
	VALUES
	    ('f1a2b3c4-d5e6-7890-aaaa-111111111111', 'INV-2026-00001', 'aaaaaaaa-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-7890-1234-567890abcdef', 150000.00, 'PAID',      '2026-01-05 10:00:00'),
	    ('f2a2b3c4-d5e6-7890-bbbb-222222222222', 'INV-2026-00002', 'aaaaaaaa-0000-0000-0000-000000000002', 'b2c3d4e5-f678-9012-3456-7890abcdef12', 500000.00, 'PAID',      '2026-01-12 14:30:00'),
	    ('f3a2b3c4-d5e6-7890-cccc-333333333333', 'INV-2026-00003', 'aaaaaaaa-0000-0000-0000-000000000003', 'c3d4e5f6-7890-1234-5678-90abcdef1234', 1200000.00,'PENDING',   '2026-02-01 09:15:00'),
	    ('f4a2b3c4-d5e6-7890-dddd-444444444444', 'INV-2026-00004', 'aaaaaaaa-0000-0000-0000-000000000001', 'b2c3d4e5-f678-9012-3456-7890abcdef12', 500000.00, 'CANCELLED', '2026-02-10 11:00:00'),
	    ('f5a2b3c4-d5e6-7890-eeee-555555555555', 'INV-2026-00005', 'aaaaaaaa-0000-0000-0000-000000000004', 'a1b2c3d4-e5f6-7890-1234-567890abcdef', 150000.00, 'PAID',      '2026-02-18 16:45:00'),
	    ('f6a2b3c4-d5e6-7890-ffff-666666666666', 'INV-2026-00006', 'aaaaaaaa-0000-0000-0000-000000000005', 'c3d4e5f6-7890-1234-5678-90abcdef1234', 1200000.00,'PENDING',   '2026-03-01 08:00:00'),
	    ('f7a2b3c4-d5e6-7890-aaaa-777777777777', 'INV-2026-00007', 'aaaaaaaa-0000-0000-0000-000000000002', 'a1b2c3d4-e5f6-7890-1234-567890abcdef', 150000.00, 'EXPIRED',   '2026-03-05 13:20:00'),
	    ('f8a2b3c4-d5e6-7890-bbbb-888888888888', 'INV-2026-00008', 'aaaaaaaa-0000-0000-0000-000000000003', 'b2c3d4e5-f678-9012-3456-7890abcdef12', 500000.00, 'PAID',      '2026-03-08 17:00:00')
	ON CONFLICT (id) DO NOTHING;
	`

	log.Println("Menjalankan Seeder Orders (Dummy)...")
	_, err = pool.Exec(ctx, ordersSeederSQL)
	if err != nil {
		log.Printf("[WARN] Gagal seed dummy orders: %v", err)
	}

	log.Println("Migrasi dan Seeding selesai! (Subscription)")
}
