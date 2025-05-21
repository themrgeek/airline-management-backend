package main

import (
	"fmt"
	"log"
	"net/http"

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
	r := router.SetupRoutes(db)

	// Start server
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
