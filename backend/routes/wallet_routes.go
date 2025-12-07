package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupWalletRoutes(router *gin.Engine) {
	wallet := router.Group("/api/wallet")
	{
		// Public routes
		wallet.GET("/validate/:walletId", controllers.ValidateWalletID)
		wallet.GET("/info/:walletId", controllers.GetWalletByID)

		// Protected routes
		wallet.POST("/generate", middleware.AuthRequired(), controllers.GenerateWallet)
		wallet.GET("/my-wallet", middleware.AuthRequired(), controllers.GetWallet)
		wallet.GET("/export-key", middleware.AuthRequired(), controllers.ExportPrivateKey)

		// Beneficiary routes
		wallet.GET("/beneficiaries", middleware.AuthRequired(), controllers.GetBeneficiaries)
		wallet.POST("/beneficiaries", middleware.AuthRequired(), controllers.AddBeneficiary)
		wallet.DELETE("/beneficiaries/:id", middleware.AuthRequired(), controllers.DeleteBeneficiary)
	}
}
