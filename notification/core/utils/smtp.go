package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

// ============================================================
// EmailSender — abstraction over Brevo HTTP API and SMTP
// ============================================================

// EmailSender is the interface used by the notification service.
type EmailSender interface {
	SendEmail(to, subject, body string) error
	SendHTMLEmail(to, subject, htmlBody string) error
}

// NewEmailSender returns Brevo-based sender if BREVO_API_KEY is set,
// otherwise falls back to SMTP.
func NewEmailSender() EmailSender {
	// Priority: Brevo > SMTP
	if key := getEnv("BREVO_API_KEY"); key != "" {
		from := getEnv("BREVO_FROM_EMAIL", getEnv("SMTP_FROM", "noreply@thinktala.com"))
		fromName := getEnv("BREVO_FROM_NAME", "ThinkTala")
		log.Printf("[EMAIL] Using Brevo HTTP API (from: %s <%s>)", fromName, from)
		return &BrevoClient{APIKey: key, FromEmail: from, FromName: fromName}
	}
	log.Println("[EMAIL] Using SMTP (fallback)")
	return NewSMTPClient()
}

// ============================================================
// Brevo (ex-Sendinblue) HTTP API Client — 300 emails/day free
// No domain required, just verify sender email
// ============================================================

type BrevoClient struct {
	APIKey    string
	FromEmail string
	FromName  string
}

type brevoPayload struct {
	Sender      brevoContact   `json:"sender"`
	To          []brevoContact `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent,omitempty"`
	TextContent string         `json:"textContent,omitempty"`
}

type brevoContact struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

func (b *BrevoClient) SendEmail(to, subject, body string) error {
	return b.send(brevoPayload{
		Sender:      brevoContact{Email: b.FromEmail, Name: b.FromName},
		To:          []brevoContact{{Email: to}},
		Subject:     subject,
		TextContent: body,
	})
}

func (b *BrevoClient) SendHTMLEmail(to, subject, htmlBody string) error {
	return b.send(brevoPayload{
		Sender:      brevoContact{Email: b.FromEmail, Name: b.FromName},
		To:          []brevoContact{{Email: to}},
		Subject:     subject,
		HTMLContent: htmlBody,
	})
}

func (b *BrevoClient) send(payload brevoPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("brevo marshal: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("brevo request: %w", err)
	}
	req.Header.Set("api-key", b.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("brevo send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("brevo error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// ============================================================
// SMTP Client (legacy, for local dev)
// ============================================================

// SMTPClient holds SMTP configuration.
type SMTPClient struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// NewSMTPClient creates a new SMTPClient from environment variables.
func NewSMTPClient() *SMTPClient {
	return &SMTPClient{
		Host:     getEnv("SMTP_HOST"),
		Port:     getEnv("SMTP_PORT", "587"),
		Username: getEnv("SMTP_USER"),
		Password: getEnv("SMTP_PASS"),
		From:     getEnv("SMTP_FROM"),
	}
}

// SendEmail sends a plain text email.
func (s *SMTPClient) SendEmail(to, subject, body string) error {
	if s.Host == "" {
		return fmt.Errorf("SMTP configuration missing")
	}

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s\r\n",
		to, subject, body,
	))

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}

// SendHTMLEmail sends an HTML email.
func (s *SMTPClient) SendHTMLEmail(to, subject, htmlBody string) error {
	if s.Host == "" {
		return fmt.Errorf("SMTP configuration missing")
	}

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s\r\n",
		to, subject, htmlBody,
	))

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}

func getEnv(key string, fallback ...string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}
