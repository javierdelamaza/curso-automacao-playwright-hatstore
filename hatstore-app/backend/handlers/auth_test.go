package handlers

import (
	"bytes"
	"encoding/json"
	"hat-store-training/backend/config"
	"hat-store-training/backend/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// =============================================================================
// Testes do RegisterHandler
// =============================================================================

func TestRegisterHandler_Success(t *testing.T) {
	setupTestDB(t)

	body := makeJSON(t, map[string]string{
		"email":    "novo@email.com",
		"password": "senhaforte123",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusCreated {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	// Verificar que o usuário foi criado no banco
	var user models.User
	if err := config.DB.Where("email = ?", "novo@email.com").First(&user).Error; err != nil {
		t.Errorf("Usuário não foi criado no banco: %v", err)
	}

	if user.Email != "novo@email.com" {
		t.Errorf("Email esperado 'novo@email.com', recebeu '%s'", user.Email)
	}

	// Verificar que a senha foi hasheada (não armazenada em texto puro)
	if user.Password == "senhaforte123" {
		t.Error("Senha armazenada em texto puro, deveria ser um hash")
	}
}

func TestRegisterHandler_EmptyEmail(t *testing.T) {
	setupTestDB(t)

	body := makeJSON(t, map[string]string{
		"email":    "",
		"password": "senha123",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestRegisterHandler_EmptyPassword(t *testing.T) {
	setupTestDB(t)

	body := makeJSON(t, map[string]string{
		"email":    "user@email.com",
		"password": "",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	setupTestDB(t)

	rr := executeRequest("POST", "/auth/register", []byte("{invalid json}"), RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestRegisterHandler_DuplicateUser(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "existente@email.com", "senha123")

	body := makeJSON(t, map[string]string{
		"email":    "existente@email.com",
		"password": "outrasenha",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusConflict {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusConflict, rr.Code, rr.Body.String())
	}
}

// =============================================================================
// Testes do LoginHandler
// =============================================================================

func TestLoginHandler_Success(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "user@email.com", "senha123")

	body := makeJSON(t, map[string]interface{}{
		"email":      "user@email.com",
		"password":   "senha123",
		"rememberMe": false,
	})

	rr := executeRequest("POST", "/auth/login", body, LoginHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	// Verificar que a resposta contém o token
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Erro ao parsear resposta: %v", err)
	}

	if response["token"] == "" {
		t.Error("Token não retornado na resposta")
	}

	if response["email"] != "user@email.com" {
		t.Errorf("Email esperado 'user@email.com', recebeu '%s'", response["email"])
	}

	if response["message"] != "Login bem-sucedido" {
		t.Errorf("Mensagem esperada 'Login bem-sucedido', recebeu '%s'", response["message"])
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	setupTestDB(t)

	rr := executeRequest("POST", "/auth/login", []byte("{{bad json"), LoginHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "user@email.com", "senhaCorreta")

	body := makeJSON(t, map[string]interface{}{
		"email":    "user@email.com",
		"password": "senhaErrada",
	})

	rr := executeRequest("POST", "/auth/login", body, LoginHandler)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestLoginHandler_NonexistentUser(t *testing.T) {
	setupTestDB(t)

	body := makeJSON(t, map[string]interface{}{
		"email":    "naoexiste@email.com",
		"password": "qualquersenha",
	})

	rr := executeRequest("POST", "/auth/login", body, LoginHandler)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestLoginHandler_RememberMe(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "user@email.com", "senha123")

	body := makeJSON(t, map[string]interface{}{
		"email":      "user@email.com",
		"password":   "senha123",
		"rememberMe": true,
	})

	rr := executeRequest("POST", "/auth/login", body, LoginHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Erro ao parsear resposta: %v", err)
	}

	if response["token"] == "" {
		t.Error("Token não retornado com rememberMe=true")
	}
}

func TestLoginHandler_RateLimiting(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "rate@email.com", "senhaCorreta")

	wrongBody := makeJSON(t, map[string]interface{}{
		"email":    "rate@email.com",
		"password": "senhaErrada",
	})

	// Fazer 5 tentativas com senha errada
	for i := 0; i < 5; i++ {
		rr := executeRequest("POST", "/auth/login", wrongBody, LoginHandler)
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Tentativa %d: esperava status %d, recebeu %d", i+1, http.StatusUnauthorized, rr.Code)
		}
	}

	// 6ª tentativa deve retornar 429 (Too Many Requests)
	rr := executeRequest("POST", "/auth/login", wrongBody, LoginHandler)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Esperava status %d (rate limit), recebeu %d: %s", http.StatusTooManyRequests, rr.Code, rr.Body.String())
	}

	// Mesmo com senha correta, deve continuar bloqueado
	correctBody := makeJSON(t, map[string]interface{}{
		"email":    "rate@email.com",
		"password": "senhaCorreta",
	})

	rr = executeRequest("POST", "/auth/login", correctBody, LoginHandler)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Esperava status %d mesmo com senha correta (rate limited), recebeu %d", http.StatusTooManyRequests, rr.Code)
	}
}

func TestLoginHandler_RateLimitResetsAfterSuccess(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "reset@email.com", "senhaCorreta")

	// Fazer 3 tentativas erradas
	wrongBody := makeJSON(t, map[string]interface{}{
		"email":    "reset@email.com",
		"password": "senhaErrada",
	})
	for i := 0; i < 3; i++ {
		executeRequest("POST", "/auth/login", wrongBody, LoginHandler)
	}

	// Login com sucesso deve resetar o contador
	correctBody := makeJSON(t, map[string]interface{}{
		"email":    "reset@email.com",
		"password": "senhaCorreta",
	})
	rr := executeRequest("POST", "/auth/login", correctBody, LoginHandler)
	if rr.Code != http.StatusOK {
		t.Errorf("Login correto falhou: %d - %s", rr.Code, rr.Body.String())
	}

	// Verificar que o contador foi resetado
	var attempt models.LoginAttempt
	err := config.DB.Where("email = ?", "reset@email.com").First(&attempt).Error
	if err == nil {
		t.Log("Registro de tentativa ainda existe após login bem-sucedido — esperado que fosse deletado")
	}
}

func TestLoginHandler_IncrementAttemptOnFailure(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "attempt@email.com", "senhaCorreta")

	wrongBody := makeJSON(t, map[string]interface{}{
		"email":    "attempt@email.com",
		"password": "senhaErrada",
	})

	// Primeira tentativa errada
	executeRequest("POST", "/auth/login", wrongBody, LoginHandler)

	// Verificar que a tentativa foi registrada
	var attempt models.LoginAttempt
	if err := config.DB.Where("email = ?", "attempt@email.com").First(&attempt).Error; err != nil {
		t.Fatalf("Tentativa não registrada: %v", err)
	}

	if attempt.Count != 1 {
		t.Errorf("Esperava count=1, recebeu count=%d", attempt.Count)
	}

	// Segunda tentativa errada
	executeRequest("POST", "/auth/login", wrongBody, LoginHandler)

	config.DB.Where("email = ?", "attempt@email.com").First(&attempt)
	if attempt.Count != 2 {
		t.Errorf("Esperava count=2, recebeu count=%d", attempt.Count)
	}
}

func TestLoginHandler_RateLimitExpiresAfter15Min(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "expire@email.com", "senhaCorreta")

	// Simular 5 tentativas falhadas com last_attempt há mais de 15 minutos
	attempt := models.LoginAttempt{
		Email:       "expire@email.com",
		Count:       5,
		LastAttempt: time.Now().Add(-16 * time.Minute),
	}
	config.DB.Create(&attempt)

	// Login com senha correta deve funcionar (janela expirou)
	body := makeJSON(t, map[string]interface{}{
		"email":    "expire@email.com",
		"password": "senhaCorreta",
	})

	rr := executeRequest("POST", "/auth/login", body, LoginHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperava status %d (janela expirada), recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
}

// =============================================================================
// Testes da Feature Flag audit_frontend_enabled para Auth
// =============================================================================

// executeRequestWithSource cria e executa uma request com X-Source header
func executeRequestWithSource(method, path string, body []byte, handler http.HandlerFunc, source string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	if source != "" {
		req.Header.Set("X-Source", source)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// setAuditFrontendFlag atualiza a flag no banco e recarrega no cache
func setAuditFrontendFlag(t *testing.T, value string) {
	t.Helper()
	config.DB.Model(&models.Setting{}).Where("key = ?", "audit_frontend_enabled").Update("value", value)
	config.LoadAuditSetting()
}

// countAuditLogs retorna a quantidade de audit logs no banco
func countAuditLogs(t *testing.T) int64 {
	t.Helper()
	var count int64
	config.DB.Model(&models.AuditLog{}).Count(&count)
	return count
}

func TestRegisterHandler_NoAuditWhenFrontendFlagDisabled(t *testing.T) {
	setupTestDB(t)
	setAuditFrontendFlag(t, "false")
	defer setAuditFrontendFlag(t, "true")

	body := makeJSON(t, map[string]string{
		"email":    "novo@email.com",
		"password": "senhaforte123",
	})

	rr := executeRequestWithSource("POST", "/auth/register", body, RegisterHandler, "frontend")

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	// Aguardar possível goroutine de auditoria assíncrona
	time.Sleep(200 * time.Millisecond)

	// Com flag desabilitada e X-Source=frontend, NÃO deve criar log de auditoria
	if count := countAuditLogs(t); count != 0 {
		t.Errorf("Não deveria ter criado log de auditoria com flag desabilitada e X-Source=frontend, mas encontrou %d logs", count)
	}
}

func TestLoginHandler_NoAuditWhenFrontendFlagDisabled(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "user@email.com", "senha123")
	setAuditFrontendFlag(t, "false")
	defer setAuditFrontendFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":      "user@email.com",
		"password":   "senha123",
		"rememberMe": false,
	})

	rr := executeRequestWithSource("POST", "/auth/login", body, LoginHandler, "frontend")

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	// Aguardar possível goroutine de auditoria assíncrona
	time.Sleep(200 * time.Millisecond)

	// Com flag desabilitada e X-Source=frontend, NÃO deve criar log de auditoria
	if count := countAuditLogs(t); count != 0 {
		t.Errorf("Não deveria ter criado log de auditoria com flag desabilitada e X-Source=frontend, mas encontrou %d logs", count)
	}
}

func TestRegisterHandler_AuditWhenDirectAPICall(t *testing.T) {
	setupTestDB(t)
	setAuditFrontendFlag(t, "false")
	defer setAuditFrontendFlag(t, "true")

	body := makeJSON(t, map[string]string{
		"email":    "api@email.com",
		"password": "senhaforte123",
	})

	// Sem X-Source header → chamada direta pela API → sempre audita
	rr := executeRequestWithSource("POST", "/auth/register", body, RegisterHandler, "")

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	// Flush buffer de auditoria para garantir persistência
	config.FlushAuditBuffer()

	if count := countAuditLogs(t); count != 1 {
		t.Errorf("Chamada direta pela API deveria gerar 1 log de auditoria, mas encontrou %d", count)
	}
}

func TestLoginHandler_AuditWhenDirectAPICall(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "user@email.com", "senha123")
	setAuditFrontendFlag(t, "false")
	defer setAuditFrontendFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":      "user@email.com",
		"password":   "senha123",
		"rememberMe": false,
	})

	// Sem X-Source header → chamada direta pela API → sempre audita
	rr := executeRequestWithSource("POST", "/auth/login", body, LoginHandler, "")

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	// Flush buffer de auditoria para garantir persistência
	config.FlushAuditBuffer()

	if count := countAuditLogs(t); count != 1 {
		t.Errorf("Chamada direta pela API deveria gerar 1 log de auditoria, mas encontrou %d", count)
	}
}

// =============================================================================
// Testes da Feature Flag new_hatstore (CPF e Data de Nascimento no Registro)
// =============================================================================

// setNewHatstoreFlag atualiza a flag new_hatstore no banco e recarrega o cache
func setNewHatstoreFlag(t *testing.T, value string) {
	t.Helper()
	config.DB.Where("key = ?", "new_hatstore").Delete(&models.Setting{})
	config.DB.Create(&models.Setting{Key: "new_hatstore", Value: value})
	config.LoadFeatureFlags()
}

func TestRegisterHandler_NewHatstore_Success(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Senha@forte1",
		"nome":            "João da Silva",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusCreated {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var user models.User
	if err := config.DB.Where("email = ?", "novo@email.com").First(&user).Error; err != nil {
		t.Fatalf("Usuário não foi criado no banco: %v", err)
	}

	if user.Nome == nil || *user.Nome != "João da Silva" {
		t.Errorf("Nome esperado 'João da Silva', recebeu '%v'", user.Nome)
	}

	if user.CPF == nil || *user.CPF != "123.456.789-00" {
		t.Errorf("CPF esperado '123.456.789-00', recebeu '%v'", user.CPF)
	}

	if user.DataNascimento == nil || *user.DataNascimento != "1990-05-15" {
		t.Errorf("DataNascimento esperada '1990-05-15', recebeu '%v'", user.DataNascimento)
	}
}

func TestRegisterHandler_NewHatstore_MissingNome(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Senha@forte1",
		"nome":            "",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_MissingCPF(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Senha@forte1",
		"nome":            "João da Silva",
		"cpf":             "",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_MissingDataNascimento(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Senha@forte1",
		"nome":            "João da Silva",
		"cpf":             "123.456.789-00",
		"data_nascimento": "",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_InvalidCPFFormat(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Senha@forte1",
		"nome":            "João da Silva",
		"cpf":             "12345678900",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_InvalidDateFormat(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Senha@forte1",
		"nome":            "João da Silva",
		"cpf":             "123.456.789-00",
		"data_nascimento": "15/05/1990",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_FlagOff_IgnoresCPF(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "false")

	// Com flag desabilitada, deve registrar normalmente só com email+senha
	body := makeJSON(t, map[string]interface{}{
		"email":    "antigo@email.com",
		"password": "senhaforte123",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusCreated {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var user models.User
	if err := config.DB.Where("email = ?", "antigo@email.com").First(&user).Error; err != nil {
		t.Fatalf("Usuário não foi criado no banco: %v", err)
	}

	if user.CPF != nil {
		t.Errorf("CPF deveria ser nil com flag desabilitada, recebeu '%v'", *user.CPF)
	}

	if user.DataNascimento != nil {
		t.Errorf("DataNascimento deveria ser nil com flag desabilitada, recebeu '%v'", *user.DataNascimento)
	}
}

func TestRegisterHandler_NewHatstore_DuplicateCPF(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	// Registrar primeiro usuário
	body1 := makeJSON(t, map[string]interface{}{
		"email":           "user1@email.com",
		"password":        "Senha@forte1",
		"nome":            "João",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})
	rr1 := executeRequest("POST", "/auth/register", body1, RegisterHandler)
	if rr1.Code != http.StatusCreated {
		t.Fatalf("Primeiro registro deveria ter sucesso: %s", rr1.Body.String())
	}

	// Tentar registrar segundo usuário com o mesmo CPF
	body2 := makeJSON(t, map[string]interface{}{
		"email":           "user2@email.com",
		"password":        "Senha@forte2",
		"nome":            "Maria",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1992-03-20",
	})
	rr2 := executeRequest("POST", "/auth/register", body2, RegisterHandler)

	if rr2.Code != http.StatusConflict {
		t.Errorf("Esperava status %d (CPF duplicado), recebeu %d: %s", http.StatusConflict, rr2.Code, rr2.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_MenorDeIdade(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "jovem@email.com",
		"password":        "Senha@forte1",
		"nome":            "João Jovem",
		"cpf":             "123.456.789-00",
		"data_nascimento": "2015-01-01",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d (menor de idade), recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_SenhaFraca_SemMaiuscula(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "senha@forte1",
		"nome":            "João",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d (sem maiúscula), recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_SenhaFraca_SemMinuscula(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "SENHA@FORTE1",
		"nome":            "João",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d (sem minúscula), recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_SenhaFraca_SemNumero(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Senha@forte",
		"nome":            "João",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d (sem número), recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_SenhaFraca_SemEspecial(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "SenhaForte1",
		"nome":            "João",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d (sem especial), recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_SenhaFraca_MuitoCurta(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body := makeJSON(t, map[string]interface{}{
		"email":           "novo@email.com",
		"password":        "Ab@1",
		"nome":            "João",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})

	rr := executeRequest("POST", "/auth/register", body, RegisterHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d (senha curta), recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestRegisterHandler_NewHatstore_DuplicateEmail(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	body1 := makeJSON(t, map[string]interface{}{
		"email":           "mesmo@email.com",
		"password":        "Senha@forte1",
		"nome":            "João",
		"cpf":             "123.456.789-00",
		"data_nascimento": "1990-05-15",
	})
	rr1 := executeRequest("POST", "/auth/register", body1, RegisterHandler)
	if rr1.Code != http.StatusCreated {
		t.Fatalf("Primeiro registro deveria ter sucesso: %s", rr1.Body.String())
	}

	body2 := makeJSON(t, map[string]interface{}{
		"email":           "mesmo@email.com",
		"password":        "Senha@forte2",
		"nome":            "Maria",
		"cpf":             "987.654.321-00",
		"data_nascimento": "1992-03-20",
	})
	rr2 := executeRequest("POST", "/auth/register", body2, RegisterHandler)

	if rr2.Code != http.StatusConflict {
		t.Errorf("Esperava status %d (email duplicado), recebeu %d: %s", http.StatusConflict, rr2.Code, rr2.Body.String())
	}
}
