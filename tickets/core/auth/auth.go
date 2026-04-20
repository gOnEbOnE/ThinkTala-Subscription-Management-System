package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"tickets/app/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	publicKeyOnce sync.Once
	jwtPublicKey  any
	jwtKeyErr     error
)

func parseClaimsFromAuthorization(authHeader string) (jwt.MapClaims, error) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(authHeader)), "bearer ") {
		return nil, errors.New("missing bearer token")
	}

	tokenString := strings.TrimSpace(authHeader[7:])
	if tokenString == "" {
		return nil, errors.New("empty bearer token")
	}

	pubKey, err := loadJWTPublicKey()
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid bearer token")
	}

	return claims, nil
}

func extractRoleFromClaims(claims jwt.MapClaims) string {
	if roleCode, ok := claims["role_code"].(string); ok {
		return strings.ToUpper(strings.TrimSpace(roleCode))
	}
	if role, ok := claims["role"].(string); ok {
		return strings.ToUpper(strings.TrimSpace(role))
	}

	if userMapAny, ok := claims["user"].(map[string]any); ok {
		if roleCode, ok := userMapAny["role_code"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(roleCode))
		}
		if role, ok := userMapAny["role"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(role))
		}
	}

	if userMapIface, ok := claims["user"].(map[string]interface{}); ok {
		if roleCode, ok := userMapIface["role_code"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(roleCode))
		}
		if role, ok := userMapIface["role"].(string); ok {
			return strings.ToUpper(strings.TrimSpace(role))
		}
	}

	return ""
}

func loadJWTPublicKey() (any, error) {
	publicKeyOnce.Do(func() {
		if b64 := strings.TrimSpace(os.Getenv("JWT_PUBLIC_KEY_B64")); b64 != "" {
			decoded, err := base64.StdEncoding.DecodeString(b64)
			if err == nil {
				if key, keyErr := jwt.ParseRSAPublicKeyFromPEM(decoded); keyErr == nil {
					jwtPublicKey = key
					jwtKeyErr = nil
					return
				}
			}
		}

		paths := []string{"certs/public.pem", "../users/certs/public.pem"}
		for _, p := range paths {
			pemBytes, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
			if err == nil {
				jwtPublicKey = key
				jwtKeyErr = nil
				return
			}
		}

		jwtKeyErr = errors.New("jwt public key not found")
	})

	if jwtPublicKey == nil {
		return nil, jwtKeyErr
	}

	return jwtPublicKey, nil
}

func extractRoleFromAuthorization(authHeader string) string {
	claims, err := parseClaimsFromAuthorization(authHeader)
	if err != nil {
		return ""
	}

	return extractRoleFromClaims(claims)
}

func HasSupportMonitoringRole(r *http.Request) bool {
	role := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	if role == "" {
		role = extractRoleFromAuthorization(r.Header.Get("Authorization"))
	}
	return role == "ADMIN_SUPPORT" || role == "OPERASIONAL"
}

func BuildDisplayName(rawName, rawEmail string) string {
	name := strings.TrimSpace(rawName)
	if name != "" {
		return name
	}

	email := strings.TrimSpace(rawEmail)
	if email == "" {
		return "User"
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
		return strings.TrimSpace(parts[0])
	}
	return email
}

func extractUserFromClaims(claims jwt.MapClaims) models.AuthenticatedUser {
	user := models.AuthenticatedUser{}

	if role := extractRoleFromClaims(claims); role != "" {
		user.Role = role
	}

	if id, ok := claims["user_id"].(string); ok {
		user.ID = strings.TrimSpace(id)
	}
	if email, ok := claims["email"].(string); ok {
		user.Email = strings.TrimSpace(email)
	}
	if name, ok := claims["name"].(string); ok {
		user.Name = strings.TrimSpace(name)
	}

	if userMapAny, ok := claims["user"].(map[string]any); ok {
		if id, ok := userMapAny["id"].(string); ok && user.ID == "" {
			user.ID = strings.TrimSpace(id)
		}
		if email, ok := userMapAny["email"].(string); ok && user.Email == "" {
			user.Email = strings.TrimSpace(email)
		}
		if name, ok := userMapAny["name"].(string); ok && user.Name == "" {
			user.Name = strings.TrimSpace(name)
		}
	}

	if userMapIface, ok := claims["user"].(map[string]interface{}); ok {
		if id, ok := userMapIface["id"].(string); ok && user.ID == "" {
			user.ID = strings.TrimSpace(id)
		}
		if email, ok := userMapIface["email"].(string); ok && user.Email == "" {
			user.Email = strings.TrimSpace(email)
		}
		if name, ok := userMapIface["name"].(string); ok && user.Name == "" {
			user.Name = strings.TrimSpace(name)
		}
	}

	if user.Name == "" {
		user.Name = BuildDisplayName("", user.Email)
	}

	return user
}

func GetAuthenticatedUserFromRequest(r *http.Request) models.AuthenticatedUser {
	user := models.AuthenticatedUser{
		ID:    strings.TrimSpace(r.Header.Get("X-User-ID")),
		Email: strings.TrimSpace(r.Header.Get("X-User-Email")),
		Role:  strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role"))),
	}
	user.Name = BuildDisplayName("", user.Email)

	claims, err := parseClaimsFromAuthorization(r.Header.Get("Authorization"))
	if err != nil {
		return user
	}

	fromClaims := extractUserFromClaims(claims)
	if user.ID == "" {
		user.ID = fromClaims.ID
	}
	if user.Email == "" {
		user.Email = fromClaims.Email
	}
	if user.Role == "" {
		user.Role = fromClaims.Role
	}
	if fromClaims.Name != "" {
		user.Name = fromClaims.Name
	} else {
		user.Name = BuildDisplayName("", user.Email)
	}

	return user
}

func HasClientSupportCreateRole(r *http.Request) bool {
	user := GetAuthenticatedUserFromRequest(r)
	return user.Role == "CLIENT" || user.Role == "SUPERADMIN" || user.Role == "CEO"
}
