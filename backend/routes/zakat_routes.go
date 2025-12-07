package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupZakatRoutes configures all zakat-related routes
func SetupZakatRoutes(router *gin.Engine) {
	zakat := router.Group("/api/zakat")
	{
		// Public endpoints
		zakat.GET("/settings", controllers.GetZakatSettings)
		zakat.GET("/recipients", controllers.GetZakatRecipients)

		// Protected endpoints (require authentication)
		protected := zakat.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/summary", controllers.GetZakatSummary)
			protected.POST("/calculate", controllers.CalculateZakat)
			protected.POST("/pay", controllers.PayZakat)
			protected.GET("/history", controllers.GetZakatHistory)
		}
	}
}
