package model

import "time"

type Session struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	RefreshToken string    `json:"-" gorm:"size:500;uniqueIndex;not null"` // hide from API
	TokenPrefix  string    `json:"token_prefix" gorm:"size:16;index;not null"`
	DeviceName   string    `json:"device_name" gorm:"size:255"`
	IPAddress    string    `json:"ip_address" gorm:"size:100"`
	LastActive   time.Time `json:"last_active"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
