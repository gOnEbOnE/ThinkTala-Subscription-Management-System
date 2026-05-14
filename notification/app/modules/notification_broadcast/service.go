package notification

import (
	"fmt"
	"sort"
	"strings"
)

// Service berisi business logic untuk broadcast notifications.
// Service hanya bicara dengan Repository, tidak boleh langsung ke database.
type Service struct {
	repo *Repository
}

const maxPinnedNotifications = 2

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

// Recent mengembalikan ringkasan notifikasi terbaru (news + event) untuk drawer.
func (s *Service) Recent(role, userID string, limit int) ([]RecentNotification, error) {
	if limit <= 0 || limit > 5 {
		limit = 5
	}
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		role = "client"
	}

	items := make([]RecentNotification, 0)

	events, err := s.repo.ListRecentEvents(userID, limit)
	if err != nil {
		// Log tapi jangan propagate — drawer tetap bisa tampil dengan news saja
		fmt.Printf("[RECENT] ListRecentEvents error (non-fatal): %v\n", err)
	} else {
		items = append(items, events...)
	}

	news, err := s.repo.ListRecentNews(role, userID, limit)
	if err != nil {
		// Log tapi jangan propagate — drawer tetap bisa tampil dengan events saja
		fmt.Printf("[RECENT] ListRecentNews error (non-fatal): %v\n", err)
	} else {
		items = append(items, news...)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

// MarkRead menandai satu item sebagai read.
func (s *Service) MarkRead(userID, sourceType, sourceID string) error {
	return s.repo.MarkRead(userID, sourceType, sourceID)
}

// MarkAllRead menandai semua item sebagai read untuk user tertentu.
func (s *Service) MarkAllRead(role, userID string) error {
	return s.repo.MarkAllRead(userID, role)
}

// GetByID mengambil satu notification berdasarkan ID.
func (s *Service) GetByID(id string) (Notification, error) {
	return s.repo.GetByID(id)
}

// GetPublicByID mengambil notification untuk client (aktif) dan menambah view_count.
func (s *Service) GetPublicByID(id, role, userID string) (Notification, error) {
	return s.repo.GetPublicByID(id, role, userID)
}

// Create membuat notification baru.
func (s *Service) Create(req CreateNotificationRequest, id string) error {
	title := strings.TrimSpace(req.Title)
	desc := strings.TrimSpace(req.Description)
	if desc == "" {
		desc = strings.TrimSpace(req.Message)
	}
	typeVal := strings.ToLower(strings.TrimSpace(req.Type))
	targetRole := strings.ToLower(strings.TrimSpace(req.TargetRole))

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
	if targetRole == "" {
		targetRole = "client"
	}
	if !isValidTargetRole(targetRole) {
		return fmt.Errorf("target_role tidak valid")
	}

	req.Title = title
	req.Description = desc
	req.Message = desc
	req.Type = typeVal
	req.TargetRole = targetRole

	if req.IsPinned != nil && *req.IsPinned {
		pinnedCount, err := s.repo.countPinned("")
		if err != nil {
			return fmt.Errorf("gagal mengecek slot pinned")
		}
		if pinnedCount >= maxPinnedNotifications {
			return fmt.Errorf("slot pinned penuh (maksimal %d). Nonaktifkan atau hapus salah satu news pinned terlebih dahulu", maxPinnedNotifications)
		}
	}

	return s.repo.Create(req, id)
}

// Update memperbarui notification berdasarkan ID.
func (s *Service) Update(id string, req UpdateNotificationRequest) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if req.Type != "" && !isValidType(req.Type) {
		return fmt.Errorf("type tidak valid")
	}
	if targetRole := strings.ToLower(strings.TrimSpace(req.TargetRole)); targetRole != "" {
		if !isValidTargetRole(targetRole) {
			return fmt.Errorf("target_role tidak valid")
		}
		req.TargetRole = targetRole
	}

	if strings.TrimSpace(req.Description) == "" && strings.TrimSpace(req.Message) != "" {
		req.Description = strings.TrimSpace(req.Message)
	}
	if req.Type != "" {
		req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	}

	if req.IsPinned != nil && *req.IsPinned && !existing.IsPinned {
		pinnedCount, err := s.repo.countPinned(id)
		if err != nil {
			return fmt.Errorf("gagal mengecek slot pinned")
		}
		if pinnedCount >= maxPinnedNotifications {
			return fmt.Errorf("slot pinned penuh (maksimal %d). Nonaktifkan atau hapus salah satu news pinned terlebih dahulu", maxPinnedNotifications)
		}
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

func isValidTargetRole(targetRole string) bool {
	switch strings.ToLower(strings.TrimSpace(targetRole)) {
	case "client", "client_never_bought", "client_paid_active", "client_lapsed", "client_expiring_soon":
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
