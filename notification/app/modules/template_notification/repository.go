package template

import (
	"context"
	"fmt"

	"notification/core/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository menyediakan akses query langsung ke table notification_templates.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository membuat instance Repository baru.
func NewRepository() *Repository {
	return &Repository{db: database.DB}
}

// List mengambil semua template, bisa difilter by channel dan/atau event_type.
func (r *Repository) List(channel, eventType string) ([]NotificationTemplate, error) {
	query := `
		SELECT id, name, event_type, channel, subject, content, created_at, updated_at, created_by, updated_by
		FROM notification_templates WHERE 1=1`
	args := []any{}
	i := 1
	if channel != "" {
		query += fmt.Sprintf(" AND channel = $%d", i)
		args = append(args, channel)
		i++
	}
	if eventType != "" {
		query += fmt.Sprintf(" AND event_type = $%d", i)
		args = append(args, eventType)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []NotificationTemplate
	for rows.Next() {
		var t NotificationTemplate
		rows.Scan(&t.ID, &t.Name, &t.EventType, &t.Channel, &t.Subject, &t.Content,
			&t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy)
		list = append(list, t)
	}
	if list == nil {
		list = []NotificationTemplate{}
	}
	return list, nil
}

// GetByID mengambil satu template berdasarkan ID.
func (r *Repository) GetByID(id string) (NotificationTemplate, error) {
	var t NotificationTemplate
	err := r.db.QueryRow(context.Background(), `
		SELECT id, name, event_type, channel, subject, content, created_at, updated_at, created_by, updated_by
		FROM notification_templates WHERE id = $1
	`, id).Scan(&t.ID, &t.Name, &t.EventType, &t.Channel, &t.Subject, &t.Content,
		&t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy)
	return t, err
}

// Exists mengecek apakah record dengan ID tertentu ada di database.
func (r *Repository) Exists(id string) bool {
	var exists bool
	r.db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM notification_templates WHERE id = $1)", id,
	).Scan(&exists)
	return exists
}

// Create menyimpan template baru ke database.
func (r *Repository) Create(req CreateTemplateRequest, id string) error {
	var subject *string
	if req.Subject != "" {
		subject = &req.Subject
	}
	var createdBy *string
	if req.CreatedBy != "" {
		createdBy = &req.CreatedBy
	}
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, req.Name, req.EventType, req.Channel, subject, req.Content, createdBy)
	return err
}

// Update memperbarui data template berdasarkan ID.
func (r *Repository) Update(id string, req UpdateTemplateRequest) error {
	var subject *string
	if req.Subject != "" {
		subject = &req.Subject
	}
	var updatedBy *string
	if req.UpdatedBy != "" {
		updatedBy = &req.UpdatedBy
	}
	_, err := r.db.Exec(context.Background(), `
		UPDATE notification_templates SET
			name       = COALESCE(NULLIF($2,''), name),
			event_type = COALESCE(NULLIF($3,''), event_type),
			channel    = COALESCE(NULLIF($4,''), channel),
			subject    = $5,
			content    = COALESCE(NULLIF($6,''), content),
			updated_by = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id, req.Name, req.EventType, req.Channel, subject, req.Content, updatedBy)
	return err
}

// Delete menghapus template berdasarkan ID.
func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM notification_templates WHERE id = $1", id)
	return err
}
