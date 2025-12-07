package routes

import (
	"crypto-wallet-backend/controllers"
	"crypto-wallet-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.Engine) {
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.GET("/stats", controllers.GetSystemStats)
		admin.GET("/users", controllers.GetAllUsers)
		admin.GET("/transactions", controllers.GetAllTransactions)
		admin.GET("/blocks", controllers.GetAllBlocks)
		admin.GET("/logs", controllers.GetSystemLogs)
		admin.PUT("/users/:userId/toggle-admin", controllers.ToggleUserAdmin)
		admin.DELETE("/users/:userId", controllers.DeleteUser)
	}
}
