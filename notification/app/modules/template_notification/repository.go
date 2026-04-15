package template

import (
	"context"
	"fmt"
	"time"

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

// RegisterEventType mencatat event_type ke registry (ON CONFLICT DO NOTHING).
func (r *Repository) RegisterEventType(eventType string) {
	r.db.Exec(context.Background(), `
		INSERT INTO notification_event_types (event_type)
		VALUES ($1)
		ON CONFLICT (event_type) DO NOTHING
	`, eventType)
}

// ListEventTypes mengembalikan semua event_type yang sudah terdaftar.
func (r *Repository) ListEventTypes() ([]string, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT event_type FROM notification_event_types ORDER BY registered_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var et string
		rows.Scan(&et)
		list = append(list, et)
	}
	if list == nil {
		list = []string{}
	}
	return list, nil
}

// GetByEventType mengambil template pertama yang cocok dengan event_type dan channel.
func (r *Repository) GetByEventType(eventType, channel string) (NotificationTemplate, error) {
	var t NotificationTemplate
	err := r.db.QueryRow(context.Background(), `
		SELECT id, name, event_type, channel, subject, content, created_at, updated_at, created_by, updated_by
		FROM notification_templates
		WHERE event_type = $1 AND channel = $2
		LIMIT 1
	`, eventType, channel).Scan(&t.ID, &t.Name, &t.EventType, &t.Channel, &t.Subject, &t.Content,
		&t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy)
	return t, err
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

// ── Notification Logs ─────────────────────────────────────────────────────────

// SaveLog menyimpan record log baru dengan status 'pending'.
func (r *Repository) SaveLog(id, eventType, channel, to, subject, content string) {
	r.db.Exec(context.Background(), `
		INSERT INTO notification_logs (id, event_type, channel, to_address, subject, content, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending')
	`, id, eventType, channel, to, subject, content)
}

// MarkLogSent menandai log sebagai berhasil dikirim.
func (r *Repository) MarkLogSent(id string) {
	r.db.Exec(context.Background(), `
		UPDATE notification_logs SET status='sent', sent_at=NOW(), error_msg=NULL WHERE id=$1
	`, id)
}

// MarkLogFailed menandai log gagal dan menjadwalkan retry dengan exponential backoff.
func (r *Repository) MarkLogFailed(id, errMsg string, retryCount int, nextRetryAt *time.Time) {
	r.db.Exec(context.Background(), `
		UPDATE notification_logs
		SET status='failed', error_msg=$1, retry_count=$2, next_retry_at=$3
		WHERE id=$4
	`, errMsg, retryCount, nextRetryAt, id)
}

// GetRetryableLogs mengambil log yang gagal dan siap di-retry.
func (r *Repository) GetRetryableLogs() ([]NotificationLog, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, event_type, channel, to_address, subject, content, retry_count, max_retries
		FROM notification_logs
		WHERE status = 'failed'
		  AND retry_count < max_retries
		  AND next_retry_at <= NOW()
		ORDER BY next_retry_at ASC
		LIMIT 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []NotificationLog
	for rows.Next() {
		var l NotificationLog
		var subject, content *string
		rows.Scan(&l.ID, &l.EventType, &l.Channel, &l.ToAddress, &subject, &content, &l.RetryCount, &l.MaxRetries)
		if subject != nil {
			l.Subject = *subject
		}
		if content != nil {
			l.Content = *content
		}
		list = append(list, l)
	}
	return list, nil
}

// ListLogs mengambil log notifikasi untuk monitoring, dengan filter status opsional.
func (r *Repository) ListLogs(status string, limit, offset int) ([]NotificationLog, int, error) {
	where := "WHERE 1=1"
	args := []any{}
	i := 1
	if status != "" {
		where += fmt.Sprintf(" AND status=$%d", i)
		args = append(args, status)
		i++
	}

	var total int
	r.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM notification_logs "+where, args...).Scan(&total)

	args = append(args, limit, offset)
	rows, err := r.db.Query(context.Background(), `
		SELECT id, event_type, channel, to_address, status, retry_count, max_retries,
		       next_retry_at, sent_at, error_msg, created_at
		FROM notification_logs `+where+fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, i, i+1),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []NotificationLog
	for rows.Next() {
		var l NotificationLog
		rows.Scan(&l.ID, &l.EventType, &l.Channel, &l.ToAddress, &l.Status,
			&l.RetryCount, &l.MaxRetries, &l.NextRetryAt, &l.SentAt, &l.ErrorMsg, &l.CreatedAt)
		list = append(list, l)
	}
	if list == nil {
		list = []NotificationLog{}
	}
	return list, total, nil
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
