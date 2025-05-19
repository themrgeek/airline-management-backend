package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/themrgeek/airline-management-backend/controller"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()

	// Add middleware here if needed
	router.Use(securityHeaders())

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", controller.RegisterUser)
		auth.POST("/verify-otp", controller.VerifyOTP)
	}

	return router
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
