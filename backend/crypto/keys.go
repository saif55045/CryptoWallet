package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

// KeyPair represents a public/private key pair
type KeyPair struct {
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     *ecdsa.PublicKey
	PrivateKeyHex string
	PublicKeyHex  string
}

// GenerateKeyPair generates a new ECDSA key pair using P-256 curve
func GenerateKeyPair() (*KeyPair, error) {
	// Generate private key using P-256 curve (secp256r1)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Get public key
	publicKey := &privateKey.PublicKey

	// Convert private key to hex
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// Convert public key to hex (compressed format)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	publicKeyHex := hex.EncodeToString(publicKeyBytes)

	return &KeyPair{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		PrivateKeyHex: privateKeyHex,
		PublicKeyHex:  publicKeyHex,
	}, nil
}

// GenerateWalletID creates a wallet ID by hashing the public key
func GenerateWalletID(publicKeyHex string) string {
	// First SHA-256 hash
	hash := sha256.Sum256([]byte(publicKeyHex))
	
	// Second SHA-256 hash (double hashing for extra security)
	hash2 := sha256.Sum256(hash[:])
	
	// Take first 20 bytes and convert to hex (40 characters)
	walletID := hex.EncodeToString(hash2[:20])
	
	return walletID
}

// PrivateKeyFromHex reconstructs a private key from hex string
func PrivateKeyFromHex(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	
	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	
	return privateKey, nil
}

// PublicKeyFromHex reconstructs a public key from hex string
func PublicKeyFromHex(publicKeyHex string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, err
	}
	
	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	
	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}
	
	return publicKey, nil
}

// PrivateKeyToPEM converts private key to PEM format
func PrivateKeyToPEM(privateKey *ecdsa.PrivateKey) (string, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	
	return string(pem.EncodeToMemory(pemBlock)), nil
}

// PublicKeyToPEM converts public key to PEM format
func PublicKeyToPEM(publicKey *ecdsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	
	return string(pem.EncodeToMemory(pemBlock)), nil
}
