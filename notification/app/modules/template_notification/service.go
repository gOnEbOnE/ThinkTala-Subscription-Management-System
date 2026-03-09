package template

import "errors"

// Service berisi business logic untuk notification templates.
// Validasi bisnis seperti subject wajib untuk email ada di sini.
type Service struct {
	repo *Repository
}

// NewService membuat instance Service baru.
func NewService() *Service {
	return &Service{repo: NewRepository()}
}

// List mengambil semua template, bisa difilter by channel dan event_type.
func (s *Service) List(channel, eventType string) ([]NotificationTemplate, error) {
	return s.repo.List(channel, eventType)
}

// GetByID mengambil satu template berdasarkan ID.
func (s *Service) GetByID(id string) (NotificationTemplate, error) {
	return s.repo.GetByID(id)
}

// Create membuat template baru dengan validasi bisnis.
func (s *Service) Create(req CreateTemplateRequest, id string) error {
	if req.Channel == "email" && req.Subject == "" {
		return errors.New("subject wajib diisi untuk channel email")
	}
	return s.repo.Create(req, id)
}

// Update memperbarui template dengan validasi bisnis.
func (s *Service) Update(id string, req UpdateTemplateRequest) error {
	if !s.repo.Exists(id) {
		return ErrNotFound
	}
	if req.Channel == "email" && req.Subject == "" {
		return errors.New("subject wajib diisi untuk channel email")
	}
	return s.repo.Update(id, req)
}

// Delete menghapus template berdasarkan ID.
func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

// ErrNotFound digunakan ketika template dengan ID yang diminta tidak ditemukan.
var ErrNotFound = errors.New("template tidak ditemukan")
