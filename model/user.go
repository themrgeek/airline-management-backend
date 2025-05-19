// model/user.go
package model

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"not null"`
	Email      string    `gorm:"unique;not null"`
	Phone      string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	Role       string    `gorm:"not null"` // user, pilot, admin
	IsVerified bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
