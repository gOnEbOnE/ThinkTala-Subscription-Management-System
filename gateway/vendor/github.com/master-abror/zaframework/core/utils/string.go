package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

// -- String Manipulation --

func ReplaceSpaceToUnderscore(words string) string {
	return strings.ReplaceAll(words, " ", "_")
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func IsNumber(number string) bool {
	_, err := strconv.ParseInt(number, 10, 64)
	return err == nil
}

// -- Random Generators (Secure) --

const (
	digits   = "0123456789"
	alphanum = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// GenerateRandomString generates a secure random string with given length
func GenerateRandomString(length int) string {
	return secureRandom(length, alphanum)
}

// GenerateRandomNumber generates a secure random numeric string
func GenerateRandomNumber(length int) string {
	return secureRandom(length, digits)
}

func secureRandom(length int, charset string) string {
	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))

	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			b[i] = charset[0] // Fallback (should rare)
		} else {
			b[i] = charset[n.Int64()]
		}
	}
	return string(b)
}
