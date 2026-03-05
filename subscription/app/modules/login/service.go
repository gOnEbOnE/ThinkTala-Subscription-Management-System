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

// Fungsi Utama yang dipanggil Worker
func (s *Service) ProcessLoginJob(ctx context.Context, payload any) (any, error) {
	// 1. Parse Payload
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	loginID := data["login_id"].(string)
	password := data["password"].(string)

	// Optional fields
	ip, _ := data["ip"].(string)
	browser, _ := data["browser"].(string)
	os, _ := data["os"].(string)
	lat, _ := data["latitude"].(string)
	long, _ := data["longitude"].(string)

	// 3. Find User (Query Utama)
	user, err := s.repo.FindUserByEmail(ctx, "email", loginID)
	if err != nil {
		return nil, err
	}

	// 5. Validate Password
	if hash := utils.CheckPasswordHash(password, user.Password); hash != true {
		return nil, fmt.Errorf("Login ID atau kata sandi anda salah!")
	}

	// 6. Build Current User Data (Logika Jabatan yang rumit itu)
	userData := s.buildCurrentUser(ctx, user, loginID, ip, lat, long, browser, os)

	// 7. Create Token
	token, err := utils.CreateJWT(userData, 24*time.Hour) // Menggunakan JWT Util Engine
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat token")
	}

	// 8. Encrypt Signature (Sesuai kode lama)
	token = s.encryptTokenSignature(token)

	// Tambahkan token ke userData
	userData["token"] = token

	return LoginResult{
		Token:    token,
		UserData: userData,
	}, nil
}

// --- Helper Logics ---

func (s *Service) encryptTokenSignature(token string) string {
	jwtParts := strings.Split(token, ".")
	if len(jwtParts) != 3 {
		return token
	}
	encrypted, _ := utils.Encrypt(jwtParts[2]) // Menggunakan AES engine
	return jwtParts[0] + "." + jwtParts[1] + "." + encrypted
}

func (s *Service) buildCurrentUser(ctx context.Context, user *User, inputID, ip, lat, long, browser, os string) map[string]any {
	// Determine Actual Login ID based on Account Type

	roleID := user.RoleID

	return map[string]any{
		"id":             user.ID,
		"login":          inputID,
		"login_id":       inputID,
		"fullname":       user.Fullname,
		"group_id":       user.GroupID,
		"level_id":       user.LevelID,
		"role_id":        roleID,
		"group":          user.Group,
		"level":          user.Level,
		"role":           user.Role,
		"nik":            user.Nik,
		"active":         user.IsActive,
		"photo":          user.Photo,
		"ip":             ip,
		"latitude":       lat,
		"longitude":      long,
		"browser":        browser,
		"os":             os,
		"multiple_login": user.MultipleLogin,
		"jenis_kelamin":  user.JenisKelamin,
		"jabatan":        "",
		"jenis_akun":     "",
	}
}
