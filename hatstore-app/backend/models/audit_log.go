package models

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Action    string         `gorm:"not null;index" json:"action"`
	UserEmail string         `gorm:"index" json:"user_email"`
	Details   datatypes.JSON `gorm:"type:jsonb" json:"details"`
	IPAddress string         `json:"ip_address"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica o nome da tabela
func (AuditLog) TableName() string {
	return "audit_logs"
}

// LogDetails armazena os detalhes do evento em JSON
type LogDetails struct {
	CPF       string `json:"cpf,omitempty"`
	PedidoID  int    `json:"pedido_id,omitempty"`
	Total     string `json:"total,omitempty"`
	Message   string `json:"message,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
}
