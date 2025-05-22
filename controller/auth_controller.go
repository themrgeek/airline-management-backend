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
func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userID string
	userID = r.Context().Value("userID").(string)

	// Get user from DB
	var user model.User
	err := db.QueryRow("SELECT id, name, email, role FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}
func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userID string
	userID = r.Context().Value("userID").(string)

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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Update user in database
	_, err = db.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
		user.Name, user.Email, user.Password, userID)

	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}
func DeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userID string
	userID = r.Context().Value("userID").(string)

	// Delete user from database
	_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
func GetAllUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get all users from DB
	rows, err := db.Query("SELECT id, name, email, role FROM users")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}
func GetUserByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userID string
	userID = r.Context().Value("userID").(string)

	// Get user from DB
	var user model.User
	err := db.QueryRow("SELECT id, name, email, role FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}
func UpdatePassword(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userID string
	userID = r.Context().Value("userID").(string)

	var credentials struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if strings.TrimSpace(credentials.OldPassword) == "" || strings.TrimSpace(credentials.NewPassword) == "" {
		http.Error(w, "Both old and new passwords are required", http.StatusBadRequest)
		return
	}

	if !validatePassword(credentials.NewPassword) {
		http.Error(w, "New password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Get user from DB
	var user model.User
	err := db.QueryRow("SELECT password FROM users WHERE id = ?", userID).
		Scan(&user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Compare old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.OldPassword))
	if err != nil {
		http.Error(w, "Old password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Update password in database
	_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, userID)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}
