package register

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/master-abror/zaframework/core/utils"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// dispatchNotification mengirim event ke Notification Service untuk diproses berdasarkan template.
// Jika Notification Service tidak tersedia atau template belum ada, fallback ke SMTP langsung.
func dispatchNotification(eventType, channel, to string, vars map[string]string) {

	// TODO(queue): Aktifkan Redis queue setelah worker siap di-deploy
	if err := utils.PublishNotificationEvent(eventType, "email", to, vars); err == nil {
		log.Printf("[NOTIF] Event dipublish ke queue: event=%s to=%s", eventType, to)
		return
	}
	baseURL := utils.GetEnv("NOTIFICATION_SERVICE_URL", "http://localhost:5003")

	payload := map[string]any{
		"event_type": eventType,
		"channel":    channel,
		"to":         to,
		"vars":       vars,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/api/notifications/send", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("[NOTIF] Notification service tidak tersedia (%v), fallback SMTP", err)
		fallbackSendOTPEmail(to, vars)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]any
		json.NewDecoder(resp.Body).Decode(&result)
		log.Printf("[NOTIF] Gagal kirim via template (%d: %v), fallback SMTP", resp.StatusCode, result["error"])
		fallbackSendOTPEmail(to, vars)
	}
}

// fallbackSendOTPEmail mengirim email OTP langsung via SMTP jika Notification Service tidak tersedia.
func fallbackSendOTPEmail(to string, vars map[string]string) {
	name := vars["name"]
	otp := vars["otp"]
	smtpClient := utils.NewSMTPClient()
	subject := "Kode Verifikasi ThinkNalyze"
	body := fmt.Sprintf(
		"Halo %s,\n\nKode verifikasi Anda: %s\n\nKode ini berlaku selama 5 menit.\nJangan bagikan kode ini kepada siapa pun.\n\n— ThinkNalyze Team",
		name, otp,
	)
	if err := smtpClient.SendEmail(to, subject, body); err != nil {
		log.Printf("[OTP] Gagal mengirim email ke %s: %v", to, err)
	} else {
		log.Printf("[OTP] Kode verifikasi (fallback) terkirim ke %s", to)
	}
}

// ProcessRegisterJob — dipanggil oleh worker untuk proses registrasi
func (s *Service) ProcessRegisterJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	fullName := data["full_name"].(string)
	email := data["email"].(string)
	password := data["password"].(string)
	phone := data["no_telp"].(string)
	birthdate := data["tanggal_lahir"].(string)

	// 1. Validasi tanggal lahir tidak boleh di masa depan
	dob, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return nil, fmt.Errorf("Format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
	}
	if dob.After(time.Now()) {
		return nil, fmt.Errorf("Tanggal lahir tidak boleh di masa depan")
	}

	// 2. Cek email — jika sudah terdaftar, cek statusnya
	userID, status, err := s.repo.GetUserStatusByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengecek email: %v", err)
	}

	if userID != "" {
		if status == "active" {
			// Email sudah aktif → tolak
			return nil, fmt.Errorf("Email sudah terdaftar dan aktif")
		}
		// Status inactive → hapus data lama, izinkan register ulang
		err = s.repo.DeleteInactiveUser(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("Gagal memproses registrasi ulang: %v", err)
		}
		log.Printf("[REGISTER] User inactive %s dihapus untuk registrasi ulang", email)
	}

	// 3. Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengenkripsi password")
	}

	// 4. Generate UUID
	newUUID, err := utils.CreateUUID()
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat ID user")
	}
	newUserID := newUUID.String()

	// 5. Simpan user ke database (status: inactive)
	err = s.repo.CreateUser(ctx, newUserID, fullName, email, hashedPassword, phone, birthdate)
	if err != nil {
		return nil, fmt.Errorf("Gagal menyimpan data: %v", err)
	}

	// 6. Generate OTP 6 digit
	otpCode := utils.GenerateRandomNumber(6)
	otpExpiry := time.Now().Add(5 * time.Minute)

	// 7. Simpan OTP ke database
	err = s.repo.SaveOTP(ctx, newUserID, email, otpCode, otpExpiry)
	if err != nil {
		return nil, fmt.Errorf("Gagal menyimpan kode OTP")
	}

	// 8. Kirim OTP via Notification Service (async)
	go dispatchNotification("otp_verification", "email", email, map[string]string{
		"name": fullName,
		"otp":  otpCode,
	})

	return RegisterResult{
		UserID:  newUserID,
		Email:   email,
		Message: "Cek email Anda untuk kode OTP",
	}, nil
}

// ProcessVerifyOTPJob — dipanggil oleh worker untuk verifikasi OTP
func (s *Service) ProcessVerifyOTPJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	email := data["email"].(string)
	otpCode := data["otp_code"].(string)

	// 1. Cari OTP yang valid
	otp, err := s.repo.GetValidOTP(ctx, email, otpCode)
	if err != nil {
		return nil, fmt.Errorf("Gagal memverifikasi OTP: %v", err)
	}
	if otp == nil {
		return nil, fmt.Errorf("Kode OTP tidak valid")
	}

	// 2. Cek expired
	if time.Now().After(otp.ExpiresAt) {
		return nil, fmt.Errorf("Kode OTP sudah kedaluwarsa")
	}

	// 3. Mark OTP as used
	err = s.repo.MarkOTPUsed(ctx, otp.ID)
	if err != nil {
		return nil, fmt.Errorf("Gagal memproses OTP")
	}

	// 4. Activate user
	err = s.repo.ActivateUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengaktifkan akun")
	}

	return map[string]any{
		"message": "Akun berhasil diaktifkan",
		"email":   email,
	}, nil
}

// ProcessResendOTPJob — kirim ulang OTP
func (s *Service) ProcessResendOTPJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	email := data["email"].(string)

	// 1. Cari user ID
	userID, err := s.repo.GetUserIDByEmail(ctx, email)
	if err != nil || userID == "" {
		return nil, fmt.Errorf("Email tidak ditemukan")
	}

	// 2. Generate OTP baru
	otpCode := utils.GenerateRandomNumber(6)
	otpExpiry := time.Now().Add(5 * time.Minute)

	// 3. Simpan OTP baru
	err = s.repo.SaveOTP(ctx, userID, email, otpCode, otpExpiry)
	if err != nil {
		return nil, fmt.Errorf("Gagal menyimpan kode OTP")
	}

	// 4. Kirim via Notification Service (async)
	go dispatchNotification("otp_verification", "email", email, map[string]string{
		"otp": otpCode,
	})

	return map[string]any{
		"message": "Kode OTP baru telah dikirim ke email Anda",
	}, nil
}
