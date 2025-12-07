package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

// SignData signs data using the private key and returns the signature as hex
func SignData(privateKeyHex string, data string) (string, error) {
	// Parse the private key
	privateKey, err := PrivateKeyFromHex(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// Hash the data
	hash := sha256.Sum256([]byte(data))

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// Encode r and s as hex (each 32 bytes)
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	// Pad to 32 bytes if necessary
	rPadded := make([]byte, 32)
	sPadded := make([]byte, 32)
	copy(rPadded[32-len(rBytes):], rBytes)
	copy(sPadded[32-len(sBytes):], sBytes)

	// Concatenate r and s
	signature := append(rPadded, sPadded...)
	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies a signature using the public key
func VerifySignature(publicKeyHex string, data string, signatureHex string) (bool, error) {
	// Parse the public key
	publicKey, err := PublicKeyFromHex(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Decode the signature
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	if len(signatureBytes) != 64 {
		return false, errors.New("invalid signature length")
	}

	// Extract r and s from signature
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Hash the data
	hash := sha256.Sum256([]byte(data))

	// Verify the signature
	valid := ecdsa.Verify(publicKey, hash[:], r, s)
	return valid, nil
}

// HashTransaction creates a hash of transaction data for signing
func HashTransaction(txID string, inputIndex int, amount float64, recipientWallet string) string {
	data := fmt.Sprintf("%s:%d:%.8f:%s", txID, inputIndex, amount, recipientWallet)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// HashTransactionData creates a comprehensive hash of all transaction data
func HashTransactionData(inputs []InputData, outputs []OutputData) string {
	data := ""
	
	// Add all inputs
	for _, input := range inputs {
		data += fmt.Sprintf("%s:%d:%.8f|", input.TransactionID, input.OutputIndex, input.Amount)
	}
	
	data += "->"
	
	// Add all outputs
	for _, output := range outputs {
		data += fmt.Sprintf("%s:%.8f|", output.WalletID, output.Amount)
	}
	
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// InputData represents input data for hashing
type InputData struct {
	TransactionID string
	OutputIndex   int
	Amount        float64
}

// OutputData represents output data for hashing
type OutputData struct {
	WalletID string
	Amount   float64
}

// GenerateTransactionID creates a unique transaction ID from transaction data
func GenerateTransactionID(senderWallet string, inputs []InputData, outputs []OutputData, timestamp int64) string {
	data := fmt.Sprintf("%s:%d:", senderWallet, timestamp)
	
	for _, input := range inputs {
		data += fmt.Sprintf("%s:%d:%.8f|", input.TransactionID, input.OutputIndex, input.Amount)
	}
	
	for _, output := range outputs {
		data += fmt.Sprintf("%s:%.8f|", output.WalletID, output.Amount)
	}
	
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CreateInputSignatureData creates the data string that needs to be signed for an input
func CreateInputSignatureData(txID string, inputTxID string, inputIndex int, amount float64) string {
	return fmt.Sprintf("%s:%s:%d:%.8f", txID, inputTxID, inputIndex, amount)
}
