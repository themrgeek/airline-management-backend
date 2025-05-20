package models

import "gorm.io/gorm"

type OTP struct {
	gorm.Model
	UserID uint   `gorm:"not null"`
	OTP    string `gorm:"not null"`
	Expiry int64  `gorm:"not null"` // Unix timestamp
}
