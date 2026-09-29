package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hat-store-training/backend/config"
	"hat-store-training/backend/models"

	"github.com/golang-jwt/jwt/v5"
)

// generateTestToken cria um JWT válido para testes
func generateTestToken(t *testing.T, email string) string {
	t.Helper()
	config.JWTKey = []byte("test-secret-key")
	claims := jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		t.Fatalf("Erro ao gerar token de teste: %v", err)
	}
	return tokenString
}

func TestGetProfile_Success(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "user@test.com", "senha123")

	nome := "João da Silva"
	cpf := "123.456.789-00"
	dob := "1990-05-15"
	config.DB.Model(&models.User{}).Where("email = ?", "user@test.com").Updates(map[string]interface{}{
		"nome":            nome,
		"cpf":             cpf,
		"data_nascimento": dob,
	})

	token := generateTestToken(t, "user@test.com")

	req := httptest.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	GetProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("Erro ao parsear resposta: %v", err)
	}

	if result["email"] != "user@test.com" {
		t.Errorf("Email esperado 'user@test.com', recebeu '%v'", result["email"])
	}
	if result["nome"] != "João da Silva" {
		t.Errorf("Nome esperado 'João da Silva', recebeu '%v'", result["nome"])
	}
	if result["cpf"] != "123.456.789-00" {
		t.Errorf("CPF esperado '123.456.789-00', recebeu '%v'", result["cpf"])
	}
	if result["data_nascimento"] != "1990-05-15" {
		t.Errorf("DataNascimento esperada '1990-05-15', recebeu '%v'", result["data_nascimento"])
	}
	// Senha não deve aparecer na resposta
	if _, exists := result["password"]; exists {
		t.Error("Password não deveria estar presente na resposta")
	}
}

func TestGetProfile_InvalidToken(t *testing.T) {
	setupTestDB(t)
	config.JWTKey = []byte("test-secret-key")

	req := httptest.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", "Bearer token-invalido")
	rr := httptest.NewRecorder()
	GetProfile(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestGetProfile_UserNotFound(t *testing.T) {
	setupTestDB(t)

	token := generateTestToken(t, "inexistente@test.com")

	req := httptest.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	GetProfile(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusNotFound, rr.Code)
	}
}

func TestGetProfile_NoToken(t *testing.T) {
	setupTestDB(t)
	config.JWTKey = []byte("test-secret-key")

	req := httptest.NewRequest("GET", "/api/me", nil)
	rr := httptest.NewRecorder()
	GetProfile(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusUnauthorized, rr.Code)
	}
}
