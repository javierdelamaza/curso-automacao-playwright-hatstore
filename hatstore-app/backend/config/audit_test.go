package config

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// Testes do GetClientIP
// =============================================================================

func TestGetClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.100, 10.0.0.1")

	ip := GetClientIP(req)

	if ip != "192.168.1.100" {
		t.Errorf("Esperava '192.168.1.100', recebeu '%s'", ip)
	}
}

func TestGetClientIP_XForwardedFor_Single(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.50")

	ip := GetClientIP(req)

	if ip != "10.0.0.50" {
		t.Errorf("Esperava '10.0.0.50', recebeu '%s'", ip)
	}
}

func TestGetClientIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "172.16.0.1")

	ip := GetClientIP(req)

	if ip != "172.16.0.1" {
		t.Errorf("Esperava '172.16.0.1', recebeu '%s'", ip)
	}
}

func TestGetClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.0.1:12345"

	ip := GetClientIP(req)

	if ip != "192.168.0.1" {
		t.Errorf("Esperava '192.168.0.1', recebeu '%s'", ip)
	}
}

func TestGetClientIP_Priority(t *testing.T) {
	// X-Forwarded-For deve ter prioridade sobre X-Real-IP
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "1.1.1.1")
	req.Header.Set("X-Real-IP", "2.2.2.2")
	req.RemoteAddr = "3.3.3.3:80"

	ip := GetClientIP(req)

	if ip != "1.1.1.1" {
		t.Errorf("X-Forwarded-For deveria ter prioridade, esperava '1.1.1.1', recebeu '%s'", ip)
	}
}

// =============================================================================
// Testes do ShouldAudit
// =============================================================================

func TestShouldAudit_APICallAlwaysAudits(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	// Sem header X-Source → chamada de API

	if !ShouldAudit(req) {
		t.Error("Chamadas de API (sem X-Source) devem sempre ser auditadas")
	}
}

func TestShouldAudit_FrontendEnabled(t *testing.T) {
	// Garantir que a flag está ativada
	auditFrontendEnabledLock.Lock()
	auditFrontendEnabled = true
	auditFrontendEnabledLock.Unlock()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Source", "frontend")

	if !ShouldAudit(req) {
		t.Error("Frontend deveria ser auditado quando flag está ativada")
	}
}

func TestShouldAudit_FrontendDisabled(t *testing.T) {
	// Desativar a flag
	auditFrontendEnabledLock.Lock()
	auditFrontendEnabled = false
	auditFrontendEnabledLock.Unlock()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Source", "frontend")

	if ShouldAudit(req) {
		t.Error("Frontend NÃO deveria ser auditado quando flag está desativada")
	}

	// Restaurar
	auditFrontendEnabledLock.Lock()
	auditFrontendEnabled = true
	auditFrontendEnabledLock.Unlock()
}

func TestShouldAudit_OtherSource(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Source", "mobile")

	// Qualquer coisa diferente de "frontend" deve ser auditada
	if !ShouldAudit(req) {
		t.Error("Qualquer fonte diferente de 'frontend' deveria ser auditada")
	}
}

// =============================================================================
// Testes do IsAuditFrontendEnabled
// =============================================================================

func TestIsAuditFrontendEnabled_Default(t *testing.T) {
	auditFrontendEnabledLock.Lock()
	auditFrontendEnabled = true
	auditFrontendEnabledLock.Unlock()

	if !IsAuditFrontendEnabled() {
		t.Error("Valor padrão deveria ser true")
	}
}

func TestIsAuditFrontendEnabled_Toggle(t *testing.T) {
	auditFrontendEnabledLock.Lock()
	auditFrontendEnabled = false
	auditFrontendEnabledLock.Unlock()

	if IsAuditFrontendEnabled() {
		t.Error("Deveria retornar false após ser desativado")
	}

	auditFrontendEnabledLock.Lock()
	auditFrontendEnabled = true
	auditFrontendEnabledLock.Unlock()
}

// =============================================================================
// Testes do LogAuditForRequest (verifica que ShouldAudit controla a execução)
// =============================================================================

func TestLogAuditForRequest_SkipsWhenFrontendDisabled(t *testing.T) {
	// Garantir que DB é nil para detectar se tentaria gravar
	originalDB := DB
	DB = nil
	defer func() { DB = originalDB }()

	auditFrontendEnabledLock.Lock()
	auditFrontendEnabled = false
	auditFrontendEnabledLock.Unlock()
	defer func() {
		auditFrontendEnabledLock.Lock()
		auditFrontendEnabled = true
		auditFrontendEnabledLock.Unlock()
	}()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Source", "frontend")

	// Não deve causar panic mesmo com DB nil, porque ShouldAudit retorna false
	LogAuditForRequest(req, "test_action", "test@test.com", nil, "127.0.0.1")
}

// =============================================================================
// Testes do JWTKey
// =============================================================================

func TestJWTKey_IsSet(t *testing.T) {
	if len(JWTKey) == 0 {
		t.Error("JWTKey não deveria estar vazio em ambiente de teste")
	}
}

// =============================================================================
// Testes utilitários para handler de requisições
// =============================================================================

func TestGetClientIP_EmptyHeaders(t *testing.T) {
	req := &http.Request{
		Header:     http.Header{},
		RemoteAddr: "127.0.0.1:8080",
	}

	ip := GetClientIP(req)

	if ip != "127.0.0.1" {
		t.Errorf("Esperava '127.0.0.1', recebeu '%s'", ip)
	}
}
