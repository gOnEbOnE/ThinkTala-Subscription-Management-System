package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateRSAKeyPair generates RSA private and public keys
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// SavePrivateKey saves RSA private key to PEM file
func SavePrivateKey(filename string, key *rsa.PrivateKey) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	return os.WriteFile(filename, privateKeyPem, 0600)
}

// SavePublicKey saves RSA public key to PEM file
func SavePublicKey(filename string, key *rsa.PublicKey) error {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	return os.WriteFile(filename, publicKeyPem, 0644)
}

// GenerateAndSaveKeys generates RSA key pair and saves to files
func GenerateAndSaveKeys(privatePath, publicPath string) error {
	fmt.Printf("🔐 Generating RSA key pair (2048-bit)...\n")

	privateKey, publicKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		return err
	}

	fmt.Printf("💾 Saving private key to: %s\n", privatePath)
	if err := SavePrivateKey(privatePath, privateKey); err != nil {
		return err
	}

	fmt.Printf("💾 Saving public key to: %s\n", publicPath)
	if err := SavePublicKey(publicPath, publicKey); err != nil {
		return err
	}

	fmt.Printf("✅ Written to %s\n", privatePath)
	fmt.Printf("✅ Written to %s\n\n", publicPath)

	return nil
}
