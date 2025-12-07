package controllers

import (
	"context"
	"crypto-wallet-backend/database"
	"crypto-wallet-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getZakatCalculationCollection() *mongo.Collection {
	return database.GetCollection("zakat_calculations")
}

func getZakatPaymentCollection() *mongo.Collection {
	return database.GetCollection("zakat_payments")
}

func getZakatRecipientCollection() *mongo.Collection {
	return database.GetCollection("zakat_recipients")
}

// GetZakatSettings returns the current zakat configuration
func GetZakatSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"settings": models.DefaultZakatSettings,
	})
}

// GetZakatSummary returns the user's zakat summary
func GetZakatSummary(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found. Please generate a wallet first."})
		return
	}

	// Calculate current balance
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"walletId": wallet.WalletID, "isSpent": false}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	}

	cursor, err := getUTXOCollection().Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate balance"})
		return
	}
	defer cursor.Close(ctx)

	var balance float64
	if cursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		cursor.Decode(&result)
		balance = result.Total
	}

	// Get last calculation
	var lastCalc models.ZakatCalculation
	err = getZakatCalculationCollection().FindOne(ctx,
		bson.M{"userId": objID},
		options.FindOne().SetSort(bson.M{"calculatedAt": -1})).Decode(&lastCalc)

	var lastCalcPtr *models.ZakatCalculation
	if err == nil {
		lastCalcPtr = &lastCalc
	}

	// Calculate total paid
	pipeline2 := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"userId": objID, "status": "confirmed"}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}, "count": bson.M{"$sum": 1}}}},
	}

	cursor2, _ := getZakatPaymentCollection().Aggregate(ctx, pipeline2)
	defer cursor2.Close(ctx)

	var totalPaid float64
	var paymentCount int
	if cursor2.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
			Count int     `bson:"count"`
		}
		cursor2.Decode(&result)
		totalPaid = result.Total
		paymentCount = result.Count
	}

	// Calculate zakat eligibility and amount
	nisab := models.DefaultZakatSettings.NisabInCoins
	isEligible := balance >= nisab
	var zakatDue float64
	if isEligible {
		zakatDue = balance * models.DefaultZakatSettings.ZakatRate
	}

	// Calculate next due date (1 lunar year from last calculation)
	var nextDueDate *time.Time
	if lastCalcPtr != nil && !lastCalcPtr.IsPaid {
		nextDueDate = &lastCalcPtr.ValidUntil
	}

	summary := models.ZakatSummary{
		CurrentBalance:  balance,
		NisabThreshold:  nisab,
		IsEligible:      isEligible,
		ZakatDue:        zakatDue,
		LastCalculation: lastCalcPtr,
		TotalPaid:       totalPaid,
		PaymentCount:    paymentCount,
		NextDueDate:     nextDueDate,
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

// CalculateZakat calculates the zakat due for the user
func CalculateZakat(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Calculate current balance
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"walletId": wallet.WalletID, "isSpent": false}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	}

	cursor, err := getUTXOCollection().Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate balance"})
		return
	}
	defer cursor.Close(ctx)

	var balance float64
	if cursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		cursor.Decode(&result)
		balance = result.Total
	}

	// Calculate zakat
	nisab := models.DefaultZakatSettings.NisabInCoins
	rate := models.DefaultZakatSettings.ZakatRate
	isEligible := balance >= nisab

	var zakatAmount float64
	if isEligible {
		zakatAmount = balance * rate
	}

	now := time.Now()
	validUntil := now.AddDate(0, 0, models.DefaultZakatSettings.LunarYearDays)

	calculation := models.ZakatCalculation{
		UserID:         objID,
		WalletID:       wallet.WalletID,
		TotalBalance:   balance,
		NisabThreshold: nisab,
		IsEligible:     isEligible,
		ZakatRate:      rate,
		ZakatAmount:    zakatAmount,
		CalculatedAt:   now,
		ValidUntil:     validUntil,
		IsPaid:         false,
	}

	result, err := getZakatCalculationCollection().InsertOne(ctx, calculation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save calculation"})
		return
	}

	calculation.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Zakat calculated successfully",
		"calculation": calculation,
	})
}

// PayZakat processes a zakat payment
func PayZakat(c *gin.Context) {
	userID := c.GetString("userId")

	var req models.ZakatPaymentRequest
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

	// Validate calculation ID if provided
	var calcID primitive.ObjectID
	if req.CalculationID != "" {
		calcID, err = primitive.ObjectIDFromHex(req.CalculationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calculation ID"})
			return
		}

		// Verify calculation exists and belongs to user
		var calc models.ZakatCalculation
		err = getZakatCalculationCollection().FindOne(ctx, bson.M{
			"_id":    calcID,
			"userId": objID,
		}).Decode(&calc)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Calculation not found"})
			return
		}

		if calc.IsPaid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This zakat calculation has already been paid"})
			return
		}
	}

	// Validate amount
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	// Validate recipient wallet
	recipientWallet := req.RecipientWallet
	if recipientWallet == "" {
		// Use default zakat fund wallet
		recipientWallet = models.DefaultZakatSettings.ZakatFundWallet
	}

	// Check if recipient wallet exists (unless it's the default fund)
	if recipientWallet != models.DefaultZakatSettings.ZakatFundWallet {
		var recipientWalletDoc models.Wallet
		err = getWalletCollection().FindOne(ctx, bson.M{"walletId": recipientWallet}).Decode(&recipientWalletDoc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Recipient wallet not found"})
			return
		}
	}

	// Check sender has sufficient balance
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"walletId": wallet.WalletID, "isSpent": false}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	}

	cursor, _ := getUTXOCollection().Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	var balance float64
	if cursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		cursor.Decode(&result)
		balance = result.Total
	}

	if balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance for zakat payment"})
		return
	}

	// Create the zakat payment transaction using existing transaction logic
	// First, get UTXOs to spend
	utxoCursor, err := getUTXOCollection().Find(ctx, bson.M{
		"walletId": wallet.WalletID,
		"isSpent":  false,
	}, options.Find().SetSort(bson.M{"amount": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch UTXOs"})
		return
	}
	defer utxoCursor.Close(ctx)

	var utxos []models.UTXO
	if err := utxoCursor.All(ctx, &utxos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse UTXOs"})
		return
	}

	// Select UTXOs to cover the amount
	var selectedUTXOs []models.UTXO
	var totalInput float64
	for _, utxo := range utxos {
		if totalInput >= req.Amount {
			break
		}
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalInput += utxo.Amount
	}

	if totalInput < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Create transaction
	now := time.Now()
	txID := primitive.NewObjectID().Hex()

	// Build inputs
	var inputs []models.SignedInput
	for _, utxo := range selectedUTXOs {
		inputs = append(inputs, models.SignedInput{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
			Amount:        utxo.Amount,
			PublicKey:     wallet.PublicKey,
			Signature:     "ZAKAT_PAYMENT_" + txID, // Internal transaction
		})
	}

	// Build outputs
	change := totalInput - req.Amount
	outputs := []models.TransactionOutput{
		{
			WalletID: recipientWallet,
			Amount:   req.Amount,
		},
	}

	if change > 0 {
		outputs = append(outputs, models.TransactionOutput{
			WalletID: wallet.WalletID,
			Amount:   change,
		})
	}

	transaction := models.Transaction{
		TransactionID: txID,
		Type:          models.TxTypeZakat,
		SenderWallet:  wallet.WalletID,
		Inputs:        inputs,
		Outputs:       outputs,
		TotalInput:    totalInput,
		TotalOutput:   req.Amount + change,
		Fee:           0, // No fee for zakat
		Status:        models.TxStatusPending,
		Message:       "Zakat Payment",
		Timestamp:     now,
	}

	// Start atomic transaction
	session, err := database.GetClient().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert transaction
		_, err := getTransactionCollection().InsertOne(sessCtx, transaction)
		if err != nil {
			return nil, err
		}

		// Mark UTXOs as spent
		for _, utxo := range selectedUTXOs {
			_, err := getUTXOCollection().UpdateOne(sessCtx,
				bson.M{"transactionId": utxo.TransactionID, "outputIndex": utxo.OutputIndex},
				bson.M{"$set": bson.M{"isSpent": true, "spentInTx": txID, "spentAt": now}})
			if err != nil {
				return nil, err
			}
		}

		// Create new UTXOs for outputs
		for idx, output := range outputs {
			newUTXO := models.UTXO{
				TransactionID: txID,
				OutputIndex:   idx,
				WalletID:      output.WalletID,
				Amount:        output.Amount,
				IsSpent:       false,
				IsConfirmed:   false,
				CreatedAt:     now,
			}

			// Get public key for recipient
			if output.WalletID == wallet.WalletID {
				newUTXO.PublicKey = wallet.PublicKey
			} else if output.WalletID != models.DefaultZakatSettings.ZakatFundWallet {
				var recipientWalletDoc models.Wallet
				getWalletCollection().FindOne(sessCtx, bson.M{"walletId": output.WalletID}).Decode(&recipientWalletDoc)
				newUTXO.PublicKey = recipientWalletDoc.PublicKey
			}

			_, err := getUTXOCollection().InsertOne(sessCtx, newUTXO)
			if err != nil {
				return nil, err
			}
		}

		// Create zakat payment record
		payment := models.ZakatPayment{
			UserID:          objID,
			WalletID:        wallet.WalletID,
			CalculationID:   calcID,
			Amount:          req.Amount,
			RecipientWallet: recipientWallet,
			TransactionID:   txID,
			Status:          "pending",
			PaidAt:          now,
		}
		_, err = getZakatPaymentCollection().InsertOne(sessCtx, payment)
		if err != nil {
			return nil, err
		}

		// Mark calculation as paid if provided
		if req.CalculationID != "" {
			_, err = getZakatCalculationCollection().UpdateOne(sessCtx,
				bson.M{"_id": calcID},
				bson.M{"$set": bson.M{"isPaid": true, "paidAt": now, "paymentTxId": txID}})
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process zakat payment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Zakat payment processed successfully! 🕌",
		"transactionId": txID,
		"amount":        req.Amount,
		"recipient":     recipientWallet,
	})
}

// GetZakatHistory returns the user's zakat payment history
func GetZakatHistory(c *gin.Context) {
	userID := c.GetString("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get calculations
	calcCursor, err := getZakatCalculationCollection().Find(ctx,
		bson.M{"userId": objID},
		options.Find().SetSort(bson.M{"calculatedAt": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch calculations"})
		return
	}
	defer calcCursor.Close(ctx)

	var calculations []models.ZakatCalculation
	if err := calcCursor.All(ctx, &calculations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse calculations"})
		return
	}

	// Get payments
	paymentCursor, err := getZakatPaymentCollection().Find(ctx,
		bson.M{"userId": objID},
		options.Find().SetSort(bson.M{"paidAt": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}
	defer paymentCursor.Close(ctx)

	var payments []models.ZakatPayment
	if err := paymentCursor.All(ctx, &payments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payments"})
		return
	}

	// Sync payment status with actual transaction status
	for i, payment := range payments {
		if payment.Status == "pending" && payment.TransactionID != "" {
			var tx models.Transaction
			err := getTransactionCollection().FindOne(ctx, bson.M{"transactionId": payment.TransactionID}).Decode(&tx)
			if err == nil && tx.Status == models.TxStatusConfirmed {
				// Update the payment status in database
				getZakatPaymentCollection().UpdateOne(ctx,
					bson.M{"transactionId": payment.TransactionID},
					bson.M{"$set": bson.M{"status": "confirmed"}})
				payments[i].Status = "confirmed"
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"calculations": calculations,
		"payments":     payments,
	})
}

// GetZakatRecipients returns verified zakat recipients
func GetZakatRecipients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getZakatRecipientCollection().Find(ctx, bson.M{"isVerified": true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipients"})
		return
	}
	defer cursor.Close(ctx)

	var recipients []models.ZakatRecipient
	if err := cursor.All(ctx, &recipients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse recipients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recipients": recipients,
		"categories": []string{
			models.ZakatCategoryPoor,
			models.ZakatCategoryNeedy,
			models.ZakatCategoryCollectors,
			models.ZakatCategoryNewMuslims,
			models.ZakatCategorySlaves,
			models.ZakatCategoryDebtors,
			models.ZakatCategoryFiSabilillah,
			models.ZakatCategoryTravelers,
		},
	})
}
