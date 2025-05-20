package controllers

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/themrgeek/airline-management-backend/config"
	model "github.com/themrgeek/airline-management-backend/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var twilioClient *twilio.RestClient

func InitTwilio(accountSID, authToken string) {
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateOTP() string {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

func SendOTP(phone, otp string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom("+1234567890") // Your Twilio phone number
	params.SetBody("Your OTP for Airline Management System is: " + otp)

	_, err := twilioClient.Api.CreateMessage(params)
	return err
}

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser model.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create user
	newUser := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		Phone:    user.Phone,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Generate and send OTP
	otp := GenerateOTP()
	expiry := time.Now().Add(10 * time.Minute).Unix()

	otpRecord := model.OTP{
		UserID: newUser.ID,
		OTP:    otp,
		Expiry: expiry,
	}

	if err := config.DB.Create(&otpRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create OTP record"})
		return
	}

	// In production, uncomment this to actually send OTP
	// if err := SendOTP(newUser.Phone, otp); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send OTP"})
	// 	return
	// }

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. OTP sent for verification.",
		"otp":     otp, // Remove this in production
	})
}

func VerifyOTP(c *gin.Context) {
	type VerifyRequest struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user model.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Find OTP
	var otpRecord model.OTP
	if err := config.DB.Where("user_id = ? AND otp = ?", user.ID, req.OTP).First(&otpRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if OTP is expired
	if time.Now().Unix() > otpRecord.Expiry {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
		return
	}

	// Update user as verified
	user.Verified = true
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify user"})
		return
	}

	// Delete OTP record
	config.DB.Delete(&otpRecord)

	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}

func Login(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user model.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	if !CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is verified
	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not verified. Please verify with OTP."})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key")) // Change this to a secure secret
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
