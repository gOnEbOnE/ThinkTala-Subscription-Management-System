package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

// NewEmailSender returns a Brevo sender.
// SMTP fallback is intentionally removed to enforce a single delivery channel.
func NewEmailSender() EmailSender {
	if key := getEnv("BREVO_API_KEY"); key != "" {
		from := getEnv("BREVO_FROM_EMAIL", "noreply@thinktala.com")
		fromName := getEnv("BREVO_FROM_NAME", "ThinkTala")
		log.Printf("[EMAIL] Using Brevo HTTP API (from: %s <%s>)", fromName, from)
		return &BrevoClient{APIKey: key, FromEmail: from, FromName: fromName}
	}
	log.Println("[EMAIL] BREVO_API_KEY is missing - Brevo sender disabled")
	return &MissingBrevoSender{}
}

type MissingBrevoSender struct{}

func (m *MissingBrevoSender) SendEmail(_, _, _ string) error {
	return fmt.Errorf("brevo sender not configured: BREVO_API_KEY is required")
}

func (m *MissingBrevoSender) SendHTMLEmail(_, _, _ string) error {
	return fmt.Errorf("brevo sender not configured: BREVO_API_KEY is required")
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

func getEnv(key string, fallback ...string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}
