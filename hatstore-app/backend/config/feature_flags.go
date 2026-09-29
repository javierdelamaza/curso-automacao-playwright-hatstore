package config

import (
	"log"
	"sync"
	"time"

	"hat-store-training/backend/models"
)

// Cache em memória para feature flags
var (
	filtersEnabled          = true
	filtersEnabledLock      sync.RWMutex
	platformEnabled         = true
	platformEnabledLock     sync.RWMutex
	fixPasswordAuditLog     = false
	fixPasswordAuditLogLock sync.RWMutex
	fixHatoffDiscount       = false
	fixHatoffDiscountLock   sync.RWMutex
	newHatstore             = false
	newHatstoreLock         sync.RWMutex
	flagsCacheTTL           = 24 * time.Hour
)

// allFlagKeys lista todas as keys que precisam ser carregadas do banco
var allFlagKeys = []string{
	"filters_enabled",
	"platform_enabled",
	"fix_password_audit_log",
	"fix_hatoff_discount",
	"new_hatstore",
	"audit_frontend_enabled",
}

// loadAllFlagsFromDB carrega todas as flags em uma única query
func loadAllFlagsFromDB() {
	if DB == nil {
		return
	}
	var settings []models.Setting
	if err := DB.Where("key IN ?", allFlagKeys).Find(&settings).Error; err != nil {
		log.Printf("[ERROR] Erro ao carregar feature flags do banco: %v", err)
		return
	}

	// Criar mapa para lookup rápido
	settingsMap := make(map[string]string, len(settings))
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	// Atualizar todas as flags (resetar para default se não encontrada no banco)
	filtersEnabledLock.Lock()
	if v, ok := settingsMap["filters_enabled"]; ok {
		filtersEnabled = v == "true"
	} else {
		filtersEnabled = true
	}
	filtersEnabledLock.Unlock()

	platformEnabledLock.Lock()
	if v, ok := settingsMap["platform_enabled"]; ok {
		platformEnabled = v == "true"
	} else {
		platformEnabled = true
	}
	platformEnabledLock.Unlock()

	fixPasswordAuditLogLock.Lock()
	if v, ok := settingsMap["fix_password_audit_log"]; ok {
		fixPasswordAuditLog = v == "true"
	} else {
		fixPasswordAuditLog = false
	}
	fixPasswordAuditLogLock.Unlock()

	fixHatoffDiscountLock.Lock()
	if v, ok := settingsMap["fix_hatoff_discount"]; ok {
		fixHatoffDiscount = v == "true"
	} else {
		fixHatoffDiscount = false
	}
	fixHatoffDiscountLock.Unlock()

	newHatstoreLock.Lock()
	if v, ok := settingsMap["new_hatstore"]; ok {
		newHatstore = v == "true"
	} else {
		newHatstore = false
	}
	newHatstoreLock.Unlock()

	auditFrontendEnabledLock.Lock()
	if v, ok := settingsMap["audit_frontend_enabled"]; ok {
		auditFrontendEnabled = v == "true"
	} else {
		auditFrontendEnabled = true
	}
	auditFrontendEnabledLock.Unlock()

	log.Printf("[INFO] Feature flags carregadas do banco (%d flags em 1 query)", len(settings))
}

// LoadFeatureFlags carrega as feature flags do banco para o cache em memória (chamar na inicialização)
func LoadFeatureFlags() {
	loadAllFlagsFromDB()
}

// StartFeatureFlagsSyncJob inicia uma goroutine que sincroniza as flags do banco a cada 2 minutos
func StartFeatureFlagsSyncJob() {
	go func() {
		ticker := time.NewTicker(flagsCacheTTL)
		defer ticker.Stop()
		for range ticker.C {
			loadAllFlagsFromDB()
		}
	}()
}

// IsFiltersEnabled retorna o valor cacheado da flag de filtros
func IsFiltersEnabled() bool {
	filtersEnabledLock.RLock()
	defer filtersEnabledLock.RUnlock()
	return filtersEnabled
}

// IsPlatformEnabled retorna o valor cacheado da flag de plataforma
func IsPlatformEnabled() bool {
	platformEnabledLock.RLock()
	defer platformEnabledLock.RUnlock()
	return platformEnabled
}

// IsFixPasswordAuditLogEnabled retorna o valor cacheado da flag
func IsFixPasswordAuditLogEnabled() bool {
	fixPasswordAuditLogLock.RLock()
	defer fixPasswordAuditLogLock.RUnlock()
	return fixPasswordAuditLog
}

// IsFixHatoffDiscountEnabled retorna o valor cacheado da flag
func IsFixHatoffDiscountEnabled() bool {
	fixHatoffDiscountLock.RLock()
	defer fixHatoffDiscountLock.RUnlock()
	return fixHatoffDiscount
}

// IsNewHatstoreEnabled retorna o valor cacheado da flag new_hatstore
func IsNewHatstoreEnabled() bool {
	newHatstoreLock.RLock()
	defer newHatstoreLock.RUnlock()
	return newHatstore
}
