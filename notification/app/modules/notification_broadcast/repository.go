package notification

import (
	"context"
	"time"

	"notification/core/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository menyediakan akses query langsung ke table notifications.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository membuat instance Repository baru.
func NewRepository() *Repository {
	return &Repository{db: database.DB}
}

// List mengambil semua notification, diurutkan terbaru.
func (r *Repository) List() ([]Notification, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, title, message, type, target_role, is_active, created_at, created_by, updated_at
		FROM notifications
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Notification
	for rows.Next() {
		var n Notification
		rows.Scan(&n.ID, &n.Title, &n.Message, &n.Type, &n.TargetRole,
			&n.IsActive, &n.CreatedAt, &n.CreatedBy, &n.UpdatedAt)
		list = append(list, n)
	}
	if list == nil {
		list = []Notification{}
	}
	return list, nil
}

// ListPublic mengambil notification aktif sesuai role, max 20 data.
func (r *Repository) ListPublic(role string) ([]map[string]any, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, title, message, type, target_role, created_at
		FROM notifications
		WHERE is_active = TRUE AND (target_role = 'all' OR target_role = $1)
		ORDER BY created_at DESC LIMIT 20
	`, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]any
	for rows.Next() {
		var id, title, msg, typ, target string
		var createdAt time.Time
		rows.Scan(&id, &title, &msg, &typ, &target, &createdAt)
		list = append(list, map[string]any{
			"id": id, "title": title, "message": msg,
			"type": typ, "target_role": target, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]any{}
	}
	return list, nil
}

// GetByID mengambil satu notification berdasarkan ID.
func (r *Repository) GetByID(id string) (Notification, error) {
	var n Notification
	err := r.db.QueryRow(context.Background(), `
		SELECT id, title, message, type, target_role, is_active, created_at, created_by, updated_at
		FROM notifications WHERE id = $1
	`, id).Scan(&n.ID, &n.Title, &n.Message, &n.Type, &n.TargetRole,
		&n.IsActive, &n.CreatedAt, &n.CreatedBy, &n.UpdatedAt)
	return n, err
}

// Create menyimpan notification baru ke database.
func (r *Repository) Create(req CreateNotificationRequest, id string) error {
	typ := req.Type
	if typ == "" {
		typ = "info"
	}
	targetRole := req.TargetRole
	if targetRole == "" {
		targetRole = "all"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	var createdBy *string
	if req.CreatedBy != "" {
		createdBy = &req.CreatedBy
	}

	_, err := r.db.Exec(context.Background(), `
		INSERT INTO notifications (id, title, message, type, target_role, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, req.Title, req.Message, typ, targetRole, isActive, createdBy)
	return err
}

// Update memperbarui data notification berdasarkan ID.
func (r *Repository) Update(id string, req UpdateNotificationRequest) error {
	var updatedBy *string
	if req.UpdatedBy != "" {
		updatedBy = &req.UpdatedBy
	}
	_, err := r.db.Exec(context.Background(), `
		UPDATE notifications SET
			title       = COALESCE(NULLIF($2,''), title),
			message     = COALESCE(NULLIF($3,''), message),
			type        = COALESCE(NULLIF($4,''), type),
			target_role = COALESCE(NULLIF($5,''), target_role),
			is_active   = COALESCE($6, is_active),
			updated_by  = $7,
			updated_at  = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id, req.Title, req.Message, req.Type, req.TargetRole, req.IsActive, updatedBy)
	return err
}

// Delete menghapus notification berdasarkan ID.
func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM notifications WHERE id = $1", id)
	return err
}
