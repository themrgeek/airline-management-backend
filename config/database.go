package config

import (
	"database/sql"
	"fmt"
	"os"
)

func ConnectDB() (*sql.DB, error) {
	// 1) Loading environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	// 2) Connect to the database
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	// 3) Verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// 4) Handle errors
	if err != nil {
		return nil, err
	}
	// 5) Return the database connection
	return db, nil
}
