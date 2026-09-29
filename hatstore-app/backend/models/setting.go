package models

// Setting armazena configurações da aplicação no banco de dados
type Setting struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value string `gorm:"not null" json:"value"`
}

func (Setting) TableName() string {
	return "settings"
}
