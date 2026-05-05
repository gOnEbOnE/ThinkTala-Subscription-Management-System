package reset

import (
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

	var userID string
	err := h.DB.Pool.QueryRow(r.Context(),
		`SELECT id FROM users WHERE email = $1 AND status = 'active' LIMIT 1`, email,
	).Scan(&userID)
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

	frontendURL := utils.GetEnv("FRONTEND_URL", "http://localhost:3000")
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	go func() {
		sender := utils.NewEmailSender()
		subject := "Reset Kata Sandi ThinkNalyze"
		htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#f4f6f9;font-family:'Segoe UI',Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f9;padding:40px 0;">
    <tr><td align="center">
      <table width="520" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);">
        <tr><td style="background:#1a1c2e;padding:28px 36px;text-align:center;">
          <span style="color:#ffffff;font-size:22px;font-weight:700;letter-spacing:0.5px;">ThinkNalyze</span>
        </td></tr>
        <tr><td style="padding:36px 36px 28px;">
          <p style="margin:0 0 16px;font-size:16px;color:#1a1c2e;font-weight:600;">Halo,</p>
          <p style="margin:0 0 20px;font-size:14px;color:#4a5568;line-height:1.6;">
            Kami menerima permintaan untuk mereset kata sandi akun ThinkNalyze Anda.
            Klik tombol di bawah untuk membuat kata sandi baru. Tautan ini berlaku selama <strong>15 menit</strong>.
          </p>
          <table cellpadding="0" cellspacing="0" width="100%%"><tr><td align="center" style="padding:8px 0 24px;">
            <a href="%s" style="display:inline-block;background:#4e73df;color:#ffffff;text-decoration:none;font-weight:600;font-size:15px;padding:14px 36px;border-radius:8px;">
              Reset Kata Sandi
            </a>
          </td></tr></table>
          <p style="margin:0 0 8px;font-size:13px;color:#718096;">Atau salin tautan berikut ke browser Anda:</p>
          <p style="margin:0 0 24px;font-size:12px;color:#4e73df;word-break:break-all;">%s</p>
          <hr style="border:none;border-top:1px solid #e8ecf0;margin:0 0 20px;">
          <p style="margin:0;font-size:12px;color:#a0aec0;line-height:1.6;">
            Jika Anda tidak meminta reset kata sandi, abaikan email ini — akun Anda tetap aman.
          </p>
        </td></tr>
        <tr><td style="background:#f7f9fc;padding:16px 36px;text-align:center;">
          <p style="margin:0;font-size:11px;color:#a0aec0;">&copy; 2026 ThinkNalyze. All rights reserved.</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, resetURL, resetURL)
		if err := sender.SendHTMLEmail(email, subject, htmlBody); err != nil {
			log.Printf("[RESET] failed send email to %s: %v", email, err)
		}
	}()

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
