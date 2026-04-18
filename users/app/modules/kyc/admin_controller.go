package kyc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/master-abror/zaframework/core/concurrency"
	resp "github.com/master-abror/zaframework/core/http"
)

// AdminController menangani HTTP request untuk admin KYC review
type AdminController struct {
	Dispatcher *concurrency.Dispatcher
	Response   *resp.ResponseHelper
}

// NewAdminController membuat instance baru Admin KYC Controller
func NewAdminController(d *concurrency.Dispatcher, r *resp.ResponseHelper) *AdminController {
	return &AdminController{
		Dispatcher: d,
		Response:   r,
	}
}

// List menangani GET /api/admin/kyc — daftar semua pengajuan KYC
func (c *AdminController) List(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-Role")
	if role == "" {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid", nil)
		return
	}

	// Parse query params
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// PBI-7: Validate status parameter — only valid values allowed
	if status != "" {
		upperStatus := strings.ToUpper(status)
		validStatuses := map[string]bool{"PENDING": true, "APPROVED": true, "REJECTED": true}
		if !validStatuses[upperStatus] {
			resp.ApiJSON(w, r, http.StatusBadRequest, false, "Status tidak valid. Gunakan: PENDING, APPROVED, atau REJECTED", nil)
			return
		}
		status = strings.ToLower(upperStatus)
	}

	payload := map[string]any{
		"role":   role,
		"status": status,
		"search": search,
		"page":   page,
		"limit":  limit,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_kyc_list", payload, concurrency.PriorityNormal)
	if err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(errMsg, "FORBIDDEN:") {
			resp.ApiJSON(w, r, http.StatusForbidden, false, strings.TrimPrefix(errMsg, "FORBIDDEN:"), nil)
			return
		}
		if strings.HasPrefix(errMsg, "BAD_REQUEST:") {
			resp.ApiJSON(w, r, http.StatusBadRequest, false, strings.TrimPrefix(errMsg, "BAD_REQUEST:"), nil)
			return
		}
		resp.ApiJSON(w, r, http.StatusInternalServerError, false, errMsg, nil)
		return
	}

	resp.ApiJSON(w, r, http.StatusOK, true, "Data KYC berhasil diambil", result)
}

// Detail menangani GET /api/admin/kyc/{id} — detail KYC
func (c *AdminController) Detail(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-Role")
	if role == "" {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid", nil)
		return
	}

	// Extract ID from Go 1.22+ path param, fallback to manual
	kycID := r.PathValue("id")
	if kycID == "" {
		kycID = extractKYCID(r.URL.Path)
	}
	if kycID == "" {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "ID KYC tidak valid", nil)
		return
	}

	payload := map[string]any{
		"role":   role,
		"kyc_id": kycID,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_kyc_detail", payload, concurrency.PriorityNormal)
	if err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(errMsg, "FORBIDDEN:") {
			resp.ApiJSON(w, r, http.StatusForbidden, false, strings.TrimPrefix(errMsg, "FORBIDDEN:"), nil)
			return
		}
		if strings.HasPrefix(errMsg, "NOT_FOUND:") {
			resp.ApiJSON(w, r, http.StatusNotFound, false, strings.TrimPrefix(errMsg, "NOT_FOUND:"), nil)
			return
		}
		resp.ApiJSON(w, r, http.StatusInternalServerError, false, errMsg, nil)
		return
	}

	resp.ApiJSON(w, r, http.StatusOK, true, "Detail KYC ditemukan", result)
}

// Approve menangani POST /api/admin/kyc/{id}/approve
func (c *AdminController) Approve(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-Role")
	reviewerID := r.Header.Get("X-User-ID")
	if role == "" || reviewerID == "" {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid", nil)
		return
	}

	// Extract ID from Go 1.22+ path param, fallback to manual
	kycID := r.PathValue("id")
	if kycID == "" {
		kycID = extractKYCIDFromAction(r.URL.Path, "approve")
	}
	if kycID == "" {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "ID KYC tidak valid", nil)
		return
	}

	payload := map[string]any{
		"role":        role,
		"kyc_id":      kycID,
		"reviewer_id": reviewerID,
		"action":      "approve",
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_kyc_review", payload, concurrency.PriorityHigh)
	if err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(errMsg, "FORBIDDEN:") {
			resp.ApiJSON(w, r, http.StatusForbidden, false, strings.TrimPrefix(errMsg, "FORBIDDEN:"), nil)
			return
		}
		if strings.HasPrefix(errMsg, "NOT_FOUND:") {
			resp.ApiJSON(w, r, http.StatusNotFound, false, strings.TrimPrefix(errMsg, "NOT_FOUND:"), nil)
			return
		}
		if strings.HasPrefix(errMsg, "CONFLICT:") {
			resp.ApiJSON(w, r, http.StatusConflict, false, strings.TrimPrefix(errMsg, "CONFLICT:"), nil)
			return
		}
		resp.ApiJSON(w, r, http.StatusUnprocessableEntity, false, errMsg, nil)
		return
	}

	resp.ApiJSON(w, r, http.StatusOK, true, "KYC berhasil di-approve", result)
}

// Reject menangani POST /api/admin/kyc/{id}/reject
func (c *AdminController) Reject(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-Role")
	reviewerID := r.Header.Get("X-User-ID")
	if role == "" || reviewerID == "" {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid", nil)
		return
	}

	kycID := r.PathValue("id")
	if kycID == "" {
		kycID = extractKYCIDFromAction(r.URL.Path, "reject")
	}
	if kycID == "" {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "ID KYC tidak valid", nil)
		return
	}

	// Parse JSON body for reason
	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "Format request tidak valid", nil)
		return
	}

	if strings.TrimSpace(body.Reason) == "" {
		resp.ApiJSON(w, r, http.StatusUnprocessableEntity, false, "Alasan penolakan wajib diisi", nil)
		return
	}

	payload := map[string]any{
		"role":          role,
		"kyc_id":        kycID,
		"reviewer_id":   reviewerID,
		"action":        "reject",
		"reject_reason": body.Reason,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_kyc_review", payload, concurrency.PriorityHigh)
	if err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(errMsg, "FORBIDDEN:") {
			resp.ApiJSON(w, r, http.StatusForbidden, false, strings.TrimPrefix(errMsg, "FORBIDDEN:"), nil)
			return
		}
		if strings.HasPrefix(errMsg, "NOT_FOUND:") {
			resp.ApiJSON(w, r, http.StatusNotFound, false, strings.TrimPrefix(errMsg, "NOT_FOUND:"), nil)
			return
		}
		if strings.HasPrefix(errMsg, "CONFLICT:") {
			resp.ApiJSON(w, r, http.StatusConflict, false, strings.TrimPrefix(errMsg, "CONFLICT:"), nil)
			return
		}
		resp.ApiJSON(w, r, http.StatusUnprocessableEntity, false, errMsg, nil)
		return
	}

	resp.ApiJSON(w, r, http.StatusOK, true, "KYC berhasil ditolak", result)
}

// extractKYCID extracts ID from /api/admin/kyc/{id}
func extractKYCID(path string) string {
	// /api/admin/kyc/some-uuid-here
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	// parts: ["", "api", "admin", "kyc", "{id}"]
	if len(parts) >= 5 {
		id := parts[4]
		// Pastikan bukan action (approve/reject)
		if id != "approve" && id != "reject" && id != "" {
			return id
		}
	}
	return ""
}

// extractKYCIDFromAction extracts ID from /api/admin/kyc/{id}/action
func extractKYCIDFromAction(path string, action string) string {
	// /api/admin/kyc/some-uuid-here/approve
	suffix := "/" + action
	if !strings.HasSuffix(path, suffix) {
		return ""
	}
	trimmed := strings.TrimSuffix(path, suffix)
	parts := strings.Split(trimmed, "/")
	if len(parts) >= 5 {
		return parts[4]
	}
	return ""
}

// AdminKYCRouter menangani routing untuk /api/admin/kyc/*
// Go 1.22+ supports path params via r.PathValue()
func (c *AdminController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// GET /api/admin/kyc — list
	if r.Method == http.MethodGet && (path == "/api/admin/kyc" || path == "/api/admin/kyc/") {
		c.List(w, r)
		return
	}

	// POST /api/admin/kyc/{id}/approve
	if r.Method == http.MethodPost && strings.HasSuffix(path, "/approve") {
		c.Approve(w, r)
		return
	}

	// POST /api/admin/kyc/{id}/reject
	if r.Method == http.MethodPost && strings.HasSuffix(path, "/reject") {
		c.Reject(w, r)
		return
	}

	// GET /api/admin/kyc/{id} — detail
	if r.Method == http.MethodGet {
		id := r.PathValue("id")
		if id == "" {
			id = extractKYCID(path)
		}
		if id != "" {
			c.Detail(w, r)
			return
		}
	}

	resp.ApiJSON(w, r, http.StatusNotFound, false, "Route not found", nil)
}
