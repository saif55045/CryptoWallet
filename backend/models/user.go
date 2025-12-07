package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FullName      string             `json:"fullName" bson:"fullName" binding:"required"`
	Email         string             `json:"email" bson:"email" binding:"required,email"`
	CNIC          string             `json:"cnic" bson:"cnic"`
	Password      string             `json:"password,omitempty" bson:"password"`
	WalletID      string             `json:"walletId" bson:"walletId"`
	PublicKey     string             `json:"publicKey" bson:"publicKey"`
	PrivateKey    string             `json:"privateKey,omitempty" bson:"privateKey"` // Encrypted
	Beneficiaries []string           `json:"beneficiaries" bson:"beneficiaries"`
	IsVerified    bool               `json:"isVerified" bson:"isVerified"`
	IsAdmin       bool               `json:"isAdmin" bson:"isAdmin"`
	GoogleID      string             `json:"googleId,omitempty" bson:"googleId"`
	AuthProvider  string             `json:"authProvider" bson:"authProvider"` // "email" or "google"
	OTP           string             `json:"-" bson:"otp"`
	OTPExpiry     time.Time          `json:"-" bson:"otpExpiry"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	CNIC     string `json:"cnic" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type GoogleAuthRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}
