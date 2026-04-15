package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

type Repository interface {
	FindUserByEmail(ctx context.Context, key, value string) (*User, error)
}

type userRepo struct {
	db *database.DBWrapper
}

func NewRepository(db *database.DBWrapper) Repository {
	return &userRepo{db: db}
}

// 1. Query Utama (account.za_users join groups, levels, roles)
func (r *userRepo) FindUserByEmail(ctx context.Context, key, value string) (*User, error) {
	var u User

	// FIX: Tambahkan casting '::boolean' pada u.is_active
	// Postgres otomatis mengubah "1", "true", "t", "yes", "on" menjadi BOOLEAN TRUE
	query := fmt.Sprintf(`
		SELECT 
			u.id, u.fullname, u.level_id, u.role_id, u.group_id, u.password, u.email,
			u.nrk, u.nik, u.npsn, u.nikki, u.nisn, u.no_ip, u.photo, u.multiple_login, u.jenis_kelamin, u.jenis_akun, 
            u.is_active::boolean, -- <--- PERBAIKAN DISINI
			l.name as level_name, r.name as role_name, g.name as group_name
		FROM account.za_users u
		JOIN account.za_groups g ON g.id = u.group_id
		JOIN account.za_levels l ON l.id = u.level_id
		JOIN account.za_roles r ON r.id = u.role_id
		WHERE u.%s = $1 LIMIT 1
	`, key)

	err := r.db.Pool.QueryRow(ctx, query, value).Scan(
		&u.ID, &u.Fullname, &u.LevelID, &u.RoleID, &u.GroupID, &u.Password, &u.Email,
		&u.Nrk, &u.Nik, &u.Npsn, &u.Nikki, &u.Nisn, &u.Nip, &u.Photo, &u.MultipleLogin, &u.JenisKelamin, &u.JenisAkun, &u.IsActive,
		&u.Level, &u.Role, &u.Group,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
