package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func seedDashboardClients(ctx context.Context, pool *pgxpool.Pool) {
	clientPwd := hashPwd("Client123")

	// Hapus seed lama agar id tetap sinkron dengan subscription.orders.user_id
	if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id::text LIKE 'f1010001-%'`); err != nil {
		log.Printf("[WARN] Gagal cleanup dashboard clients: %v", err)
	}

	const q = `
	INSERT INTO users (id, name, email, password, phone, group_id, level_id, role_id, status, created_at, created_by)
	VALUES ($1, $2, $3, $4, $5, '2e98c63f-5474-4506-826c-ded22b59b3dc', 2, 'df47ce1c-1455-4a20-bafe-c2b7c2ab9994', 'active', CURRENT_TIMESTAMP, '10ef7bff-4c69-4b56-aec8-ef7427601952')
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		email = EXCLUDED.email,
		password = EXCLUDED.password,
		phone = EXCLUDED.phone,
		status = 'active'
	`

	clients := []struct {
		id, name, email, phone string
	}{
		{"f1010001-0001-4000-8000-000000000001", "Budi Santoso", "budi.santoso@mail.com", "081210010001"},
		{"f1010001-0001-4000-8000-000000000002", "Siti Rahayu", "siti.rahayu@mail.com", "081210010002"},
		{"f1010001-0001-4000-8000-000000000003", "Andi Wijaya", "andi.wijaya@mail.com", "081210010003"},
		{"f1010001-0001-4000-8000-000000000004", "Dewi Lestari", "dewi.lestari@mail.com", "081210010004"},
		{"f1010001-0001-4000-8000-000000000005", "Rizky Pratama", "rizky.pratama@mail.com", "081210010005"},
		{"f1010001-0001-4000-8000-000000000006", "Nina Kusuma", "nina.kusuma@mail.com", "081210010006"},
		{"f1010001-0001-4000-8000-000000000007", "Fajar Nugroho", "fajar.nugroho@mail.com", "081210010007"},
		{"f1010001-0001-4000-8000-000000000008", "Maya Sari", "maya.sari@mail.com", "081210010008"},
		{"f1010001-0001-4000-8000-000000000009", "Doni Hermawan", "doni.hermawan@mail.com", "081210010009"},
		{"f1010001-0001-4000-8000-000000000010", "Putri Anggraini", "putri.anggraini@mail.com", "081210010010"},
		{"f1010001-0001-4000-8000-000000000011", "Agus Setiawan", "agus.setiawan@mail.com", "081210010011"},
		{"f1010001-0001-4000-8000-000000000012", "Rina Melati", "rina.melati@mail.com", "081210010012"},
		{"f1010001-0001-4000-8000-000000000013", "Eko Prasetyo", "eko.prasetyo@mail.com", "081210010013"},
		{"f1010001-0001-4000-8000-000000000014", "Linda Wijaya", "linda.wijaya@mail.com", "081210010014"},
		{"f1010001-0001-4000-8000-000000000015", "Hendra Gunawan", "hendra.gunawan@mail.com", "081210010015"},
		{"f1010001-0001-4000-8000-000000000016", "Reza Akbar", "reza.akbar@mail.com", "081210010016"},
		{"f1010001-0001-4000-8000-000000000017", "Citra Dewi", "citra.dewi@mail.com", "081210010017"},
		{"f1010001-0001-4000-8000-000000000018", "Bambang Hartono", "bambang.hartono@mail.com", "081210010018"},
		{"f1010001-0001-4000-8000-000000000019", "Wulan Sari", "wulan.sari@mail.com", "081210010019"},
		{"f1010001-0001-4000-8000-000000000020", "Yoga Pratama", "yoga.pratama@mail.com", "081210010020"},
	}

	log.Println("Menjalankan Seeder (dashboard client accounts)...")
	for _, c := range clients {
		if _, err := pool.Exec(ctx, q, c.id, c.name, c.email, clientPwd, c.phone); err != nil {
			log.Printf("[WARN] Gagal seed client %s: %v", c.email, err)
		}
	}
	log.Println("  ✓ 20 client accounts untuk management dashboard (password: Client123)")
}
