package login

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/master-abror/zaframework/core/concurrency"
	resp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/token"
	"github.com/master-abror/zaframework/core/utils"
)

func (c *Controller) ApiAuth(w http.ResponseWriter, r *http.Request) {
	// Pengecekan Content-Type
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		resp.ApiJSON(w, r, http.StatusUnsupportedMediaType, false, "Content-Type harus application/json", nil)
		return
	}

	// Batasi ukuran body (1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	// Decode payload
	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Anti field liar

	if err := decoder.Decode(&login); err != nil {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "Format JSON tidak valid atau field tidak dikenali", nil)
		return
	}

	// Validasi Input
	email := strings.TrimSpace(login.Email)
	password := strings.TrimSpace(login.Password)
	if email == "" || password == "" {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "Email dan Password wajib diisi!", nil)
		return
	}

	// Ambil data client
	browser, osName := utils.GetUserAgent(r)

	payload := map[string]any{
		"login_id":  email,
		"password":  password,
		"ip":        utils.GetClientIP(r),
		"browser":   browser,
		"os":        osName,
		"latitude":  "",
		"longitude": "",
	}

	// Dispatch ke worker/service
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "auth", payload, concurrency.PriorityHigh)
	if err != nil {
		// Menggunakan 401 Unauthorized jika gagal login
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, err.Error(), nil)
		return
	}

	// Casting result
	loginRes, ok := result.(LoginResult) // Pastikan struct LoginResult sudah di-import/didefinisikan
	if !ok {
		resp.ApiJSON(w, r, http.StatusInternalServerError, false, "Terjadi kesalahan internal (casting result)", nil)
		return
	}

	// Set otorisasi user
	var key string
	if loginRes.Token != "" {
		if key, err = token.SetUserAuthz(w, r, loginRes.Token); err != nil {
			resp.ApiJSON(w, r, http.StatusInternalServerError, false, "Gagal mengatur sesi: "+err.Error(), nil)
			return
		}
	}

	// Kirim balikan sukses
	resp.ApiJSON(w, r, http.StatusOK, true, "Otentikasi akun berhasil.", map[string]any{
		"token": key,
	})
}
