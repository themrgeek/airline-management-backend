package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/themrgeek/airline-management-backend/model"
	"github.com/themrgeek/airline-management-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// Register User
func RegisterUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user model.User
	json.NewDecoder(r.Body).Decode(&user)

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// Insert into database
	_, err := db.Exec("INSERT INTO users(name, email, password, role) VALUES(?, ?, ?, ?)",
		user.Name, user.Email, user.Password, "user") // Default role="user"

	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

// Login User
func LoginUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&credentials)

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

	// Generate JWT token (you'll implement this)
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
