package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

// ==========================================
// 1. DATA TRANSFER OBJECTS (DTO)
// ==========================================

// ==========================================
// 2. REPOSITORY INTERFACE
// ==========================================

type Repository interface {
	// Hanya fungsi untuk listing user inactive
	GetUsersInactive(ctx context.Context, page, limit int, search string) ([]UserListDTO, int, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, id string, data UserUpdateDTO) error
}

// ==========================================
// 3. REPOSITORY IMPLEMENTATION
// ==========================================

type userRepo struct {
	db *database.DBWrapper
}

func NewRepository(db *database.DBWrapper) Repository {
	return &userRepo{db: db}
}

func (r *userRepo) FindUserByID(ctx context.Context, id string) (*User, error) {
	var u User

	// PERBAIKAN:
	// Saya tambahkan COALESCE pada: multiple_login, jenis_kelamin, dan jenis_akun
	query := `
        SELECT 
            u.id, u.fullname, u.level_id, u.role_id, u.group_id, u.email,
            COALESCE(u.nrk, ''), COALESCE(u.nik, ''), COALESCE(u.npsn, ''), 
            COALESCE(u.nikki, ''), COALESCE(u.nisn, ''), COALESCE(u.no_ip, ''), 
            COALESCE(u.photo, ''), 
            
            -- PERBAIKAN DISINI (Bungkus semua string nullable)
            COALESCE(u.multiple_login, 0), 
            COALESCE(u.jenis_kelamin, ''), 
			COALESCE(u.dokumen_pendukung, ''), 
            COALESCE(u.jenis_akun, ''), 

            u.is_active::boolean,
            
            COALESCE(l.name, '') as level_name, 
            COALESCE(r.name, '') as role_name, 
            
            CASE 
                WHEN LENGTH(u.group_id) = 8 THEN COALESCE(s.nama_satuan_pendidikan, u.group_id)
                ELSE COALESCE(g.name, '-')
            END as group_name

        FROM account.za_users u
        LEFT JOIN account.za_groups g ON g.id = u.group_id
        LEFT JOIN account.sekolah s ON s.npsn = u.group_id
        LEFT JOIN account.za_levels l ON l.id = u.level_id
        LEFT JOIN account.za_roles r ON r.id = u.role_id
        WHERE u.id = $1 LIMIT 1
    `

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Fullname, &u.LevelID, &u.RoleID, &u.GroupID, &u.Email,
		&u.Nrk, &u.Nik, &u.Npsn, &u.Nikki, &u.Nisn, &u.Nip,
		&u.Photo,
		&u.MultipleLogin, &u.JenisKelamin, &u.DokumenPendukun, &u.JenisAkun, // Sekarang aman dari NULL
		&u.IsActive,
		&u.Level, &u.Role, &u.Group,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user tidak ditemukan")
		}
		return nil, err
	}
	return &u, nil
}

// GetUsersInactive mengambil data user dengan status pending/inactive
// Logika Group: Jika group_id 8 karakter (NPSN) ambil dari table sekolah, jika tidak ambil dari za_groups
func (r *userRepo) GetUsersInactive(ctx context.Context, page, limit int, search string) ([]UserListDTO, int, error) {
	var users []UserListDTO
	var total int

	// 1. Hitung Offset untuk Pagination
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// 2. Base Query
	// Kita join ke tabel sekolah (s) dan za_groups (g) sekaligus
	baseQuery := `
        FROM account.za_users u
        LEFT JOIN account.za_groups g ON g.id = u.group_id
        LEFT JOIN account.sekolah s ON s.npsn = u.group_id
        LEFT JOIN account.za_roles r ON r.id = u.role_id
        LEFT JOIN account.za_levels l ON l.id = u.level_id
        WHERE u.is_active = '0' AND u.level_id = 3 AND LENGTH(TRIM(u.group_id)) = 8
    `

	// 3. Dynamic Search Logic
	var searchArgs []interface{}
	var whereClause string

	if search != "" {
		// Kita mencari di tabel user dan nama-nama di tabel relasi (Case Insensitive)
		// Tambahkan pencarian ke s.nama_satuan_pendidikan
		whereClause = ` AND (
            u.fullname ILIKE $1 OR
            u.nik ILIKE $1 OR
            u.nrk ILIKE $1 OR
            u.nikki ILIKE $1 OR
            u.npsn ILIKE $1 OR
            g.name ILIKE $1 OR
            s.nama_satuan_pendidikan ILIKE $1 OR
            r.name ILIKE $1 OR
            l.name ILIKE $1
        )`
		searchTerm := "%" + search + "%"
		searchArgs = append(searchArgs, searchTerm)
	}

	// ---------------------------------------------------------
	// QUERY 1: Hitung Total Data (Untuk Pagination Frontend)
	// ---------------------------------------------------------
	countQuery := "SELECT count(*) " + baseQuery + whereClause

	// Eksekusi count
	err := r.db.Pool.QueryRow(ctx, countQuery, searchArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data: %w", err)
	}

	// Jika total 0, langsung return array kosong (hemat resource)
	if total == 0 {
		return []UserListDTO{}, 0, nil
	}

	// ---------------------------------------------------------
	// QUERY 2: Ambil Data Sebenarnya
	// ---------------------------------------------------------
	// Perhatikan logic CASE WHEN pada pemilihan Group Name
	selectQuery := `
        SELECT 
            u.id, 
            u.fullname, 
            COALESCE(u.nrk, '-'), 
            COALESCE(u.nik, '-'), 
            COALESCE(u.nikki, '-'), 
            COALESCE(u.npsn, '-'),
            COALESCE(u.status, 'Pending'),
            u.is_active::boolean,
            
            -- LOGIKA GROUP ID (8 Karakter = NPSN Sekolah, Selain itu = Group Biasa)
            CASE 
                WHEN LENGTH(u.group_id) = 8 THEN COALESCE(s.nama_satuan_pendidikan, u.group_id)
                ELSE COALESCE(g.name, '-')
            END as group_name,

            COALESCE(r.name, '-') as role_name,
            COALESCE(l.name, '-') as level_name
    ` + baseQuery + whereClause + ` ORDER BY u.created_at DESC LIMIT $` + fmt.Sprint(len(searchArgs)+1) + ` OFFSET $` + fmt.Sprint(len(searchArgs)+2)

	// Tambahkan limit & offset ke arguments
	args := append(searchArgs, limit, offset)

	rows, err := r.db.Pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u UserListDTO
		// Scan harus urut sesuai urutan SELECT di atas
		err := rows.Scan(
			&u.ID,
			&u.Fullname,
			&u.NRK,
			&u.NIK,
			&u.NIKKI,
			&u.NPSN,
			&u.Status,   // Status Text
			&u.IsActive, // Boolean true/false
			&u.GroupName,
			&u.RoleName,
			&u.LevelName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal scan row user: %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, id string, data UserUpdateDTO) error {
	// Query Update
	query := `
        UPDATE account.za_users 
        SET 
            fullname = $1, 
            nik = $2, 
            nrk = $3, 
            is_active = $4,
            updated_at = NOW() -- Selalu update timestamp
        WHERE id = $5
    `

	// Konversi bool ke string '1'/'0' jika database Anda pakai char(1) untuk boolean
	// Tapi karena sebelumnya kita pakai ::boolean, saya asumsikan kolomnya support true/false atau '1'/'0'
	// Postgres native boolean menerima true/false.
	// Jika tipe kolom Anda char/varchar/int, sesuaikan valuenya.
	var activeVal string = "0"
	if data.IsActive {
		activeVal = "1"
	}

	cmdTag, err := r.db.Pool.Exec(ctx, query,
		data.Fullname,
		data.NIK,
		data.NRK,
		activeVal,
		id,
	)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or no changes made")
	}

	return nil
}
