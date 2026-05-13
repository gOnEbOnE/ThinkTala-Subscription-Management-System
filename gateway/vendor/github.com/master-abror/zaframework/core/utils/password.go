package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengubah plain text password menjadi hash (enkripsi satu arah).
// Digunakan saat Register atau Reset Password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash membandingkan plain password dari input user dengan hash di database.
// Digunakan saat Login. Return true jika cocok.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
