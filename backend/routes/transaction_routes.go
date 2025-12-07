package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupTransactionRoutes sets up routes for transaction operations
func SetupTransactionRoutes(router *gin.Engine) {
	tx := router.Group("/api/transaction")
	tx.Use(middleware.AuthRequired())
	{
		// Simple send transaction (server-side signing)
		tx.POST("/send", controllers.SendTransaction)

		// Create a transaction preview (unsigned)
		tx.POST("/create", controllers.CreateTransaction)

		// Sign and broadcast a transaction
		tx.POST("/broadcast", controllers.SignAndBroadcastTransaction)

		// Get user's transaction history
		tx.GET("/my-transactions", controllers.GetMyTransactions)

		// Get transaction stats
		tx.GET("/stats", controllers.GetTransactionStats)

		// Get a specific transaction by ID
		tx.GET("/:txId", controllers.GetTransaction)
	}
}
