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
// EmailSender — abstraction over SMTP and Resend HTTP API
// ============================================================

// EmailSender is the interface used by the notification service.
type EmailSender interface {
	SendEmail(to, subject, body string) error
	SendHTMLEmail(to, subject, htmlBody string) error
}

// NewEmailSender returns Resend-based sender if RESEND_API_KEY is set,
// otherwise falls back to SMTP.
func NewEmailSender() EmailSender {
	if key := getEnv("RESEND_API_KEY"); key != "" {
		from := getEnv("RESEND_FROM", getEnv("SMTP_FROM", "noreply@thinktala.com"))
		log.Printf("[EMAIL] Using Resend HTTP API (from: %s)", from)
		return &ResendClient{APIKey: key, From: from}
	}
	log.Println("[EMAIL] Using SMTP (fallback)")
	return NewSMTPClient()
}

// ============================================================
// Resend HTTP API Client
// ============================================================

type ResendClient struct {
	APIKey string
	From   string
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

func (r *ResendClient) SendEmail(to, subject, body string) error {
	return r.send(resendPayload{
		From: r.From, To: []string{to}, Subject: subject, Text: body,
	})
}

func (r *ResendClient) SendHTMLEmail(to, subject, htmlBody string) error {
	return r.send(resendPayload{
		From: r.From, To: []string{to}, Subject: subject, HTML: htmlBody,
	})
}

func (r *ResendClient) send(payload resendPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("resend marshal: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("resend send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend error %d: %s", resp.StatusCode, string(respBody))
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
