package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupBlockchainRoutes configures all blockchain-related routes
func SetupBlockchainRoutes(router *gin.Engine) {
	blockchain := router.Group("/api/blockchain")
	{
		// Public endpoints
		blockchain.GET("/stats", controllers.GetBlockchainStats)
		blockchain.GET("/blocks", controllers.GetBlocks)
		blockchain.GET("/block/:identifier", controllers.GetBlock)
		blockchain.GET("/latest", controllers.GetLatestBlock)
		blockchain.GET("/validate", controllers.ValidateBlockchain)
		blockchain.GET("/mining-status", controllers.GetMiningStatus)

		// Protected endpoints (require authentication)
		protected := blockchain.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.POST("/genesis", controllers.CreateGenesisBlock)
			protected.POST("/mine", controllers.MineBlock)
			protected.GET("/my-blocks", controllers.GetMyMinedBlocks)
		}
	}
}
