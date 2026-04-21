package notification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"notification/core/database"

	"github.com/jackc/pgx/v5"
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

// List mengambil semua notification, diurutkan terbaru, dengan filter opsional type & status.
func (r *Repository) List(typeFilter, statusFilter string) ([]Notification, error) {
	ctx := context.Background()
	query := `
		SELECT
			id,
			title,
			COALESCE(description, message, '') AS description,
			COALESCE(message, description, '') AS message,
			type,
			target_role,
			cta_url,
			image_url,
			expiry_date,
			is_active,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM notifications
		WHERE 1=1
	`

	args := make([]any, 0)
	argPos := 1

	if typeFilter != "" {
		query += fmt.Sprintf(" AND LOWER(type) = $%d", argPos)
		args = append(args, strings.ToLower(typeFilter))
		argPos++
	}

	switch strings.ToLower(statusFilter) {
	case "active", "true", "1":
		query += " AND is_active = TRUE AND (expiry_date IS NULL OR expiry_date > NOW())"
	case "inactive", "false", "0":
		query += " AND is_active = FALSE"
	case "expired":
		query += " AND expiry_date IS NOT NULL AND expiry_date <= NOW()"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Notification
	for rows.Next() {
		var n Notification
		err = rows.Scan(
			&n.ID,
			&n.Title,
			&n.Description,
			&n.Message,
			&n.Type,
			&n.TargetRole,
			&n.CTAURL,
			&n.ImageURL,
			&n.ExpiryDate,
			&n.IsActive,
			&n.CreatedAt,
			&n.CreatedBy,
			&n.UpdatedAt,
			&n.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		n.Status = deriveStatus(n.IsActive, n.ExpiryDate)
		list = append(list, n)
	}
	if list == nil {
		list = []Notification{}
	}
	return list, nil
}

// ListPublic mengambil notification aktif sesuai audience, max 20 data.
func (r *Repository) ListPublic(role, userID string) ([]map[string]any, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role != "" && role != "client" {
		return []map[string]any{}, nil
	}

	audiences := map[string]struct{}{"client": {}}
	if segments, err := r.resolveClientAudienceSegments(strings.TrimSpace(userID)); err == nil {
		for _, segment := range segments {
			segment = strings.ToLower(strings.TrimSpace(segment))
			if segment == "" {
				continue
			}
			audiences[segment] = struct{}{}
		}
	}

	audienceKeys := make([]string, 0, len(audiences))
	for key := range audiences {
		audienceKeys = append(audienceKeys, key)
	}

	rows, err := r.db.Query(context.Background(), `
		SELECT id, title, COALESCE(description, message, '') AS description, type, target_role,
		       cta_url, image_url, expiry_date, created_at
		FROM notifications
		WHERE is_active = TRUE
		  AND (expiry_date IS NULL OR expiry_date > NOW())
		  AND LOWER(target_role) = ANY($1)
		ORDER BY created_at DESC LIMIT 20
	`, audienceKeys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]any
	for rows.Next() {
		var id, title, desc, typ, target string
		var ctaURL, imageURL *string
		var expiryDate *time.Time
		var createdAt time.Time
		err = rows.Scan(&id, &title, &desc, &typ, &target, &ctaURL, &imageURL, &expiryDate, &createdAt)
		if err != nil {
			return nil, err
		}
		list = append(list, map[string]any{
			"id":          id,
			"title":       title,
			"description": desc,
			"message":     desc,
			"type":        typ,
			"target_role": target,
			"cta_url":     ctaURL,
			"image_url":   imageURL,
			"expiry_date": expiryDate,
			"created_at":  createdAt,
		})
	}
	if list == nil {
		list = []map[string]any{}
	}
	return list, nil
}

func (r *Repository) resolveClientAudienceSegments(userID string) ([]string, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil
	}

	var isExpiringSoon bool
	err := r.db.QueryRow(context.Background(), `
		SELECT (s.end_date <= CURRENT_DATE + INTERVAL '7 days') AS is_expiring_soon
		FROM subscription.subscriptions s
		WHERE s.user_id = $1
		  AND s.status = 'ACTIVE'
		  AND s.end_date >= CURRENT_DATE
		ORDER BY s.end_date ASC, s.created_at DESC
		LIMIT 1
	`, userID).Scan(&isExpiringSoon)
	if err == nil {
		if isExpiringSoon {
			return []string{"client_expiring_soon"}, nil
		}
		return []string{"client_paid_active"}, nil
	}
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	var orderStatus string
	err = r.db.QueryRow(context.Background(), `
		SELECT COALESCE(status, '')
		FROM subscription.orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID).Scan(&orderStatus)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	orderStatus = strings.ToUpper(strings.TrimSpace(orderStatus))
	if orderStatus == "PENDING_PAYMENT" || orderStatus == "CANCELLED" {
		return []string{"client_never_bought"}, nil
	}

	var latestSubscriptionExpired bool
	err = r.db.QueryRow(context.Background(), `
		SELECT (status IN ('EXPIRED', 'CANCELLED') OR end_date < CURRENT_DATE) AS is_expired
		FROM subscription.subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID).Scan(&latestSubscriptionExpired)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
	} else if latestSubscriptionExpired {
		return []string{"client_lapsed"}, nil
	}

	if orderStatus == "PAID" {
		return []string{"client_paid_active"}, nil
	}
	if orderStatus == "" {
		return []string{"client_never_bought"}, nil
	}
	return []string{"client_never_bought"}, nil
}

// GetByID mengambil satu notification berdasarkan ID.
func (r *Repository) GetByID(id string) (Notification, error) {
	var n Notification
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id,
			title,
			COALESCE(description, message, '') AS description,
			COALESCE(message, description, '') AS message,
			type,
			target_role,
			cta_url,
			image_url,
			expiry_date,
			is_active,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM notifications WHERE id = $1
	`, id).Scan(
		&n.ID,
		&n.Title,
		&n.Description,
		&n.Message,
		&n.Type,
		&n.TargetRole,
		&n.CTAURL,
		&n.ImageURL,
		&n.ExpiryDate,
		&n.IsActive,
		&n.CreatedAt,
		&n.CreatedBy,
		&n.UpdatedAt,
		&n.UpdatedBy,
	)
	if err == nil {
		n.Status = deriveStatus(n.IsActive, n.ExpiryDate)
	}
	return n, err
}

// Create menyimpan notification baru ke database.
func (r *Repository) Create(req CreateNotificationRequest, id string) error {
	desc := strings.TrimSpace(req.Description)
	if desc == "" {
		desc = strings.TrimSpace(req.Message)
	}

	typ := strings.ToLower(strings.TrimSpace(req.Type))
	targetRole := strings.ToLower(strings.TrimSpace(req.TargetRole))
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	if typ == "" {
		typ = "info"
	}
	if targetRole == "" {
		targetRole = "client"
	}

	var ctaURL *string
	if v := strings.TrimSpace(req.CTAURL); v != "" {
		ctaURL = &v
	}
	var imageURL *string
	if v := strings.TrimSpace(req.ImageURL); v != "" {
		imageURL = &v
	}

	var expiryAt *time.Time
	if strings.TrimSpace(req.ExpiryDate) != "" {
		t, err := parseFlexibleTime(req.ExpiryDate)
		if err != nil {
			return err
		}
		expiryAt = &t
	}

	var createdBy *string
	if req.CreatedBy != "" {
		createdBy = &req.CreatedBy
	}

	_, err := r.db.Exec(context.Background(), `
		INSERT INTO notifications
			(id, title, message, description, type, target_role, cta_url, image_url, expiry_date, is_active, created_by)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, id, req.Title, desc, desc, typ, targetRole, ctaURL, imageURL, expiryAt, isActive, createdBy)
	return err
}

// Update memperbarui data notification berdasarkan ID.
func (r *Repository) Update(id string, req UpdateNotificationRequest) error {
	desc := strings.TrimSpace(req.Description)
	if desc == "" {
		desc = strings.TrimSpace(req.Message)
	}

	typ := strings.ToLower(strings.TrimSpace(req.Type))
	targetRole := strings.ToLower(strings.TrimSpace(req.TargetRole))
	ctaURL := strings.TrimSpace(req.CTAURL)
	imageURL := strings.TrimSpace(req.ImageURL)

	var expiryAt *time.Time
	if strings.TrimSpace(req.ExpiryDate) != "" {
		t, err := parseFlexibleTime(req.ExpiryDate)
		if err != nil {
			return err
		}
		expiryAt = &t
	}

	var updatedBy *string
	if req.UpdatedBy != "" {
		updatedBy = &req.UpdatedBy
	}

	res, err := r.db.Exec(context.Background(), `
		UPDATE notifications SET
			title       = COALESCE(NULLIF($2,''), title),
			description = COALESCE(NULLIF($3,''), description),
			message     = COALESCE(NULLIF($3,''), message),
			type        = COALESCE(NULLIF($4,''), type),
			target_role = COALESCE(NULLIF($5,''), target_role),
			cta_url     = COALESCE(NULLIF($6,''), cta_url),
			image_url   = COALESCE(NULLIF($7,''), image_url),
			expiry_date = COALESCE($8, expiry_date),
			is_active   = COALESCE($9, is_active),
			updated_by  = $10,
			updated_at  = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id, req.Title, desc, typ, targetRole, ctaURL, imageURL, expiryAt, req.IsActive, updatedBy)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

// Delete menghapus notification berdasarkan ID.
func (r *Repository) Delete(id string) error {
	res, err := r.db.Exec(context.Background(), "DELETE FROM notifications WHERE id = $1", id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func deriveStatus(isActive bool, expiryDate *time.Time) string {
	if expiryDate != nil && expiryDate.Before(time.Now()) {
		return "expired"
	}
	if isActive {
		return "active"
	}
	return "inactive"
}

func parseFlexibleTime(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	layouts := []string{time.RFC3339, "2006-01-02T15:04", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, trimmed); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("format expiry_date tidak valid")
}
