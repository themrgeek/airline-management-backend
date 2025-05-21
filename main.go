package main

import (
	"fmt"
	"log"
	"os"

	"github.com/themrgeek/airline-management-backend/config"
	"github.com/themrgeek/airline-management-backend/model"
	"github.com/themrgeek/airline-management-backend/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("DB connection failed")
	}
	defer db.Close()

	// Create tables
	model.CreateUserTable(db)

	// Setup router
	handler := router.SetupRoutes(db)

	// Get certificate and key file paths from environment variables
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile == "" || keyFile == "" {
		log.Fatal("TLS certificate and key file paths must be specified in .env file")
	}

	// Start HTTPS server
	fmt.Println("Server starting on :443 with HTTPS")
	if err := router.StartServer(handler, certFile, keyFile); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
