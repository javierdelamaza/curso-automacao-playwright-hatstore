// backend/models/user.go
package models

import "time"

type User struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	Email          string    `gorm:"uniqueIndex;not null" json:"email"`
	Password       string    `gorm:"not null" json:"-"` // não retorna em JSON
	Nome           *string   `json:"nome,omitempty"`
	CPF            *string   `gorm:"uniqueIndex" json:"cpf,omitempty"`
	DataNascimento *string   `json:"data_nascimento,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}
