// model/otp.go
package model

import "time"

type OTP struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"not null"`
	OTPCode        string    `gorm:"not null"`
	ExpirationTime time.Time `gorm:"not null"`
}
