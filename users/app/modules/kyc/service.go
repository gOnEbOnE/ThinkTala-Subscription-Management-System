package kyc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/master-abror/zaframework/core/utils"
)

// Service menangani business logic KYC
type Service struct {
	repo Repository
}

// NewService membuat instance baru KYC Service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// dispatchNotification mengirim event KYC ke Notification Service.
// Urutan: 1) Redis queue → 2) HTTP langsung → 3) SMTP fallback.
func dispatchKYCNotification(eventType, to string, vars map[string]string) {
	// 1. Coba Redis queue terlebih dahulu (async, reliable)
	if err := utils.PublishNotificationEvent(eventType, "email", to, vars); err == nil {
		log.Printf("[KYC NOTIF] Event dipublish ke queue: event=%s to=%s", eventType, to)
		return
	}

	// 2. Fallback: HTTP langsung ke notification service
	baseURL := utils.GetEnv("NOTIFICATION_SERVICE_URL", "http://localhost:5003")
	payload := map[string]any{
		"event_type": eventType,
		"channel":    "email",
		"to":         to,
		"vars":       vars,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/api/notifications/send", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("[KYC NOTIF] Notification service tidak tersedia (%v), fallback SMTP", err)
		// Derive status dari eventType agar fallback email sesuai
		statusFallback := "approved"
		if strings.Contains(eventType, "reject") {
			statusFallback = "rejected"
		}
		sendKYCFallbackEmail(to, vars["full_name"], statusFallback, vars["reject_reason"])
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]any
		json.NewDecoder(resp.Body).Decode(&result)
		log.Printf("[KYC NOTIF] Gagal kirim via template (%d), fallback SMTP", resp.StatusCode)
		statusFallback := "approved"
		if strings.Contains(eventType, "reject") {
			statusFallback = "rejected"
		}
		sendKYCFallbackEmail(to, vars["full_name"], statusFallback, vars["reject_reason"])
	}
}

// sendKYCFallbackEmail mengirim email KYC langsung via SMTP jika Notification Service tidak tersedia.
func sendKYCFallbackEmail(to, fullName, status, rejectReason string) {
	smtp := utils.NewSMTPClient()
	var subject, body string

	if status == "approved" {
		subject = "✓ Verifikasi KYC Berhasil - ThinkNalyze"
		body = fmt.Sprintf(
			"Halo %s,\n\n"+
				"Selamat! Verifikasi identitas (KYC) Anda telah disetujui.\n\n"+
				"Anda sekarang memiliki akses penuh ke semua fitur ThinkNalyze, termasuk:\n"+
				"• Market Insight\n"+
				"• Deep Scanner\n"+
				"• Ask Nizza AI Assistant\n"+
				"• Trading Features\n\n"+
				"Silakan login ke dashboard Anda untuk mulai menggunakan layanan kami.\n\n"+
				"Terima kasih telah mempercayai ThinkNalyze.\n\n"+
				"Salam,\n"+
				"Tim ThinkNalyze",
			fullName,
		)
	} else if status == "rejected" {
		subject = "⚠️ Verifikasi KYC Memerlukan Perbaikan - ThinkNalyze"
		body = fmt.Sprintf(
			"Halo %s,\n\n"+
				"Mohon maaf, verifikasi identitas (KYC) Anda memerlukan perbaikan.\n\n"+
				"ALASAN PENOLAKAN:\n"+
				"%s\n\n"+
				"LANGKAH SELANJUTNYA:\n"+
				"1. Login ke dashboard Anda\n"+
				"2. Buka menu KYC Verification\n"+
				"3. Klik tombol \"Perbaiki Data KYC\"\n"+
				"4. Perbaiki data sesuai alasan penolakan di atas\n"+
				"5. Submit ulang untuk direview kembali\n\n"+
				"Jika ada pertanyaan, silakan hubungi tim support kami.\n\n"+
				"Salam,\n"+
				"Tim ThinkNalyze",
			fullName,
			rejectReason,
		)
	} else {
		return
	}

	if err := smtp.SendEmail(to, subject, body); err != nil {
		log.Printf("[KYC NOTIFICATION] Gagal mengirim fallback email ke %s: %v", to, err)
	} else {
		log.Printf("[KYC NOTIFICATION] Fallback email terkirim ke %s", to)
	}
}

// ProcessKYCSubmitJob — handler untuk dispatcher job "kyc_submit"
// Dipanggil oleh concurrency dispatcher dari controller
func (s *Service) ProcessKYCSubmitJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	userID, _ := data["user_id"].(string)
	fullName, _ := data["full_name"].(string)
	nik, _ := data["nik"].(string)
	address, _ := data["address"].(string)
	birthdate, _ := data["birthdate"].(string)
	phone, _ := data["phone"].(string)
	ktpImage, _ := data["ktp_image"].(string)

	// ========== VALIDASI ==========

	// 1. Field wajib tidak boleh kosong
	fullName = strings.TrimSpace(fullName)
	nik = strings.TrimSpace(nik)
	address = strings.TrimSpace(address)
	birthdate = strings.TrimSpace(birthdate)
	phone = strings.TrimSpace(phone)

	if fullName == "" || nik == "" || address == "" || birthdate == "" || phone == "" || ktpImage == "" {
		return nil, fmt.Errorf("semua field wajib diisi (nama, NIK, alamat, tanggal lahir, telepon, KTP)")
	}

	// 2. Validasi NIK harus 16 digit angka
	nikRegex := regexp.MustCompile(`^\d{16}$`)
	if !nikRegex.MatchString(nik) {
		return nil, fmt.Errorf("NIK harus 16 digit angka")
	}

	// 3. Validasi tanggal lahir (format YYYY-MM-DD)
	_, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return nil, fmt.Errorf("format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
	}

	// 4. Validasi nomor telepon (10-15 digit, boleh diawali +)
	phoneRegex := regexp.MustCompile(`^\+?\d{10,15}$`)
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	if !phoneRegex.MatchString(cleanPhone) {
		return nil, fmt.Errorf("nomor telepon tidak valid (10-15 digit)")
	}

	// 5. Cek apakah user sudah punya KYC pending/approved
	hasPending, err := s.repo.UserHasPendingKYC(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengecek status KYC: %v", err)
	}
	if hasPending {
		return nil, fmt.Errorf("DUPLICATE_KYC:Anda sudah memiliki pengajuan KYC yang aktif")
	}

	// 6. Cek apakah user punya KYC yang sudah di-reject → boleh re-submit
	existingKYC, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengecek data KYC: %v", err)
	}

	isResubmit := existingKYC != nil && existingKYC.Status == "rejected"

	if !isResubmit {
		// Hanya cek duplikat NIK kalau bukan re-submit dari user yang sama
		nikExistsForOther, err := s.repo.NIKExistsForOtherUser(ctx, nik, userID)
		if err != nil {
			return nil, fmt.Errorf("gagal mengecek NIK: %v", err)
		}
		if nikExistsForOther {
			return nil, fmt.Errorf("DUPLICATE_NIK:NIK sudah terdaftar di sistem")
		}
	}

	// ========== SIMPAN KE DATABASE ==========
	if isResubmit {
		// UPDATE record lama: reset ke pending, hapus rejection reason
		if err := s.repo.UpdateResubmission(ctx, existingKYC.ID, fullName, nik, address, birthdate, cleanPhone, ktpImage); err != nil {
			return nil, fmt.Errorf("gagal memperbarui data KYC: %v", err)
		}
		return KYCSubmitResult{
			ID:      existingKYC.ID,
			Status:  "pending",
			Message: "Dokumen KYC berhasil diperbarui dan sedang dalam proses verifikasi ulang",
		}, nil
	}

	// Buat submission baru
	submission := &KYCSubmission{
		ID:        uuid.New().String(),
		UserID:    userID,
		FullName:  fullName,
		NIK:       nik,
		Address:   address,
		Birthdate: birthdate,
		Phone:     cleanPhone,
		KTPImage:  ktpImage,
		Status:    "pending",
	}

	if err := s.repo.CreateSubmission(ctx, submission); err != nil {
		return nil, fmt.Errorf("gagal menyimpan data KYC: %v", err)
	}

	return KYCSubmitResult{
		ID:      submission.ID,
		Status:  "pending",
		Message: "Dokumen KYC berhasil dikirim dan sedang dalam proses verifikasi",
	}, nil
}

// ProcessKYCStatusJob — handler untuk dispatcher job "kyc_status"
func (s *Service) ProcessKYCStatusJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	userID, _ := data["user_id"].(string)
	if userID == "" {
		return nil, fmt.Errorf("user_id wajib diisi")
	}

	sub, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data KYC: %v", err)
	}

	if sub == nil {
		return nil, nil // Belum pernah submit KYC
	}

	return KYCStatusResult{
		ID:           sub.ID,
		FullName:     sub.FullName,
		NIK:          sub.NIK,
		Address:      sub.Address,
		Birthdate:    sub.Birthdate,
		Phone:        sub.Phone,
		KTPImage:     sub.KTPImage,
		Status:       sub.Status,
		RejectReason: sub.RejectReason,
		ReviewedAt:   sub.ReviewedAt,
		CreatedAt:    sub.CreatedAt,
	}, nil
}

// ========== ADMIN SERVICE METHODS ==========

// ProcessAdminKYCListJob — handler untuk dispatcher job "admin_kyc_list"
func (s *Service) ProcessAdminKYCListJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	role, _ := data["role"].(string)
	if role != "COMPLIANCE" && role != "ADMIN_KYC" && role != "SUPERADMIN" {
		return nil, fmt.Errorf("FORBIDDEN:Anda tidak memiliki akses ke halaman ini")
	}

	status, _ := data["status"].(string)
	search, _ := data["search"].(string)
	page, _ := data["page"].(int)
	limit, _ := data["limit"].(int)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	items, total, err := s.repo.ListAll(ctx, status, search, page, limit)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data KYC: %v", err)
	}

	if items == nil {
		items = []KYCListItem{}
	}

	return map[string]any{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	}, nil
}

// ProcessAdminKYCDetailJob — handler untuk dispatcher job "admin_kyc_detail"
func (s *Service) ProcessAdminKYCDetailJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	role, _ := data["role"].(string)
	if role != "COMPLIANCE" && role != "ADMIN_KYC" && role != "SUPERADMIN" {
		return nil, fmt.Errorf("FORBIDDEN:Anda tidak memiliki akses")
	}

	kycID, _ := data["kyc_id"].(string)
	if kycID == "" {
		return nil, fmt.Errorf("kyc_id wajib diisi")
	}

	detail, err := s.repo.GetDetailByID(ctx, kycID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail KYC: %v", err)
	}

	if detail == nil {
		return nil, fmt.Errorf("NOT_FOUND:Data KYC tidak ditemukan")
	}

	return detail, nil
}

// ProcessAdminKYCReviewJob — handler untuk dispatcher job "admin_kyc_review"
func (s *Service) ProcessAdminKYCReviewJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	role, _ := data["role"].(string)
	if role != "COMPLIANCE" && role != "ADMIN_KYC" && role != "SUPERADMIN" {
		return nil, fmt.Errorf("FORBIDDEN:Anda tidak memiliki akses")
	}

	kycID, _ := data["kyc_id"].(string)
	reviewerID, _ := data["reviewer_id"].(string)
	action, _ := data["action"].(string)
	rejectReason, _ := data["reject_reason"].(string)

	if kycID == "" {
		return nil, fmt.Errorf("kyc_id wajib diisi")
	}

	// Cek KYC exists dan masih pending
	detail, err := s.repo.GetDetailByID(ctx, kycID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data KYC: %v", err)
	}
	if detail == nil {
		return nil, fmt.Errorf("NOT_FOUND:Data KYC tidak ditemukan")
	}
	if detail.Status != "pending" {
		return nil, fmt.Errorf("CONFLICT:KYC ini sudah di-review sebelumnya (status: %s)", detail.Status)
	}

	var newStatus string
	switch action {
	case "approve":
		newStatus = "approved"
	case "reject":
		newStatus = "rejected"
		if strings.TrimSpace(rejectReason) == "" {
			return nil, fmt.Errorf("Alasan penolakan wajib diisi")
		}
	default:
		return nil, fmt.Errorf("action tidak valid (gunakan 'approve' atau 'reject')")
	}

	if err := s.repo.UpdateStatus(ctx, kycID, newStatus, reviewerID, rejectReason); err != nil {
		return nil, fmt.Errorf("gagal mengupdate status KYC: %v", err)
	}

	// Kirim email notification ke user (async)
	s.sendKYCNotification(detail.Email, detail.FullName, newStatus, rejectReason)

	msg := "KYC berhasil di-approve"
	if newStatus == "rejected" {
		msg = "KYC berhasil ditolak"
	}

	return map[string]any{
		"id":      kycID,
		"status":  newStatus,
		"message": msg,
	}, nil
}

// sendKYCNotification mengirim email notifikasi ke user saat KYC di-approve/reject.
// dispatchKYCNotification sudah menangani fallback SMTP secara internal —
// jangan panggil sendKYCFallbackEmail di sini agar tidak terjadi double message.
func (s *Service) sendKYCNotification(email, fullName, status, rejectReason string) {
	go func() {
		switch status {
		case "approved":
			dispatchKYCNotification("kyc_approved", email, map[string]string{
				"full_name": fullName,
			})
		case "rejected":
			dispatchKYCNotification("kyc_rejected", email, map[string]string{
				"full_name":     fullName,
				"reject_reason": rejectReason,
			})
		}
	}()
}
