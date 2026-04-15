package login

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/master-abror/zaframework/core/utils"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ProcessLoginJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	loginID := data["login_id"].(string)
	password := data["password"].(string)

	ip, _ := data["ip"].(string)
	browser, _ := data["browser"].(string)
	os, _ := data["os"].(string)
	lat, _ := data["latitude"].(string)
	long, _ := data["longitude"].(string)

	user, err := s.repo.FindUser(ctx, "email", loginID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("Email atau Password tidak cocok")
	}

	// Cek status user harus 'active'
	if strings.ToLower(user.Status) != "active" {
		return nil, fmt.Errorf("Akun belum aktif. Silakan verifikasi email terlebih dahulu.")
	}

	if hash := utils.CheckPasswordHash(password, user.Password); hash != true {
		return nil, fmt.Errorf("Email atau Password tidak cocok")
	}

	userData := map[string]any{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"group_id":   user.GroupID,
		"level_id":   user.LevelID,
		"role_id":    user.RoleID,
		"group":      user.Group,
		"group_code": user.GroupCode,
		"level":      user.Level,
		"level_code": user.LevelCode,
		"role":       user.Role,
		"role_code":  user.RoleCode,
		"status":     user.Status,
		"photo":      user.Photo,
		"ip":         ip,
		"latitude":   lat,
		"longitude":  long,
		"browser":    browser,
		"os":         os,
	}

	// Pastikan role_code tidak kosong dan dalam format uppercase
	roleCode := strings.ToUpper(strings.TrimSpace(userData["role_code"].(string)))
	userData["role_code"] = roleCode

	token, err := utils.CreateJWT(userData, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat token")
	}

	// ✅ JANGAN ENCRYPT TOKEN SIGNATURE! Simpan RAW JWT
	// DIHAPUS: token = utils.EncryptTokenSignature(token)
	userData["token"] = token

	return LoginResult{
		Token:    token,
		UserData: userData,
	}, nil
}

// ProcessAssumeRoleJob — dipanggil oleh worker untuk proses assume role (SUPERADMIN only)
func (s *Service) ProcessAssumeRoleJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	userEmail := data["user_email"].(string)
	targetRoleCode := strings.ToUpper(data["target_role_code"].(string))

	// 1. Cari user di database
	user, err := s.repo.FindUser(ctx, "email", userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("User tidak ditemukan")
	}

	// 2. Pastikan user adalah SUPERADMIN
	if strings.ToUpper(user.LevelCode) != "SUPERADMIN" {
		return nil, fmt.Errorf("Hanya Super Admin yang dapat melakukan simulasi role")
	}

	// 3. Validasi target role
	validRoles := map[string]bool{
		"CEO":         true,
		"OPERASIONAL": true,
		"COMPLIANCE":  true,
		"CLIENT":      true,
	}
	if !validRoles[targetRoleCode] {
		return nil, fmt.Errorf("Role '%s' tidak valid untuk simulasi", targetRoleCode)
	}

	// 4. Cari role di database
	role, err := s.repo.FindRoleByCode(ctx, targetRoleCode)
	if err != nil || role == nil {
		return nil, fmt.Errorf("Role '%s' tidak ditemukan di database", targetRoleCode)
	}

	// 5. Buat JWT baru dengan role yang di-assume tapi tetap simpan original level
	userData := map[string]any{
		"id":                 user.ID,
		"name":               user.Name,
		"email":              user.Email,
		"group_id":           user.GroupID,
		"level_id":           user.LevelID,
		"role_id":            role.ID, // Role yang di-assume
		"group":              user.Group,
		"group_code":         user.GroupCode,
		"level":              user.Level,     // Tetap SUPERADMIN
		"level_code":         user.LevelCode, // Tetap SUPERADMIN
		"role":               role.Name,      // Role yang di-assume
		"role_code":          role.Code,      // Role yang di-assume
		"status":             user.Status,
		"photo":              user.Photo,
		"assumed_role":       true,
		"original_role":      user.Role,
		"original_role_code": user.RoleCode,
	}

	roleCode := strings.ToUpper(strings.TrimSpace(role.Code))
	userData["role_code"] = roleCode

	token, err := utils.CreateJWT(userData, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat token")
	}

	// Tentukan redirect URL berdasarkan assumed role
	var redirectURL string
	switch roleCode {
	case "OPERASIONAL", "CEO":
		redirectURL = "/ops/dashboard"
	case "COMPLIANCE":
		redirectURL = "/compliance/dashboard"
	case "CLIENT":
		redirectURL = "/client/dashboard"
	default:
		redirectURL = "/ops/dashboard"
	}

	return AssumeRoleResult{
		Token:       token,
		RedirectURL: redirectURL,
		AssumedRole: role.Code,
	}, nil
}
