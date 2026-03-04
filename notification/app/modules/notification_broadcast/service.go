package notification

// Service berisi business logic untuk broadcast notifications.
// Service hanya bicara dengan Repository, tidak boleh langsung ke database.
type Service struct {
	repo *Repository
}

// NewService membuat instance Service baru.
func NewService() *Service {
	return &Service{repo: NewRepository()}
}

// List mengambil semua notification.
func (s *Service) List() ([]Notification, error) {
	return s.repo.List()
}

// ListPublic mengambil notification aktif untuk role tertentu.
func (s *Service) ListPublic(role string) ([]map[string]any, error) {
	return s.repo.ListPublic(role)
}

// GetByID mengambil satu notification berdasarkan ID.
func (s *Service) GetByID(id string) (Notification, error) {
	return s.repo.GetByID(id)
}

// Create membuat notification baru.
func (s *Service) Create(req CreateNotificationRequest, id string) error {
	return s.repo.Create(req, id)
}

// Update memperbarui notification berdasarkan ID.
func (s *Service) Update(id string, req UpdateNotificationRequest) error {
	return s.repo.Update(id, req)
}

// Delete menghapus notification berdasarkan ID.
func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
