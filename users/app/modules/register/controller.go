package register

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/utils"
)

type Controller struct {
	Dispatcher *concurrency.Dispatcher
	Response   *ehttp.ResponseHelper
}

func NewController(d *concurrency.Dispatcher, r *ehttp.ResponseHelper) *Controller {
	return &Controller{
		Dispatcher: d,
		Response:   r,
	}
}

// Register — GET /account/register (render halaman)
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"AppName": "ThinkNalyze",
		"Year":    time.Now().Year(),
	}
	c.Response.View(w, r, "public/views/register/page.html", "Buat Akun | ThinkNalyze", data)
}

// VerifyOTPPage — GET /account/verify-otp (render halaman OTP)
func (c *Controller) VerifyOTPPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"AppName": "ThinkNalyze",
		"Year":    time.Now().Year(),
	}
	c.Response.View(w, r, "public/views/register/verify-otp.html", "Verifikasi OTP | ThinkNalyze", data)
}

// Submit — POST /api/auth/register (proses registrasi)
func (c *Controller) Submit(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		ehttp.ApiJSON(w, r, http.StatusUnsupportedMediaType, false, "Content-Type harus application/json", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input RegisterInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Format JSON tidak valid", nil)
		return
	}

	// Validasi struct
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Password = strings.TrimSpace(input.Password)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Birthdate = strings.TrimSpace(input.Birthdate)

	// Validasi semua field required
	if input.FullName == "" || input.Email == "" || input.Password == "" || input.Phone == "" || input.Birthdate == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Semua field wajib diisi", nil)
		return
	}

	// Validasi email format
	if !utils.IsValidEmail(input.Email) {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Format email tidak valid", nil)
		return
	}

	// Validasi password minimal 8 karakter
	if len(input.Password) < 8 {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Password minimal 8 karakter", nil)
		return
	}

	// Dispatch ke worker
	payload := map[string]any{
		"full_name":     input.FullName,
		"email":         input.Email,
		"password":      input.Password,
		"no_telp":       input.Phone,
		"tanggal_lahir": input.Birthdate,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "register", payload, concurrency.PriorityHigh)
	if err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	regResult, ok := result.(RegisterResult)
	if !ok {
		ehttp.ApiJSON(w, r, http.StatusInternalServerError, false, "Terjadi kesalahan internal", nil)
		return
	}

	ehttp.ApiJSON(w, r, http.StatusCreated, true, "Cek email Anda untuk kode OTP", map[string]any{
		"user_id": regResult.UserID,
		"email":   regResult.Email,
	})
}

// VerifyOTP — POST /api/auth/verify-otp (verifikasi OTP)
func (c *Controller) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		ehttp.ApiJSON(w, r, http.StatusUnsupportedMediaType, false, "Content-Type harus application/json", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input OTPVerifyInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Format JSON tidak valid", nil)
		return
	}

	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.OTPCode = strings.TrimSpace(input.OTPCode)

	if input.Email == "" || input.OTPCode == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Email dan kode OTP wajib diisi", nil)
		return
	}

	if len(input.OTPCode) != 6 {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Kode OTP harus 6 digit", nil)
		return
	}

	payload := map[string]any{
		"email":    input.Email,
		"otp_code": input.OTPCode,
	}

	_, err := c.Dispatcher.DispatchAndWait(r.Context(), "verify_otp", payload, concurrency.PriorityHigh)
	if err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	ehttp.ApiJSON(w, r, http.StatusOK, true, "Akun berhasil diaktifkan", nil)
}

// ResendOTP — POST /api/auth/resend-otp (kirim ulang OTP)
func (c *Controller) ResendOTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		ehttp.ApiJSON(w, r, http.StatusUnsupportedMediaType, false, "Content-Type harus application/json", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Format JSON tidak valid", nil)
		return
	}

	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Email wajib diisi", nil)
		return
	}

	payload := map[string]any{
		"email": input.Email,
	}

	_, err := c.Dispatcher.DispatchAndWait(r.Context(), "resend_otp", payload, concurrency.PriorityHigh)
	if err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	ehttp.ApiJSON(w, r, http.StatusOK, true, "Kode OTP baru telah dikirim ke email Anda", nil)
}
