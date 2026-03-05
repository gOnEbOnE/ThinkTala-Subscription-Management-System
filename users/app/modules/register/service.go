package register

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/master-abror/zaframework/core/utils"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
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

	// 8. Kirim OTP via email (async)
	go func() {
		smtp := utils.NewSMTPClient()
		subject := "Kode Verifikasi ThinkNalyze"
		body := fmt.Sprintf(
			"Halo %s,\n\nKode verifikasi Anda: %s\n\nKode ini berlaku selama 5 menit.\nJangan bagikan kode ini kepada siapa pun.\n\n— ThinkNalyze Team",
			fullName, otpCode,
		)
		if err := smtp.SendEmail(email, subject, body); err != nil {
			log.Printf("[OTP] Gagal mengirim email ke %s: %v", email, err)
		} else {
			log.Printf("[OTP] Kode verifikasi terkirim ke %s", email)
		}
	}()

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

	// 4. Kirim email
	go func() {
		smtp := utils.NewSMTPClient()
		subject := "Kode Verifikasi Baru - ThinkNalyze"
		body := fmt.Sprintf(
			"Kode verifikasi baru Anda: %s\n\nKode ini berlaku selama 5 menit.\n\n— ThinkNalyze Team",
			otpCode,
		)
		if err := smtp.SendEmail(email, subject, body); err != nil {
			log.Printf("[OTP] Gagal mengirim ulang email ke %s: %v", email, err)
		} else {
			log.Printf("[OTP] Kode verifikasi ulang terkirim ke %s", email)
		}
	}()

	return map[string]any{
		"message": "Kode OTP baru telah dikirim ke email Anda",
	}, nil
}
