package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupLogRoutes configures all logging and reporting routes
func SetupLogRoutes(router *gin.Engine) {
	// Activity Logs routes
	logs := router.Group("/api/logs")
	logs.Use(middleware.AuthRequired())
	{
		logs.GET("/", controllers.GetActivityLogs)
		logs.GET("/stats", controllers.GetActivityStats)
	}

	// Reports routes
	reports := router.Group("/api/reports")
	reports.Use(middleware.AuthRequired())
	{
		reports.POST("/transactions", controllers.GenerateTransactionReport)
		reports.GET("/wallet", controllers.GenerateWalletReport)
	}
}
