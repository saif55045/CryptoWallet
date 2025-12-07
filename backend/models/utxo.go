package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UTXO represents an Unspent Transaction Output
// This is the fundamental building block for tracking balances in a UTXO-based system
type UTXO struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TransactionID   string             `json:"transactionId" bson:"transactionId"`     // ID of the transaction that created this UTXO
	OutputIndex     int                `json:"outputIndex" bson:"outputIndex"`         // Index of this output in the transaction
	WalletID        string             `json:"walletId" bson:"walletId"`               // Owner's wallet ID
	Amount          float64            `json:"amount" bson:"amount"`                   // Amount in this UTXO
	PublicKey       string             `json:"publicKey" bson:"publicKey"`             // Owner's public key for verification
	IsSpent         bool               `json:"isSpent" bson:"isSpent"`                 // Whether this UTXO has been spent
	SpentInTx       string             `json:"spentInTx,omitempty" bson:"spentInTx"`   // Transaction ID that spent this UTXO
	BlockHash       string             `json:"blockHash,omitempty" bson:"blockHash"`   // Block hash where this UTXO was confirmed
	IsConfirmed     bool               `json:"isConfirmed" bson:"isConfirmed"`         // Whether this UTXO is confirmed in a block
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	SpentAt         *time.Time         `json:"spentAt,omitempty" bson:"spentAt"`
}

// UTXOInput represents a reference to a UTXO being spent in a transaction
type UTXOInput struct {
	TransactionID string  `json:"transactionId" bson:"transactionId"` // Reference to the UTXO's transaction
	OutputIndex   int     `json:"outputIndex" bson:"outputIndex"`     // Index of the output being spent
	Amount        float64 `json:"amount" bson:"amount"`               // Amount being spent
	Signature     string  `json:"signature" bson:"signature"`         // Digital signature proving ownership
}

// UTXOOutput represents a new output being created in a transaction
type UTXOOutput struct {
	WalletID  string  `json:"walletId" bson:"walletId"`   // Recipient's wallet ID
	Amount    float64 `json:"amount" bson:"amount"`       // Amount being sent
	PublicKey string  `json:"publicKey" bson:"publicKey"` // Recipient's public key
}

// UTXOSet represents a collection of UTXOs for balance calculation
type UTXOSet struct {
	WalletID     string  `json:"walletId"`
	TotalBalance float64 `json:"totalBalance"`
	UTXOs        []UTXO  `json:"utxos"`
	UTXOCount    int     `json:"utxoCount"`
}

// BalanceResponse is returned when querying a wallet's balance
type BalanceResponse struct {
	WalletID         string  `json:"walletId"`
	Balance          float64 `json:"balance"`
	ConfirmedBalance float64 `json:"confirmedBalance"`
	PendingBalance   float64 `json:"pendingBalance"`
	UTXOCount        int     `json:"utxoCount"`
}

// CoinbaseUTXO creates a new UTXO for mining rewards or initial distribution
type CoinbaseRequest struct {
	WalletID string  `json:"walletId" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Reason   string  `json:"reason"` // "mining_reward", "initial_distribution", etc.
}
