package token

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/master-abror/zaframework/core/session"
	"github.com/master-abror/zaframework/core/utils"
)

type CurrentUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Photo     *string   `json:"photo"`
	GroupID   string    `json:"group_id"`
	LevelID   string    `json:"level_id"`
	RoleID    string    `json:"role_id"`
	Status    string    `json:"status"`
	Level     string    `json:"level"`
	Role      string    `json:"role"`
	Group     string    `json:"group"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *string   `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by"`
}

func SetUserAuthz(w http.ResponseWriter, r *http.Request, token_value string) (string, error) {
	if !utils.IsRedisEnabled() {
		_, _ = session.Set(w, r, "token", token_value)
		return token_value, nil
	}

	token_key, err := utils.CreateUUID()
	if err != nil {
		return "", err
	}
	token_key_encrypt, err := utils.Encrypt(token_key.String())
	if err != nil {
		return "", err
	}
	_, err = session.Set(w, r, "token", token_key_encrypt)
	if err != nil {
		return "", err
	}

	ctx := r.Context()

	v := utils.GetEnv("SESSION_LIFETIME", "86400")
	sec, err := strconv.Atoi(v)
	if err != nil {
		return "", err
	}

	sessionLifetime := time.Duration(sec) * time.Second

	// ✅ SIMPAN RAW TOKEN (TANPA ENKRIPSI)
	err = utils.RedisSet(ctx, token_key.String(), token_value, sessionLifetime)
	if err != nil {
		return "", err
	}

	return token_key_encrypt, nil
}

func GetUserAuthz(r *http.Request, key string) (*CurrentUser, error) {
	token_key_encrypt := session.Get(r, key)

	val := token_key_encrypt

	tokenStr, ok := val.(string)
	if !ok || tokenStr == "" {
		return nil, errors.New("token_key_encrypt is missing or not string")
	}

	if !utils.IsRedisEnabled() {
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			return nil, err
		}

		u := claims.User
		photoStr, _ := u["photo"].(string)

		currentUser := CurrentUser{
			ID:      u["id"].(string),
			Name:    u["name"].(string),
			Email:   u["email"].(string),
			Role:    u["role"].(string),
			RoleID:  u["role_id"].(string),
			Level:   u["level"].(string),
			LevelID: u["level_id"].(string),
			Group:   u["group"].(string),
			GroupID: u["group_id"].(string),
			Status:  u["status"].(string),
			Photo:   &photoStr,
		}

		return &currentUser, nil
	}

	token_key_decrypt, err := utils.Decrypt(tokenStr)
	if err != nil {
		claims, jwtErr := utils.ValidateJWT(tokenStr)
		if jwtErr != nil {
			return nil, err
		}

		u := claims.User
		photoStr, _ := u["photo"].(string)

		currentUser := CurrentUser{
			ID:      u["id"].(string),
			Name:    u["name"].(string),
			Email:   u["email"].(string),
			Role:    u["role"].(string),
			RoleID:  u["role_id"].(string),
			Level:   u["level"].(string),
			LevelID: u["level_id"].(string),
			Group:   u["group"].(string),
			GroupID: u["group_id"].(string),
			Status:  u["status"].(string),
			Photo:   &photoStr,
		}

		return &currentUser, nil
	}

	ctx := r.Context()

	// ✅ AMBIL RAW TOKEN DARI REDIS
	token_raw, err := utils.RedisGet(ctx, string(token_key_decrypt))
	if err != nil {
		return nil, err
	}

	// ✅ VALIDASI JWT LANGSUNG
	claims, err := utils.ValidateJWT(token_raw)
	if err != nil {
		return nil, err
	}

	v := utils.GetEnv("JWT_EXPIRED")
	sec, err := strconv.Atoi(v)
	if err != nil {
		return nil, err
	}

	sessionLifetime := time.Duration(sec) * time.Second

	// ✅ REFRESH DAN SIMPAN RAW TOKEN
	newToken, _, err := utils.RefreshJWT(token_raw, sessionLifetime)
	if err != nil {
		return nil, err
	}

	err = utils.RedisSet(ctx, string(token_key_decrypt), newToken, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	u := claims.User

	photoStr, _ := u["photo"].(string)

	currentUser := CurrentUser{
		ID:      u["id"].(string),
		Name:    u["name"].(string),
		Email:   u["email"].(string),
		Role:    u["role"].(string),
		RoleID:  u["role_id"].(string),
		Level:   u["level"].(string),
		LevelID: u["level_id"].(string),
		Group:   u["group"].(string),
		GroupID: u["group_id"].(string),
		Status:  u["status"].(string),
		Photo:   &photoStr,
	}

	return &currentUser, nil
}
