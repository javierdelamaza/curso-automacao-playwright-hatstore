package config

import (
	"testing"
)

// =============================================================================
// Testes das Feature Flags
// =============================================================================

func TestIsFiltersEnabled_Default(t *testing.T) {
	filtersEnabledLock.Lock()
	filtersEnabled = true
	filtersEnabledLock.Unlock()

	if !IsFiltersEnabled() {
		t.Error("Filtros deveriam estar habilitados por padrão")
	}
}

func TestIsFiltersEnabled_Disabled(t *testing.T) {
	filtersEnabledLock.Lock()
	filtersEnabled = false
	filtersEnabledLock.Unlock()
	defer func() {
		filtersEnabledLock.Lock()
		filtersEnabled = true
		filtersEnabledLock.Unlock()
	}()

	if IsFiltersEnabled() {
		t.Error("Filtros deveriam estar desabilitados após serem setados para false")
	}
}

func TestIsPlatformEnabled_Default(t *testing.T) {
	platformEnabledLock.Lock()
	platformEnabled = true
	platformEnabledLock.Unlock()

	if !IsPlatformEnabled() {
		t.Error("Plataforma deveria estar habilitada por padrão")
	}
}

func TestIsPlatformEnabled_Disabled(t *testing.T) {
	platformEnabledLock.Lock()
	platformEnabled = false
	platformEnabledLock.Unlock()
	defer func() {
		platformEnabledLock.Lock()
		platformEnabled = true
		platformEnabledLock.Unlock()
	}()

	if IsPlatformEnabled() {
		t.Error("Plataforma deveria estar desabilitada após ser setada para false")
	}
}

func TestIsFiltersEnabled_ConcurrentSafe(t *testing.T) {
	// Testa que leitura/escrita concorrente não causa race condition
	done := make(chan bool, 10)

	for i := 0; i < 5; i++ {
		go func() {
			_ = IsFiltersEnabled()
			done <- true
		}()
	}
	for i := 0; i < 5; i++ {
		go func(val bool) {
			filtersEnabledLock.Lock()
			filtersEnabled = val
			filtersEnabledLock.Unlock()
			done <- true
		}(i%2 == 0)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Restaurar
	filtersEnabledLock.Lock()
	filtersEnabled = true
	filtersEnabledLock.Unlock()
}

func TestIsPlatformEnabled_ConcurrentSafe(t *testing.T) {
	done := make(chan bool, 10)

	for i := 0; i < 5; i++ {
		go func() {
			_ = IsPlatformEnabled()
			done <- true
		}()
	}
	for i := 0; i < 5; i++ {
		go func(val bool) {
			platformEnabledLock.Lock()
			platformEnabled = val
			platformEnabledLock.Unlock()
			done <- true
		}(i%2 == 0)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Restaurar
	platformEnabledLock.Lock()
	platformEnabled = true
	platformEnabledLock.Unlock()
}
