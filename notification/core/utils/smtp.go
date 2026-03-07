package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

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
