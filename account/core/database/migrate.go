package database

import (
	"context"
	"log"
	"strings"
)

func MigrateAndSeed(db interface{}) {
	ctx := context.Background()

	// 1. Perbaikan Casting:
	// app.DB adalah *DBWrapper, jadi kita harus cast ke situ dulu
	wrapper, ok := db.(*DBWrapper)
	if !ok {
		log.Fatalf("[FATAL] Input MigrateAndSeed bukan *database.DBWrapper")
	}

	// Ambil pool aslinya dari dalam wrapper
	pool := wrapper.Pool
	if pool == nil {
		log.Printf("[WARN] Database pool is nil, skipping migration & seeding")
		return
	}

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
		photo VARCHAR(255) NULL,
		group_id UUID REFERENCES groups(id) ON DELETE SET NULL,
		level_id INT REFERENCES levels(id) ON DELETE SET NULL,
		role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
		status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('banned', 'active', 'inactive')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by UUID REFERENCES users(id) ON DELETE SET NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_by UUID REFERENCES users(id) ON DELETE SET NULL
	);
	`

	log.Println("Menjalankan Migrasi PostgreSQL...")
	_, err := pool.Exec(ctx, migrationSQL)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "already exists") {
			log.Printf("[WARN] Migrasi race condition (safe to ignore): %v", err)
		} else {
			log.Fatalf("Gagal menjalankan migrasi: %v", err)
		}
	}

	// 2. Query Seeder (Insert jika belum ada / UPSERT)
	seederSQL := `
	-- Seed Level
	INSERT INTO levels (id, name, code) 
	VALUES (1, 'Super Admin', 'SUPERADMIN') 
	ON CONFLICT (id) DO NOTHING;

	-- Seed Group
	INSERT INTO groups (id, name, code) 
	VALUES ('1e98c63f-5474-4506-826c-ded22b59b3db', 'Owner', 'Owner') 
	ON CONFLICT (id) DO NOTHING;

	-- Seed Role
	INSERT INTO roles (id, name, code, group_id) 
	VALUES ('cf47ce1c-1455-4a20-bafe-c2b7c2ab9993', 'CEO', 'CEO', '1e98c63f-5474-4506-826c-ded22b59b3db') 
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
	ON CONFLICT (email) DO NOTHING; -- Mencegah duplikasi berdasarkan email
	`

	log.Println("Menjalankan Seeder...")
	_, err = pool.Exec(ctx, seederSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan seeder: %v", err)
	}

	log.Println("Migrasi dan Seeding selesai!")
}
