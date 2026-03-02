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
		return nil, fmt.Errorf("Login ID atau kata sandi anda salah!")
	}

	if hash := utils.CheckPasswordHash(password, user.Password); hash != true {
		return nil, fmt.Errorf("Login ID atau kata sandi anda salah!")
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
