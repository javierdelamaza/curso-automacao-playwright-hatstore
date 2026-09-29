package models

import "time"

type Hat struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Nome       string    `gorm:"not null" json:"nome"`
	Price      float64   `gorm:"not null" json:"price"`
	Quantidade int       `gorm:"not null;default:0" json:"quantidade"`
	Categoria  string    `gorm:"not null" json:"categoria"` // nacional, importado, crescer
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
