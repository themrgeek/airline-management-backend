package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	controllers "github.com/themrgeek/airline-management-backend/controller"
)

func InitializeRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// Simple welcome message
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `
			<!DOCTYPE html>
			<html>
			<head><title>Airline System</title></head>
			<body>
				<h1>Welcome to Airline System</h1>
				<p>Total users: %d</p>
			</body>
			</html>
		`, getTotalUsers(db))
	})

	return router
}

func getTotalUsers(db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
func SetupAuthRoutes(r *gin.Engine, db *sql.DB) {
	authController := controllers.AuthController{DB: db}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/signup", authController.Signup)
	}
}
