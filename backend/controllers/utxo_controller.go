package controllers

import (
	"context"
	"crypto-wallet-backend/database"
	"crypto-wallet-backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getUTXOCollection() *mongo.Collection {
	return database.GetCollection("utxos")
}

// GetBalance calculates the balance for a wallet from its UTXOs
func GetBalance(c *gin.Context) {
	walletID := c.Param("walletId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate wallet exists
	var wallet models.Wallet
	err := getWalletCollection().FindOne(ctx, bson.M{"walletId": walletID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Wallet ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get all unspent UTXOs for this wallet
	cursor, err := getUTXOCollection().Find(ctx, bson.M{
		"walletId": walletID,
		"isSpent":  false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch UTXOs"})
		return
	}
	defer cursor.Close(ctx)

	var utxos []models.UTXO
	if err := cursor.All(ctx, &utxos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse UTXOs"})
		return
	}

	// Calculate balances
	var totalBalance, confirmedBalance, pendingBalance float64
	for _, utxo := range utxos {
		totalBalance += utxo.Amount
		if utxo.IsConfirmed {
			confirmedBalance += utxo.Amount
		} else {
			pendingBalance += utxo.Amount
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": models.BalanceResponse{
			WalletID:         walletID,
			Balance:          totalBalance,
			ConfirmedBalance: confirmedBalance,
			PendingBalance:   pendingBalance,
			UTXOCount:        len(utxos),
		},
	})
}

// GetMyBalance gets the balance for the authenticated user's wallet
func GetMyBalance(c *gin.Context) {
	userID := c.GetString("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user's wallet
	var wallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found. Please generate a wallet first."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get all unspent UTXOs for this wallet
	cursor, err := getUTXOCollection().Find(ctx, bson.M{
		"walletId": wallet.WalletID,
		"isSpent":  false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch UTXOs"})
		return
	}
	defer cursor.Close(ctx)

	var utxos []models.UTXO
	if err := cursor.All(ctx, &utxos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse UTXOs"})
		return
	}

	// Calculate balances
	var totalBalance, confirmedBalance, pendingBalance float64
	for _, utxo := range utxos {
		totalBalance += utxo.Amount
		if utxo.IsConfirmed {
			confirmedBalance += utxo.Amount
		} else {
			pendingBalance += utxo.Amount
		}
	}

	// Update cached balance in wallet
	_, _ = getWalletCollection().UpdateOne(
		ctx,
		bson.M{"_id": wallet.ID},
		bson.M{"$set": bson.M{"balance": totalBalance, "updatedAt": time.Now()}},
	)

	c.JSON(http.StatusOK, gin.H{
		"balance": models.BalanceResponse{
			WalletID:         wallet.WalletID,
			Balance:          totalBalance,
			ConfirmedBalance: confirmedBalance,
			PendingBalance:   pendingBalance,
			UTXOCount:        len(utxos),
		},
	})
}

// GetUTXOs returns all UTXOs for a wallet
func GetUTXOs(c *gin.Context) {
	walletID := c.Param("walletId")
	includeSpent := c.Query("includeSpent") == "true"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{"walletId": walletID}
	if !includeSpent {
		filter["isSpent"] = false
	}

	// Get UTXOs with sorting
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := getUTXOCollection().Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch UTXOs"})
		return
	}
	defer cursor.Close(ctx)

	var utxos []models.UTXO
	if err := cursor.All(ctx, &utxos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse UTXOs"})
		return
	}

	if utxos == nil {
		utxos = []models.UTXO{}
	}

	// Calculate total
	var total float64
	for _, utxo := range utxos {
		if !utxo.IsSpent {
			total += utxo.Amount
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"utxos":        utxos,
		"count":        len(utxos),
		"totalBalance": total,
	})
}

// GetMyUTXOs returns all UTXOs for the authenticated user's wallet
func GetMyUTXOs(c *gin.Context) {
	userID := c.GetString("userId")
	includeSpent := c.Query("includeSpent") == "true"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user's wallet
	var wallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Build filter
	filter := bson.M{"walletId": wallet.WalletID}
	if !includeSpent {
		filter["isSpent"] = false
	}

	// Get UTXOs
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := getUTXOCollection().Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch UTXOs"})
		return
	}
	defer cursor.Close(ctx)

	var utxos []models.UTXO
	if err := cursor.All(ctx, &utxos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse UTXOs"})
		return
	}

	if utxos == nil {
		utxos = []models.UTXO{}
	}

	// Calculate totals
	var totalBalance, confirmedBalance float64
	for _, utxo := range utxos {
		if !utxo.IsSpent {
			totalBalance += utxo.Amount
			if utxo.IsConfirmed {
				confirmedBalance += utxo.Amount
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"utxos":            utxos,
		"count":            len(utxos),
		"totalBalance":     totalBalance,
		"confirmedBalance": confirmedBalance,
	})
}

// CreateCoinbaseUTXO creates a new UTXO (for mining rewards or initial distribution)
// This is used to add funds to the system
func CreateCoinbaseUTXO(c *gin.Context) {
	var req models.CoinbaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate wallet exists
	var wallet models.Wallet
	err := getWalletCollection().FindOne(ctx, bson.M{"walletId": req.WalletID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Wallet ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate unique transaction ID for coinbase
	txID := fmt.Sprintf("coinbase_%s", uuid.New().String())

	// Create UTXO
	utxo := models.UTXO{
		ID:            primitive.NewObjectID(),
		TransactionID: txID,
		OutputIndex:   0,
		WalletID:      req.WalletID,
		Amount:        req.Amount,
		PublicKey:     wallet.PublicKey,
		IsSpent:       false,
		IsConfirmed:   true, // Coinbase UTXOs are immediately confirmed
		CreatedAt:     time.Now(),
	}

	_, err = getUTXOCollection().InsertOne(ctx, utxo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create UTXO"})
		return
	}

	// Update wallet cached balance
	_, _ = getWalletCollection().UpdateOne(
		ctx,
		bson.M{"_id": wallet.ID},
		bson.M{
			"$inc": bson.M{"balance": req.Amount},
			"$set": bson.M{"updatedAt": time.Now()},
		},
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Coinbase UTXO created successfully",
		"utxo":    utxo,
		"reason":  req.Reason,
	})
}

// SelectUTXOsForAmount selects optimal UTXOs to cover a specific amount
// Uses a greedy algorithm to minimize the number of inputs
func SelectUTXOsForAmount(walletID string, amount float64) ([]models.UTXO, float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all unspent, confirmed UTXOs for the wallet, sorted by amount descending
	opts := options.Find().SetSort(bson.D{{Key: "amount", Value: -1}})
	cursor, err := getUTXOCollection().Find(ctx, bson.M{
		"walletId":    walletID,
		"isSpent":     false,
		"isConfirmed": true,
	}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var allUTXOs []models.UTXO
	if err := cursor.All(ctx, &allUTXOs); err != nil {
		return nil, 0, err
	}

	// Greedy selection
	var selectedUTXOs []models.UTXO
	var totalSelected float64

	for _, utxo := range allUTXOs {
		if totalSelected >= amount {
			break
		}
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalSelected += utxo.Amount
	}

	if totalSelected < amount {
		return nil, totalSelected, fmt.Errorf("insufficient balance: have %.8f, need %.8f", totalSelected, amount)
	}

	return selectedUTXOs, totalSelected, nil
}

// MarkUTXOsAsSpent marks UTXOs as spent (called during transaction processing)
func MarkUTXOsAsSpent(utxoIDs []primitive.ObjectID, spentInTx string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	_, err := getUTXOCollection().UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": utxoIDs}},
		bson.M{
			"$set": bson.M{
				"isSpent":   true,
				"spentInTx": spentInTx,
				"spentAt":   now,
			},
		},
	)

	return err
}

// CreateUTXO creates a new UTXO (internal function)
func CreateUTXO(utxo models.UTXO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := getUTXOCollection().InsertOne(ctx, utxo)
	return err
}

// GetUTXOStats returns statistics about UTXOs in the system
func GetUTXOStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Count total UTXOs
	totalCount, err := getUTXOCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count UTXOs"})
		return
	}

	// Count unspent UTXOs
	unspentCount, err := getUTXOCollection().CountDocuments(ctx, bson.M{"isSpent": false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count unspent UTXOs"})
		return
	}

	// Count spent UTXOs
	spentCount := totalCount - unspentCount

	// Calculate total value of unspent UTXOs
	pipeline := []bson.M{
		{"$match": bson.M{"isSpent": false}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	cursor, err := getUTXOCollection().Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total"})
		return
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse result"})
		return
	}

	var totalValue float64
	if len(result) > 0 {
		if val, ok := result[0]["total"].(float64); ok {
			totalValue = val
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"totalUTXOs":       totalCount,
			"unspentUTXOs":     unspentCount,
			"spentUTXOs":       spentCount,
			"totalCirculating": totalValue,
		},
	})
}
