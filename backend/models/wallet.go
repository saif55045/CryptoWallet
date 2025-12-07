package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Wallet represents a user's cryptocurrency wallet
type Wallet struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"userId" bson:"userId"`
	WalletID        string             `json:"walletId" bson:"walletId"`          // Unique wallet address (hash of public key)
	PublicKey       string             `json:"publicKey" bson:"publicKey"`        // Public key in hex format
	PrivateKey      string             `json:"-" bson:"privateKey"`               // Encrypted private key (never sent to client)
	Balance         float64            `json:"balance" bson:"balance"`            // Cached balance (calculated from UTXOs)
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// WalletResponse is used when sending wallet info to client (no private key)
type WalletResponse struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"userId"`
	WalletID  string             `json:"walletId"`
	PublicKey string             `json:"publicKey"`
	Balance   float64            `json:"balance"`
	CreatedAt time.Time          `json:"createdAt"`
}

// Beneficiary represents a saved wallet address for quick transfers
type Beneficiary struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Name      string             `json:"name" bson:"name"`
	WalletID  string             `json:"walletId" bson:"walletId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

// AddBeneficiaryRequest is used to add a new beneficiary
type AddBeneficiaryRequest struct {
	Name     string `json:"name" binding:"required"`
	WalletID string `json:"walletId" binding:"required"`
}
