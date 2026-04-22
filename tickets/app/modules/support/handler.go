package support

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tickets/app/models"
	"tickets/core/auth"
	"tickets/core/config"
	corehttp "tickets/core/http"
	"tickets/core/utils"
)

const (
	maxSupportAttachmentBytes = 5 * 1024 * 1024
	maxCreateTicketBodyBytes  = maxSupportAttachmentBytes + (512 * 1024)
)

func parseCreateSupportTicketInput(r *http.Request) (models.CreateSupportTicketInput, error) {
	input := models.CreateSupportTicketInput{}
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
		if !utils.IsAllowedImageMIME(detectedMIME) {
			return input, errors.New("format gambar tidak didukung")
		}

		input.AttachmentName = utils.SanitizeUploadFileName(fileHeader.Filename)
		input.AttachmentMIME = detectedMIME
		input.AttachmentData = attachmentData
		return input, nil
	}

	var payload models.CreateSupportTicketRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxCreateTicketBodyBytes))
	if err := decoder.Decode(&payload); err != nil {
		return input, errors.New("invalid request body")
	}

	input.Title = strings.TrimSpace(payload.Title)
	input.Category = strings.TrimSpace(payload.Category)
	input.Description = strings.TrimSpace(payload.Description)

	return input, nil
}

func HandleAdminSupportTickets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.HasSupportMonitoringRole(r) {
			corehttp.WriteJSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		if db == nil {
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "database unavailable")
			return
		}

		if ticketID, action, ok := utils.ParseAdminSupportTicketPath(r.URL.Path); ok {
			switch {
			case action == "" && r.Method == http.MethodGet:
				handleGetAdminSupportTicketDetail(w, r, db, ticketID)
			case action == "replies" && r.Method == http.MethodPost:
				handleCreateAdminSupportTicketReply(w, r, db, ticketID)
			default:
				corehttp.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}

		if r.URL.Path != "/api/admin/support/tickets" && r.URL.Path != "/api/admin/support/tickets/" {
			corehttp.WriteJSONError(w, http.StatusNotFound, "not found")
			return
		}

		if r.Method != http.MethodGet {
			corehttp.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		search := strings.TrimSpace(r.URL.Query().Get("search"))
		status := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("status")))
		sortParam := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("sort")))
		ticketID := strings.TrimSpace(r.URL.Query().Get("ticket_id"))

		page := utils.ParsePositiveInt(r.URL.Query().Get("page"), 1)
		perPage := utils.ParsePositiveInt(r.URL.Query().Get("per_page"), 6)
		if perPage > 100 {
			perPage = 100
		}

		if status != "" {
			normalizedStatus, validStatus := utils.NormalizeSupportTicketStatus(status)
			if !validStatus {
				corehttp.WriteJSONError(w, http.StatusBadRequest, "invalid status filter")
				return
			}
			status = normalizedStatus
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
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to fetch support tickets")
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
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to fetch support tickets")
			return
		}
		defer rows.Close()

		result := make([]models.SupportTicket, 0)
		for rows.Next() {
			var ticket models.SupportTicket
			if scanErr := rows.Scan(&ticket.ID, &ticket.UserName, &ticket.UserEmail, &ticket.CreatedAt, &ticket.Title, &ticket.Status, &ticket.Category, &ticket.Description, &ticket.HasAttachment); scanErr != nil {
				log.Printf("[ERROR] Failed scanning support ticket: %v", scanErr)
				corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to parse support tickets")
				return
			}
			ticket.CreatedAt = config.NormalizeNaiveTimestampToAppTZ(ticket.CreatedAt)
			ticket.AttachmentURL = utils.BuildAdminAttachmentURL(ticket.ID, ticket.HasAttachment)
			result = append(result, ticket)
		}

		if err = rows.Err(); err != nil {
			log.Printf("[ERROR] Rows error on support tickets: %v", err)
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to read support tickets")
			return
		}

		totalPages := 0
		if total > 0 {
			totalPages = (total + perPage - 1) / perPage
		}

		resp := models.SupportTicketsResponse{
			Data: result,
			Pagination: models.SupportTicketPagination{
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

func handleGetAdminSupportTicketDetail(w http.ResponseWriter, r *http.Request, db *sql.DB, ticketID string) {
	ticketID = strings.TrimSpace(ticketID)
	if ticketID == "" {
		corehttp.WriteJSONError(w, http.StatusBadRequest, "ticket id tidak valid")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var ticket models.SupportTicket
	err := db.QueryRowContext(ctx, `
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
		WHERE st.id = $1
		LIMIT 1
	`, ticketID).Scan(
		&ticket.ID,
		&ticket.UserName,
		&ticket.UserEmail,
		&ticket.CreatedAt,
		&ticket.Title,
		&ticket.Status,
		&ticket.Category,
		&ticket.Description,
		&ticket.HasAttachment,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			corehttp.WriteJSONError(w, http.StatusNotFound, "ticket tidak ditemukan")
			return
		}
		log.Printf("[ERROR] Failed querying support ticket detail: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to fetch support ticket detail")
		return
	}

	ticket.CreatedAt = config.NormalizeNaiveTimestampToAppTZ(ticket.CreatedAt)
	ticket.AttachmentURL = utils.BuildAdminAttachmentURL(ticket.ID, ticket.HasAttachment)

	replyRows, err := db.QueryContext(ctx, `
		SELECT id, ticket_id, admin_id, message, created_at
		FROM support_ticket_replies
		WHERE ticket_id = $1
		ORDER BY created_at ASC, id ASC
	`, ticketID)
	if err != nil {
		log.Printf("[ERROR] Failed querying support ticket replies: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to fetch support ticket replies")
		return
	}
	defer replyRows.Close()

	replies := make([]models.SupportTicketReply, 0)
	for replyRows.Next() {
		var reply models.SupportTicketReply
		if scanErr := replyRows.Scan(&reply.ID, &reply.TicketID, &reply.AdminID, &reply.Message, &reply.CreatedAt); scanErr != nil {
			log.Printf("[ERROR] Failed scanning support ticket reply: %v", scanErr)
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to parse support ticket replies")
			return
		}
		reply.CreatedAt = config.NormalizeNaiveTimestampToAppTZ(reply.CreatedAt)
		replies = append(replies, reply)
	}

	if err = replyRows.Err(); err != nil {
		log.Printf("[ERROR] Rows error on support ticket replies: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to read support ticket replies")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SupportTicketDetailResponse{
		Ticket:  ticket,
		Replies: replies,
	})
}

func handleCreateAdminSupportTicketReply(w http.ResponseWriter, r *http.Request, db *sql.DB, ticketID string) {
	ticketID = strings.TrimSpace(ticketID)
	if ticketID == "" {
		corehttp.WriteJSONError(w, http.StatusBadRequest, "ticket id tidak valid")
		return
	}

	admin := auth.GetAuthenticatedUserFromRequest(r)
	if strings.TrimSpace(admin.ID) == "" {
		corehttp.WriteJSONError(w, http.StatusBadRequest, "admin tidak valid")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 256*1024)

	var payload models.CreateSupportTicketReplyRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		corehttp.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	message := strings.TrimSpace(payload.Message)
	if message == "" {
		corehttp.WriteJSONError(w, http.StatusBadRequest, "message wajib diisi")
		return
	}
	if len(message) > 5000 {
		corehttp.WriteJSONError(w, http.StatusBadRequest, "message maksimal 5000 karakter")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[ERROR] Failed starting transaction for support reply: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to save support reply")
		return
	}
	defer tx.Rollback()

	var currentStatus string
	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM support_tickets
		WHERE id = $1
		FOR UPDATE
	`, ticketID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			corehttp.WriteJSONError(w, http.StatusNotFound, "ticket tidak ditemukan")
			return
		}
		log.Printf("[ERROR] Failed locking support ticket for reply: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to save support reply")
		return
	}

	normalizedCurrentStatus, validCurrentStatus := utils.NormalizeSupportTicketStatus(currentStatus)
	if !validCurrentStatus {
		normalizedCurrentStatus = "ON PROCESS"
	}

	if normalizedCurrentStatus == "DONE" {
		corehttp.WriteJSONError(w, http.StatusForbidden, "ticket sudah selesai dan tidak dapat dijawab kembali")
		return
	}

	nextStatus := normalizedCurrentStatus
	if requestedStatus := strings.TrimSpace(payload.Status); requestedStatus != "" {
		normalizedRequestedStatus, validRequestedStatus := utils.NormalizeSupportTicketStatus(requestedStatus)
		if !validRequestedStatus {
			corehttp.WriteJSONError(w, http.StatusBadRequest, "status tidak valid")
			return
		}
		nextStatus = normalizedRequestedStatus
	}

		replyID := utils.NewUUID()

	createdAt := config.NormalizeNaiveTimestampToAppTZ(time.Now().In(config.SupportAppLocation()))

	_, err = tx.ExecContext(ctx, `
		INSERT INTO support_ticket_replies (id, ticket_id, admin_id, message, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, replyID, ticketID, admin.ID, message, createdAt)
	if err != nil {
		log.Printf("[ERROR] Failed inserting support ticket reply: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to save support reply")
		return
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE support_tickets
		SET status = $1
		WHERE id = $2
	`, nextStatus, ticketID)
	if err != nil {
		log.Printf("[ERROR] Failed updating support ticket status: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to update support ticket status")
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[ERROR] Failed committing support ticket reply: %v", err)
		corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to save support reply")
		return
	}

	reply := models.SupportTicketReply{
		ID:        replyID,
		TicketID:  ticketID,
		AdminID:   admin.ID,
		Message:   message,
		CreatedAt: createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"reply":         reply,
		"ticket_status": nextStatus,
	})
}

func HandleAdminSupportTicketAttachment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			corehttp.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if !auth.HasSupportMonitoringRole(r) {
			corehttp.WriteJSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		if db == nil {
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "database unavailable")
			return
		}

		ticketID := strings.TrimSpace(r.URL.Query().Get("ticket_id"))
		if ticketID == "" {
			corehttp.WriteJSONError(w, http.StatusBadRequest, "ticket_id wajib diisi")
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
				corehttp.WriteJSONError(w, http.StatusNotFound, "ticket tidak ditemukan")
				return
			}
			log.Printf("[ERROR] Failed fetching support ticket attachment: %v", err)
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "gagal memuat lampiran")
			return
		}

		if len(attachmentData) == 0 {
			corehttp.WriteJSONError(w, http.StatusNotFound, "lampiran tidak ditemukan")
			return
		}

		contentType := strings.ToLower(strings.TrimSpace(attachmentMIME.String))
		if !utils.IsAllowedImageMIME(contentType) {
			contentType = strings.ToLower(strings.TrimSpace(http.DetectContentType(attachmentData)))
		}
		if !utils.IsAllowedImageMIME(contentType) {
			contentType = "application/octet-stream"
		}

		filename := utils.SanitizeUploadFileName(attachmentName.String)
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

func HandleCreateUserSupportTicket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			corehttp.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if !auth.HasClientSupportCreateRole(r) {
			corehttp.WriteJSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		if db == nil {
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "database unavailable")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxCreateTicketBodyBytes)

		input, err := parseCreateSupportTicketInput(r)
		if err != nil {
			corehttp.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		title := strings.TrimSpace(input.Title)
		description := strings.TrimSpace(input.Description)
		category, validCategory := utils.NormalizeTicketCategory(input.Category)

		if len(title) < 5 {
			corehttp.WriteJSONError(w, http.StatusBadRequest, "title minimal 5 karakter")
			return
		}
		if !validCategory {
			corehttp.WriteJSONError(w, http.StatusBadRequest, "kategori tidak valid")
			return
		}
		if description == "" {
			corehttp.WriteJSONError(w, http.StatusBadRequest, "description wajib diisi")
			return
		}

		user := auth.GetAuthenticatedUserFromRequest(r)
		if strings.TrimSpace(user.ID) == "" {
			corehttp.WriteJSONError(w, http.StatusBadRequest, "user tidak valid")
			return
		}

		ticketID := utils.NewUUID()

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		createdAt := config.NormalizeNaiveTimestampToAppTZ(time.Now().In(config.SupportAppLocation()))
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
			user.ID, auth.BuildDisplayName(user.Name, user.Email), userEmailValue,
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
		`, ticketID, user.ID, auth.BuildDisplayName(user.Name, user.Email), reporterEmail, title, category, description, createdAt, attachmentName, attachmentMIME, input.AttachmentData)
		if err != nil {
			log.Printf("[ERROR] Failed creating support ticket: %v", err)
			corehttp.WriteJSONError(w, http.StatusInternalServerError, "failed to create ticket")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":             ticketID,
			"user_id":        user.ID,
			"user_name":      auth.BuildDisplayName(user.Name, user.Email),
			"user_email":     reporterEmail,
			"title":          title,
			"category":       category,
			"description":    description,
			"status":         "ON PROCESS",
			"created_at":     config.NormalizeNaiveTimestampToAppTZ(createdAt),
			"has_attachment": len(input.AttachmentData) > 0,
			"attachment_url": utils.BuildAdminAttachmentURL(ticketID, len(input.AttachmentData) > 0),
		})
	}
}
