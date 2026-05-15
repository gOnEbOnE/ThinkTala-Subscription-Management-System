package notification

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

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
			is_pinned,
			view_count,
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
			&n.IsPinned,
			&n.ViewCount,
			&n.CreatedAt,
			&n.CreatedBy,
			&n.UpdatedAt,
			&n.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		n.Status = deriveStatus(n.IsActive, n.ExpiryDate)
		applyDetailFlags(&n)
		list = append(list, n)
	}
	if list == nil {
		list = []Notification{}
	}
	return list, nil
}

// ListPublic mengambil notification aktif sesuai audience, max 20 data.
func (r *Repository) ListPublic(role, userID string) ([]map[string]any, error) {
	audienceKeys, err := r.buildAudienceKeys(role, userID)
	if err != nil {
		return nil, err
	}
	if len(audienceKeys) == 0 {
		return []map[string]any{}, nil
	}

	rows, err := r.db.Query(context.Background(), `
		SELECT id, title, COALESCE(description, message, '') AS description, type, target_role,
		       cta_url, image_url, expiry_date, created_at, is_pinned
		FROM notifications
		WHERE is_active = TRUE
		  AND (expiry_date IS NULL OR expiry_date > NOW())
		  AND LOWER(target_role) = ANY($1)
		ORDER BY is_pinned DESC, created_at DESC LIMIT 20
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
		var isPinned bool
		err = rows.Scan(&id, &title, &desc, &typ, &target, &ctaURL, &imageURL, &expiryDate, &createdAt, &isPinned)
		if err != nil {
			return nil, err
		}
		plainDesc := normalizePlainText(desc)
		isLong := isLongContent(plainDesc)
		excerpt := buildExcerpt(plainDesc)
		list = append(list, map[string]any{
			"id":          id,
			"title":       title,
			"description": excerpt,
			"message":     excerpt,
			"excerpt":     excerpt,
			"type":        typ,
			"target_role": target,
			"cta_url":     ctaURL,
			"image_url":   imageURL,
			"is_pinned":   isPinned,
			"has_detail":  isLong,
			"is_long_content": isLong,
			"expiry_date": expiryDate,
			"created_at":  createdAt,
		})
	}
	if list == nil {
		list = []map[string]any{}
	}
	return list, nil
}

// ListRecentNews mengambil ringkasan news terbaru untuk drawer.
func (r *Repository) ListRecentNews(role, userID string, limit int) ([]RecentNotification, error) {
	audienceKeys, err := r.buildAudienceKeys(role, userID)
	if err != nil {
		return nil, err
	}
	if len(audienceKeys) == 0 {
		return []RecentNotification{}, nil
	}
	if limit <= 0 {
		limit = 5
	}

	rows, err := r.db.Query(context.Background(), `
		SELECT n.id,
			   n.title,
			   COALESCE(n.description, n.message, '') AS description,
			   n.type,
			   n.cta_url,
			   n.image_url,
			   n.created_at,
			   EXISTS (
				   SELECT 1 FROM notification_reads nr
				   WHERE nr.user_id = $2
					 AND nr.source_type = 'news'
					 AND nr.source_id = n.id::text
			   ) AS is_read
		FROM notifications n
		WHERE n.is_active = TRUE
		  AND (n.expiry_date IS NULL OR n.expiry_date > NOW())
		  AND LOWER(n.target_role) = ANY($1)
		ORDER BY n.created_at DESC
		LIMIT $3
	`, audienceKeys, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]RecentNotification, 0)
	for rows.Next() {
		var id, title, desc, typ string
		var ctaURL, imageURL *string
		var createdAt time.Time
		var isRead bool
		if err := rows.Scan(&id, &title, &desc, &typ, &ctaURL, &imageURL, &createdAt, &isRead); err != nil {
			return nil, err
		}
		if strings.TrimSpace(title) == "" {
			title = "Update Terbaru"
		}
		excerpt := buildExcerpt(normalizePlainText(desc))
		list = append(list, RecentNotification{
			ID:        id,
			Source:    "news",
			Title:     title,
			Body:      excerpt,
			Type:      typ,
			CTAURL:    ctaURL,
			ImageURL:  imageURL,
			CreatedAt: createdAt,
			IsRead:    isRead,
		})
	}
	return list, nil
}

// ListRecentEvents mengambil ringkasan log event terbaru untuk drawer.
func (r *Repository) ListRecentEvents(userID string, limit int) ([]RecentNotification, error) {
	if strings.TrimSpace(userID) == "" {
		return []RecentNotification{}, nil
	}
	if limit <= 0 {
		limit = 5
	}

	rows, err := r.db.Query(context.Background(), `
		SELECT l.id,
			   l.event_type,
			   l.channel,
			   COALESCE(l.subject, '') AS subject,
			   COALESCE(l.content, '') AS content,
			   l.status,
			   COALESCE(l.sent_at, l.created_at) AS created_at,
			   EXISTS (
				   SELECT 1 FROM notification_reads nr
				   WHERE nr.user_id = $1
					 AND nr.source_type = 'event'
					 AND nr.source_id = l.id
			   ) AS is_read
		FROM notification_logs l
		WHERE l.user_id = $1
		ORDER BY COALESCE(l.sent_at, l.created_at) DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]RecentNotification, 0)
	for rows.Next() {
		var id, eventType, channel, subject, content, status string
		var createdAt time.Time
		var isRead bool
		if err := rows.Scan(&id, &eventType, &channel, &subject, &content, &status, &createdAt, &isRead); err != nil {
			return nil, err
		}
		title := strings.TrimSpace(subject)
		if title == "" {
			title = strings.ReplaceAll(eventType, "_", " ")
		}
		excerpt := buildExcerpt(normalizePlainText(content))
		if strings.TrimSpace(excerpt) == "" {
			excerpt = strings.ToUpper(strings.ReplaceAll(eventType, "_", " "))
		}
		list = append(list, RecentNotification{
			ID:        id,
			Source:    "event",
			Title:     title,
			Body:      excerpt,
			EventType: eventType,
			Channel:   channel,
			Type:      "event",
			Status:    status,
			CreatedAt: createdAt,
			IsRead:    isRead,
		})
	}
	return list, nil
}

// MarkRead menandai satu item (news/event) sebagai read untuk user tertentu.
func (r *Repository) MarkRead(userID, sourceType, sourceID string) error {
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO notification_reads (user_id, source_type, source_id, read_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, source_type, source_id) DO NOTHING
	`, userID, sourceType, sourceID)
	return err
}

// MarkAllRead menandai semua item (news/event) sebagai read untuk user tertentu.
func (r *Repository) MarkAllRead(userID, role string) error {
	if strings.TrimSpace(userID) == "" {
		return fmt.Errorf("user_id wajib diisi")
	}

	// Event notifications (per user)
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO notification_reads (user_id, source_type, source_id, read_at)
		SELECT $1::varchar, 'event', id::varchar, NOW()
		FROM notification_logs
		WHERE user_id = $1::varchar
		ON CONFLICT (user_id, source_type, source_id) DO NOTHING
	`, userID)
	if err != nil {
		return err
	}

	// News notifications (per audience)
	audienceKeys, err := r.buildAudienceKeys(role, userID)
	if err != nil {
		return err
	}
	if len(audienceKeys) == 0 {
		return nil
	}

	_, err = r.db.Exec(context.Background(), `
		INSERT INTO notification_reads (user_id, source_type, source_id, read_at)
		SELECT $1::varchar, 'news', n.id::varchar, NOW()
		FROM notifications n
		WHERE n.is_active = TRUE
		  AND (n.expiry_date IS NULL OR n.expiry_date > NOW())
		  AND LOWER(n.target_role) = ANY($2::varchar[])
		ON CONFLICT (user_id, source_type, source_id) DO NOTHING
	`, userID, audienceKeys)
	return err
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
	// If not pgx.ErrNoRows, log but DO NOT propagate — fall through to next check.
	// This makes the function resilient to cross-DB/schema access errors.
	if err != pgx.ErrNoRows {
		log.Printf("[NOTIF AUDIENCE] subscription query error for user %s (non-fatal): %v", userID, err)
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
		// Cross-schema query unavailable — return nil (caller will use default ["client","all"])
		log.Printf("[NOTIF AUDIENCE] orders query error for user %s (non-fatal): %v", userID, err)
		return nil, nil
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
			log.Printf("[NOTIF AUDIENCE] expiry query error for user %s (non-fatal): %v", userID, err)
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

func (r *Repository) buildAudienceKeys(role, userID string) ([]string, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	log.Printf("[AUDIENCE] input role=%q userID=%q", role, userID)
	if role != "" && role != "client" {
		log.Printf("[AUDIENCE] role=%q is not 'client' → returning empty keys", role)
		return []string{}, nil
	}

	// Include 'all' so legacy notifications with target_role='all' are delivered
	// to client audiences as well.
	audiences := map[string]struct{}{"client": {}, "all": {}}
	if segments, err := r.resolveClientAudienceSegments(strings.TrimSpace(userID)); err == nil {
		for _, segment := range segments {
			segment = strings.ToLower(strings.TrimSpace(segment))
			if segment == "" {
				continue
			}
			audiences[segment] = struct{}{}
		}
	} else if err != nil {
		return nil, err
	}

	audienceKeys := make([]string, 0, len(audiences))
	for key := range audiences {
		audienceKeys = append(audienceKeys, key)
	}
	log.Printf("[AUDIENCE] resolved keys=%v", audienceKeys)
	return audienceKeys, nil
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
			is_pinned,
			view_count,
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
		&n.IsPinned,
		&n.ViewCount,
		&n.CreatedAt,
		&n.CreatedBy,
		&n.UpdatedAt,
		&n.UpdatedBy,
	)
	if err == nil {
		n.Status = deriveStatus(n.IsActive, n.ExpiryDate)
		applyDetailFlags(&n)
	}
	return n, err
}

// GetPublicByID mengambil notification untuk client dan menambah view_count.
func (r *Repository) GetPublicByID(id, role, userID string) (Notification, error) {
	audienceKeys, err := r.buildAudienceKeys(role, userID)
	if err != nil {
		return Notification{}, err
	}
	if len(audienceKeys) == 0 {
		return Notification{}, pgx.ErrNoRows
	}

	var n Notification
	err = r.db.QueryRow(context.Background(), `
		UPDATE notifications
		SET view_count = view_count + 1
		WHERE id = $1
		  AND is_active = TRUE
		  AND (expiry_date IS NULL OR expiry_date > NOW())
		  AND LOWER(target_role) = ANY($2)
		RETURNING
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
			is_pinned,
			view_count,
			created_at,
			created_by,
			updated_at,
			updated_by
	`, id, audienceKeys).Scan(
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
		&n.IsPinned,
		&n.ViewCount,
		&n.CreatedAt,
		&n.CreatedBy,
		&n.UpdatedAt,
		&n.UpdatedBy,
	)
	if err == nil {
		n.Status = deriveStatus(n.IsActive, n.ExpiryDate)
		applyDetailFlags(&n)
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
	isPinned := false
	if req.IsPinned != nil {
		isPinned = *req.IsPinned
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
			(id, title, message, description, type, target_role, cta_url, image_url, expiry_date, is_active, is_pinned, created_by)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, id, req.Title, desc, desc, typ, targetRole, ctaURL, imageURL, expiryAt, isActive, isPinned, createdBy)
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
			is_pinned   = COALESCE($10, is_pinned),
			updated_by  = $11,
			updated_at  = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id, req.Title, desc, typ, targetRole, ctaURL, imageURL, expiryAt, req.IsActive, req.IsPinned, updatedBy)
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

func (r *Repository) countPinned(excludeID string) (int, error) {
	ctx := context.Background()
	var count int

	if strings.TrimSpace(excludeID) == "" {
		err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE is_pinned = TRUE`).Scan(&count)
		return count, err
	}

	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE is_pinned = TRUE AND id <> $1`, excludeID).Scan(&count)
	return count, err
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

const (
	excerptMaxRunes = 180
)

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func normalizePlainText(value string) string {
	clean := htmlTagRe.ReplaceAllString(value, "")
	clean = strings.ReplaceAll(clean, "\u00a0", " ")
	clean = strings.TrimSpace(clean)
	if clean == "" {
		return ""
	}
	return strings.Join(strings.Fields(clean), " ")
}

func buildExcerpt(value string) string {
	plain := normalizePlainText(value)
	if plain == "" {
		return ""
	}
	if utf8.RuneCountInString(plain) <= excerptMaxRunes {
		return plain
	}
	runes := []rune(plain)
	trimmed := strings.TrimSpace(string(runes[:excerptMaxRunes]))
	return strings.TrimRight(trimmed, " .,") + "..."
}

func isLongContent(value string) bool {
	plain := normalizePlainText(value)
	return utf8.RuneCountInString(plain) > excerptMaxRunes
}

func applyDetailFlags(n *Notification) {
	plain := normalizePlainText(n.Description)
	isLong := isLongContent(plain)
	n.IsLong = isLong
	n.HasDetail = isLong
	n.Excerpt = buildExcerpt(plain)
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
