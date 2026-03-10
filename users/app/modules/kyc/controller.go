package kyc

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/master-abror/zaframework/core/concurrency"
	resp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/utils"
)

// Controller menangani HTTP request untuk fitur KYC
type Controller struct {
	Dispatcher *concurrency.Dispatcher
	Response   *resp.ResponseHelper
}

// NewController membuat instance baru KYC Controller
func NewController(d *concurrency.Dispatcher, r *resp.ResponseHelper) *Controller {
	return &Controller{
		Dispatcher: d,
		Response:   r,
	}
}

// Submit menangani POST /api/kyc/submit — multipart/form-data
func (c *Controller) Submit(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil user_id dari cookie token (auth middleware sudah validasi)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		// Fallback: coba ambil dari context (jika pakai AuthMiddleware)
		if ctxUser := r.Context().Value("current_user"); ctxUser != nil {
			if cu, ok := ctxUser.(resp.CurrentUser); ok {
				userID = cu.ID
			}
		}
	}

	if userID == "" {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid, silakan login kembali", nil)
		return
	}

	// 2. Validasi & upload file KTP (max 2MB, hanya jpg/png/pdf)
	// → Uploads to Cloudflare R2 when R2_* env vars are set; falls back to local disk otherwise.
	allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf"}
	maxSize := int64(2 * 1024 * 1024) // 2MB

	filename, err := utils.UploadFileToR2(r, "ktp_image", "kyc", allowedExts, maxSize)
	if err != nil {
		errMsg := err.Error()
		// Cek apakah error karena ukuran/format file
		if strings.Contains(errMsg, "terlalu besar") || strings.Contains(errMsg, "tidak diizinkan") || strings.Contains(errMsg, "bukan gambar") {
			resp.ApiJSON(w, r, http.StatusUnprocessableEntity, false, errMsg, nil)
			return
		}
		resp.ApiJSON(w, r, http.StatusUnprocessableEntity, false, "Gagal mengupload file KTP: "+errMsg, nil)
		return
	}

	// 3. Ambil field form lainnya
	fullName := strings.TrimSpace(r.FormValue("full_name"))
	nik := strings.TrimSpace(r.FormValue("nik"))
	address := strings.TrimSpace(r.FormValue("address"))
	birthdate := strings.TrimSpace(r.FormValue("birthdate"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	// 4. Dispatch ke service via worker pool
	payload := map[string]any{
		"user_id":   userID,
		"full_name": fullName,
		"nik":       nik,
		"address":   address,
		"birthdate": birthdate,
		"phone":     phone,
		"ktp_image": filename,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "kyc_submit", payload, concurrency.PriorityHigh)
	if err != nil {
		errMsg := err.Error()

		// Cek error duplikat NIK → 409 Conflict
		if strings.HasPrefix(errMsg, "DUPLICATE_NIK:") {
			resp.ApiJSON(w, r, http.StatusConflict, false, strings.TrimPrefix(errMsg, "DUPLICATE_NIK:"), nil)
			return
		}

		// Cek error duplikat KYC submission
		if strings.HasPrefix(errMsg, "DUPLICATE_KYC:") {
			resp.ApiJSON(w, r, http.StatusConflict, false, strings.TrimPrefix(errMsg, "DUPLICATE_KYC:"), nil)
			return
		}

		// Error validasi → 422
		resp.ApiJSON(w, r, http.StatusUnprocessableEntity, false, errMsg, nil)
		return
	}

	// 5. Sukses → 201 Created
	kycResult := result.(KYCSubmitResult)
	resp.ApiJSON(w, r, http.StatusCreated, true, kycResult.Message, map[string]any{
		"id":     kycResult.ID,
		"status": kycResult.Status,
	})
}

// Status menangani GET /api/kyc/status — cek status KYC user
func (c *Controller) Status(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		if ctxUser := r.Context().Value("current_user"); ctxUser != nil {
			if cu, ok := ctxUser.(resp.CurrentUser); ok {
				userID = cu.ID
			}
		}
	}

	if userID == "" {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid, silakan login kembali", nil)
		return
	}

	payload := map[string]any{
		"user_id": userID,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "kyc_status", payload, concurrency.PriorityNormal)
	if err != nil {
		fmt.Printf("[KYC STATUS ERROR] %v\n", err)
		resp.ApiJSON(w, r, http.StatusInternalServerError, false, "Gagal mengambil status KYC", nil)
		return
	}

	// Jika user belum pernah submit KYC, kembalikan 404 Not Found
	// agar FE bisa redirect ke form KYC baru
	if result == nil {
		resp.ApiJSON(w, r, http.StatusNotFound, false, "Belum ada pengajuan KYC", nil)
		return
	}

	kycStatus := result.(KYCStatusResult)
	
	// Response wajib memuat field status dan rejection_reason
	resp.ApiJSON(w, r, http.StatusOK, true, "Data KYC ditemukan", map[string]any{
		"status":          kycStatus.Status,
		"rejection_reason": kycStatus.RejectReason,
		"kyc":             kycStatus,
	})
}
