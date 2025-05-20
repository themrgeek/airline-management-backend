package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt NullTime  `json:"deleted_at,omitempty"`
}

type UserAuth struct {
	ID           int      `json:"-"`
	UserID       int      `json:"-"`
	PasswordHash string   `json:"-"`
	LastLogin    NullTime `json:"last_login,omitempty"`
}

type NullTime struct {
	sql.NullTime
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserByEmail retrieves a user with auth data by email
func GetUserByEmail(db *sql.DB, email string) (*User, *UserAuth, error) {
	var user User
	var auth UserAuth

	query := `
		SELECT u.user_id, u.email, u.created_at, u.deleted_at,
		       ua.auth_id, ua.password_hash, ua.last_login
		FROM users u
		JOIN user_auth ua ON u.user_id = ua.user_id
		WHERE u.email = ? AND ua.auth_provider = 'local'`

	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.CreatedAt, &user.DeletedAt,
		&auth.ID, &auth.PasswordHash, &auth.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, err
	}

	return &user, &auth, nil
}
