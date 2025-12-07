package controllers

import (
	"context"
	"crypto-wallet-backend/database"
	"crypto-wallet-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getActivityLogCollection() *mongo.Collection {
	return database.GetCollection("activity_logs")
}

// LogActivity creates an activity log entry
func LogActivity(userID primitive.ObjectID, walletID string, activity models.ActivityType, description string, details map[string]interface{}, status string, c *gin.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ipAddress := ""
	userAgent := ""
	if c != nil {
		ipAddress = c.ClientIP()
		userAgent = c.GetHeader("User-Agent")
	}

	log := models.ActivityLog{
		UserID:      userID,
		WalletID:    walletID,
		Activity:    activity,
		Description: description,
		Details:     details,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Status:      status,
		Timestamp:   time.Now(),
	}

	_, err := getActivityLogCollection().InsertOne(ctx, log)
	return err
}

// GetActivityLogs returns activity logs for the authenticated user
func GetActivityLogs(c *gin.Context) {
	userID := c.GetString("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	activityType := c.Query("activity")
	status := c.Query("status")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	skip := (page - 1) * limit

	// Build filter
	filter := bson.M{"userId": objID}

	if activityType != "" {
		filter["activity"] = activityType
	}
	if status != "" {
		filter["status"] = status
	}

	// Date range filter
	if startDateStr != "" || endDateStr != "" {
		dateFilter := bson.M{}
		if startDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err == nil {
				dateFilter["$gte"] = startDate
			}
		}
		if endDateStr != "" {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err == nil {
				// Add 1 day to include the end date
				dateFilter["$lte"] = endDate.Add(24 * time.Hour)
			}
		}
		if len(dateFilter) > 0 {
			filter["timestamp"] = dateFilter
		}
	}

	// Get total count
	total, _ := getActivityLogCollection().CountDocuments(ctx, filter)

	// Get logs
	cursor, err := getActivityLogCollection().Find(ctx, filter,
		options.Find().
			SetSort(bson.M{"timestamp": -1}).
			SetSkip(int64(skip)).
			SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity logs"})
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
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GenerateTransactionReport generates a transaction report for a time period
func GenerateTransactionReport(c *gin.Context) {
	userID := c.GetString("userId")

	var req models.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Determine date range
	var startDate, endDate time.Time
	now := time.Now()

	switch req.ReportType {
	case "daily":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "weekly":
		// Start from beginning of the week (Sunday)
		weekday := int(now.Weekday())
		startDate = time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(7 * 24 * time.Hour)
	case "monthly":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
	case "yearly":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0)
	case "custom":
		if req.StartDate != "" {
			startDate, _ = time.Parse("2006-01-02", req.StartDate)
		} else {
			startDate = now.AddDate(0, -1, 0) // Default to last month
		}
		if req.EndDate != "" {
			endDate, _ = time.Parse("2006-01-02", req.EndDate)
			endDate = endDate.Add(24 * time.Hour) // Include end date
		} else {
			endDate = now
		}
	default:
		// Default to monthly
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
	}

	// Query transactions
	txFilter := bson.M{
		"timestamp": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
		"$or": []bson.M{
			{"senderWallet": wallet.WalletID},
			{"outputs.walletId": wallet.WalletID},
		},
	}

	cursor, err := getTransactionCollection().Find(ctx, txFilter)
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

	// Process transactions
	var totalSent, totalReceived, totalFees float64
	var sentTxs, receivedTxs []models.TransactionSummary

	for _, tx := range transactions {
		if tx.SenderWallet == wallet.WalletID {
			// Sent transaction
			var sentAmount float64
			var counterparty string
			for _, output := range tx.Outputs {
				if output.WalletID != wallet.WalletID {
					sentAmount += output.Amount
					counterparty = output.WalletID
				}
			}
			totalSent += sentAmount
			totalFees += tx.Fee

			sentTxs = append(sentTxs, models.TransactionSummary{
				TransactionID: tx.TransactionID,
				Type:          "sent",
				Amount:        sentAmount,
				Counterparty:  counterparty,
				Status:        string(tx.Status),
				Timestamp:     tx.Timestamp,
				BlockHash:     tx.BlockHash,
			})
		}

		// Check if received
		for _, output := range tx.Outputs {
			if output.WalletID == wallet.WalletID && tx.SenderWallet != wallet.WalletID {
				totalReceived += output.Amount
				receivedTxs = append(receivedTxs, models.TransactionSummary{
					TransactionID: tx.TransactionID,
					Type:          "received",
					Amount:        output.Amount,
					Counterparty:  tx.SenderWallet,
					Status:        string(tx.Status),
					Timestamp:     tx.Timestamp,
					BlockHash:     tx.BlockHash,
				})
			}
		}
	}

	// Get mining rewards for period
	blockFilter := bson.M{
		"minerWalletId": wallet.WalletID,
		"timestamp": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}
	blockCursor, _ := getBlockCollection().Find(ctx, blockFilter)
	defer blockCursor.Close(ctx)

	var blocksMined int
	var miningRewards float64
	for blockCursor.Next(ctx) {
		var block models.Block
		if blockCursor.Decode(&block) == nil {
			blocksMined++
			miningRewards += block.MiningReward
		}
	}

	// Get zakat payments for period
	zakatFilter := bson.M{
		"userId": objID,
		"paidAt": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}
	zakatCursor, _ := getZakatPaymentCollection().Find(ctx, zakatFilter)
	defer zakatCursor.Close(ctx)

	var zakatPaid float64
	var zakatPayments int
	for zakatCursor.Next(ctx) {
		var payment models.ZakatPayment
		if zakatCursor.Decode(&payment) == nil {
			zakatPayments++
			zakatPaid += payment.Amount
		}
	}

	report := models.TransactionReport{
		UserID:               objID,
		WalletID:             wallet.WalletID,
		ReportType:           req.ReportType,
		StartDate:            startDate,
		EndDate:              endDate,
		GeneratedAt:          now,
		TotalSent:            totalSent,
		TotalReceived:        totalReceived,
		TotalFees:            totalFees,
		NetChange:            totalReceived + miningRewards - totalSent - totalFees - zakatPaid,
		TransactionCount:     len(transactions),
		SentTransactions:     sentTxs,
		ReceivedTransactions: receivedTxs,
		BlocksMined:          blocksMined,
		MiningRewards:        miningRewards,
		ZakatPaid:            zakatPaid,
		ZakatPayments:        zakatPayments,
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GenerateWalletReport generates a comprehensive wallet report
func GenerateWalletReport(c *gin.Context) {
	userID := c.GetString("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	now := time.Now()

	// Calculate balance
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"walletId": wallet.WalletID, "isSpent": false}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}, "count": bson.M{"$sum": 1}}}},
	}

	cursor, _ := getUTXOCollection().Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	var currentBalance float64
	var unspentUTXOs int
	if cursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
			Count int     `bson:"count"`
		}
		cursor.Decode(&result)
		currentBalance = result.Total
		unspentUTXOs = result.Count
	}

	// Count total and spent UTXOs
	totalUTXOs, _ := getUTXOCollection().CountDocuments(ctx, bson.M{"walletId": wallet.WalletID})
	spentUTXOs, _ := getUTXOCollection().CountDocuments(ctx, bson.M{"walletId": wallet.WalletID, "isSpent": true})

	// Pending balance (unconfirmed UTXOs)
	pendingPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"walletId": wallet.WalletID, "isSpent": false, "isConfirmed": false}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	}
	pendingCursor, _ := getUTXOCollection().Aggregate(ctx, pendingPipeline)
	defer pendingCursor.Close(ctx)

	var pendingBalance float64
	if pendingCursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		pendingCursor.Decode(&result)
		pendingBalance = result.Total
	}

	// Transaction stats
	sentCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{"senderWallet": wallet.WalletID})
	
	// Received transactions (where wallet is in outputs but not sender)
	receivedCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{
		"outputs.walletId": wallet.WalletID,
		"senderWallet":     bson.M{"$ne": wallet.WalletID},
	})

	// Total sent amount
	sentPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"senderWallet": wallet.WalletID}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$totalOutput"}}}},
	}
	sentCursor, _ := getTransactionCollection().Aggregate(ctx, sentPipeline)
	defer sentCursor.Close(ctx)

	var totalSentAmount float64
	if sentCursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		sentCursor.Decode(&result)
		totalSentAmount = result.Total
	}

	// Mining stats
	blocksMined, _ := getBlockCollection().CountDocuments(ctx, bson.M{"minerWalletId": wallet.WalletID})
	
	miningPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"minerWalletId": wallet.WalletID}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$miningReward"}}}},
	}
	miningCursor, _ := getBlockCollection().Aggregate(ctx, miningPipeline)
	defer miningCursor.Close(ctx)

	var totalMiningRewards float64
	if miningCursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		miningCursor.Decode(&result)
		totalMiningRewards = result.Total
	}

	// Zakat stats
	nisab := models.DefaultZakatSettings.NisabInCoins
	zakatEligible := currentBalance >= nisab
	var zakatDue float64
	if zakatEligible {
		zakatDue = currentBalance * models.DefaultZakatSettings.ZakatRate
	}

	zakatPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"userId": objID, "status": "confirmed"}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	}
	zakatCursor, _ := getZakatPaymentCollection().Aggregate(ctx, zakatPipeline)
	defer zakatCursor.Close(ctx)

	var totalZakatPaid float64
	if zakatCursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		zakatCursor.Decode(&result)
		totalZakatPaid = result.Total
	}

	// Last activity
	var lastLog models.ActivityLog
	err = getActivityLogCollection().FindOne(ctx,
		bson.M{"userId": objID},
		options.FindOne().SetSort(bson.M{"timestamp": -1})).Decode(&lastLog)

	var lastActivity *time.Time
	if err == nil {
		lastActivity = &lastLog.Timestamp
	}

	report := models.WalletReport{
		UserID:              objID,
		WalletID:            wallet.WalletID,
		GeneratedAt:         now,
		CurrentBalance:      currentBalance,
		AvailableBalance:    currentBalance - pendingBalance,
		PendingBalance:      pendingBalance,
		TotalUTXOs:          int(totalUTXOs),
		SpentUTXOs:          int(spentUTXOs),
		UnspentUTXOs:        unspentUTXOs,
		TotalTransactions:   int(sentCount + receivedCount),
		SentTransactions:    int(sentCount),
		ReceivedTransactions: int(receivedCount),
		TotalSentAmount:     totalSentAmount,
		TotalReceivedAmount: currentBalance + totalSentAmount - totalMiningRewards, // Approximation
		BlocksMined:         int(blocksMined),
		TotalMiningRewards:  totalMiningRewards,
		ZakatEligible:       zakatEligible,
		ZakatDue:            zakatDue,
		TotalZakatPaid:      totalZakatPaid,
		LastActivity:        lastActivity,
		WalletCreatedAt:     wallet.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetActivityStats returns activity statistics
func GetActivityStats(c *gin.Context) {
	userID := c.GetString("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Count activities by type
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"userId": objID}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$activity", "count": bson.M{"$sum": 1}}}},
	}

	cursor, err := getActivityLogCollection().Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate stats"})
		return
	}
	defer cursor.Close(ctx)

	activityCounts := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if cursor.Decode(&result) == nil {
			activityCounts[result.ID] = result.Count
		}
	}

	// Get last 7 days activity
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	recentCount, _ := getActivityLogCollection().CountDocuments(ctx, bson.M{
		"userId":    objID,
		"timestamp": bson.M{"$gte": sevenDaysAgo},
	})

	totalCount, _ := getActivityLogCollection().CountDocuments(ctx, bson.M{"userId": objID})

	c.JSON(http.StatusOK, gin.H{
		"activityCounts":   activityCounts,
		"recentActivities": recentCount,
		"totalActivities":  totalCount,
	})
}
