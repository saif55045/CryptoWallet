package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUTXORoutes(router *gin.Engine) {
	utxo := router.Group("/api/utxo")
	{
		// Public routes
		utxo.GET("/balance/:walletId", controllers.GetBalance)
		utxo.GET("/list/:walletId", controllers.GetUTXOs)
		utxo.GET("/stats", controllers.GetUTXOStats)

		// Protected routes
		utxo.GET("/my-balance", middleware.AuthRequired(), controllers.GetMyBalance)
		utxo.GET("/my-utxos", middleware.AuthRequired(), controllers.GetMyUTXOs)
		
		// Admin/System routes (for adding funds to the system)
		utxo.POST("/coinbase", middleware.AuthRequired(), controllers.CreateCoinbaseUTXO)
	}
}
