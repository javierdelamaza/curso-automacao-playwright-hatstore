package models

import "time"

type Pedido struct {
	ID        int          `gorm:"primaryKey" json:"id"`
	Nome      string       `gorm:"not null" json:"nome"`
	CPF       string       `gorm:"not null;index" json:"cpf"`
	Email     string       `gorm:"not null" json:"email"`
	Telefone  string       `gorm:"not null" json:"telefone"`
	Endereco  string       `gorm:"not null" json:"endereco"`
	Numero    string       `gorm:"not null" json:"numero"`
	Bairro    string       `gorm:"not null" json:"bairro"`
	CEP       string       `gorm:"not null" json:"cep"`
	Cidade    string       `gorm:"not null" json:"cidade"`
	UF        string       `gorm:"not null" json:"uf"`
	Pagamento string       `gorm:"not null" json:"pagamento"` // "pix" ou "boleto"
	Itens     []PedidoItem `gorm:"foreignKey:PedidoID;constraint:OnDelete:CASCADE" json:"itens"`
	Total     float64      `gorm:"not null" json:"total"`
	Frete     float64      `gorm:"not null;default:0" json:"frete"`
	Cupom     string       `json:"cupom"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

type PedidoItem struct {
	ID         int     `gorm:"primaryKey" json:"id"`
	PedidoID   int     `gorm:"not null;index" json:"pedido_id"`
	HatID      int     `gorm:"not null" json:"hat_id"`
	Nome       string  `gorm:"not null" json:"nome"`
	Price      float64 `gorm:"not null" json:"price"`
	Quantidade int     `gorm:"not null" json:"quantidade"`
}

type CupomHatOffUsado struct {
	CPF       string    `gorm:"primaryKey" json:"cpf"`
	Usado     bool      `gorm:"not null;default:true" json:"usado"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
