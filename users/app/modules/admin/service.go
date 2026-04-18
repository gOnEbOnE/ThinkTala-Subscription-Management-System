package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/master-abror/zaframework/core/utils"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// dispatchNotification mengirim event ke Notification Service untuk diproses.
// Menggunakan pattern yang sama dengan register/service.go
func dispatchNotification(eventType, channel, to string, vars map[string]string) {
	// Coba kirim via Redis queue terlebih dahulu
	if err := utils.PublishNotificationEvent(eventType, "email", to, vars); err == nil {
		log.Printf("[ADMIN-NOTIF] Event dipublish ke queue: event=%s to=%s", eventType, to)
		return
	}

	// Fallback: kirim langsung ke Notification Service via HTTP
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
		log.Printf("[ADMIN-NOTIF] Notification service tidak tersedia (%v), fallback SMTP", err)
		fallbackSendCredentialsEmail(to, vars)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ADMIN-NOTIF] Gagal kirim via template (status %d), fallback SMTP", resp.StatusCode)
		fallbackSendCredentialsEmail(to, vars)
	}
}

// fallbackSendCredentialsEmail mengirim email kredensial langsung via SMTP
func fallbackSendCredentialsEmail(to string, vars map[string]string) {
	name := vars["name"]
	email := vars["email"]
	password := vars["password"]
	role := vars["role"]

	smtpClient := utils.NewSMTPClient()
	subject := "Akun ThinkNalyze Anda Telah Dibuat"
	body := fmt.Sprintf(
		"Halo %s,\n\nAkun ThinkNalyze Anda telah berhasil dibuat oleh administrator.\n\nBerikut kredensial login Anda:\nEmail: %s\nPassword: %s\nRole: %s\n\nSilakan login di ThinkNalyze dan segera ubah password Anda.\n\n-- ThinkNalyze Team",
		name, email, password, role,
	)
	if err := smtpClient.SendEmail(to, subject, body); err != nil {
		log.Printf("[ADMIN-NOTIF] Gagal mengirim email ke %s: %v", to, err)
	} else {
		log.Printf("[ADMIN-NOTIF] Email kredensial (fallback) terkirim ke %s", to)
	}
}

// ProcessCreateUserJob — dipanggil oleh worker untuk proses pembuatan user internal
func (s *Service) ProcessCreateUserJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	fullName := data["full_name"].(string)
	email := data["email"].(string)
	password := data["password"].(string)
	role := strings.ToUpper(data["role"].(string))
	adminEmail, _ := data["admin_email"].(string)

	// 1. Validasi role
	if !AllowedRoles[role] {
		log.Printf("[ADMIN] action=CREATE_USER admin=%s target=%s role=%s result=FAILED reason=invalid_role", adminEmail, email, role)
		return nil, fmt.Errorf("Role '%s' tidak valid. Role yang diizinkan: OPERASIONAL, COMPLIANCE, MANAJEMEN, ADMIN_CS", role)
	}

	// 2. Cek email unik
	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		log.Printf("[ADMIN] action=CREATE_USER admin=%s target=%s role=%s result=FAILED reason=db_error", adminEmail, email, role)
		return nil, fmt.Errorf("Gagal mengecek email: %v", err)
	}
	if exists {
		log.Printf("[ADMIN] action=CREATE_USER admin=%s target=%s role=%s result=FAILED reason=email_exists", adminEmail, email, role)
		return nil, fmt.Errorf("Email sudah terdaftar")
	}

	// 3. Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("[ADMIN] action=CREATE_USER admin=%s target=%s role=%s result=FAILED reason=hash_error", adminEmail, email, role)
		return nil, fmt.Errorf("Gagal mengenkripsi password")
	}

	// 4. Cari role ID di database
	roleID, err := s.repo.FindRoleByCode(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("Gagal mencari role: %v", err)
	}
	if roleID == "" {
		log.Printf("[ADMIN] action=CREATE_USER admin=%s target=%s role=%s result=FAILED reason=role_not_found", adminEmail, email, role)
		return nil, fmt.Errorf("Role '%s' tidak ditemukan di database", role)
	}

	// 5. Ambil default group dan level
	groupID, levelID, err := s.repo.FindDefaultGroupAndLevel(ctx)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil konfigurasi default: %v", err)
	}

	// 6. Generate UUID
	newUUID, err := utils.CreateUUID()
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat ID user")
	}
	newUserID := newUUID.String()

	// 7. Simpan user ke database (status: active)
	err = s.repo.CreateUser(ctx, newUserID, fullName, email, hashedPassword, roleID, groupID, levelID)
	if err != nil {
		log.Printf("[ADMIN] action=CREATE_USER admin=%s target=%s role=%s result=FAILED reason=insert_error err=%v", adminEmail, email, role, err)
		return nil, fmt.Errorf("Gagal menyimpan data user: %v", err)
	}

	// 8. Audit log
	log.Printf("[ADMIN] action=CREATE_USER admin=%s target=%s role=%s result=SUCCESS user_id=%s", adminEmail, email, role, newUserID)

	// 9. Dispatch USER_CREATED event ke Notification Service (async)
	plainPassword := password // simpan sebelum di-hash untuk dikirim via email
	go dispatchNotification("USER_CREATED", "email", email, map[string]string{
		"name":     fullName,
		"email":    email,
		"password": plainPassword,
		"role":     role,
	})

	return CreateUserResult{
		ID:        newUserID,
		FullName:  fullName,
		Email:     email,
		Role:      role,
		Status:    "active",
		CreatedAt: time.Now(),
	}, nil
}

// ProcessGetUsersJob — dipanggil oleh worker untuk mengambil daftar akun internal
func (s *Service) ProcessGetUsersJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(GetUsersParams)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	// Validate / Default pagination
	if data.Page < 1 {
		data.Page = 1
	}
	if data.PerPage < 1 {
		data.PerPage = 20
	}

	// Fetch data
	users, err := s.repo.GetInternalUsers(ctx, data)
	if err != nil {
		log.Printf("[ADMIN] GetUsers failed: %v", err)
		return nil, fmt.Errorf("Gagal mengambil data user")
	}

	// Fetch count
	total, err := s.repo.CountInternalUsers(ctx, data)
	if err != nil {
		log.Printf("[ADMIN] CountUsers failed: %v", err)
		return nil, fmt.Errorf("Gagal menghitung jumlah user")
	}

	if users == nil {
		users = []UserListItem{}
	}

	return GetUsersResponse{
		Data:    users,
		Total:   total,
		Page:    data.Page,
		PerPage: data.PerPage,
	}, nil
}

