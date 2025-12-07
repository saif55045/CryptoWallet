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
)

func getWalletCollection() *mongo.Collection {
	return database.GetCollection("wallets")
}

func getBeneficiaryCollection() *mongo.Collection {
	return database.GetCollection("beneficiaries")
}

// GenerateWallet creates a new wallet with key pair for a user
func GenerateWallet(c *gin.Context) {
	userID := c.GetString("userId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert userID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user already has a wallet
	var existingWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&existingWallet)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet already exists for this user"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate key pair
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate key pair"})
		return
	}

	// Generate wallet ID from public key
	walletID := crypto.GenerateWalletID(keyPair.PublicKeyHex)

	// Encrypt private key before storing
	encryptedPrivateKey, err := crypto.EncryptPrivateKey(keyPair.PrivateKeyHex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt private key"})
		return
	}

	// Create wallet
	wallet := models.Wallet{
		ID:         primitive.NewObjectID(),
		UserID:     objID,
		WalletID:   walletID,
		PublicKey:  keyPair.PublicKeyHex,
		PrivateKey: encryptedPrivateKey,
		Balance:    0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert wallet into database
	_, err = getWalletCollection().InsertOne(ctx, wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// Update user with wallet ID and public key
	_, err = getUserCollection().UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"walletId":  walletID,
				"publicKey": keyPair.PublicKeyHex,
				"updatedAt": time.Now(),
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Wallet created successfully",
		"wallet": models.WalletResponse{
			ID:        wallet.ID,
			UserID:    wallet.UserID,
			WalletID:  wallet.WalletID,
			PublicKey: wallet.PublicKey,
			Balance:   wallet.Balance,
			CreatedAt: wallet.CreatedAt,
		},
	})
}

// GetWallet returns the wallet for the authenticated user
func GetWallet(c *gin.Context) {
	userID := c.GetString("userId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"wallet": models.WalletResponse{
			ID:        wallet.ID,
			UserID:    wallet.UserID,
			WalletID:  wallet.WalletID,
			PublicKey: wallet.PublicKey,
			Balance:   wallet.Balance,
			CreatedAt: wallet.CreatedAt,
		},
	})
}

// GetWalletByID returns wallet info by wallet ID (public endpoint)
func GetWalletByID(c *gin.Context) {
	walletID := c.Param("walletId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	// Return only public info
	c.JSON(http.StatusOK, gin.H{
		"walletId":  wallet.WalletID,
		"publicKey": wallet.PublicKey,
		"balance":   wallet.Balance,
	})
}

// ValidateWalletID checks if a wallet ID exists
func ValidateWalletID(c *gin.Context) {
	walletID := c.Param("walletId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wallet models.Wallet
	err := getWalletCollection().FindOne(ctx, bson.M{"walletId": walletID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"valid": false, "error": "Invalid Wallet ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "walletId": wallet.WalletID})
}

// ExportPrivateKey returns the decrypted private key (use with caution!)
func ExportPrivateKey(c *gin.Context) {
	userID := c.GetString("userId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

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

	// Decrypt private key
	privateKeyHex, err := crypto.DecryptPrivateKey(wallet.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt private key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"warning":    "Keep your private key secure! Never share it with anyone.",
		"privateKey": privateKeyHex,
	})
}

// AddBeneficiary adds a new beneficiary to the user's list
func AddBeneficiary(c *gin.Context) {
	userID := c.GetString("userId")
	
	var req models.AddBeneficiaryRequest
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

	// Validate wallet ID exists
	var wallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"walletId": req.WalletID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Wallet ID - wallet does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if beneficiary already exists
	var existingBeneficiary models.Beneficiary
	err = getBeneficiaryCollection().FindOne(ctx, bson.M{
		"userId":   objID,
		"walletId": req.WalletID,
	}).Decode(&existingBeneficiary)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Beneficiary already exists"})
		return
	}

	// Create beneficiary
	beneficiary := models.Beneficiary{
		ID:        primitive.NewObjectID(),
		UserID:    objID,
		Name:      req.Name,
		WalletID:  req.WalletID,
		CreatedAt: time.Now(),
	}

	_, err = getBeneficiaryCollection().InsertOne(ctx, beneficiary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add beneficiary"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Beneficiary added successfully",
		"beneficiary": beneficiary,
	})
}

// GetBeneficiaries returns all beneficiaries for the user
func GetBeneficiaries(c *gin.Context) {
	userID := c.GetString("userId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cursor, err := getBeneficiaryCollection().Find(ctx, bson.M{"userId": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer cursor.Close(ctx)

	var beneficiaries []models.Beneficiary
	if err := cursor.All(ctx, &beneficiaries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch beneficiaries"})
		return
	}

	if beneficiaries == nil {
		beneficiaries = []models.Beneficiary{}
	}

	c.JSON(http.StatusOK, gin.H{
		"beneficiaries": beneficiaries,
	})
}

// DeleteBeneficiary removes a beneficiary from the user's list
func DeleteBeneficiary(c *gin.Context) {
	userID := c.GetString("userId")
	beneficiaryID := c.Param("id")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	beneficiaryObjID, err := primitive.ObjectIDFromHex(beneficiaryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid beneficiary ID"})
		return
	}

	result, err := getBeneficiaryCollection().DeleteOne(ctx, bson.M{
		"_id":    beneficiaryObjID,
		"userId": userObjID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete beneficiary"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Beneficiary not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Beneficiary deleted successfully",
	})
}
