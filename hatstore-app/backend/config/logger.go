package config

import (
	"log"
	"os"
)

// isDebugEnabled verifica se debug está habilitado
func isDebugEnabled() bool {
	return os.Getenv("DEBUG") == "true"
}

// LogDebug imprime logs de debug apenas se DEBUG=true
func LogDebug(message string, args ...interface{}) {
	if isDebugEnabled() {
		log.Printf("[DEBUG] "+message, args...)
	}
}

// LogDebugf imprime logs formatados de debug apenas se DEBUG=true
func LogDebugf(format string, args ...interface{}) {
	if isDebugEnabled() {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// LogInfo imprime logs de informação (sempre)
func LogInfo(message string, args ...interface{}) {
	log.Printf("[INFO] "+message, args...)
}

// LogWarn imprime logs de aviso (sempre)
func LogWarn(message string, args ...interface{}) {
	log.Printf("[WARN] "+message, args...)
}

// LogError imprime logs de erro (sempre)
func LogError(message string, args ...interface{}) {
	log.Printf("[ERROR] "+message, args...)
}

// LogSection imprime um separador de seção de debug
func LogSection(title string) {
	if isDebugEnabled() {
		log.Printf("[DEBUG] ========== %s ==========", title)
	}
}
