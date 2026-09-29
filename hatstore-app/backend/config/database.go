package config

import (
	"log"
	"os"
	"time"

	"hat-store-training/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

// InitDatabase estabelece conexão com PostgreSQL
func InitDatabase() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Desabilita prepared statements (evita erro com PgBouncer/Neon)
	}), &gorm.Config{})
	if err != nil {
		log.Printf("[ERROR] Erro ao conectar ao banco de dados: %v", err)
		return err
	}

	DB = database
	log.Println("[INFO] Conectado com sucesso ao PostgreSQL via Neon")

	// Configurar connection pool para otimizar uso no Neon serverless
	sqlDB, err := database.DB()
	if err != nil {
		log.Printf("[ERROR] Erro ao obter sql.DB: %v", err)
		return err
	}
	sqlDB.SetMaxOpenConns(3)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(30 * time.Second)
	log.Println("[INFO] Connection pool configurado: maxOpen=3, maxIdle=1, maxLifetime=5m, maxIdleTime=30s")

	// Auto-migrar models apenas fora de produção
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := DB.AutoMigrate(
			&models.User{},
			&models.Hat{},
			&models.Pedido{},
			&models.PedidoItem{},
			&models.CupomHatOffUsado{},
			&models.LoginAttempt{},
			&models.AuditLog{},
			&models.Setting{},
		); err != nil {
			log.Printf("[ERROR] Erro ao executar migrations: %v", err)
			return err
		}
		log.Println("[INFO] Migrations executadas com sucesso")
	} else {
		log.Println("[INFO] AutoMigrate ignorado em produção (usar migrations manuais)")
	}

	// Seed: garantir que as settings padrão existem (batch único com ON CONFLICT)
	defaultSettings := []models.Setting{
		{Key: "audit_frontend_enabled", Value: "true"},
		{Key: "filters_enabled", Value: "true"},
		{Key: "platform_enabled", Value: "true"},
	}
	DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&defaultSettings)
	log.Println("[INFO] Settings padrão verificadas (1 query batch)")

	return nil
}

// GetDB retorna a instância global do banco de dados
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database não foi inicializado. Chame InitDatabase() primeiro.")
	}
	return DB
}
