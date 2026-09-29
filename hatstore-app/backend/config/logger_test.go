package config

import (
	"os"
	"testing"
)

// =============================================================================
// Testes do Logger
// =============================================================================

func TestIsDebugEnabled_True(t *testing.T) {
	os.Setenv("DEBUG", "true")
	defer os.Setenv("DEBUG", "")

	if !isDebugEnabled() {
		t.Error("isDebugEnabled deveria retornar true quando DEBUG=true")
	}
}

func TestIsDebugEnabled_False(t *testing.T) {
	os.Setenv("DEBUG", "")

	if isDebugEnabled() {
		t.Error("isDebugEnabled deveria retornar false quando DEBUG não está definido")
	}
}

func TestIsDebugEnabled_OtherValue(t *testing.T) {
	os.Setenv("DEBUG", "false")
	defer os.Setenv("DEBUG", "")

	if isDebugEnabled() {
		t.Error("isDebugEnabled deveria retornar false quando DEBUG=false")
	}
}

// Testar que as funções de log não causam panic
func TestLogFunctions_NoPanic(t *testing.T) {
	os.Setenv("DEBUG", "true")
	defer os.Setenv("DEBUG", "")

	// Nenhuma dessas deve causar panic
	LogDebug("teste debug %s", "param")
	LogDebugf("teste debugf %d", 42)
	LogInfo("teste info %s", "param")
	LogWarn("teste warn %s", "param")
	LogError("teste error %s", "param")
	LogSection("Seção de Teste")
}

func TestLogDebug_DoesNothingWhenDisabled(t *testing.T) {
	os.Setenv("DEBUG", "")

	// Não deve causar panic mesmo com debug desativado
	LogDebug("teste %s", "param")
	LogDebugf("teste %d", 42)
	LogSection("Seção de Teste")
}
