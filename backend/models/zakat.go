package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ZakatCalculation represents a zakat calculation for a user
type ZakatCalculation struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"userId" bson:"userId"`
	WalletID        string             `json:"walletId" bson:"walletId"`
	TotalBalance    float64            `json:"totalBalance" bson:"totalBalance"`       // Total balance at calculation time
	NisabThreshold  float64            `json:"nisabThreshold" bson:"nisabThreshold"`   // Nisab value in coins
	IsEligible      bool               `json:"isEligible" bson:"isEligible"`           // Whether balance meets nisab
	ZakatRate       float64            `json:"zakatRate" bson:"zakatRate"`             // Usually 2.5%
	ZakatAmount     float64            `json:"zakatAmount" bson:"zakatAmount"`         // Amount of zakat due
	CalculatedAt    time.Time          `json:"calculatedAt" bson:"calculatedAt"`
	ValidUntil      time.Time          `json:"validUntil" bson:"validUntil"`           // Calculation valid for 1 lunar year
	IsPaid          bool               `json:"isPaid" bson:"isPaid"`
	PaidAt          *time.Time         `json:"paidAt,omitempty" bson:"paidAt,omitempty"`
	PaymentTxID     string             `json:"paymentTxId,omitempty" bson:"paymentTxId,omitempty"`
}

// ZakatPayment represents a zakat payment record
type ZakatPayment struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"userId" bson:"userId"`
	WalletID        string             `json:"walletId" bson:"walletId"`
	CalculationID   primitive.ObjectID `json:"calculationId" bson:"calculationId"`
	Amount          float64            `json:"amount" bson:"amount"`
	RecipientWallet string             `json:"recipientWallet" bson:"recipientWallet"` // Zakat fund wallet
	TransactionID   string             `json:"transactionId" bson:"transactionId"`
	Status          string             `json:"status" bson:"status"`                   // pending, confirmed, failed
	PaidAt          time.Time          `json:"paidAt" bson:"paidAt"`
	ConfirmedAt     *time.Time         `json:"confirmedAt,omitempty" bson:"confirmedAt,omitempty"`
	BlockHash       string             `json:"blockHash,omitempty" bson:"blockHash,omitempty"`
}

// ZakatSettings holds the zakat configuration
type ZakatSettings struct {
	NisabInGold      float64 `json:"nisabInGold"`      // Nisab threshold in grams of gold (87.48g)
	NisabInSilver    float64 `json:"nisabInSilver"`    // Nisab threshold in grams of silver (612.36g)
	NisabInCoins     float64 `json:"nisabInCoins"`     // Nisab threshold in our coins
	ZakatRate        float64 `json:"zakatRate"`        // Standard rate 2.5%
	LunarYearDays    int     `json:"lunarYearDays"`    // ~354 days
	ZakatFundWallet  string  `json:"zakatFundWallet"`  // Official zakat collection wallet
}

// Default Zakat Settings
var DefaultZakatSettings = ZakatSettings{
	NisabInGold:     87.48,    // 87.48 grams of gold
	NisabInSilver:   612.36,   // 612.36 grams of silver
	NisabInCoins:    1000.0,   // 1000 coins as nisab threshold for our system
	ZakatRate:       0.025,    // 2.5%
	LunarYearDays:   354,      // Lunar year
	ZakatFundWallet: "ZAKAT_FUND_OFFICIAL", // Will be set to actual wallet
}

// ZakatSummary provides a summary of user's zakat status
type ZakatSummary struct {
	CurrentBalance   float64            `json:"currentBalance"`
	NisabThreshold   float64            `json:"nisabThreshold"`
	IsEligible       bool               `json:"isEligible"`
	ZakatDue         float64            `json:"zakatDue"`
	LastCalculation  *ZakatCalculation  `json:"lastCalculation,omitempty"`
	TotalPaid        float64            `json:"totalPaid"`
	PaymentCount     int                `json:"paymentCount"`
	NextDueDate      *time.Time         `json:"nextDueDate,omitempty"`
}

// ZakatRecipient represents a verified zakat recipient/organization
type ZakatRecipient struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	WalletID    string             `json:"walletId" bson:"walletId"`
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category"` // poor, needy, zakat-collectors, etc.
	IsVerified  bool               `json:"isVerified" bson:"isVerified"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

// Zakat recipient categories (8 categories as per Islamic law)
const (
	ZakatCategoryPoor           = "poor"           // Al-Fuqara
	ZakatCategoryNeedy          = "needy"          // Al-Masakin
	ZakatCategoryCollectors     = "collectors"     // Al-Amilina Alayha
	ZakatCategoryNewMuslims     = "new-muslims"    // Al-Mu'allafatu Qulubuhum
	ZakatCategorySlaves         = "slaves"         // Fir-Riqab
	ZakatCategoryDebtors        = "debtors"        // Al-Gharimin
	ZakatCategoryFiSabilillah   = "fi-sabilillah"  // Fi Sabilillah
	ZakatCategoryTravelers      = "travelers"      // Ibn as-Sabil
)

// ZakatPaymentRequest is the request body for paying zakat
type ZakatPaymentRequest struct {
	CalculationID   string  `json:"calculationId"`
	Amount          float64 `json:"amount"`
	RecipientWallet string  `json:"recipientWallet"`
}
