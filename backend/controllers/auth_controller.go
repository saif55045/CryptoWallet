package controllers

import (
	"context"
	"crypto-wallet-backend/database"
	"crypto-wallet-backend/models"
	"crypto-wallet-backend/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return database.GetCollection("users")
}

// Signup creates a new user account and sends OTP
func Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser models.User
	err := getUserCollection().FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check CNIC uniqueness
	err = getUserCollection().FindOne(ctx, bson.M{"cnic": req.CNIC}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this CNIC already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate OTP
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(10 * time.Minute)

	log.Printf("🔐 [Signup] Generated OTP for %s: %s (Expires at: %v)", req.Email, otp, otpExpiry)

	// Create user
	user := models.User{
		ID:           primitive.NewObjectID(),
		FullName:     req.FullName,
		Email:        req.Email,
		CNIC:         req.CNIC,
		Password:     hashedPassword,
		IsVerified:   false,
		OTP:          otp,
		OTPExpiry:    otpExpiry,
		Beneficiaries: []string{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = getUserCollection().InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send OTP email
	if err := utils.SendOTPEmail(user.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please verify your email with the OTP sent.",
		"email":   user.Email,
	})
}

// VerifyOTP verifies the OTP and activates the user account
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ [VerifyOTP] Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("🔍 [VerifyOTP] Request received - Email: %s, OTP: %s", req.Email, req.OTP)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user
	var user models.User
	err := getUserCollection().FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		log.Printf("❌ [VerifyOTP] User not found: %s, Error: %v", req.Email, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Printf("✅ [VerifyOTP] User found - Email: %s, IsVerified: %v", user.Email, user.IsVerified)
	log.Printf("🔑 [VerifyOTP] Stored OTP: %s, Received OTP: %s", user.OTP, req.OTP)
	log.Printf("⏰ [VerifyOTP] OTP Expiry: %v, Current Time: %v", user.OTPExpiry, time.Now())

	// Check if already verified
	if user.IsVerified {
		log.Printf("⚠️ [VerifyOTP] User already verified: %s", user.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already verified"})
		return
	}

	// Check OTP expiry
	if time.Now().After(user.OTPExpiry) {
		log.Printf("⚠️ [VerifyOTP] OTP expired for: %s", user.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired. Please request a new one."})
		return
	}

	// Verify OTP
	if user.OTP != req.OTP {
		log.Printf("❌ [VerifyOTP] Invalid OTP - Expected: %s, Got: %s", user.OTP, req.OTP)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	log.Printf("✅ [VerifyOTP] OTP matched! Verifying user: %s", user.Email)

	log.Printf("✅ [VerifyOTP] OTP matched! Verifying user: %s", user.Email)

	// Update user as verified
	update := bson.M{
		"$set": bson.M{
			"isVerified": true,
			"updatedAt":  time.Now(),
		},
		"$unset": bson.M{
			"otp":       "",
			"otpExpiry": "",
		},
	}

	_, err = getUserCollection().UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		log.Printf("❌ [VerifyOTP] Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
		return
	}

	log.Printf("🎉 [VerifyOTP] User successfully verified: %s", user.Email)

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"token":   token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"email":    user.Email,
			"fullName": user.FullName,
			"isAdmin":  user.IsAdmin,
		},
	})
}

// ResendOTP resends the OTP to the user's email
func ResendOTP(c *gin.Context) {
	var req models.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user
	var user models.User
	err := getUserCollection().FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already verified
	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already verified"})
		return
	}

	// Generate new OTP
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(10 * time.Minute)

	// Update user with new OTP
	update := bson.M{
		"$set": bson.M{
			"otp":       otp,
			"otpExpiry": otpExpiry,
			"updatedAt": time.Now(),
		},
	}

	_, err = getUserCollection().UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP"})
		return
	}

	// Send OTP email
	if err := utils.SendOTPEmail(user.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
	})
}

// Login authenticates a user
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user
	var user models.User
	err := getUserCollection().FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if verified
	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email first"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"email":    user.Email,
			"fullName": user.FullName,
			"walletId": user.WalletID,
			"isAdmin":  user.IsAdmin,
		},
	})
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	userID := c.GetString("userId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	err = getUserCollection().FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't send password and private key
	user.Password = ""
	user.PrivateKey = ""

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GoogleAuth handles Google OAuth authentication
func GoogleAuth(c *gin.Context) {
	var req models.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify Google token by calling Google's tokeninfo endpoint
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify Google token"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google token"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Google response"})
		return
	}

	var googleUser struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.Unmarshal(body, &googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Google user info"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists with this Google ID or email
	var user models.User
	err = getUserCollection().FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"googleId": googleUser.Sub},
			{"email": googleUser.Email},
		},
	}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		// Create new user
		user = models.User{
			ID:            primitive.NewObjectID(),
			FullName:      googleUser.Name,
			Email:         googleUser.Email,
			GoogleID:      googleUser.Sub,
			AuthProvider:  "google",
			IsVerified:    true, // Google accounts are pre-verified
			IsAdmin:       false,
			Beneficiaries: []string{},
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		_, err = getUserCollection().InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	} else {
		// Update existing user's Google ID if not set
		if user.GoogleID == "" {
			getUserCollection().UpdateOne(ctx,
				bson.M{"_id": user.ID},
				bson.M{"$set": bson.M{
					"googleId":     googleUser.Sub,
					"authProvider": "google",
					"isVerified":   true,
					"updatedAt":    time.Now(),
				}})
		}
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Google login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"email":    user.Email,
			"fullName": user.FullName,
			"walletId": user.WalletID,
			"isAdmin":  user.IsAdmin,
		},
		"needsWallet": user.WalletID == "",
	})
}
