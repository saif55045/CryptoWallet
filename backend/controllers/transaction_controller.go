package controllers

import (
	"context"
	"crypto-wallet-backend/crypto"
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

func getTransactionCollection() *mongo.Collection {
	return database.GetCollection("transactions")
}

// CreateTransaction creates a new unsigned transaction (preview)
func CreateTransaction(c *gin.Context) {
	userID := c.GetString("userId")

	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get sender's wallet
	var senderWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&senderWallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found. Please generate a wallet first."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Validate recipient wallet exists
	var recipientWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"walletId": req.RecipientWalletID}).Decode(&recipientWallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipient wallet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Cannot send to yourself
	if senderWallet.WalletID == req.RecipientWalletID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send to your own wallet"})
		return
	}

	// Get sender's UTXOs (unspent only)
	cursor, err := getUTXOCollection().Find(ctx, bson.M{
		"walletId": senderWallet.WalletID,
		"isSpent":  false,
	}, options.Find().SetSort(bson.M{"amount": -1})) // Sort by amount descending
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

	// Calculate total available balance
	var totalAvailable float64
	for _, utxo := range utxos {
		totalAvailable += utxo.Amount
	}

	if totalAvailable < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Insufficient balance",
			"available": totalAvailable,
			"requested": req.Amount,
		})
		return
	}

	// Select UTXOs for transaction (greedy algorithm)
	var selectedUTXOs []models.UTXO
	var totalInput float64
	for _, utxo := range utxos {
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalInput += utxo.Amount
		if totalInput >= req.Amount {
			break
		}
	}

	// Calculate change
	change := totalInput - req.Amount
	fee := 0.0 // For now, no transaction fees

	// Build transaction inputs
	var inputs []models.SignedInput
	var inputDataForHash []crypto.InputData
	for _, utxo := range selectedUTXOs {
		inputs = append(inputs, models.SignedInput{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
			Amount:        utxo.Amount,
			PublicKey:     senderWallet.PublicKey,
			Signature:     "", // Will be filled after signing
		})
		inputDataForHash = append(inputDataForHash, crypto.InputData{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
			Amount:        utxo.Amount,
		})
	}

	// Build transaction outputs
	var outputs []models.TransactionOutput
	var outputDataForHash []crypto.OutputData

	// Output to recipient
	outputs = append(outputs, models.TransactionOutput{
		WalletID:  req.RecipientWalletID,
		Amount:    req.Amount,
		PublicKey: recipientWallet.PublicKey,
	})
	outputDataForHash = append(outputDataForHash, crypto.OutputData{
		WalletID: req.RecipientWalletID,
		Amount:   req.Amount,
	})

	// Change output (if any)
	if change > 0 {
		outputs = append(outputs, models.TransactionOutput{
			WalletID:  senderWallet.WalletID,
			Amount:    change,
			PublicKey: senderWallet.PublicKey,
		})
		outputDataForHash = append(outputDataForHash, crypto.OutputData{
			WalletID: senderWallet.WalletID,
			Amount:   change,
		})
	}

	// Generate transaction ID
	timestamp := time.Now()
	txID := crypto.GenerateTransactionID(senderWallet.WalletID, inputDataForHash, outputDataForHash, timestamp.Unix())

	// Generate data to sign for each input
	var dataToSign []string
	for i, input := range inputs {
		signData := crypto.CreateInputSignatureData(txID, input.TransactionID, input.OutputIndex, input.Amount)
		dataToSign = append(dataToSign, signData)
		inputs[i].Signature = "" // Placeholder for signature
	}

	// Create transaction preview
	preview := models.TransactionPreview{
		TransactionID:     txID,
		Inputs:            inputs,
		Outputs:           outputs,
		TotalInput:        totalInput,
		TotalOutput:       req.Amount + change,
		Fee:               fee,
		DataToSign:        dataToSign,
		RecipientWalletID: req.RecipientWalletID,
		Amount:            req.Amount,
		Change:            change,
	}

	c.JSON(http.StatusOK, gin.H{
		"preview": preview,
		"message": "Transaction created. Please sign each input with your private key.",
	})
}

// SignAndBroadcastTransaction signs and broadcasts a transaction
func SignAndBroadcastTransaction(c *gin.Context) {
	userID := c.GetString("userId")

	var req struct {
		TransactionID     string                    `json:"transactionId" binding:"required"`
		RecipientWalletID string                    `json:"recipientWalletId" binding:"required"`
		Amount            float64                   `json:"amount" binding:"required"`
		Signatures        []models.InputSignature   `json:"signatures" binding:"required"`
		Message           string                    `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get sender's wallet
	var senderWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&senderWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Get recipient's wallet
	var recipientWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"walletId": req.RecipientWalletID}).Decode(&recipientWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient wallet not found"})
		return
	}

	// Get sender's UTXOs
	cursor, err := getUTXOCollection().Find(ctx, bson.M{
		"walletId": senderWallet.WalletID,
		"isSpent":  false,
	}, options.Find().SetSort(bson.M{"amount": -1}))
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

	// Select UTXOs for transaction
	var selectedUTXOs []models.UTXO
	var totalInput float64
	for _, utxo := range utxos {
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalInput += utxo.Amount
		if totalInput >= req.Amount {
			break
		}
	}

	if totalInput < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	change := totalInput - req.Amount

	// Verify we have signatures for all inputs
	if len(req.Signatures) != len(selectedUTXOs) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Signature count mismatch",
			"expected": len(selectedUTXOs),
			"received": len(req.Signatures),
		})
		return
	}

	// Build and verify inputs
	var inputs []models.SignedInput
	var inputDataForHash []crypto.InputData
	var outputDataForHash []crypto.OutputData

	// Add output data for hash
	outputDataForHash = append(outputDataForHash, crypto.OutputData{
		WalletID: req.RecipientWalletID,
		Amount:   req.Amount,
	})
	if change > 0 {
		outputDataForHash = append(outputDataForHash, crypto.OutputData{
			WalletID: senderWallet.WalletID,
			Amount:   change,
		})
	}

	for _, utxo := range selectedUTXOs {
		inputDataForHash = append(inputDataForHash, crypto.InputData{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
			Amount:        utxo.Amount,
		})
	}

	// Use the provided transaction ID
	txID := req.TransactionID

	// Verify each signature
	for i, utxo := range selectedUTXOs {
		signData := crypto.CreateInputSignatureData(txID, utxo.TransactionID, utxo.OutputIndex, utxo.Amount)
		
		// Find the signature for this input
		var signature string
		for _, sig := range req.Signatures {
			if sig.InputIndex == i {
				signature = sig.Signature
				break
			}
		}

		if signature == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature for input", "inputIndex": i})
			return
		}

		// Verify the signature
		valid, err := crypto.VerifySignature(senderWallet.PublicKey, signData, signature)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to verify signature", "details": err.Error()})
			return
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature", "inputIndex": i})
			return
		}

		inputs = append(inputs, models.SignedInput{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
			Amount:        utxo.Amount,
			PublicKey:     senderWallet.PublicKey,
			Signature:     signature,
		})
	}

	// Build outputs
	var outputs []models.TransactionOutput
	outputs = append(outputs, models.TransactionOutput{
		WalletID:  req.RecipientWalletID,
		Amount:    req.Amount,
		PublicKey: recipientWallet.PublicKey,
	})
	if change > 0 {
		outputs = append(outputs, models.TransactionOutput{
			WalletID:  senderWallet.WalletID,
			Amount:    change,
			PublicKey: senderWallet.PublicKey,
		})
	}

	// Create the transaction
	timestamp := time.Now()
	transaction := models.Transaction{
		TransactionID: txID,
		Type:          models.TxTypeTransfer,
		Inputs:        inputs,
		Outputs:       outputs,
		TotalInput:    totalInput,
		TotalOutput:   req.Amount + change,
		Fee:           0,
		SenderWallet:  senderWallet.WalletID,
		Status:        models.TxStatusPending,
		Timestamp:     timestamp,
		Message:       req.Message,
	}

	// Start a session for atomic operations
	session, err := database.GetClient().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database session"})
		return
	}
	defer session.EndSession(ctx)

	// Execute transaction atomically
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Mark input UTXOs as spent
		now := time.Now()
		for _, utxo := range selectedUTXOs {
			_, err := getUTXOCollection().UpdateOne(sessCtx, bson.M{
				"transactionId": utxo.TransactionID,
				"outputIndex":   utxo.OutputIndex,
				"isSpent":       false,
			}, bson.M{
				"$set": bson.M{
					"isSpent":   true,
					"spentInTx": txID,
					"spentAt":   now,
				},
			})
			if err != nil {
				return nil, err
			}
		}

		// Create new UTXOs for outputs
		for i, output := range outputs {
			newUTXO := models.UTXO{
				TransactionID: txID,
				OutputIndex:   i,
				WalletID:      output.WalletID,
				Amount:        output.Amount,
				PublicKey:     output.PublicKey,
				IsSpent:       false,
				IsConfirmed:   false, // Will be confirmed when included in a block
				CreatedAt:     now,
			}
			_, err := getUTXOCollection().InsertOne(sessCtx, newUTXO)
			if err != nil {
				return nil, err
			}
		}

		// Save the transaction
		_, err := getTransactionCollection().InsertOne(sessCtx, transaction)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
		"message":     "Transaction broadcast successfully!",
	})
}

// GetMyTransactions gets all transactions for the authenticated user
func GetMyTransactions(c *gin.Context) {
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Find transactions where user is sender or receiver
	cursor, err := getTransactionCollection().Find(ctx, bson.M{
		"$or": []bson.M{
			{"senderWallet": wallet.WalletID},
			{"outputs.walletId": wallet.WalletID},
		},
	}, options.Find().SetSort(bson.M{"timestamp": -1}))
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

	// Convert to history items
	var history []models.TransactionHistoryItem
	for _, tx := range transactions {
		item := models.TransactionHistoryItem{
			TransactionID: tx.TransactionID,
			Type:          tx.Type,
			Fee:           tx.Fee,
			Status:        tx.Status,
			Timestamp:     tx.Timestamp,
			Message:       tx.Message,
		}

		if tx.SenderWallet == wallet.WalletID {
			// User sent this transaction
			item.Direction = "sent"
			// Find the recipient (first output that's not change)
			for _, output := range tx.Outputs {
				if output.WalletID != wallet.WalletID {
					item.Counterparty = output.WalletID
					item.Amount = output.Amount
					break
				}
			}
		} else {
			// User received this transaction
			item.Direction = "received"
			item.Counterparty = tx.SenderWallet
			// Find how much user received
			for _, output := range tx.Outputs {
				if output.WalletID == wallet.WalletID {
					item.Amount = output.Amount
					break
				}
			}
		}

		history = append(history, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": history,
		"count":        len(history),
	})
}

// GetTransaction gets a specific transaction by ID
func GetTransaction(c *gin.Context) {
	txID := c.Param("txId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var transaction models.Transaction
	err := getTransactionCollection().FindOne(ctx, bson.M{"transactionId": txID}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
	})
}

// GetTransactionStats returns transaction statistics
func GetTransactionStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Count transactions by status
	pendingCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{"status": models.TxStatusPending})
	confirmedCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{"status": models.TxStatusConfirmed})
	totalCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total":     totalCount,
			"pending":   pendingCount,
			"confirmed": confirmedCount,
		},
	})
}

// SendTransaction creates, signs (server-side), and broadcasts a transaction in one step
// This is a simpler approach for the wallet application
func SendTransaction(c *gin.Context) {
	userID := c.GetString("userId")

	var req struct {
		RecipientWalletID string  `json:"recipientWalletId" binding:"required"`
		Amount            float64 `json:"amount" binding:"required,gt=0"`
		Message           string  `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get sender's wallet
	var senderWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&senderWallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found. Please generate a wallet first."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Validate recipient wallet exists
	var recipientWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"walletId": req.RecipientWalletID}).Decode(&recipientWallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipient wallet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Cannot send to yourself
	if senderWallet.WalletID == req.RecipientWalletID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send to your own wallet"})
		return
	}

	// Get sender's UTXOs (unspent only)
	cursor, err := getUTXOCollection().Find(ctx, bson.M{
		"walletId": senderWallet.WalletID,
		"isSpent":  false,
	}, options.Find().SetSort(bson.M{"amount": -1}))
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

	// Calculate total available balance
	var totalAvailable float64
	for _, utxo := range utxos {
		totalAvailable += utxo.Amount
	}

	if totalAvailable < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Insufficient balance",
			"available": totalAvailable,
			"requested": req.Amount,
		})
		return
	}

	// Select UTXOs for transaction (greedy algorithm)
	var selectedUTXOs []models.UTXO
	var totalInput float64
	for _, utxo := range utxos {
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalInput += utxo.Amount
		if totalInput >= req.Amount {
			break
		}
	}

	change := totalInput - req.Amount

	// Decrypt sender's private key for signing
	privateKeyHex, err := crypto.DecryptPrivateKey(senderWallet.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt private key"})
		return
	}

	// Build inputs and outputs for transaction ID generation
	var inputDataForHash []crypto.InputData
	var outputDataForHash []crypto.OutputData

	for _, utxo := range selectedUTXOs {
		inputDataForHash = append(inputDataForHash, crypto.InputData{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
			Amount:        utxo.Amount,
		})
	}

	outputDataForHash = append(outputDataForHash, crypto.OutputData{
		WalletID: req.RecipientWalletID,
		Amount:   req.Amount,
	})
	if change > 0 {
		outputDataForHash = append(outputDataForHash, crypto.OutputData{
			WalletID: senderWallet.WalletID,
			Amount:   change,
		})
	}

	// Generate transaction ID
	timestamp := time.Now()
	txID := crypto.GenerateTransactionID(senderWallet.WalletID, inputDataForHash, outputDataForHash, timestamp.Unix())

	// Build and sign inputs
	var inputs []models.SignedInput
	for _, utxo := range selectedUTXOs {
		signData := crypto.CreateInputSignatureData(txID, utxo.TransactionID, utxo.OutputIndex, utxo.Amount)
		signature, err := crypto.SignData(privateKeyHex, signData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign transaction"})
			return
		}

		inputs = append(inputs, models.SignedInput{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
			Amount:        utxo.Amount,
			PublicKey:     senderWallet.PublicKey,
			Signature:     signature,
		})
	}

	// Build outputs
	var outputs []models.TransactionOutput
	outputs = append(outputs, models.TransactionOutput{
		WalletID:  req.RecipientWalletID,
		Amount:    req.Amount,
		PublicKey: recipientWallet.PublicKey,
	})
	if change > 0 {
		outputs = append(outputs, models.TransactionOutput{
			WalletID:  senderWallet.WalletID,
			Amount:    change,
			PublicKey: senderWallet.PublicKey,
		})
	}

	// Create the transaction
	transaction := models.Transaction{
		TransactionID: txID,
		Type:          models.TxTypeTransfer,
		Inputs:        inputs,
		Outputs:       outputs,
		TotalInput:    totalInput,
		TotalOutput:   req.Amount + change,
		Fee:           0,
		SenderWallet:  senderWallet.WalletID,
		Status:        models.TxStatusPending,
		Timestamp:     timestamp,
		Message:       req.Message,
	}

	// Start a session for atomic operations
	session, err := database.GetClient().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database session"})
		return
	}
	defer session.EndSession(ctx)

	// Execute transaction atomically
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Mark input UTXOs as spent
		now := time.Now()
		for _, utxo := range selectedUTXOs {
			_, err := getUTXOCollection().UpdateOne(sessCtx, bson.M{
				"transactionId": utxo.TransactionID,
				"outputIndex":   utxo.OutputIndex,
				"isSpent":       false,
			}, bson.M{
				"$set": bson.M{
					"isSpent":   true,
					"spentInTx": txID,
					"spentAt":   now,
				},
			})
			if err != nil {
				return nil, err
			}
		}

		// Create new UTXOs for outputs
		for i, output := range outputs {
			newUTXO := models.UTXO{
				TransactionID: txID,
				OutputIndex:   i,
				WalletID:      output.WalletID,
				Amount:        output.Amount,
				PublicKey:     output.PublicKey,
				IsSpent:       false,
				IsConfirmed:   false,
				CreatedAt:     now,
			}
			_, err := getUTXOCollection().InsertOne(sessCtx, newUTXO)
			if err != nil {
				return nil, err
			}
		}

		// Save the transaction
		_, err := getTransactionCollection().InsertOne(sessCtx, transaction)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
		"message":     "Transaction sent successfully!",
	})
}
