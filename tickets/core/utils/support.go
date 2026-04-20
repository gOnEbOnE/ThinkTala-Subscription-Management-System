package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func NormalizeTicketCategory(raw string) (string, bool) {
	cleaned := strings.ToUpper(strings.TrimSpace(raw))
	cleaned = strings.ReplaceAll(cleaned, "-", "_")
	cleaned = strings.ReplaceAll(cleaned, " ", "_")

	switch cleaned {
	case "MASALAH_TEKNIS", "TEKNIS", "TEKNIK", "TECHNICAL", "MASALAHTECHNIS":
		return "MASALAH_TEKNIS", true
	case "PEMBAYARAN", "PAYMENT":
		return "PEMBAYARAN", true
	case "AKUN", "ACCOUNT":
		return "AKUN", true
	case "LANGGANAN", "SUBSCRIPTION", "PAKET", "MEMBERSHIP":
		return "LANGGANAN", true
	case "NOTIFIKASI", "NOTIFICATION", "NOTIF":
		return "NOTIFIKASI", true
	case "LAINNYA", "OTHER", "GENERAL", "UMUM":
		return "LAINNYA", true
	default:
		return "", false
	}
}

func NewUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	hexID := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexID[0:8],
		hexID[8:12],
		hexID[12:16],
		hexID[16:20],
		hexID[20:32]), nil
}

func ParsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func NormalizeSupportTicketStatus(raw string) (string, bool) {
	cleaned := strings.ToUpper(strings.TrimSpace(raw))
	cleaned = strings.ReplaceAll(cleaned, "_", " ")

	switch cleaned {
	case "ON PROCESS", "DONE":
		return cleaned, true
	default:
		return "", false
	}
}

func ParseAdminSupportTicketPath(path string) (ticketID, action string, ok bool) {
	const prefix = "/api/admin/support/tickets/"
	if !strings.HasPrefix(path, prefix) {
		return "", "", false
	}

	suffix := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if suffix == "" {
		return "", "", false
	}

	parts := strings.Split(suffix, "/")
	if len(parts) == 1 {
		id := strings.TrimSpace(parts[0])
		return id, "", id != ""
	}

	if len(parts) == 2 && strings.EqualFold(strings.TrimSpace(parts[1]), "replies") {
		id := strings.TrimSpace(parts[0])
		return id, "replies", id != ""
	}

	return "", "", false
}

func IsAllowedImageMIME(contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/jpg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

func SanitizeUploadFileName(fileName string) string {
	baseName := filepath.Base(strings.TrimSpace(fileName))
	if baseName == "" || baseName == "." || baseName == string(filepath.Separator) {
		return "attachment"
	}
	if len(baseName) > 160 {
		baseName = baseName[:160]
	}
	return strings.ReplaceAll(baseName, "\x00", "")
}

func BuildAdminAttachmentURL(ticketID string, hasAttachment bool) string {
	if !hasAttachment || strings.TrimSpace(ticketID) == "" {
		return ""
	}
	return fmt.Sprintf("/api/admin/support/tickets/attachment?ticket_id=%s", ticketID)
}
