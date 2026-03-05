package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

type Repository interface {
	FindUser(ctx context.Context, key, value string) (*User, error)
}

type userRepo struct {
	db *database.DBWrapper
}

func NewRepository(db *database.DBWrapper) Repository {
	return &userRepo{db: db}
}

func (r *userRepo) FindUser(ctx context.Context, key, value string) (*User, error) {
	var u User

	// Query disesuaikan dengan kolom tabel 'users' yang baru
	query := fmt.Sprintf(`
        SELECT 
            u.id, u.name, u.email, u.password, u.photo, 
            u.group_id, u.level_id, u.role_id, u.status,
            u.created_at, u.created_by, u.updated_at, u.updated_by,
            l.name as level_name, 
            r.name as role_name, 
            g.name as group_name
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
