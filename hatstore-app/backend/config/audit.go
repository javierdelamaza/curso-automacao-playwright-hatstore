package config

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"hat-store-training/backend/models"

	"gorm.io/datatypes"
)

// Cache em memória para a flag de auditoria do frontend
var (
	auditFrontendEnabled     = true
	auditFrontendEnabledLock sync.RWMutex
)

// Buffer para batch insert de audit logs (reduz queries individuais)
var auditBufCh = make(chan models.AuditLog, 200)

// StartAuditBufferJob inicia goroutine que acumula audit logs e faz flush em batch
func StartAuditBufferJob() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		batch := make([]models.AuditLog, 0, 20)
		for {
			select {
			case entry := <-auditBufCh:
				batch = append(batch, entry)
				if len(batch) >= 20 {
					flushAuditBatch(batch)
					batch = batch[:0]
				}
			case <-ticker.C:
				if len(batch) > 0 {
					flushAuditBatch(batch)
					batch = batch[:0]
				}
			}
		}
	}()
}

func flushAuditBatch(batch []models.AuditLog) {
	if DB == nil || len(batch) == 0 {
		return
	}
	if err := DB.Create(&batch).Error; err != nil {
		log.Printf("[ERROR] Erro ao flush audit batch (%d logs): %v", len(batch), err)
	}
}

// FlushAuditBuffer drena o canal de audit e insere tudo no banco (útil para testes)
func FlushAuditBuffer() {
	var batch []models.AuditLog
	for {
		select {
		case entry := <-auditBufCh:
			batch = append(batch, entry)
		default:
			if len(batch) > 0 {
				flushAuditBatch(batch)
			}
			return
		}
	}
}

// DrainAuditBuffer descarta todos os eventos pendentes no canal (para reset entre testes)
func DrainAuditBuffer() {
	for {
		select {
		case <-auditBufCh:
		default:
			return
		}
	}
}

// LoadAuditSetting carrega a flag do banco para o cache em memória
// Delega para loadAllFlagsFromDB() que carrega todas as flags em uma única query
func LoadAuditSetting() {
	loadAllFlagsFromDB()
}

// StartAuditSettingSyncJob mantido para compatibilidade, mas não faz mais nada
// (sincronização agora unificada em StartFeatureFlagsSyncJob)
func StartAuditSettingSyncJob() {
	// Noop - sincronizado por StartFeatureFlagsSyncJob()
}

// IsAuditFrontendEnabled retorna o valor cacheado da flag
func IsAuditFrontendEnabled() bool {
	auditFrontendEnabledLock.RLock()
	defer auditFrontendEnabledLock.RUnlock()
	return auditFrontendEnabled
}

// ShouldAudit verifica se a auditoria deve ser feita para esta requisição.
// - Chamadas diretas pela API (sem header X-Source) → sempre audita
// - Chamadas do frontend (X-Source: frontend) → só audita se a flag estiver ativada
func ShouldAudit(r *http.Request) bool {
	source := r.Header.Get("X-Source")
	if source == "frontend" {
		return IsAuditFrontendEnabled()
	}
	// Chamada direta pela API → sempre audita
	return true
}

// LogAuditForRequest é um wrapper que verifica ShouldAudit antes de registrar (assíncrono)
func LogAuditForRequest(r *http.Request, action string, userEmail string, details interface{}, ipAddress string) {
	source := r.Header.Get("X-Source")
	shouldAudit := ShouldAudit(r)
	LogDebugf("LogAuditForRequest: action=%s, X-Source=%s, auditFrontendEnabled=%v, shouldAudit=%v", action, source, IsAuditFrontendEnabled(), shouldAudit)
	if !shouldAudit {
		return
	}
	LogAuditAsync(action, userEmail, details, ipAddress)
}

// LogAuditAsync envia um evento de auditoria para o buffer de batch insert (não bloqueia)
func LogAuditAsync(action string, userEmail string, details interface{}, ipAddress string) {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		log.Printf("[ERROR] Erro ao converter detalhes para JSON: %v", err)
		detailsJSON = []byte("{}")
	}

	auditLog := models.AuditLog{
		Action:    action,
		UserEmail: userEmail,
		Details:   datatypes.JSON(detailsJSON),
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	select {
	case auditBufCh <- auditLog:
	default:
		log.Printf("[WARN] Audit buffer cheio, descartando log: %s", action)
	}
}

// logAuditWithRetry tenta registrar auditoria com retry automático
func logAuditWithRetry(action string, userEmail string, details interface{}, ipAddress string, retries int) bool {
	for attempt := 1; attempt <= retries; attempt++ {
		if logAuditInternal(action, userEmail, details, ipAddress) {
			return true
		}
		if attempt < retries {
			time.Sleep(100 * time.Millisecond)
		}
	}
	log.Printf("[ERROR] Falha ao registrar auditoria após %d tentativas: action=%s, email=%s", retries, action, userEmail)
	return false
}

// logAuditInternal faz o registro real da auditoria
func logAuditInternal(action string, userEmail string, details interface{}, ipAddress string) bool {
	if DB == nil {
		log.Printf("[ERROR] Database não inicializado ao registrar auditoria: %s", action)
		return false
	}

	// Converter details para JSON
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		log.Printf("[ERROR] Erro ao converter detalhes para JSON: %v", err)
		detailsJSON = []byte("{}")
	}

	auditLog := models.AuditLog{
		Action:    action,
		UserEmail: userEmail,
		Details:   datatypes.JSON(detailsJSON),
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := DB.Create(&auditLog)

	if result.Error != nil {
		log.Printf("[ERROR] Erro ao registrar auditoria: action=%s, email=%s, erro=%v", action, userEmail, result.Error)
		return false
	}

	return true
}

// LogAudit registra um evento de auditoria no banco de dados (síncrono)
func LogAudit(action string, userEmail string, details interface{}, ipAddress string) {
	success := logAuditWithRetry(action, userEmail, details, ipAddress, 3)
	if !success {
		log.Printf("[ERROR] LogAudit falhou após retries | action=%s, email=%s", action, userEmail)
	}
}

// GetClientIP extrai o IP do cliente da requisição
func GetClientIP(r *http.Request) string {
	// Verificar X-Forwarded-For (proxy/load balancer)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Verificar X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Usar RemoteAddr como fallback
	return strings.Split(r.RemoteAddr, ":")[0]
}

// CleanOldAuditLogs deleta logs com mais de 30 dias (em lotes para evitar locks longos)
func CleanOldAuditLogs() {
	if DB == nil {
		log.Printf("[WARN] Database não inicializado, não foi possível limpar logs de auditoria")
		return
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	// Deletar em lotes de 1000 para não travar a tabela
	for {
		result := DB.Where("created_at < ?", thirtyDaysAgo).Limit(1000).Delete(&models.AuditLog{})
		if result.Error != nil {
			log.Printf("[ERROR] Erro ao limpar logs de auditoria antigos: %v", result.Error)
			return
		}
		if result.RowsAffected == 0 {
			break
		}
		log.Printf("[INFO] Removidos %d logs de auditoria antigos", result.RowsAffected)
	}

	log.Printf("[INFO] Limpeza de logs de auditoria concluída")
}
