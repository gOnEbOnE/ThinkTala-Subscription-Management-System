package notification

import (
	"fmt"
	"strings"
)

// Service berisi business logic untuk broadcast notifications.
// Service hanya bicara dengan Repository, tidak boleh langsung ke database.
type Service struct {
	repo *Repository
}

// NewService membuat instance Service baru.
func NewService() *Service {
	return &Service{repo: NewRepository()}
}

// List mengambil semua notification dengan filter opsional.
func (s *Service) List(typeFilter, statusFilter string) ([]Notification, error) {
	if typeFilter != "" && typeFilter != "all" && !isValidType(typeFilter) {
		return nil, fmt.Errorf("type tidak valid")
	}
	if statusFilter != "" && statusFilter != "all" && !isValidStatusFilter(statusFilter) {
		return nil, fmt.Errorf("status filter tidak valid")
	}
	return s.repo.List(strings.ToLower(typeFilter), strings.ToLower(statusFilter))
}

// ListPublic mengambil notification aktif untuk role tertentu.
func (s *Service) ListPublic(role, userID string) ([]map[string]any, error) {
	return s.repo.ListPublic(role, userID)
}

// GetByID mengambil satu notification berdasarkan ID.
func (s *Service) GetByID(id string) (Notification, error) {
	return s.repo.GetByID(id)
}

// Create membuat notification baru.
func (s *Service) Create(req CreateNotificationRequest, id string) error {
	title := strings.TrimSpace(req.Title)
	desc := strings.TrimSpace(req.Description)
	if desc == "" {
		desc = strings.TrimSpace(req.Message)
	}
	typeVal := strings.ToLower(strings.TrimSpace(req.Type))

	if title == "" {
		return fmt.Errorf("title wajib diisi")
	}
	if desc == "" {
		return fmt.Errorf("description wajib diisi")
	}
	if typeVal == "" {
		return fmt.Errorf("type wajib diisi")
	}
	if !isValidType(typeVal) {
		return fmt.Errorf("type tidak valid")
	}

	req.Title = title
	req.Description = desc
	req.Message = desc
	req.Type = typeVal

	return s.repo.Create(req, id)
}

// Update memperbarui notification berdasarkan ID.
func (s *Service) Update(id string, req UpdateNotificationRequest) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return err
	}

	if req.Type != "" && !isValidType(req.Type) {
		return fmt.Errorf("type tidak valid")
	}

	if strings.TrimSpace(req.Description) == "" && strings.TrimSpace(req.Message) != "" {
		req.Description = strings.TrimSpace(req.Message)
	}
	if req.Type != "" {
		req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	}

	return s.repo.Update(id, req)
}

// Delete menghapus notification berdasarkan ID.
func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func isValidType(typeVal string) bool {
	switch strings.ToLower(strings.TrimSpace(typeVal)) {
	case "system", "promo", "warning", "info", "analysis", "education", "event":
		return true
	default:
		return false
	}
}

func isValidStatusFilter(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active", "inactive", "expired", "true", "false", "1", "0":
		return true
	default:
		return false
	}
}
