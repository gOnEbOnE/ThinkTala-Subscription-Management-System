package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type supportTicket struct {
	ID            string    `json:"id"`
	UserName      string    `json:"user_name"`
	UserEmail     string    `json:"user_email,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	Category      string    `json:"category,omitempty"`
	Description   string    `json:"description,omitempty"`
	HasAttachment bool      `json:"has_attachment"`
	AttachmentURL string    `json:"attachment_url,omitempty"`
}

type createSupportTicketRequest struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type authenticatedUser struct {
	ID    string
	Name  string
	Email string
	Role  string
}

type supportTicketPagination struct {
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type supportTicketsResponse struct {
	Data       []supportTicket         `json:"data"`
	Pagination supportTicketPagination `json:"pagination"`
}

var (
	publicKeyOnce sync.Once
	jwtPublicKey  any
	jwtKeyErr     error
	appTZOnce     sync.Once
	appTZLoc      *time.Location
)

const (
	maxSupportAttachmentBytes = 5 * 1024 * 1024
	maxCreateTicketBodyBytes  = maxSupportAttachmentBytes + (512 * 1024)
)

func envOrDefault(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func supportAppLocation() *time.Location {
	appTZOnce.Do(func() {
		tzName := envOrDefault("APP_TIMEZONE", envOrDefault("read_db_timezone", "Asia/Jakarta"))
		loc, err := time.LoadLocation(tzName)
		if err != nil {
			log.Printf("[WARNING] Invalid timezone %q, fallback to Asia/Jakarta: %v", tzName, err)
			loc = time.FixedZone("Asia/Jakarta", 7*60*60)
		}
		appTZLoc = loc
	})

	if appTZLoc == nil {
		return time.FixedZone("Asia/Jakarta", 7*60*60)
	}

	return appTZLoc
}

func normalizeNaiveTimestampToAppTZ(ts time.Time) time.Time {
	if ts.IsZero() {
		return ts
	}

	loc := supportAppLocation()
	year, month, day := ts.Date()
	hour, minute, second := ts.Clock()

	// Re-attach app timezone to DB timestamp fields without shifting wall-clock time.
	return time.Date(year, month, day, hour, minute, second, ts.Nanosecond(), loc)
}

func openDB() (*sql.DB, error) {
	host := envOrDefault("read_db_host", envOrDefault("DB_HOST", "localhost"))
	port := envOrDefault("read_db_port", envOrDefault("DB_PORT", "5432"))
	user := envOrDefault("read_db_user", envOrDefault("DB_USER", "postgres"))
	pass := envOrDefault("read_db_pass", envOrDefault("DB_PASS", "postgres"))
	name := envOrDefault("read_db_name", envOrDefault("DB_NAME", "postgres"))
	sslMode := envOrDefault("read_db_ssl_mode", "disable")
	timeZone := envOrDefault("read_db_timezone", "Asia/Jakarta")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, pass, name, sslMode, timeZone)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func migrateSupportTicketsTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS support_tickets (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			reporter_name VARCHAR(255) NOT NULL DEFAULT '',
			reporter_email VARCHAR(255) NOT NULL DEFAULT '',
			title TEXT NOT NULL,
			category VARCHAR(50) NOT NULL DEFAULT 'MASALAH_TEKNIS',
			description TEXT NOT NULL DEFAULT '',
			status VARCHAR(20) NOT NULL DEFAULT 'ON PROCESS',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		ALTER TABLE support_tickets DROP CONSTRAINT IF EXISTS support_tickets_user_id_fkey;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255);
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS reporter_name VARCHAR(255) NOT NULL DEFAULT '';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS reporter_email VARCHAR(255) NOT NULL DEFAULT '';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS category VARCHAR(50) NOT NULL DEFAULT 'MASALAH_TEKNIS';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS attachment_name VARCHAR(255);
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS attachment_mime VARCHAR(100);
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS attachment_data BYTEA;

		CREATE INDEX IF NOT EXISTS idx_support_tickets_created_at ON support_tickets(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_support_tickets_user_id ON support_tickets(user_id);
		CREATE INDEX IF NOT EXISTS idx_support_tickets_status ON support_tickets(status);
	`)

	return err
}

func parseClaimsFromAuthorization(authHeader string) (jwt.MapClaims, error) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(authHeader)), "bearer ") {
		return nil, errors.New("missing bearer token")
	}

	tokenString := strings.TrimSpace(authHeader[7:])
	if tokenString == "" {
		return nil, errors.New("empty bearer token")
	}

	pubKey, err := loadJWTPublicKey()
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid bearer token")
	}

	return claims, nil
}

func extractRoleFromClaims(claims jwt.MapClaims) string {
	if roleCode, ok := claims["role_code"].(string); ok {
		return strings.ToUpper(strings.TrimSpace(roleCode))
	}
	if role, ok := claims["role"].(string); ok {
		return strings.ToUpper(strings.TrimSpace(role))
	}

	if userMapAny, ok := claims["user"].(map[string]any); ok {
		if roleCode, ok := userMapAny["role_code"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(roleCode))
		}
		if role, ok := userMapAny["role"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(role))
		}
	}

	if userMapIface, ok := claims["user"].(map[string]interface{}); ok {
		if roleCode, ok := userMapIface["role_code"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(roleCode))
		}
		if role, ok := userMapIface["role"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(role))
		}
	}

	return ""
}

func loadJWTPublicKey() (any, error) {
	publicKeyOnce.Do(func() {
		if b64 := strings.TrimSpace(os.Getenv("JWT_PUBLIC_KEY_B64")); b64 != "" {
			decoded, err := base64.StdEncoding.DecodeString(b64)
			if err == nil {
				if key, keyErr := jwt.ParseRSAPublicKeyFromPEM(decoded); keyErr == nil {
					jwtPublicKey = key
					jwtKeyErr = nil
					return
				}
			}
		}

		paths := []string{"certs/public.pem", "../users/certs/public.pem"}
		for _, p := range paths {
			pemBytes, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
			if err == nil {
				jwtPublicKey = key
				jwtKeyErr = nil
				return
			}
		}

		jwtKeyErr = errors.New("jwt public key not found")
	})

	if jwtPublicKey == nil {
		return nil, jwtKeyErr
	}

	return jwtPublicKey, nil
}

func extractRoleFromAuthorization(authHeader string) string {
	claims, err := parseClaimsFromAuthorization(authHeader)
	if err != nil {
		return ""
	}

	return extractRoleFromClaims(claims)
}

func hasSupportMonitoringRole(r *http.Request) bool {
	role := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	if role == "" {
		role = extractRoleFromAuthorization(r.Header.Get("Authorization"))
	}
	return role == "ADMIN_SUPPORT" || role == "OPERASIONAL"
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func buildDisplayName(rawName, rawEmail string) string {
	name := strings.TrimSpace(rawName)
	if name != "" {
		return name
	}

	email := strings.TrimSpace(rawEmail)
	if email == "" {
		return "User"
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
		return strings.TrimSpace(parts[0])
	}
	return email
}

func extractUserFromClaims(claims jwt.MapClaims) authenticatedUser {
	user := authenticatedUser{}

	if role := extractRoleFromClaims(claims); role != "" {
		user.Role = role
	}

	if id, ok := claims["user_id"].(string); ok {
		user.ID = strings.TrimSpace(id)
	}
	if email, ok := claims["email"].(string); ok {
		user.Email = strings.TrimSpace(email)
	}
	if name, ok := claims["name"].(string); ok {
		user.Name = strings.TrimSpace(name)
	}

	if userMapAny, ok := claims["user"].(map[string]any); ok {
		if id, ok := userMapAny["id"].(string); ok && user.ID == "" {
			user.ID = strings.TrimSpace(id)
		}
		if email, ok := userMapAny["email"].(string); ok && user.Email == "" {
			user.Email = strings.TrimSpace(email)
		}
		if name, ok := userMapAny["name"].(string); ok && user.Name == "" {
			user.Name = strings.TrimSpace(name)
		}
	}

	if userMapIface, ok := claims["user"].(map[string]interface{}); ok {
		if id, ok := userMapIface["id"].(string); ok && user.ID == "" {
			user.ID = strings.TrimSpace(id)
		}
		if email, ok := userMapIface["email"].(string); ok && user.Email == "" {
			user.Email = strings.TrimSpace(email)
		}
		if name, ok := userMapIface["name"].(string); ok && user.Name == "" {
			user.Name = strings.TrimSpace(name)
		}
	}

	if user.Name == "" {
		user.Name = buildDisplayName("", user.Email)
	}

	return user
}

func getAuthenticatedUserFromRequest(r *http.Request) authenticatedUser {
	user := authenticatedUser{
		ID:    strings.TrimSpace(r.Header.Get("X-User-ID")),
		Email: strings.TrimSpace(r.Header.Get("X-User-Email")),
		Role:  strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role"))),
	}
	user.Name = buildDisplayName("", user.Email)

	claims, err := parseClaimsFromAuthorization(r.Header.Get("Authorization"))
	if err != nil {
		return user
	}

	fromClaims := extractUserFromClaims(claims)
	if user.ID == "" {
		user.ID = fromClaims.ID
	}
	if user.Email == "" {
		user.Email = fromClaims.Email
	}
	if user.Role == "" {
		user.Role = fromClaims.Role
	}
	if fromClaims.Name != "" {
		user.Name = fromClaims.Name
	} else {
		user.Name = buildDisplayName("", user.Email)
	}

	return user
}

func hasClientSupportCreateRole(r *http.Request) bool {
	user := getAuthenticatedUserFromRequest(r)
	return user.Role == "CLIENT" || user.Role == "SUPERADMIN" || user.Role == "CEO"
}

func normalizeTicketCategory(raw string) (string, bool) {
	cleaned := strings.ToUpper(strings.TrimSpace(raw))
	cleaned = strings.ReplaceAll(cleaned, "-", "_")
	cleaned = strings.ReplaceAll(cleaned, " ", "_")

	switch cleaned {
	case "MASALAH_TEKNIS", "TEKNIS", "TEKNIK", "TECHNICAL", "MASALAHTECHNIS":
		return "MASALAH_TEKNIS", true
	case "PEMBAYARAN", "PAYMENT":
		return "PEMBAYARAN", true
	case "AKUN", "ACCOUNT":
		return "AKUN", true
	default:
		return "", false
	}
}

func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	hexID := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexID[0:8],
		hexID[8:12],
		hexID[12:16],
		hexID[16:20],
		hexID[20:32]), nil
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func isAllowedImageMIME(contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/jpg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

func sanitizeUploadFileName(fileName string) string {
	baseName := filepath.Base(strings.TrimSpace(fileName))
	if baseName == "" || baseName == "." || baseName == string(filepath.Separator) {
		return "attachment"
	}
	if len(baseName) > 160 {
		baseName = baseName[:160]
	}
	return strings.ReplaceAll(baseName, "\x00", "")
}

func buildAdminAttachmentURL(ticketID string, hasAttachment bool) string {
	if !hasAttachment || strings.TrimSpace(ticketID) == "" {
		return ""
	}
	return fmt.Sprintf("/api/admin/support/tickets/attachment?ticket_id=%s", ticketID)
}

type createSupportTicketInput struct {
	Title          string
	Category       string
	Description    string
	AttachmentName string
	AttachmentMIME string
	AttachmentData []byte
}

func parseCreateSupportTicketInput(r *http.Request) (createSupportTicketInput, error) {
	input := createSupportTicketInput{}
	contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(maxCreateTicketBodyBytes); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "request body too large") {
				return input, errors.New("ukuran gambar maksimal 5 MB")
			}
			return input, errors.New("invalid request body")
		}

		input.Title = strings.TrimSpace(r.FormValue("title"))
		input.Category = strings.TrimSpace(r.FormValue("category"))
		input.Description = strings.TrimSpace(r.FormValue("description"))

		file, fileHeader, err := r.FormFile("attachment")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				return input, nil
			}
			return input, errors.New("gagal membaca file gambar")
		}
		defer file.Close()

		if fileHeader.Size > maxSupportAttachmentBytes && fileHeader.Size > 0 {
			return input, errors.New("ukuran gambar maksimal 5 MB")
		}

		attachmentData, err := io.ReadAll(io.LimitReader(file, maxSupportAttachmentBytes+1))
		if err != nil {
			return input, errors.New("gagal membaca file gambar")
		}

		if len(attachmentData) == 0 {
			return input, errors.New("file gambar kosong")
		}
		if len(attachmentData) > maxSupportAttachmentBytes {
			return input, errors.New("ukuran gambar maksimal 5 MB")
		}

		detectedMIME := strings.ToLower(strings.TrimSpace(http.DetectContentType(attachmentData)))
		if !isAllowedImageMIME(detectedMIME) {
			return input, errors.New("format gambar tidak didukung")
		}

		input.AttachmentName = sanitizeUploadFileName(fileHeader.Filename)
		input.AttachmentMIME = detectedMIME
		input.AttachmentData = attachmentData
		return input, nil
	}

	var payload createSupportTicketRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxCreateTicketBodyBytes))
	if err := decoder.Decode(&payload); err != nil {
		return input, errors.New("invalid request body")
	}

	input.Title = strings.TrimSpace(payload.Title)
	input.Category = strings.TrimSpace(payload.Category)
	input.Description = strings.TrimSpace(payload.Description)

	return input, nil
}

func handleAdminSupportTickets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if !hasSupportMonitoringRole(r) {
			writeJSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database unavailable")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		search := strings.TrimSpace(r.URL.Query().Get("search"))
		status := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("status")))
		sortParam := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("sort")))
		ticketID := strings.TrimSpace(r.URL.Query().Get("ticket_id"))

		page := parsePositiveInt(r.URL.Query().Get("page"), 1)
		perPage := parsePositiveInt(r.URL.Query().Get("per_page"), 6)
		if perPage > 100 {
			perPage = 100
		}

		if status != "" && status != "ON PROCESS" && status != "DONE" {
			writeJSONError(w, http.StatusBadRequest, "invalid status filter")
			return
		}

		sortOrder := "DESC"
		if sortParam == "oldest" || sortParam == "asc" {
			sortOrder = "ASC"
		}

		whereClauses := make([]string, 0)
		args := make([]any, 0)

		if ticketID != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("st.id = $%d", len(args)+1))
			args = append(args, ticketID)
		}

		if status != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("st.status = $%d", len(args)+1))
			args = append(args, status)
		}

		if search != "" {
			whereClauses = append(whereClauses,
				fmt.Sprintf("(st.title ILIKE $%d OR COALESCE(NULLIF(st.reporter_name, ''), COALESCE(u.name, '')) ILIKE $%d OR COALESCE(NULLIF(st.reporter_email, ''), COALESCE(u.email, '')) ILIKE $%d)", len(args)+1, len(args)+1, len(args)+1))
			args = append(args, "%"+search+"%")
		}

		whereSQL := ""
		if len(whereClauses) > 0 {
			whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
		}

		countQuery := `
			SELECT COUNT(1)
			FROM support_tickets st
			LEFT JOIN users u ON u.id = st.user_id
		` + whereSQL

		var total int
		if err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
			log.Printf("[ERROR] Failed counting support tickets: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch support tickets")
			return
		}

		offset := (page - 1) * perPage

		query := fmt.Sprintf(`
			SELECT st.id,
			       COALESCE(NULLIF(st.reporter_name, ''), COALESCE(u.name, '')) AS user_name,
			       COALESCE(NULLIF(st.reporter_email, ''), COALESCE(u.email, '')) AS user_email,
			       st.created_at,
			       st.title,
			       st.status,
			       st.category,
			       st.description,
			       CASE WHEN st.attachment_data IS NULL OR OCTET_LENGTH(st.attachment_data) = 0 THEN FALSE ELSE TRUE END AS has_attachment
			FROM support_tickets st
			LEFT JOIN users u ON u.id = st.user_id
			%s
			ORDER BY st.created_at %s
			LIMIT $%d OFFSET $%d
		`, whereSQL, sortOrder, len(args)+1, len(args)+2)

		queryArgs := append(append([]any{}, args...), perPage, offset)

		rows, err := db.QueryContext(ctx, query, queryArgs...)
		if err != nil {
			log.Printf("[ERROR] Failed querying support tickets: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch support tickets")
			return
		}
		defer rows.Close()

		result := make([]supportTicket, 0)
		for rows.Next() {
			var ticket supportTicket
			if scanErr := rows.Scan(&ticket.ID, &ticket.UserName, &ticket.UserEmail, &ticket.CreatedAt, &ticket.Title, &ticket.Status, &ticket.Category, &ticket.Description, &ticket.HasAttachment); scanErr != nil {
				log.Printf("[ERROR] Failed scanning support ticket: %v", scanErr)
				writeJSONError(w, http.StatusInternalServerError, "failed to parse support tickets")
				return
			}
			ticket.CreatedAt = normalizeNaiveTimestampToAppTZ(ticket.CreatedAt)
			ticket.AttachmentURL = buildAdminAttachmentURL(ticket.ID, ticket.HasAttachment)
			result = append(result, ticket)
		}

		if err = rows.Err(); err != nil {
			log.Printf("[ERROR] Rows error on support tickets: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to read support tickets")
			return
		}

		totalPages := 0
		if total > 0 {
			totalPages = (total + perPage - 1) / perPage
		}

		resp := supportTicketsResponse{
			Data: result,
			Pagination: supportTicketPagination{
				Page:       page,
				PerPage:    perPage,
				Total:      total,
				TotalPages: totalPages,
				HasNext:    page < totalPages,
				HasPrev:    page > 1 && totalPages > 0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func handleAdminSupportTicketAttachment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if !hasSupportMonitoringRole(r) {
			writeJSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database unavailable")
			return
		}

		ticketID := strings.TrimSpace(r.URL.Query().Get("ticket_id"))
		if ticketID == "" {
			writeJSONError(w, http.StatusBadRequest, "ticket_id wajib diisi")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var attachmentName sql.NullString
		var attachmentMIME sql.NullString
		var attachmentData []byte

		err := db.QueryRowContext(ctx, `
			SELECT attachment_name, attachment_mime, attachment_data
			FROM support_tickets
			WHERE id = $1
		`, ticketID).Scan(&attachmentName, &attachmentMIME, &attachmentData)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, http.StatusNotFound, "ticket tidak ditemukan")
				return
			}
			log.Printf("[ERROR] Failed fetching support ticket attachment: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "gagal memuat lampiran")
			return
		}

		if len(attachmentData) == 0 {
			writeJSONError(w, http.StatusNotFound, "lampiran tidak ditemukan")
			return
		}

		contentType := strings.ToLower(strings.TrimSpace(attachmentMIME.String))
		if !isAllowedImageMIME(contentType) {
			contentType = strings.ToLower(strings.TrimSpace(http.DetectContentType(attachmentData)))
		}
		if !isAllowedImageMIME(contentType) {
			contentType = "application/octet-stream"
		}

		filename := sanitizeUploadFileName(attachmentName.String)
		if filename == "" {
			filename = "attachment"
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(len(attachmentData)))
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(attachmentData)
	}
}

func handleCreateUserSupportTicket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if !hasClientSupportCreateRole(r) {
			writeJSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database unavailable")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxCreateTicketBodyBytes)

		input, err := parseCreateSupportTicketInput(r)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		title := strings.TrimSpace(input.Title)
		description := strings.TrimSpace(input.Description)
		category, validCategory := normalizeTicketCategory(input.Category)

		if len(title) < 5 {
			writeJSONError(w, http.StatusBadRequest, "title minimal 5 karakter")
			return
		}
		if !validCategory {
			writeJSONError(w, http.StatusBadRequest, "kategori tidak valid")
			return
		}
		if description == "" {
			writeJSONError(w, http.StatusBadRequest, "description wajib diisi")
			return
		}

		user := getAuthenticatedUserFromRequest(r)
		if strings.TrimSpace(user.ID) == "" {
			writeJSONError(w, http.StatusBadRequest, "user tidak valid")
			return
		}

		ticketID, err := newUUID()
		if err != nil {
			log.Printf("[ERROR] Failed generating ticket id: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to create ticket")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		createdAt := normalizeNaiveTimestampToAppTZ(time.Now().In(supportAppLocation()))
		reporterEmail := strings.TrimSpace(user.Email)

		userEmailValue := sql.NullString{}
		if reporterEmail != "" {
			userEmailValue = sql.NullString{String: reporterEmail, Valid: true}
		}

		_, _ = db.ExecContext(ctx,
			`INSERT INTO users (id, name, email)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (id) DO UPDATE SET
			 	name = EXCLUDED.name,
			 	email = COALESCE(EXCLUDED.email, users.email)`,
			user.ID, buildDisplayName(user.Name, user.Email), userEmailValue,
		)

		attachmentName := sql.NullString{}
		if strings.TrimSpace(input.AttachmentName) != "" {
			attachmentName = sql.NullString{String: strings.TrimSpace(input.AttachmentName), Valid: true}
		}

		attachmentMIME := sql.NullString{}
		if strings.TrimSpace(input.AttachmentMIME) != "" {
			attachmentMIME = sql.NullString{String: strings.TrimSpace(input.AttachmentMIME), Valid: true}
		}

		_, err = db.ExecContext(ctx, `
			INSERT INTO support_tickets (
				id,
				user_id,
				reporter_name,
				reporter_email,
				title,
				category,
				description,
				status,
				created_at,
				attachment_name,
				attachment_mime,
				attachment_data
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'ON PROCESS', $8, $9, $10, $11)
		`, ticketID, user.ID, buildDisplayName(user.Name, user.Email), reporterEmail, title, category, description, createdAt, attachmentName, attachmentMIME, input.AttachmentData)
		if err != nil {
			log.Printf("[ERROR] Failed creating support ticket: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to create ticket")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":             ticketID,
			"user_id":        user.ID,
			"user_name":      buildDisplayName(user.Name, user.Email),
			"user_email":     reporterEmail,
			"title":          title,
			"category":       category,
			"description":    description,
			"status":         "ON PROCESS",
			"created_at":     normalizeNaiveTimestampToAppTZ(createdAt),
			"has_attachment": len(input.AttachmentData) > 0,
			"attachment_url": buildAdminAttachmentURL(ticketID, len(input.AttachmentData) > 0),
		})
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2004"
	}

	db, dbErr := openDB()
	if dbErr != nil {
		log.Printf("[WARNING] Failed to connect database: %v", dbErr)
	} else {
		defer db.Close()
		if err := migrateSupportTicketsTable(db); err != nil {
			log.Printf("[WARNING] Failed to migrate support_tickets table: %v", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/support/tickets", handleAdminSupportTickets(db))
	mux.HandleFunc("/api/admin/support/tickets/", handleAdminSupportTickets(db))
	mux.HandleFunc("/api/admin/support/tickets/attachment", handleAdminSupportTicketAttachment(db))
	mux.HandleFunc("/api/admin/support/tickets/attachment/", handleAdminSupportTicketAttachment(db))
	mux.HandleFunc("/api/user/support/tickets", handleCreateUserSupportTicket(db))
	mux.HandleFunc("/api/user/support/tickets/", handleCreateUserSupportTicket(db))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service": "Tickets Service",
			"status":  "Healthy",
			"path":    r.URL.Path,
		})
	})

	fmt.Printf("Tickets service is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
