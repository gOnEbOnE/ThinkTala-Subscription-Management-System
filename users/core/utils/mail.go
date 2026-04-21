package utils

import (
	"fmt"
	"net/smtp"
)

// SMTPClient struct
type SMTPClient struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// NewSMTPClient creates new instance from Env
func NewSMTPClient() *SMTPClient {
	return &SMTPClient{
		Host:     GetEnv("SMTP_HOST"),
		Port:     GetEnv("SMTP_PORT", "587"),
		Username: GetEnv("SMTP_USER"),
		Password: GetEnv("SMTP_PASS"),
		From:     GetEnv("SMTP_FROM"),
	}
}

// SendEmail sends plain text email
func (s *SMTPClient) SendEmail(to, subject, body string) error {
	if s.Host == "" {
		return fmt.Errorf("SMTP configuration missing")
	}

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}
