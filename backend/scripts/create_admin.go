package main

import (
	"context"
	"crypto-wallet-backend/database"
	"crypto-wallet-backend/models"
	"crypto-wallet-backend/utils"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Disconnect()

	// Admin user details
	fmt.Println("=== Create Admin User ===")
	
	email := "admin@cryptovault.com"
	password := "Admin@123"
	fullName := "System Admin"
	cnic := "00000-0000000-0"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := database.GetCollection("users")

	// Check if admin already exists
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		// User exists, just make them admin
		fmt.Println("User already exists. Updating to admin...")
		_, err = userCollection.UpdateOne(ctx,
			bson.M{"email": email},
			bson.M{"$set": bson.M{
				"isAdmin":    true,
				"isVerified": true,
				"updatedAt":  time.Now(),
			}})
		if err != nil {
			log.Fatal("Failed to update user:", err)
		}
		fmt.Println("✅ User updated to admin successfully!")
		fmt.Printf("\nLogin credentials:\n")
		fmt.Printf("Email: %s\n", email)
		fmt.Printf("Password: %s\n", password)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Create admin user
	adminUser := models.User{
		ID:            primitive.NewObjectID(),
		FullName:      fullName,
		Email:         email,
		CNIC:          cnic,
		Password:      hashedPassword,
		IsVerified:    true,
		IsAdmin:       true,
		AuthProvider:  "email",
		Beneficiaries: []string{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = userCollection.InsertOne(ctx, adminUser)
	if err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	fmt.Println("✅ Admin user created successfully!")
	fmt.Printf("\nLogin credentials:\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
	fmt.Println("\nNote: You'll need to generate a wallet after logging in.")
}
