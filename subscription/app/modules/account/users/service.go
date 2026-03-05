package users

import (
	"context"
	"fmt"
	"math"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetInactiveUsersService(ctx context.Context, payload any) (any, error) {
	// 1. Validasi Input Dasar
	// Hindari page 0 atau negatif, dan limit aneh
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	page, _ := data["page"].(int)
	limit, _ := data["limit"].(int)
	search, _ := data["search"].(string)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 12 // Default sesuai dropdown frontend
	}
	// Batasi limit agar tidak memberatkan server (opsional)
	if limit > 100 {
		limit = 100
	}

	// 2. Panggil Repository
	// Kita minta Repo mengambil data + total count
	users, total, err := s.repo.GetUsersInactive(ctx, page, limit, search)
	if err != nil {
		return nil, err
	}

	// 3. Logika Bisnis: Hitung Pagination Metadata
	// Rumus: Total Halaman = Ceil(Total Data / Limit)
	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	// 4. Build Response Data
	// Kita susun struktur map persis seperti yang diharapkan script JS Frontend (users.html)
	// Format: { data: [...], meta: { ... } }
	result := map[string]any{
		"data": users, // Slice UserListDTO dari repo
		"meta": map[string]any{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}

	return result, nil
}

func (s *Service) GetDetailUserService(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("ID user tidak boleh kosong")
	}

	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// 1. Fungsi Utama
func (s *Service) UpdateUserService(ctx context.Context, id string, data UserUpdateDTO) error {
	if id == "" {
		return fmt.Errorf("ID tidak boleh kosong")
	}
	// Validasi lain bisa ditaruh disini (misal NIK harus angka)

	return s.repo.UpdateUser(ctx, id, data)
}

// 2. Wrapper untuk Worker
func (s *Service) ProcessUpdateUserJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	id, _ := data["id"].(string)
	fullname, _ := data["fullname"].(string)
	nik, _ := data["nik"].(string)
	nrk, _ := data["nrk"].(string)

	// Handle IsActive (bisa datang sebagai bool atau string "1"/"0")
	isActive := false
	if v, ok := data["is_active"].(bool); ok {
		isActive = v
	} else if v, ok := data["is_active"].(string); ok {
		isActive = (v == "1" || v == "true")
	}

	dto := UserUpdateDTO{
		Fullname: fullname,
		NIK:      nik,
		NRK:      nrk,
		IsActive: isActive,
	}

	err := s.UpdateUserService(ctx, id, dto)
	return nil, err // Return nil jika sukses
}

// 2. Wrapper untuk Worker/Dispatcher
func (s *Service) ProcessGetDetailUserJob(ctx context.Context, payload any) (any, error) {
	// Parse Payload
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}

	// Ambil ID
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("missing or invalid 'id' in payload")
	}

	// Panggil Logic Utama
	return s.GetDetailUserService(ctx, id)
}
