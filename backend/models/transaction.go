package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransactionStatus represents the current state of a transaction
type TransactionStatus string

const (
	TxStatusPending   TransactionStatus = "pending"
	TxStatusConfirmed TransactionStatus = "confirmed"
	TxStatusFailed    TransactionStatus = "failed"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TxTypeTransfer TransactionType = "transfer"
	TxTypeCoinbase TransactionType = "coinbase"
	TxTypeZakat    TransactionType = "zakat"
)

// SignedInput represents a UTXO input with its signature
type SignedInput struct {
	TransactionID string  `json:"transactionId" bson:"transactionId"` // Reference to source UTXO's transaction
	OutputIndex   int     `json:"outputIndex" bson:"outputIndex"`     // Index in the source transaction
	Amount        float64 `json:"amount" bson:"amount"`               // Amount from this input
	PublicKey     string  `json:"publicKey" bson:"publicKey"`         // Sender's public key
	Signature     string  `json:"signature" bson:"signature"`         // Digital signature proving ownership
}

// TransactionOutput represents an output in a transaction
type TransactionOutput struct {
	WalletID  string  `json:"walletId" bson:"walletId"`   // Recipient's wallet ID
	Amount    float64 `json:"amount" bson:"amount"`       // Amount to send
	PublicKey string  `json:"publicKey" bson:"publicKey"` // Recipient's public key
}

// Transaction represents a complete blockchain transaction
type Transaction struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	TransactionID string              `json:"transactionId" bson:"transactionId"` // Unique hash of the transaction
	Type          TransactionType     `json:"type" bson:"type"`
	Inputs        []SignedInput       `json:"inputs" bson:"inputs"`
	Outputs       []TransactionOutput `json:"outputs" bson:"outputs"`
	TotalInput    float64             `json:"totalInput" bson:"totalInput"`
	TotalOutput   float64             `json:"totalOutput" bson:"totalOutput"`
	Fee           float64             `json:"fee" bson:"fee"` // TotalInput - TotalOutput (goes to miners)
	SenderWallet  string              `json:"senderWallet" bson:"senderWallet"`
	Status        TransactionStatus   `json:"status" bson:"status"`
	BlockHash     string              `json:"blockHash,omitempty" bson:"blockHash"`
	BlockHeight   int64               `json:"blockHeight,omitempty" bson:"blockHeight"`
	Timestamp     time.Time           `json:"timestamp" bson:"timestamp"`
	ConfirmedAt   *time.Time          `json:"confirmedAt,omitempty" bson:"confirmedAt"`
	Message       string              `json:"message,omitempty" bson:"message"` // Optional memo
}

// CreateTransactionRequest is used when creating a new transaction
type CreateTransactionRequest struct {
	RecipientWalletID string  `json:"recipientWalletId" binding:"required"`
	Amount            float64 `json:"amount" binding:"required,gt=0"`
	Message           string  `json:"message"`
}

// SignTransactionRequest contains the data to sign for a transaction
type SignTransactionRequest struct {
	TransactionID string        `json:"transactionId" binding:"required"`
	Signatures    []InputSignature `json:"signatures" binding:"required"`
}

// InputSignature pairs an input index with its signature
type InputSignature struct {
	InputIndex int    `json:"inputIndex"`
	Signature  string `json:"signature"`
}

// TransactionPreview shows what a transaction will look like before signing
type TransactionPreview struct {
	TransactionID     string              `json:"transactionId"`
	Inputs            []SignedInput       `json:"inputs"`
	Outputs           []TransactionOutput `json:"outputs"`
	TotalInput        float64             `json:"totalInput"`
	TotalOutput       float64             `json:"totalOutput"`
	Fee               float64             `json:"fee"`
	DataToSign        []string            `json:"dataToSign"` // Hash of each input to sign
	RecipientWalletID string              `json:"recipientWalletId"`
	Amount            float64             `json:"amount"`
	Change            float64             `json:"change"`
}

// TransactionResponse is returned after transaction operations
type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
	Message     string      `json:"message"`
}

// TransactionHistoryItem represents a transaction in the history list
type TransactionHistoryItem struct {
	TransactionID string            `json:"transactionId"`
	Type          TransactionType   `json:"type"`
	Direction     string            `json:"direction"` // "sent" or "received"
	Amount        float64           `json:"amount"`
	Fee           float64           `json:"fee"`
	Counterparty  string            `json:"counterparty"` // Other party's wallet ID
	Status        TransactionStatus `json:"status"`
	Timestamp     time.Time         `json:"timestamp"`
	Message       string            `json:"message,omitempty"`
}
