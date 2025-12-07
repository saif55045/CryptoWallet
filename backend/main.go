package main

import (
	"crypto-wallet-backend/config"
	"crypto-wallet-backend/database"
	"crypto-wallet-backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Disconnect()

	// Set Gin mode
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(config.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
	})

	// Setup routes
	routes.SetupAuthRoutes(router)
	routes.SetupWalletRoutes(router)
	routes.SetupUTXORoutes(router)
	routes.SetupTransactionRoutes(router)
	routes.SetupBlockchainRoutes(router)
	routes.SetupZakatRoutes(router)
	routes.SetupLogRoutes(router)
	routes.SetupAdminRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	// Bind to 0.0.0.0 for Render/Docker compatibility
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
