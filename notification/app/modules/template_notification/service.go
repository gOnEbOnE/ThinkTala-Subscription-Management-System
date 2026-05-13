package template

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"notification/core/utils"

	"github.com/google/uuid"
)

// Service berisi business logic untuk notification templates.
// Validasi bisnis seperti subject wajib untuk email ada di sini.
type Service struct {
	repo *Repository
}

// NewService membuat instance Service baru.
func NewService() *Service {
	return &Service{repo: NewRepository()}
}

// List mengambil semua template, bisa difilter by channel dan event_type.
func (s *Service) List(channel, eventType string) ([]NotificationTemplate, error) {
	return s.repo.List(channel, eventType)
}

// GetByID mengambil satu template berdasarkan ID.
func (s *Service) GetByID(id string) (NotificationTemplate, error) {
	return s.repo.GetByID(id)
}

// Create membuat template baru dengan validasi bisnis.
func (s *Service) Create(req CreateTemplateRequest, id string) error {
	if req.Channel == "email" && req.Subject == "" {
		return errors.New("subject wajib diisi untuk channel email")
	}
	return s.repo.Create(req, id)
}

// Update memperbarui template dengan validasi bisnis.
func (s *Service) Update(id string, req UpdateTemplateRequest) error {
	if !s.repo.Exists(id) {
		return ErrNotFound
	}
	if req.Channel == "email" && req.Subject == "" {
		return errors.New("subject wajib diisi untuk channel email")
	}
	return s.repo.Update(id, req)
}

// Delete menghapus template berdasarkan ID.
func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

// Send mencari template berdasarkan event_type+channel, render placeholder, simpan log, lalu kirim.
func (s *Service) Send(req SendRequest) error {
	log.Printf("[NOTIF SEND] Received event_type=%s channel=%s to=%s", req.EventType, req.Channel, req.To)

	// Catat event_type ke registry agar muncul di dropdown
	s.repo.RegisterEventType(req.EventType)

	tpl, err := s.repo.GetByEventType(req.EventType, req.Channel)
	if err != nil {
		log.Printf("[NOTIF SEND] Template NOT found for event_type=%s channel=%s: %v", req.EventType, req.Channel, err)
		return fmt.Errorf("template tidak ditemukan untuk event_type=%s channel=%s", req.EventType, req.Channel)
	}
	log.Printf("[NOTIF SEND] Template found for event_type=%s, sending to %s", req.EventType, req.To)

	subject := ""
	if tpl.Subject != nil {
		subject = *tpl.Subject
	}
	content := tpl.Content

	for k, v := range req.Vars {
		placeholder := "{{" + k + "}}"
		subject = strings.ReplaceAll(subject, placeholder, v)
		content = strings.ReplaceAll(content, placeholder, v)
	}

	// Simpan log sebelum kirim (status pending → akan diupdate setelah send)
	logID := uuid.New().String()
	s.repo.SaveLog(logID, req.EventType, req.Channel, req.To, subject, content)

	sendErr := s.doSend(req.Channel, req.To, subject, content)

	if sendErr == nil {
		s.repo.MarkLogSent(logID)
		log.Printf("[NOTIF] Terkirim ke %s (event: %s)", req.To, req.EventType)
	} else {
		nextRetry := time.Now().Add(backoffDuration(0))
		s.repo.MarkLogFailed(logID, sendErr.Error(), 0, &nextRetry)
		log.Printf("[NOTIF] Gagal kirim ke %s: %v — dijadwalkan retry", req.To, sendErr)
	}

	return sendErr
}

// doSend melakukan pengiriman aktual tanpa menyentuh log.
func (s *Service) doSend(channel, to, subject, content string) error {
	switch channel {
	case "email":
		sender := utils.NewEmailSender()
		trimmed := strings.TrimSpace(content)
		isHTML := len(trimmed) > 0 && trimmed[0] == '<'
		if isHTML {
			return sender.SendHTMLEmail(to, subject, content)
		}
		return sender.SendEmail(to, subject, content)
	default:
		return fmt.Errorf("channel '%s' belum didukung", channel)
	}
}

// backoffDuration mengembalikan durasi tunggu sebelum retry ke-n.
// Skema: retry 0→1m, 1→5m, 2→30m
func backoffDuration(retryCount int) time.Duration {
	switch retryCount {
	case 0:
		return 1 * time.Minute
	case 1:
		return 5 * time.Minute
	case 2:
		return 30 * time.Minute
	default:
		return 1 * time.Hour
	}
}

// StartRetryWorker menjalankan loop retry di background.
// Dipanggil sekali saat startup dari main.go.
func (s *Service) StartRetryWorker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	log.Println("[RETRY WORKER] Started — interval: 30s")
	for {
		select {
		case <-ctx.Done():
			log.Println("[RETRY WORKER] Stopped")
			return
		case <-ticker.C:
			s.processRetries()
		}
	}
}

// processRetries memproses semua log dengan status failed yang sudah waktunya di-retry.
func (s *Service) processRetries() {
	retryable, err := s.repo.GetRetryableLogs()
	if err != nil || len(retryable) == 0 {
		return
	}
	log.Printf("[RETRY WORKER] Processing %d retryable notification(s)", len(retryable))
	for _, l := range retryable {
		sendErr := s.doSend(l.Channel, l.ToAddress, l.Subject, l.Content)
		if sendErr == nil {
			s.repo.MarkLogSent(l.ID)
			log.Printf("[RETRY WORKER] Berhasil retry ke %s (event: %s, attempt: %d)", l.ToAddress, l.EventType, l.RetryCount+1)
		} else {
			newCount := l.RetryCount + 1
			if newCount >= l.MaxRetries {
				// Tandai permanently failed (retry_count = max_retries)
				s.repo.MarkLogFailed(l.ID, sendErr.Error(), newCount, nil)
				log.Printf("[RETRY WORKER] Permanently failed ke %s setelah %d percobaan", l.ToAddress, newCount)
			} else {
				nextRetry := time.Now().Add(backoffDuration(newCount))
				s.repo.MarkLogFailed(l.ID, sendErr.Error(), newCount, &nextRetry)
				log.Printf("[RETRY WORKER] Retry %d/%d gagal ke %s, jadwal ulang %v", newCount, l.MaxRetries, l.ToAddress, nextRetry.Format(time.RFC3339))
			}
		}
	}
}

// GetLogs mengambil log untuk monitoring.
func (s *Service) GetLogs(status string, limit, offset int) ([]NotificationLog, int, error) {
	return s.repo.ListLogs(status, limit, offset)
}

// ListEventTypes mengembalikan daftar event_type yang sudah terdaftar.
func (s *Service) ListEventTypes() ([]string, error) {
	return s.repo.ListEventTypes()
}

// ErrNotFound digunakan ketika template dengan ID yang diminta tidak ditemukan.
var ErrNotFound = errors.New("template tidak ditemukan")
