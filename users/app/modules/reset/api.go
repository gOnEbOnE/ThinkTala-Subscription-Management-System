package reset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/master-abror/zaframework/core/database"
	"github.com/master-abror/zaframework/core/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5"
)

// APIHandler handles password reset API endpoints.
type APIHandler struct {
	DB *database.DBWrapper
}

func NewAPIHandler(db *database.DBWrapper) *APIHandler {
	return &APIHandler{DB: db}
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

// dispatchResetNotification sends a password_reset event to the notification service.
// Pattern mirrors dispatchKYCNotification in kyc/service.go exactly.
func dispatchResetNotification(to string, vars map[string]string) {
	// 1. Try Redis queue first
	if err := utils.PublishNotificationEvent("password_reset", "email", to, vars); err == nil {
		log.Printf("[RESET NOTIF] Event published to queue: to=%s", to)
		return
	}

	// 2. HTTP fallback to notification service
	baseURL := utils.GetEnv("NOTIFICATION_SERVICE_URL", "http://localhost:5003")
	payload := map[string]any{
		"event_type": "password_reset",
		"channel":    "email",
		"to":         to,
		"vars":       vars,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/api/notifications/send", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("[RESET NOTIF] Notification service unavailable: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]any
		json.NewDecoder(resp.Body).Decode(&result)
		log.Printf("[RESET NOTIF] Failed via notification service (%d: %v)", resp.StatusCode, result["error"])
	}
}

// ForgotPassword — POST /api/auth/forgot-password
func (h *APIHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(body.Email))
	if email == "" || !utils.IsValidEmail(email) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format email tidak valid"})
		return
	}

	// Always respond the same to prevent email enumeration
	successMsg := map[string]string{"message": "Jika email terdaftar, instruksi telah dikirim"}

	if h.DB == nil || h.DB.Pool == nil {
		writeJSON(w, http.StatusOK, successMsg)
		return
	}

	var userID, fullName string
	err := h.DB.Pool.QueryRow(r.Context(),
		`SELECT id, COALESCE(NULLIF(full_name, ''), NULLIF(name, ''), split_part(email, '@', 1)) FROM users WHERE email = $1 AND status = 'active' LIMIT 1`,
		email,
	).Scan(&userID, &fullName)
	if err != nil {
		// email not found — still return 200
		writeJSON(w, http.StatusOK, successMsg)
		return
	}

	// Rate-limit: if an active token already exists, silently return success
	var activeCount int
	countErr := h.DB.Pool.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM password_reset_tokens WHERE user_id = $1 AND expires_at > NOW() AND used_at IS NULL`,
		userID,
	).Scan(&activeCount)
	if countErr == nil && activeCount > 0 {
		writeJSON(w, http.StatusOK, successMsg)
		return
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(15 * time.Minute)

	_, err = h.DB.Pool.Exec(r.Context(),
		`INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	if err != nil {
		log.Printf("[RESET] failed insert token for user %s: %v", userID, err)
		writeJSON(w, http.StatusOK, successMsg)
		return
	}

	frontendURL := utils.GetEnv("FRONTEND_URL", "https://propensuy-thinknalyze.vercel.app")
	resetURL := fmt.Sprintf("%s/account/reset?token=%s", frontendURL, token)
	go dispatchResetNotification(email, map[string]string{
		"full_name": fullName,
		"reset_url": resetURL,
	})

	writeJSON(w, http.StatusOK, successMsg)
}

// ValidateResetToken — GET /api/auth/reset-password/validate?token=...
func (h *APIHandler) ValidateResetToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	token := strings.TrimSpace(r.URL.Query().Get("token"))
	invalid := map[string]string{"error": "Tautan tidak valid atau sudah kadaluarsa. Silakan minta tautan reset baru."}

	if token == "" {
		writeJSON(w, http.StatusBadRequest, invalid)
		return
	}

	if h.DB == nil || h.DB.Pool == nil {
		writeJSON(w, http.StatusBadRequest, invalid)
		return
	}

	var exists bool
	err := h.DB.Pool.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM password_reset_tokens WHERE token = $1 AND expires_at > NOW() AND used_at IS NULL)`,
		token,
	).Scan(&exists)

	if err != nil || !exists {
		writeJSON(w, http.StatusBadRequest, invalid)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Token valid"})
}

// ResetPassword — POST /api/auth/reset-password
func (h *APIHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
		return
	}

	if len(body.NewPassword) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Password minimal 8 karakter"})
		return
	}

	token := strings.TrimSpace(body.Token)
	invalid := map[string]string{"error": "Tautan tidak valid atau sudah kadaluarsa. Silakan minta tautan reset baru."}

	if token == "" {
		writeJSON(w, http.StatusBadRequest, invalid)
		return
	}

	if h.DB == nil || h.DB.Pool == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	// Re-validate token with server-side expiry/used check
	var tokenID, userID string
	err := h.DB.Pool.QueryRow(r.Context(),
		`SELECT id, user_id FROM password_reset_tokens WHERE token = $1 AND expires_at > NOW() AND used_at IS NULL LIMIT 1`,
		token,
	).Scan(&tokenID, &userID)

	if err == pgx.ErrNoRows || err != nil {
		writeJSON(w, http.StatusBadRequest, invalid)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memproses password"})
		return
	}

	_, err = h.DB.Pool.Exec(r.Context(),
		`UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`,
		string(hashed), userID,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengubah password"})
		return
	}

	_, _ = h.DB.Pool.Exec(r.Context(),
		`UPDATE password_reset_tokens SET used_at = NOW() WHERE token = $1`, token,
	)

	writeJSON(w, http.StatusOK, map[string]string{"message": "Password berhasil diubah"})
}
