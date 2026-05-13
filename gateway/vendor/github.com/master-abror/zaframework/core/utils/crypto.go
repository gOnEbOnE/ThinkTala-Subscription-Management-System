package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// Config Encryption (Bisa diambil dari Env nanti)
// Pastikan Key 32 bytes (AES-256) dan IV 16 bytes
var (
	// Default fallback keys (Sebaiknya di-override via Env saat init)
	defaultKey = []byte("ab2409mel1203za1611ar1606rau0101")
	defaultIV  = []byte("r4u1ft1t@hmul1kh")
)

// -- AES Encryption --

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil // Avoid panic
	}
	return src[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Decrypt(encrypted string) ([]byte, error) {
	key := []byte(GetEnv("APP_KEY", string(defaultKey)))
	iv := []byte(GetEnv("APP_IV", string(defaultIV)))

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	ciphertext = PKCS5UnPadding(ciphertext)
	if ciphertext == nil {
		return nil, errors.New("unpadding failed")
	}

	return ciphertext, nil
}

func Encrypt(plaintext string) (string, error) {
	key := []byte(GetEnv("APP_KEY", string(defaultKey)))
	iv := []byte(GetEnv("APP_IV", string(defaultIV)))

	plainBytes := []byte(plaintext)
	plainBytes = PKCS5Padding(plainBytes, aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainBytes)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
