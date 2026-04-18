package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

type Repository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	FindRoleByCode(ctx context.Context, code string) (string, error)
	FindDefaultGroupAndLevel(ctx context.Context) (groupID string, levelID string, err error)
	CreateUser(ctx context.Context, id, name, email, hashedPassword, roleID, groupID, levelID string) error
	GetInternalUsers(ctx context.Context, params GetUsersParams) ([]UserListItem, error)
	CountInternalUsers(ctx context.Context, params GetUsersParams) (int, error)
}

type adminRepo struct {
	db *database.DBWrapper
}

func NewRepository(db *database.DBWrapper) Repository {
	return &adminRepo{db: db}
}

// EmailExists mengecek apakah email sudah terdaftar di database
func (r *adminRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM users WHERE email = $1", email,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindRoleByCode mengambil role ID berdasarkan role code
func (r *adminRepo) FindRoleByCode(ctx context.Context, code string) (string, error) {
	var roleID string
	err := r.db.Pool.QueryRow(ctx,
		"SELECT id FROM roles WHERE UPPER(code) = UPPER($1)", code,
	).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return roleID, nil
}

// FindDefaultGroupAndLevel mengambil default group dan level untuk user internal
// Menggunakan group 'INTERNAL' dan level 'STAFF' sebagai default.
// Fallback ke group/level pertama jika tidak ditemukan.
func (r *adminRepo) FindDefaultGroupAndLevel(ctx context.Context) (string, string, error) {
	var groupID, levelID string

	// Cari group INTERNAL, fallback ke group pertama
	err := r.db.Pool.QueryRow(ctx,
		"SELECT id FROM groups WHERE UPPER(code) = 'INTERNAL' LIMIT 1",
	).Scan(&groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Fallback: ambil group pertama yang ada
			err = r.db.Pool.QueryRow(ctx, "SELECT id FROM groups ORDER BY id LIMIT 1").Scan(&groupID)
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}

	// Cari level STAFF, fallback ke level pertama
	err = r.db.Pool.QueryRow(ctx,
		"SELECT id FROM levels WHERE UPPER(code) = 'STAFF' LIMIT 1",
	).Scan(&levelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Fallback: ambil level pertama yang ada
			err = r.db.Pool.QueryRow(ctx, "SELECT id FROM levels ORDER BY id LIMIT 1").Scan(&levelID)
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}

	return groupID, levelID, nil
}

// CreateUser menyimpan user baru ke database dengan status 'active'
func (r *adminRepo) CreateUser(ctx context.Context, id, name, email, hashedPassword, roleID, groupID, levelID string) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO users (id, name, email, password, group_id, level_id, role_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, id, name, email, hashedPassword, groupID, levelID, roleID)
	return err
}

// buildInternalUsersQuery constructs the query and arguments for listing/counting users
func buildInternalUsersQuery(isCount bool, params GetUsersParams) (string, []any) {
	selectClause := "SELECT u.id, u.name, u.email, r.code as role, u.status, u.created_at, u.updated_at"
	if isCount {
		selectClause = "SELECT COUNT(*)"
	}

	query := selectClause + `
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE 1=1
	`

	args := []any{}
	argID := 1

	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.email ILIKE $%d)", argID, argID+1)
		args = append(args, searchTerm, searchTerm)
		argID += 2
	}

	if params.Role != "" {
		query += fmt.Sprintf(" AND UPPER(r.code) = UPPER($%d)", argID)
		args = append(args, params.Role)
		argID++
	}

	if params.Status != "" {
		query += fmt.Sprintf(" AND UPPER(u.status) = UPPER($%d)", argID)
		args = append(args, params.Status)
		argID++
	}

	return query, args
}

// GetInternalUsers mengambil daftar akun internal berdasarkan filter & paginasi
func (r *adminRepo) GetInternalUsers(ctx context.Context, params GetUsersParams) ([]UserListItem, error) {
	query, args := buildInternalUsersQuery(false, params)

	// Sorting and pagination
	query += " ORDER BY u.created_at DESC"
	
	if params.PerPage > 0 {
		offset := (params.Page - 1) * params.PerPage
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.PerPage, offset)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserListItem
	for rows.Next() {
		var u UserListItem
		if err := rows.Scan(&u.UserID, &u.FullName, &u.Email, &u.Role, &u.Status, &u.CreatedAt, &u.LastLoginAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// CountInternalUsers menghitung total baris untuk paginasi
func (r *adminRepo) CountInternalUsers(ctx context.Context, params GetUsersParams) (int, error) {
	query, args := buildInternalUsersQuery(true, params)
	
	var count int
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
