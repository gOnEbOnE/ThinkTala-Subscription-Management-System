package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

// CreateUser — POST /api/admin/users (buat user internal oleh Superadmin)
func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Cek role: hanya SUPERADMIN yang boleh akses
	// Cek X-User-Role langsung ATAU X-User-Level (untuk kasus assume-role)
	userRole := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	userLevel := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Level")))
	if userRole != "SUPERADMIN" && userLevel != "SUPERADMIN" {
		ehttp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak: hanya Superadmin yang dapat membuat akun", nil)
		return
	}

	// 2. Validasi Content-Type
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		ehttp.ApiJSON(w, r, http.StatusUnsupportedMediaType, false, "Content-Type harus application/json", nil)
		return
	}

	// 3. Parse JSON body
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input CreateUserInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Format JSON tidak valid", nil)
		return
	}

	// 4. Sanitize inputs
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Password = strings.TrimSpace(input.Password)
	input.Role = strings.TrimSpace(strings.ToUpper(input.Role))

	// 5. Validasi: semua field wajib diisi
	if input.FullName == "" || input.Email == "" || input.Password == "" || input.Role == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Semua field wajib diisi (full_name, email, password, role)", nil)
		return
	}

	// 6. Validasi email format
	if !utils.IsValidEmail(input.Email) {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Format email tidak valid", nil)
		return
	}

	// 7. Validasi password minimal 8 karakter
	if len(input.Password) < 8 {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Password minimal 8 karakter", nil)
		return
	}

	// 8. Validasi role harus ada di daftar yang diizinkan
	if !AllowedRoles[input.Role] {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Role tidak valid. Role yang diizinkan: OPERASIONAL, COMPLIANCE, MANAJEMEN, ADMIN_CS", nil)
		return
	}

	// 9. Dispatch ke worker
	adminEmail := strings.TrimSpace(r.Header.Get("X-User-Email"))

	payload := map[string]any{
		"full_name":   input.FullName,
		"email":       input.Email,
		"password":    input.Password,
		"role":        input.Role,
		"admin_email": adminEmail,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_create_user", payload, concurrency.PriorityHigh)
	if err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	createResult, ok := result.(CreateUserResult)
	if !ok {
		ehttp.ApiJSON(w, r, http.StatusInternalServerError, false, "Terjadi kesalahan internal", nil)
		return
	}

	ehttp.ApiJSON(w, r, http.StatusCreated, true, "Akun berhasil dibuat", map[string]any{
		"id":         createResult.ID,
		"full_name":  createResult.FullName,
		"email":      createResult.Email,
		"role":       createResult.Role,
		"status":     createResult.Status,
		"created_at": createResult.CreatedAt,
	})
}

// GetUsers — GET /api/admin/users (daftar user internal untuk Superadmin)
func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	// 1. Cek role: hanya SUPERADMIN
	userRole := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	userLevel := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Level")))
	if userRole != "SUPERADMIN" && userLevel != "SUPERADMIN" {
		ehttp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak: hanya Superadmin yang dapat mengakses data akun internal", nil)
		return
	}

	// 2. Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	
	params := GetUsersParams{
		Page:    page,
		PerPage: perPage,
		Role:    strings.TrimSpace(r.URL.Query().Get("role")),
		Status:  strings.TrimSpace(r.URL.Query().Get("status")),
		Search:  strings.TrimSpace(r.URL.Query().Get("search")),
	}

	// 3. Dispatch ke worker
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_get_users", params, concurrency.PriorityHigh)
	if err != nil {
		ehttp.ApiJSON(w, r, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	resp, ok := result.(GetUsersResponse)
	if !ok {
		ehttp.ApiJSON(w, r, http.StatusInternalServerError, false, "Terjadi kesalahan internal (invalid response format)", nil)
		return
	}

	// 4. Return success
	ehttp.ApiJSON(w, r, http.StatusOK, true, "Berhasil mengambil data", resp)
}

// GetUserDetail — GET /api/admin/users/{id} (detail satu user untuk Superadmin, PBI-53)
func (c *Controller) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	// 1. Cek role: hanya SUPERADMIN
	userRole := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	userLevel := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Level")))
	if userRole != "SUPERADMIN" && userLevel != "SUPERADMIN" {
		ehttp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak: hanya Superadmin yang dapat mengakses data akun internal", nil)
		return
	}

	// 2. Extract user ID from URL path
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "User ID diperlukan", nil)
		return
	}

	// 3. Dispatch ke worker
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_get_user_detail", userID, concurrency.PriorityHigh)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ehttp.ApiJSON(w, r, http.StatusNotFound, false, "Akun tidak ditemukan", nil)
			return
		}
		ehttp.ApiJSON(w, r, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	// 4. Return success
	ehttp.ApiJSON(w, r, http.StatusOK, true, "Berhasil mengambil data user", result)
}

// EditUser — PUT /api/admin/users/{id} (edit user internal oleh Superadmin, PBI-54)
func (c *Controller) EditUser(w http.ResponseWriter, r *http.Request) {
	// 1. Cek role: hanya SUPERADMIN
	userRole := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	userLevel := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Level")))
	if userRole != "SUPERADMIN" && userLevel != "SUPERADMIN" {
		ehttp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak: hanya Superadmin yang dapat mengedit akun internal", nil)
		return
	}

	// 2. Extract user ID
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "User ID diperlukan", nil)
		return
	}

	// 3. Parse body
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		ehttp.ApiJSON(w, r, http.StatusUnsupportedMediaType, false, "Content-Type harus application/json", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input EditUserInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Format JSON tidak valid", nil)
		return
	}

	// 4. Validate: at least one field
	hasName := input.FullName != nil
	hasRole := input.Role != nil

	if !hasName && !hasRole {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Minimal satu field (full_name atau role) harus diisi", nil)
		return
	}

	// 5. Validate full_name if provided
	if hasName {
		trimmed := strings.TrimSpace(*input.FullName)
		if trimmed == "" {
			ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Nama lengkap tidak boleh kosong", nil)
			return
		}
		input.FullName = &trimmed
	}

	// 6. Validate role if provided
	if hasRole {
		roleVal := strings.TrimSpace(strings.ToUpper(*input.Role))
		if !AllowedRoles[roleVal] {
			ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Role tidak valid. Role yang diizinkan: OPERASIONAL, COMPLIANCE, MANAJEMEN, ADMIN_CS", nil)
			return
		}
		input.Role = &roleVal
	}

	// 7. Build payload
	payload := map[string]any{
		"user_id":      userID,
		"performed_by": strings.TrimSpace(r.Header.Get("X-User-ID")),
	}
	if hasName {
		payload["full_name"] = *input.FullName
	}
	if hasRole {
		payload["role"] = *input.Role
	}

	// 8. Dispatch
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_edit_user", payload, concurrency.PriorityHigh)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ehttp.ApiJSON(w, r, http.StatusNotFound, false, "Akun tidak ditemukan", nil)
			return
		}
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	ehttp.ApiJSON(w, r, http.StatusOK, true, "Akun berhasil diperbarui", result)
}

// DeactivateUser — PATCH /api/admin/users/{id}/deactivate (PBI-55)
func (c *Controller) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Cek role: hanya SUPERADMIN
	userRole := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	userLevel := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Level")))
	if userRole != "SUPERADMIN" && userLevel != "SUPERADMIN" {
		ehttp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak: hanya Superadmin yang dapat menonaktifkan akun", nil)
		return
	}

	// 2. Extract user ID
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "User ID diperlukan", nil)
		return
	}

	// 3. Build payload
	payload := map[string]any{
		"user_id":      userID,
		"performed_by": strings.TrimSpace(r.Header.Get("X-User-ID")),
	}

	// 4. Dispatch
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_deactivate_user", payload, concurrency.PriorityHigh)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ehttp.ApiJSON(w, r, http.StatusNotFound, false, "Akun tidak ditemukan", nil)
			return
		}
		if err.Error() == "SUPERADMIN_PROTECTED" {
			ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Tidak dapat menonaktifkan akun Superadmin", nil)
			return
		}
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	ehttp.ApiJSON(w, r, http.StatusOK, true, "Akun berhasil dinonaktifkan", result)
}

// ReactivateUser — PATCH /api/admin/users/{id}/reactivate (PBI-55)
func (c *Controller) ReactivateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Cek role: hanya SUPERADMIN
	userRole := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	userLevel := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Level")))
	if userRole != "SUPERADMIN" && userLevel != "SUPERADMIN" {
		ehttp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak: hanya Superadmin yang dapat mengaktifkan kembali akun", nil)
		return
	}

	// 2. Extract user ID
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "User ID diperlukan", nil)
		return
	}

	// 3. Build payload
	payload := map[string]any{
		"user_id":      userID,
		"performed_by": strings.TrimSpace(r.Header.Get("X-User-ID")),
	}

	// 4. Dispatch
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_reactivate_user", payload, concurrency.PriorityHigh)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ehttp.ApiJSON(w, r, http.StatusNotFound, false, "Akun tidak ditemukan", nil)
			return
		}
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	ehttp.ApiJSON(w, r, http.StatusOK, true, "Akun berhasil diaktifkan kembali", result)
}

