package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// EmailSender abstracts over Brevo HTTP API.
type EmailSender interface {
	SendEmail(to, subject, body string) error
	SendHTMLEmail(to, subject, htmlBody string) error
}

// NewEmailSender returns a Brevo sender configured from env, or a no-op sender
// that logs when BREVO_API_KEY is not set.
func NewEmailSender() EmailSender {
	if key := GetEnv("BREVO_API_KEY"); key != "" {
		from := GetEnv("BREVO_FROM_EMAIL", "noreply@thinktala.com")
		fromName := GetEnv("BREVO_FROM_NAME", "ThinkTala")
		log.Printf("[EMAIL] Using Brevo HTTP API (from: %s <%s>)", fromName, from)
		return &BrevoClient{APIKey: key, FromEmail: from, FromName: fromName}
	}
	log.Println("[EMAIL] BREVO_API_KEY missing — email delivery disabled")
	return &missingBrevoSender{}
}

type missingBrevoSender struct{}

func (m *missingBrevoSender) SendEmail(_, _, _ string) error {
	return fmt.Errorf("email sender not configured: BREVO_API_KEY is required")
}
func (m *missingBrevoSender) SendHTMLEmail(_, _, _ string) error {
	return fmt.Errorf("email sender not configured: BREVO_API_KEY is required")
}

// BrevoClient sends email via Brevo (ex-Sendinblue) HTTP API.
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
