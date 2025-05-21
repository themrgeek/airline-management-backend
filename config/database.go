package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB establishes a connection to the MySQL database using environment variables
// Returns a database connection pool and error if any
func ConnectDB() (*sql.DB, error) {
	// Get database configuration from environment variables
	dbConfig := struct {
		user     string
		password string
		host     string
		port     string
		name     string
	}{
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		name:     os.Getenv("DB_NAME"),
	}

	// Validate required environment variables
	if dbConfig.user == "" || dbConfig.password == "" || dbConfig.host == "" ||
		dbConfig.port == "" || dbConfig.name == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	// Create the connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=30s",
		dbConfig.user, dbConfig.password, dbConfig.host, dbConfig.port, dbConfig.name)

	// Open database connection
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	if err = db.Ping(); err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return db, nil
}
