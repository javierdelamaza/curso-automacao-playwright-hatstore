package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hat-store-training/backend/config"
	"hat-store-training/backend/models"

	"github.com/glebarez/sqlite"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// setupTestDB cria um banco SQLite em memória e executa as migrations
func setupTestDB(t *testing.T) {
	t.Helper()

	// Aguardar goroutines assíncronas de testes anteriores (ex: LogAuditAsync)
	time.Sleep(50 * time.Millisecond)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Erro ao criar banco de teste: %v", err)
	}

	// SQLite :memory: é per-connection. Limitar a 1 conexão
	// garante que todas as operações usam a mesma base de dados.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Erro ao obter sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	// Executar migrations ANTES de atribuir config.DB para evitar
	// que goroutines assíncronas acessem o DB sem tabelas
	err = db.AutoMigrate(
		&models.User{},
		&models.Hat{},
		&models.Pedido{},
		&models.PedidoItem{},
		&models.CupomHatOffUsado{},
		&models.LoginAttempt{},
		&models.AuditLog{},
		&models.Setting{},
	)
	if err != nil {
		t.Fatalf("Erro ao executar migrations: %v", err)
	}

	// Seed: setting de auditoria (necessário para LogAuditForRequest)
	db.Create(&models.Setting{Key: "audit_frontend_enabled", Value: "true"})

	// Atribuir DB global somente após setup completo
	config.DB = db

	// Invalidar cache de chapéus e descartar audit events de testes anteriores
	InvalidateHatsCache()
	config.DrainAuditBuffer()
}

// seedHats insere chapéus de teste no banco
func seedHats(t *testing.T) []models.Hat {
	t.Helper()
	hats := []models.Hat{
		{Nome: "Boné Nacional", Price: 49.90, Quantidade: 10, Categoria: "nacional"},
		{Nome: "Chapéu Importado", Price: 129.90, Quantidade: 5, Categoria: "importado"},
		{Nome: "Boné Crescer", Price: 79.90, Quantidade: 0, Categoria: "crescer"},
		{Nome: "Cartola Importada", Price: 199.90, Quantidade: 3, Categoria: "importado"},
	}
	for i := range hats {
		if err := config.DB.Create(&hats[i]).Error; err != nil {
			t.Fatalf("Erro ao seed hat: %v", err)
		}
	}
	return hats
}

// seedUser cria um usuário de teste no banco e retorna o email
func seedUser(t *testing.T, email, password string) {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("Erro ao gerar hash: %v", err)
	}
	user := models.User{Email: email, Password: string(hashed)}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("Erro ao seed user: %v", err)
	}
}

// makeJSON converte um valor para bytes JSON
func makeJSON(t *testing.T, v interface{}) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Erro ao serializar JSON: %v", err)
	}
	return data
}

// executeRequest cria e executa uma request contra um handler
func executeRequest(method, path string, body []byte, handler http.HandlerFunc) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// executeRequestWithVars executa uma request com variáveis de rota (mux.Vars)
func executeRequestWithVars(method, path string, body []byte, handler http.HandlerFunc, vars map[string]string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// validPedidoPayload retorna um payload de pedido válido para testes
func validPedidoPayload(hatID uint) map[string]interface{} {
	return map[string]interface{}{
		"nome":      "João da Silva",
		"cpf":       "123.456.789-00",
		"email":     "joao@email.com",
		"telefone":  "(11) 99999-0000",
		"endereco":  "Rua dos Testes",
		"numero":    "123",
		"bairro":    "Centro",
		"cep":       "01234-567",
		"cidade":    "São Paulo",
		"uf":        "SP",
		"pagamento": "pix",
		"itens": []map[string]interface{}{
			{
				"id":         hatID,
				"nome":       "Boné Nacional",
				"price":      49.90,
				"quantidade": 1,
			},
		},
	}
}
