package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

type Repository interface {
	FindUser(ctx context.Context, key, value string) (*User, error)
	FindRoleByCode(ctx context.Context, code string) (*RoleInfo, error)
	GetAllRoles(ctx context.Context) ([]RoleInfo, error)
}

type userRepo struct {
	db *database.DBWrapper
}

func (r *userRepo) ensureDBReady() error {
	if r == nil || r.db == nil || r.db.Pool == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return nil
}

func NewRepository(db *database.DBWrapper) Repository {
	return &userRepo{db: db}
}

func (r *userRepo) FindUser(ctx context.Context, key, value string) (*User, error) {
	if err := r.ensureDBReady(); err != nil {
		return nil, err
	}

	var u User

	query := fmt.Sprintf(`
        SELECT 
            u.id, u.name, u.email, u.password, u.photo, 
            u.group_id, u.level_id, u.role_id, u.status,
            u.created_at, u.created_by, u.updated_at, u.updated_by,
            l.name as level_name, l.code as level_code,
            r.name as role_name, r.code as role_code,
            g.name as group_name, g.code as group_code
        FROM users u
        JOIN groups g ON g.id = u.group_id
        JOIN levels l ON l.id = u.level_id
        JOIN roles r ON r.id = u.role_id
        WHERE u.%s = $1 LIMIT 1
    `, key)

	err := r.db.Pool.QueryRow(ctx, query, value).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Photo,
		&u.GroupID, &u.LevelID, &u.RoleID, &u.Status,
		&u.CreatedAt, &u.CreatedBy, &u.UpdatedAt, &u.UpdatedBy,
		&u.Level, &u.LevelCode,
		&u.Role, &u.RoleCode,
		&u.Group, &u.GroupCode,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindRoleByCode mengambil role berdasarkan code
func (r *userRepo) FindRoleByCode(ctx context.Context, code string) (*RoleInfo, error) {
	if err := r.ensureDBReady(); err != nil {
		return nil, err
	}

	var role RoleInfo
	err := r.db.Pool.QueryRow(ctx,
		"SELECT id, name, code FROM roles WHERE UPPER(code) = UPPER($1)", code,
	).Scan(&role.ID, &role.Name, &role.Code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetAllRoles mengambil semua role yang tersedia
func (r *userRepo) GetAllRoles(ctx context.Context) ([]RoleInfo, error) {
	if err := r.ensureDBReady(); err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, "SELECT id, name, code FROM roles ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []RoleInfo
	for rows.Next() {
		var role RoleInfo
		if err := rows.Scan(&role.ID, &role.Name, &role.Code); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
