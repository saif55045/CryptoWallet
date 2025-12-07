package controllers

import (
	"context"
	"crypto-wallet-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllUsers returns all users (admin only)
func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getUserCollection().Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse users"})
		return
	}

	// Remove sensitive data
	for i := range users {
		users[i].Password = ""
		users[i].PrivateKey = ""
		users[i].OTP = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GetSystemStats returns system-wide statistics (admin only)
func GetSystemStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Count users
	userCount, _ := getUserCollection().CountDocuments(ctx, bson.M{})
	verifiedUsers, _ := getUserCollection().CountDocuments(ctx, bson.M{"isVerified": true})

	// Count wallets
	walletCount, _ := getWalletCollection().CountDocuments(ctx, bson.M{})

	// Count transactions
	txCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{})
	pendingTx, _ := getTransactionCollection().CountDocuments(ctx, bson.M{"status": "pending"})

	// Count blocks
	blockCount, _ := getBlockCollection().CountDocuments(ctx, bson.M{})

	// Count UTXOs
	utxoCount, _ := getUTXOCollection().CountDocuments(ctx, bson.M{})
	unspentUtxos, _ := getUTXOCollection().CountDocuments(ctx, bson.M{"isSpent": false})

	// Calculate total balance in system
	pipeline := []bson.M{
		{"$match": bson.M{"isSpent": false}},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}
	cursor, err := getUTXOCollection().Aggregate(ctx, pipeline)
	var totalBalance float64
	if err == nil {
		var results []bson.M
		if cursor.All(ctx, &results) == nil && len(results) > 0 {
			if total, ok := results[0]["total"].(float64); ok {
				totalBalance = total
			}
		}
	}

	// Get zakat stats
	zakatCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{"type": "zakat"})

	c.JSON(http.StatusOK, gin.H{
		"users": gin.H{
			"total":    userCount,
			"verified": verifiedUsers,
		},
		"wallets": walletCount,
		"transactions": gin.H{
			"total":   txCount,
			"pending": pendingTx,
		},
		"blocks":       blockCount,
		"utxos": gin.H{
			"total":   utxoCount,
			"unspent": unspentUtxos,
		},
		"totalBalance":  totalBalance,
		"zakatPayments": zakatCount,
	})
}

// GetAllTransactions returns all transactions (admin only)
func GetAllTransactions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getTransactionCollection().Find(ctx, bson.M{}, 
		options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(100))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	})
}

// GetAllBlocks returns all blocks (admin only)
func GetAllBlocks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getBlockCollection().Find(ctx, bson.M{}, 
		options.Find().SetSort(bson.M{"index": -1}).SetLimit(100))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocks"})
		return
	}
	defer cursor.Close(ctx)

	var blocks []models.Block
	if err := cursor.All(ctx, &blocks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse blocks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blocks": blocks,
		"count":  len(blocks),
	})
}

// ToggleUserAdmin toggles admin status for a user
func ToggleUserAdmin(c *gin.Context) {
	userID := c.Param("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user
	var user models.User
	err = getUserCollection().FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Toggle admin status
	newAdminStatus := !user.IsAdmin
	_, err = getUserCollection().UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"isAdmin": newAdminStatus, "updatedAt": time.Now()}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User admin status updated",
		"isAdmin": newAdminStatus,
	})
}

// DeleteUser deletes a user (admin only)
func DeleteUser(c *gin.Context) {
	userID := c.Param("userId")
	currentUserID := c.GetString("userId")

	if userID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result, err := getUserCollection().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetSystemLogs returns system logs (admin only)
func GetSystemLogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getActivityLogCollection().Find(ctx, bson.M{},
		options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(200))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	defer cursor.Close(ctx)

	var logs []models.ActivityLog
	if err := cursor.All(ctx, &logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": len(logs),
	})
}
