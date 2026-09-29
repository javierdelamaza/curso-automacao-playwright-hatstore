// backend/models/login_attempt.go
package models

import "time"

type LoginAttempt struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Email       string    `gorm:"index" json:"email"`
	Count       int       `json:"count"`
	LastAttempt time.Time `json:"last_attempt"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}
