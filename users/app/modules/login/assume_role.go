package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/master-abror/zaframework/core/concurrency"
	resp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/token"
	"github.com/master-abror/zaframework/core/utils"
)

type assumeRoleActor struct {
	ID      string
	Email   string
	Level   string
	LevelID string
}

func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case nil:
		return ""
	default:
		return fmt.Sprint(v)
	}
}

func resolveClaimsFromTokenCookie(cookieValue string) (*utils.AppClaims, error) {
	if strings.TrimSpace(cookieValue) == "" {
		return nil, errors.New("empty token cookie")
	}

	if decryptedKey, err := utils.Decrypt(cookieValue); err == nil {
		if jwtRaw, redisErr := utils.RedisGet(context.Background(), string(decryptedKey)); redisErr == nil && strings.TrimSpace(jwtRaw) != "" {
			return utils.ValidateJWT(jwtRaw)
		}
	}

	return utils.ValidateJWT(cookieValue)
}

func resolveAssumeRoleActor(r *http.Request) (*assumeRoleActor, error) {
	if currentUser, err := token.GetUserAuthz(r, "token"); err == nil && currentUser != nil {
		return &assumeRoleActor{
			ID:      currentUser.ID,
			Email:   currentUser.Email,
			Level:   currentUser.Level,
			LevelID: currentUser.LevelID,
		}, nil
	}

	tokenCookie, err := r.Cookie("token")
	if err != nil || strings.TrimSpace(tokenCookie.Value) == "" {
		return nil, errors.New("token cookie not found")
	}

	claims, err := resolveClaimsFromTokenCookie(tokenCookie.Value)
	if err != nil {
		return nil, err
	}

	user := claims.User
	return &assumeRoleActor{
		ID:      toString(user["id"]),
		Email:   toString(user["email"]),
		Level:   toString(user["level_code"]),
		LevelID: toString(user["level_id"]),
	}, nil
}

func isSuperAdminActor(actor *assumeRoleActor) bool {
	if actor == nil {
		return false
	}
	return strings.ToUpper(strings.TrimSpace(actor.Level)) == "SUPERADMIN" || strings.TrimSpace(actor.LevelID) == "1"
}

// AssumeRole — POST /api/auth/assume-role
// Hanya SUPERADMIN yang bisa assume role lain.
// Mengubah JWT claims dengan role target, update Redis session.
func (c *Controller) AssumeRole(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		resp.ApiJSON(w, r, http.StatusUnsupportedMediaType, false, "Content-Type harus application/json", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	// 1. Validate current user is SUPERADMIN
	currentUser, err := resolveAssumeRoleActor(r)
	if err != nil || currentUser == nil {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid, silakan login kembali", nil)
		return
	}

	if !isSuperAdminActor(currentUser) {
		resp.ApiJSON(w, r, http.StatusForbidden, false, "Hanya Super Admin yang dapat melakukan simulasi role", nil)
		return
	}

	// 2. Parse input
	var input AssumeRoleInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "Format JSON tidak valid", nil)
		return
	}

	targetRoleCode := strings.ToUpper(strings.TrimSpace(input.TargetRoleCode))
	if targetRoleCode == "" {
		resp.ApiJSON(w, r, http.StatusBadRequest, false, "target_role_code wajib diisi", nil)
		return
	}

	// 3. Dispatch to worker
	payload := map[string]any{
		"user_id":          currentUser.ID,
		"user_email":       currentUser.Email,
		"target_role_code": targetRoleCode,
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "assume_role", payload, concurrency.PriorityHigh)
	if err != nil {
		fmt.Printf("[ASSUME-ROLE ERROR] %v\n", err)
		resp.ApiJSON(w, r, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	assumeRes, ok := result.(AssumeRoleResult)
	if !ok {
		resp.ApiJSON(w, r, http.StatusInternalServerError, false, "Terjadi kesalahan internal", nil)
		return
	}

	// 4. Update session in Redis with new JWT
	if assumeRes.Token != "" {
		tokenKey, err := token.SetUserAuthz(w, r, assumeRes.Token)
		if err != nil {
			resp.ApiJSON(w, r, http.StatusInternalServerError, false, "Gagal mengupdate sesi: "+err.Error(), nil)
			return
		}

		// Set cookie for gateway
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenKey,
			Path:     "/",
			Domain:   "",
			MaxAge:   86400,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		fmt.Printf("[ASSUME-ROLE] User %s assumed role %s\n", currentUser.Email, assumeRes.AssumedRole)
	}

	resp.ApiJSON(w, r, http.StatusOK, true, "Role berhasil diubah ke "+assumeRes.AssumedRole, map[string]any{
		"assumed_role": assumeRes.AssumedRole,
		"redirect_url": assumeRes.RedirectURL,
	})
}

// GetRoles — GET /api/auth/roles
// Mengembalikan daftar role yang tersedia (hanya untuk SUPERADMIN)
func (c *Controller) GetRoles(w http.ResponseWriter, r *http.Request) {
	currentUser, err := resolveAssumeRoleActor(r)
	if err != nil || currentUser == nil {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid", nil)
		return
	}

	if !isSuperAdminActor(currentUser) {
		resp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak", nil)
		return
	}

	// We return hardcoded simulatable roles
	roles := []map[string]string{
		{"code": "CEO", "name": "CEO"},
		{"code": "OPERASIONAL", "name": "Operasional"},
		{"code": "COMPLIANCE", "name": "Compliance"},
		{"code": "MANAGEMENT", "name": "Management"},
		{"code": "CLIENT", "name": "Client"},
		{"code": "ADMIN_SUPPORT", "name": "Customer Support"},
	}

	resp.ApiJSON(w, r, http.StatusOK, true, "Daftar role tersedia", map[string]any{
		"roles": roles,
	})
}
