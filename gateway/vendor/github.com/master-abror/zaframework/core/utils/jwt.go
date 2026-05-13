package utils

import (
	"crypto/rsa"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

// InitJWTLoadKeys memuat file private.pem dan public.pem dari path
func InitJWTLoadKeys(privatePath, publicPath string) error {
	// 1. Load Private Key
	privBytes, err := os.ReadFile(privatePath)
	if err != nil {
		return err
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return err
	}

	// 2. Load Public Key
	pubBytes, err := os.ReadFile(publicPath)
	if err != nil {
		return err
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return err
	}

	return nil
}

// Struct Custom Claims (Payload Token)
type AppClaims struct {
	User map[string]any `json:"user"`
	jwt.RegisteredClaims
}

// CreateToken membuat JWT Signed dengan Private Key
func CreateJWT(payload map[string]any, duration time.Duration) (string, error) {
	if signKey == nil {
		return "", errors.New("JWT keys not initialized")
	}

	claims := AppClaims{
		User: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "za-framework",
		},
	}

	// Gunakan SigningMethodRS256 (Asymmetric)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signKey)
}

// ValidateToken memverifikasi token dengan Public Key
func ValidateJWT(tokenString string) (*AppClaims, error) {
	if verifyKey == nil {
		return nil, errors.New("JWT keys not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi Algoritma harus RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshJWT(tokenString string, duration time.Duration) (string, *AppClaims, error) {
	if signKey == nil || verifyKey == nil {
		return "", nil, errors.New("JWT keys not initialized")
	}

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", nil, err
	}

	now := time.Now()

	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(duration))

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	newToken, err := token.SignedString(signKey)
	if err != nil {
		return "", nil, err
	}

	return newToken, claims, nil
}

func EncryptTokenSignature(token string) string {
	jwtParts := strings.Split(token, ".")
	if len(jwtParts) != 3 {
		return token
	}
	encrypted, _ := Encrypt(jwtParts[2]) // Menggunakan AES engine
	return jwtParts[0] + "." + jwtParts[1] + "." + encrypted
}

// =================================================================
// TAMBAHAN BARU (UNTUK SSO LEGACY)
// =================================================================

// CreateHS256Token membuat token HMAC (Symmetric) untuk keperluan redirect SSO
// Fungsi ini dibutuhkan karena SSO menggunakan kunci simetris lama, bukan RSA
func CreateHS256Token(key []byte, userID, state string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"state":   state,
		"exp":     time.Now().Add(duration).Unix(),
	}

	// Gunakan SigningMethodHS256 (Symmetric)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign dengan key []byte langsung
	return token.SignedString(key)
}
