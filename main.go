package main

import (
	"log"

	"github.com/themrgeek/airline-management-backend/config"
	controller "github.com/themrgeek/airline-management-backend/controller"
	model "github.com/themrgeek/airline-management-backend/model"
	"github.com/themrgeek/airline-management-backend/router"
)

func main() {
	// Initialize database
	config.ConnectDB()

	// Initialize Twilio (replace with your actual credentials)
	controller.InitTwilio("your-account-sid", "your-auth-token")

	// Migrate models
	db := config.DB
	db.AutoMigrate(&model.User{}, &model.OTP{})

	// Setup router
	r := router.SetupRouter()

	// Start server
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
