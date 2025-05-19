package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/themrgeek/airline-management-backend/config"
	"github.com/themrgeek/airline-management-backend/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check existing user
	var existing model.User
	if config.DB.Where("email = ? OR phone = ?", user.Email, user.Phone).First(&existing).RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	// Generate OTP
	otp := generateOTP()
	expiration := time.Now().Add(5 * time.Minute)

	// Save OTP
	otpRecord := model.OTP{
		UserID:         user.ID,
		OTPCode:        otp,
		ExpirationTime: expiration,
	}
	if err := config.DB.Create(&otpRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP creation failed"})
		return
	}

	// Send OTP via SMS (mock implementation)
	sendOTP(user.Phone, otp)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered. Check phone for OTP"})
}

func VerifyOTP(c *gin.Context) {
	type VerifyRequest struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by phone
	var user model.User
	if config.DB.Where("phone = ?", req.Phone).First(&user).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validate OTP
	var otp model.OTP
	if config.DB.Where("user_id = ? AND otp_code = ? AND expiration_time > ?",
		user.ID, req.OTP, time.Now()).First(&otp).RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Mark user as verified
	user.IsVerified = true
	config.DB.Save(&user)

	// Delete OTP
	config.DB.Delete(&otp)

	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func sendOTP(phone, otp string) {
	// Implement actual SMS sending logic here
	// This is a mock implementation
	println("Sending OTP", otp, "to", phone)
}
