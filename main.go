package main

import (
	"log"

	"github.com/themrgeek/airline-management-backend/config"
	model "github.com/themrgeek/airline-management-backend/model"
	"github.com/themrgeek/airline-management-backend/router"
)

func main() {
	// Initialize database connection
	config.ConnectDB()

	// Auto migrate models
	if err := config.DB.AutoMigrate(&model.User{}, &model.OTP{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup router
	r := router.SetupRouter()

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
