package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

// NewUUID menghasilkan UUID v7 (Sortable by Time).
// Lebih cepat untuk indexing database dibanding UUID v4 biasa.
// Contoh: "018d2a6e-..."

// RandomString menghasilkan string acak hex sepanjang n * 2 karakter.
// Aman untuk token sesi, OTP, atau nama file.
// Contoh n=16 -> output 32 karakter
func RandomString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// Slugify mengubah string menjadi format URL-friendly.
// Contoh: "Halo Dunia 2026!" -> "halo-dunia-2026"
func Slugify(s string) string {
	// 1. Lowercase
	s = strings.ToLower(s)
	// 2. Ganti spasi dengan dash
	s = strings.ReplaceAll(s, " ", "-")
	// 3. Hapus karakter non-alphanumeric (kecuali dash)
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")
	// 4. Hapus dash berulang
	reg2, _ := regexp.Compile("-+")
	s = reg2.ReplaceAllString(s, "-")
	// 5. Trim dash di awal/akhir
	return strings.Trim(s, "-")
}
