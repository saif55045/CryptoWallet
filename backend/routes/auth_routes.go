package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/verify-otp", controllers.VerifyOTP)
		auth.POST("/resend-otp", controllers.ResendOTP)
		auth.POST("/login", controllers.Login)
		auth.POST("/google", controllers.GoogleAuth)
		
		// Protected routes
		auth.GET("/profile", middleware.AuthRequired(), controllers.GetProfile)
	}
}
