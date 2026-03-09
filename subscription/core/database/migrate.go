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

	log.Println("Migrasi dan Seeding selesai! (Subscription)")
}
