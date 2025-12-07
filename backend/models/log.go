package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityType represents the type of user activity
type ActivityType string

const (
	ActivityLogin           ActivityType = "login"
	ActivityLogout          ActivityType = "logout"
	ActivityWalletGenerate  ActivityType = "wallet_generate"
	ActivityTransactionSend ActivityType = "transaction_send"
	ActivityTransactionRecv ActivityType = "transaction_receive"
	ActivityMiningStart     ActivityType = "mining_start"
	ActivityMiningSuccess   ActivityType = "mining_success"
	ActivityZakatCalculate  ActivityType = "zakat_calculate"
	ActivityZakatPay        ActivityType = "zakat_pay"
	ActivityBeneficiaryAdd  ActivityType = "beneficiary_add"
	ActivityProfileUpdate   ActivityType = "profile_update"
	ActivityExportKey       ActivityType = "export_key"
)

// ActivityLog represents a user activity log entry
type ActivityLog struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	WalletID    string             `json:"walletId,omitempty" bson:"walletId,omitempty"`
	Activity    ActivityType       `json:"activity" bson:"activity"`
	Description string             `json:"description" bson:"description"`
	Details     map[string]interface{} `json:"details,omitempty" bson:"details,omitempty"`
	IPAddress   string             `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	UserAgent   string             `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	Status      string             `json:"status" bson:"status"` // success, failed
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
}

// TransactionReport represents a transaction report for a time period
type TransactionReport struct {
	UserID           primitive.ObjectID   `json:"userId" bson:"userId"`
	WalletID         string               `json:"walletId" bson:"walletId"`
	ReportType       string               `json:"reportType" bson:"reportType"` // daily, weekly, monthly, yearly, custom
	StartDate        time.Time            `json:"startDate" bson:"startDate"`
	EndDate          time.Time            `json:"endDate" bson:"endDate"`
	GeneratedAt      time.Time            `json:"generatedAt" bson:"generatedAt"`
	
	// Summary
	TotalSent        float64              `json:"totalSent" bson:"totalSent"`
	TotalReceived    float64              `json:"totalReceived" bson:"totalReceived"`
	TotalFees        float64              `json:"totalFees" bson:"totalFees"`
	NetChange        float64              `json:"netChange" bson:"netChange"`
	TransactionCount int                  `json:"transactionCount" bson:"transactionCount"`
	
	// Breakdown
	SentTransactions     []TransactionSummary `json:"sentTransactions" bson:"sentTransactions"`
	ReceivedTransactions []TransactionSummary `json:"receivedTransactions" bson:"receivedTransactions"`
	
	// Mining
	BlocksMined      int                  `json:"blocksMined" bson:"blocksMined"`
	MiningRewards    float64              `json:"miningRewards" bson:"miningRewards"`
	
	// Zakat
	ZakatPaid        float64              `json:"zakatPaid" bson:"zakatPaid"`
	ZakatPayments    int                  `json:"zakatPayments" bson:"zakatPayments"`
}

// TransactionSummary is a simplified transaction for reports
type TransactionSummary struct {
	TransactionID string    `json:"transactionId" bson:"transactionId"`
	Type          string    `json:"type" bson:"type"` // sent, received, mining, zakat
	Amount        float64   `json:"amount" bson:"amount"`
	Counterparty  string    `json:"counterparty" bson:"counterparty"` // Other wallet involved
	Status        string    `json:"status" bson:"status"`
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
	BlockHash     string    `json:"blockHash,omitempty" bson:"blockHash,omitempty"`
}

// WalletReport provides a comprehensive wallet overview
type WalletReport struct {
	UserID          primitive.ObjectID `json:"userId" bson:"userId"`
	WalletID        string             `json:"walletId" bson:"walletId"`
	GeneratedAt     time.Time          `json:"generatedAt" bson:"generatedAt"`
	
	// Balance Info
	CurrentBalance  float64            `json:"currentBalance" bson:"currentBalance"`
	AvailableBalance float64           `json:"availableBalance" bson:"availableBalance"`
	PendingBalance  float64            `json:"pendingBalance" bson:"pendingBalance"`
	
	// UTXO Info
	TotalUTXOs      int                `json:"totalUtxos" bson:"totalUtxos"`
	SpentUTXOs      int                `json:"spentUtxos" bson:"spentUtxos"`
	UnspentUTXOs    int                `json:"unspentUtxos" bson:"unspentUtxos"`
	
	// Transaction Stats
	TotalTransactions   int            `json:"totalTransactions" bson:"totalTransactions"`
	SentTransactions    int            `json:"sentTransactions" bson:"sentTransactions"`
	ReceivedTransactions int           `json:"receivedTransactions" bson:"receivedTransactions"`
	TotalSentAmount     float64        `json:"totalSentAmount" bson:"totalSentAmount"`
	TotalReceivedAmount float64        `json:"totalReceivedAmount" bson:"totalReceivedAmount"`
	
	// Mining Stats
	BlocksMined     int                `json:"blocksMined" bson:"blocksMined"`
	TotalMiningRewards float64         `json:"totalMiningRewards" bson:"totalMiningRewards"`
	
	// Zakat Stats
	ZakatEligible   bool               `json:"zakatEligible" bson:"zakatEligible"`
	ZakatDue        float64            `json:"zakatDue" bson:"zakatDue"`
	TotalZakatPaid  float64            `json:"totalZakatPaid" bson:"totalZakatPaid"`
	
	// Activity
	LastActivity    *time.Time         `json:"lastActivity,omitempty" bson:"lastActivity,omitempty"`
	WalletCreatedAt time.Time          `json:"walletCreatedAt" bson:"walletCreatedAt"`
}

// ReportRequest is the request body for generating reports
type ReportRequest struct {
	ReportType string `json:"reportType"` // daily, weekly, monthly, yearly, custom
	StartDate  string `json:"startDate,omitempty"`
	EndDate    string `json:"endDate,omitempty"`
}

// ActivityLogFilter for querying activity logs
type ActivityLogFilter struct {
	Activity  string `json:"activity,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Status    string `json:"status,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}
