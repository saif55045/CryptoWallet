package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// getEncryptionKey derives a 32-byte key from JWT_SECRET
func getEncryptionKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-encryption-key-change-in-production"
	}
	
	// Use SHA-256 to derive a 32-byte key
	hash := sha256.Sum256([]byte(secret))
	return hash[:]
}

// EncryptPrivateKey encrypts the private key using AES-GCM
func EncryptPrivateKey(privateKeyHex string) (string, error) {
	key := getEncryptionKey()
	plaintext := []byte(privateKeyHex)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and append nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	
	// Return as base64 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey decrypts the encrypted private key
func DecryptPrivateKey(encryptedKey string) (string, error) {
	key := getEncryptionKey()
	
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
