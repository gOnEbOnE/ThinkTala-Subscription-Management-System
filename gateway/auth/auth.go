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

func resolveClaimsFromCookie(cookieValue string) (*utils.AppClaims, error) {
	if strings.TrimSpace(cookieValue) == "" {
		return nil, errors.New("empty token cookie")
	}

	// Legacy flow: cookie contains encrypted Redis key -> resolve raw JWT from Redis.
	if decryptedKey, err := utils.Decrypt(cookieValue); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if jwtRaw, redisErr := utils.RedisGet(ctx, string(decryptedKey)); redisErr == nil {
			return utils.ValidateJWT(jwtRaw)
		}
	}

	// Fallback local flow: cookie contains raw JWT directly.
	return utils.ValidateJWT(cookieValue)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("[AUTH] Path: %s", path)

		if isPublicPath(path) {
			next.ServeHTTP(w, r)
			return
		}

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			log.Printf("[AUTH] No token cookie found")
			http.Redirect(w, r, "/account/login", http.StatusSeeOther)
			return
		}

		claims, err := resolveClaimsFromCookie(tokenCookie.Value)
		if err != nil {
			log.Printf("[AUTH] Invalid token: %v", err)
			http.Redirect(w, r, "/account/login", http.StatusSeeOther)
			return
		}

		userRole := ""
		if u, ok := claims.User["role_code"].(string); ok {
			userRole = strings.ToUpper(strings.TrimSpace(u))
		}
		if userRole == "" {
			if u, ok := claims.User["role"].(string); ok {
				userRole = strings.ToUpper(strings.TrimSpace(u))
			}
		}

		levelCode, _ := claims.User["level_code"].(string)
		if strings.EqualFold(strings.TrimSpace(levelCode), "SUPERADMIN") && !isAssumedRole(claims.User["assumed_role"]) {
			userRole = "SUPERADMIN"
		}

		if !isRoleAllowed(path, userRole) {
			log.Printf("[AUTH] Role %s not allowed for path %s", userRole, path)
			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user_role", userRole)
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
	role = strings.ToUpper(strings.TrimSpace(role))

	// Package dashboard is restricted strictly to MANAGEMENT/ADMIN only.
	if strings.HasPrefix(path, "/api/dashboard/packages") ||
		strings.HasPrefix(path, "/api/dashboard/package/") ||
		strings.HasPrefix(path, "/management/dashboard-packages") ||
		strings.HasPrefix(path, "/dashboard/packages/") {
		return role == "MANAGEMENT" || role == "ADMIN" || role == "SUPERADMIN"
	}

	// Management dashboard and APIs are restricted to MANAGEMENT/ADMIN/SUPERADMIN only.
	if strings.HasPrefix(path, "/api/dashboard/") ||
		strings.HasPrefix(path, "/management/") ||
		strings.HasPrefix(path, "/dashboard/customer/") {
		return role == "MANAGEMENT" || role == "ADMIN" || role == "SUPERADMIN"
	}

	// Customer Support access
	if strings.HasPrefix(path, "/ops/tickets") || strings.HasPrefix(path, "/ops/support-ticket-detail") {
		return role == "ADMIN_SUPPORT"
	}

	if strings.HasPrefix(path, "/api/admin/support/tickets") {
		return role == "ADMIN_SUPPORT"
	}

	// SuperAdmin & CEO can access everything
	if role == "SUPERADMIN" || role == "CEO" {
		return true
	}

	// /api/admin/users → SUPERADMIN only (user management)
	if strings.HasPrefix(path, "/api/admin/users") {
		return role == "SUPERADMIN"
	}

	if strings.HasPrefix(path, "/api/admin/kyc") {
		return role == "COMPLIANCE" || role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/api/admin/") {
		return role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/api/subscription/") {
		return role == "CLIENT" || role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/api/kyc/") {
		return role == "CLIENT" || role == "COMPLIANCE" || role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/api/notifications") {
		return role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/api/help/") {
		return role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/api/operational/") {
		return role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/api/ops/") {
		return role == "OPERASIONAL"
	}

	if strings.HasPrefix(path, "/client/") && role == "CLIENT" {
		return true
	}

	if strings.HasPrefix(path, "/support/") && role == "CLIENT" {
		return true
	}

	if strings.HasPrefix(path, "/ops/") && role == "OPERASIONAL" {
		return true
	}

	if strings.HasPrefix(path, "/compliance/") && role == "COMPLIANCE" {
		return true
	}

	return false
}

func CheckRoleAccess(path string, role string) bool {
	return isRoleAllowed(path, role)
}

func isAssumedRole(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return strings.EqualFold(strings.TrimSpace(x), "true")
	default:
		return false
	}
}

type TokenUser struct {
	ID          string
	Email       string
	Role        string
	RoleCode    string
	LevelCode   string
	AssumedRole bool
}

func GetUserFromToken(r *http.Request) (*TokenUser, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Value == "" {
		return nil, errors.New("token cookie not found")
	}

	claims, err := resolveClaimsFromCookie(tokenCookie.Value)
	if err != nil {
		return nil, err
	}

	u := claims.User

	id, _ := u["id"].(string)
	email, _ := u["email"].(string)
	role, _ := u["role"].(string)
	roleCode, _ := u["role_code"].(string)
	levelCode, _ := u["level_code"].(string)
	assumedRole := isAssumedRole(u["assumed_role"])

	if roleCode == "" {
		roleCode = strings.ToUpper(strings.TrimSpace(role))
	}

	// Level SUPERADMIN harus membawa role_code SUPERADMIN, kecuali saat sedang assume role.
	if strings.EqualFold(strings.TrimSpace(levelCode), "SUPERADMIN") && !assumedRole {
		roleCode = "SUPERADMIN"
		if strings.TrimSpace(role) == "" || strings.EqualFold(strings.TrimSpace(role), "CEO") {
			role = "Super Admin"
		}
	}

	return &TokenUser{
		ID:          id,
		Email:       email,
		Role:        role,
		RoleCode:    strings.ToUpper(strings.TrimSpace(roleCode)),
		LevelCode:   strings.ToUpper(strings.TrimSpace(levelCode)),
		AssumedRole: assumedRole,
	}, nil
}
