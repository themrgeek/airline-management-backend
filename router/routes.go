package router

import (
	controller "github.com/themrgeek/airline-management-backend/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/verify", controller.VerifyOTP)
		auth.POST("/login", controller.Login)
	}

	return r
}
