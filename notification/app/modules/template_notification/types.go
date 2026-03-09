package template

import "time"

// NotificationTemplate adalah model untuk template pesan notifikasi
// yang di-mapping ke event_type tertentu (misal: otp_verification, user_kyc_approved).
type NotificationTemplate struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	EventType string    `json:"event_type"` // otp_verification | user_register | user_kyc_approved | ...
	Channel   string    `json:"channel"`    // email | telegram
	Subject   *string   `json:"subject,omitempty"`
	Content   string    `json:"content"` // Mendukung placeholder: {{name}}, {{otp}}, dst.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *string   `json:"created_by,omitempty"`
	UpdatedBy *string   `json:"updated_by,omitempty"`
}

// NotificationLog adalah record setiap percobaan pengiriman notifikasi.
type NotificationLog struct {
	ID          string     `json:"id"`
	EventType   string     `json:"event_type"`
	Channel     string     `json:"channel"`
	ToAddress   string     `json:"to"`
	Subject     string     `json:"subject,omitempty"`
	Content     string     `json:"content,omitempty"`
	Status      string     `json:"status"` // pending | sent | failed
	RetryCount  int        `json:"retry_count"`
	MaxRetries  int        `json:"max_retries"`
	NextRetryAt *time.Time `json:"next_retry_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	ErrorMsg    *string    `json:"error_msg,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CreateTemplateRequest adalah payload untuk membuat template baru.
type CreateTemplateRequest struct {
	Name      string `json:"name"       binding:"required"`
	EventType string `json:"event_type" binding:"required"`
	Channel   string `json:"channel"    binding:"required"`
	Subject   string `json:"subject"`
	Content   string `json:"content"    binding:"required"`
	CreatedBy string `json:"created_by"`
}

// SendRequest adalah payload untuk mengirim notifikasi berdasarkan event_type.
// Notification Service akan mencari template, render placeholder, lalu kirim via channel.
type SendRequest struct {
	EventType string            `json:"event_type" binding:"required"` // otp_verification, user_kyc_approved, ...
	Channel   string            `json:"channel"    binding:"required"` // email | telegram
	To        string            `json:"to"         binding:"required"` // alamat tujuan (email address, dll.)
	Vars      map[string]string `json:"vars"`                          // placeholder values: {"name": "Rifqi", "otp": "123456"}
}

// UpdateTemplateRequest adalah payload untuk update template.
// Validasi wajib sama dengan CreateTemplateRequest (PBI-30).
type UpdateTemplateRequest struct {
	Name      string `json:"name"       binding:"required"`
	EventType string `json:"event_type" binding:"required"`
	Channel   string `json:"channel"    binding:"required"`
	Subject   string `json:"subject"`
	Content   string `json:"content"    binding:"required"`
	UpdatedBy string `json:"updated_by"`
}
