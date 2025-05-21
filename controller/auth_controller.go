package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/themrgeek/airline-management-backend/model"
	"github.com/themrgeek/airline-management-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// validateEmail checks if email is valid
func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePassword checks if password meets requirements
func validatePassword(password string) bool {
	return len(password) >= 8
}

// Register User
func RegisterUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if strings.TrimSpace(user.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if !validateEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if !validatePassword(user.Password) {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert into database
	_, err = db.Exec("INSERT INTO users(name, email, password, role) VALUES(?, ?, ?, ?)",
		user.Name, user.Email, user.Password, "user")

	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login User
func LoginUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if !validateEmail(credentials.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(credentials.Password) == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Get user from DB
	var user model.User
	err := db.QueryRow("SELECT id, password, role FROM users WHERE email = ?", credentials.Email).
		Scan(&user.ID, &user.Password, &user.Role)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
