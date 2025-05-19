// config/database.go
package config

import (
	"fmt"
	"os"

	"github.com/themrgeek/airline-management-backend/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = db
	DB.AutoMigrate(&model.User{}, &model.OTP{})
}
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database connection")
	}
	sqlDB.Close()
}
