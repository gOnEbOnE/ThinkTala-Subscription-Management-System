package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/master-abror/zaframework/core/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("[AUTH] Path: %s", path)

		// Skip public paths
		if isPublicPath(path) {
			next.ServeHTTP(w, r)
			return
		}

		// ✅ STEP 1: Get token from cookie
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			log.Printf("[AUTH] No token cookie found")
			http.Redirect(w, r, "/account/login", http.StatusSeeOther)
			return
		}

		encryptedKey := tokenCookie.Value

		// ✅ STEP 2: Decrypt key to get UUID
		decryptedKey, err := utils.Decrypt(encryptedKey)
		if err != nil {
			log.Printf("[AUTH] Failed to decrypt token: %v", err)
			http.Redirect(w, r, "/account/login", http.StatusSeeOther)
			return
		}

		// ✅ STEP 3: Get RAW JWT from Redis using UUID key
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		jwtRaw, err := utils.RedisGet(ctx, string(decryptedKey))
		if err != nil {
			log.Printf("[AUTH] Token error: %v", err)
			http.Redirect(w, r, "/account/login", http.StatusSeeOther)
			return
		}

		// ✅ STEP 4: Validate JWT
		claims, err := utils.ValidateJWT(jwtRaw)
		if err != nil {
			log.Printf("[AUTH] Invalid token: %v", err)
			http.Redirect(w, r, "/account/login", http.StatusSeeOther)
			return
		}

		// ✅ STEP 5: Check role authorization
		userRole := ""
		if u, ok := claims.User["role"].(string); ok {
			userRole = strings.ToUpper(u)
		}

		if !isRoleAllowed(path, userRole) {
			log.Printf("[AUTH] Role %s not allowed for path %s", userRole, path)
			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
			return
		}

		// ✅ STEP 6: Add user info to context
		ctx = context.WithValue(r.Context(), "user_role", userRole)
		ctx = context.WithValue(ctx, "user_id", claims.User["id"])
		ctx = context.WithValue(ctx, "user_email", claims.User["email"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isPublicPath(path string) bool {
	publicPaths := []string{
		"/account/login",
		"/account/register",
		"/account/verify-otp",
		"/api/auth",
		"/api/auth/",
		"/assets/",
		"/favicon.ico",
		"/",
	}
	for _, p := range publicPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func isRoleAllowed(path string, role string) bool {
	// Normalize role
	role = strings.ToUpper(strings.TrimSpace(role))

	// SuperAdmin & CEO can access everything
	if role == "SUPERADMIN" || role == "CEO" {
		return true
	}

	// /api/admin/* → Admin-only: CEO, SUPERADMIN, OPERASIONAL
	if strings.HasPrefix(path, "/api/admin/") {
		return role == "OPERASIONAL"
	}

	// /api/subscription/catalog → CLIENT, OPERASIONAL, and above can view
	if strings.HasPrefix(path, "/api/subscription/") {
		return role == "CLIENT" || role == "OPERASIONAL"
	}

	// /client/* → CLIENT role
	if strings.HasPrefix(path, "/client/") && role == "CLIENT" {
		return true
	}

	// /ops/* → OPERASIONAL role
	if strings.HasPrefix(path, "/ops/") && role == "OPERASIONAL" {
		return true
	}

	// /compliance/* → COMPLIANCE role
	if strings.HasPrefix(path, "/compliance/") && role == "COMPLIANCE" {
		return true
	}

	return false
}

// CheckRoleAccess is an exported wrapper for role/path authorization.
func CheckRoleAccess(path string, role string) bool {
	return isRoleAllowed(path, role)
}

type TokenUser struct {
	ID       string
	Email    string
	Role     string
	RoleCode string
}

// GetUserFromToken reads auth cookie, resolves JWT from Redis, validates it,
// then returns minimal user info used by gateway routing.
func GetUserFromToken(r *http.Request) (*TokenUser, error) {
	// 1) Cookie: encrypted redis key
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Value == "" {
		return nil, errors.New("token cookie not found")
	}

	// 2) Decrypt key (UUID/string key for redis lookup)
	decryptedKey, err := utils.Decrypt(tokenCookie.Value)
	if err != nil {
		return nil, err
	}

	// 3) Read raw JWT from redis
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	jwtRaw, err := utils.RedisGet(ctx, string(decryptedKey))
	if err != nil {
		return nil, err
	}

	// 4) Validate JWT (RSA)
	claims, err := utils.ValidateJWT(jwtRaw)
	if err != nil {
		return nil, err
	}

	u := claims.User

	id, _ := u["id"].(string)
	email, _ := u["email"].(string)
	role, _ := u["role"].(string)
	roleCode, _ := u["role_code"].(string)

	// Fallback if role_code is absent
	if roleCode == "" {
		roleCode = strings.ToUpper(strings.TrimSpace(role))
	}

	return &TokenUser{
		ID:       id,
		Email:    email,
		Role:     role,
		RoleCode: strings.ToUpper(strings.TrimSpace(roleCode)),
	}, nil
}
