package handlers

import (
	"encoding/json"
	"fmt"
	"hat-store-training/backend/config"
	"hat-store-training/backend/models"
	"net/http"
	"testing"
)

// =============================================================================
// Testes do GetHats
// =============================================================================

func TestGetHats_ReturnsAll(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats", nil, GetHats)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var hats []HatComEstoque
	if err := json.Unmarshal(rr.Body.Bytes(), &hats); err != nil {
		t.Fatalf("Erro ao parsear resposta: %v", err)
	}

	if len(hats) != 4 {
		t.Errorf("Esperava 4 chapéus, recebeu %d", len(hats))
	}
}

func TestGetHats_TemEstoqueField(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats", nil, GetHats)

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	// Verificar o campo temEstoque
	for _, hat := range hats {
		if hat.Nome == "Boné Crescer" && hat.TemEstoque {
			t.Error("Boné Crescer tem quantidade 0, temEstoque deveria ser false")
		}
		if hat.Nome == "Boné Nacional" && !hat.TemEstoque {
			t.Error("Boné Nacional tem quantidade 10, temEstoque deveria ser true")
		}
	}
}

func TestGetHats_FilterByCategoria(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats?categoria=importado", nil, GetHats)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	if len(hats) != 2 {
		t.Errorf("Esperava 2 chapéus importados, recebeu %d", len(hats))
	}

	for _, hat := range hats {
		if hat.Categoria != "importado" {
			t.Errorf("Esperava categoria 'importado', recebeu '%s'", hat.Categoria)
		}
	}
}

func TestGetHats_FilterByMultipleCategorias(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats?categoria=nacional,crescer", nil, GetHats)

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	if len(hats) != 2 {
		t.Errorf("Esperava 2 chapéus (nacional + crescer), recebeu %d", len(hats))
	}

	for _, hat := range hats {
		if hat.Categoria != "nacional" && hat.Categoria != "crescer" {
			t.Errorf("Categoria inesperada: %s", hat.Categoria)
		}
	}
}

func TestGetHats_FilterByPriceMin(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats?min=100", nil, GetHats)

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	for _, hat := range hats {
		if hat.Price < 100 {
			t.Errorf("Chapéu %s tem preço %.2f, deveria ser >= 100", hat.Nome, hat.Price)
		}
	}
}

func TestGetHats_FilterByPriceMax(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats?max=80", nil, GetHats)

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	for _, hat := range hats {
		if hat.Price > 80 {
			t.Errorf("Chapéu %s tem preço %.2f, deveria ser <= 80", hat.Nome, hat.Price)
		}
	}
}

func TestGetHats_FilterByPriceRange(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats?min=50&max=130", nil, GetHats)

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	for _, hat := range hats {
		if hat.Price < 50 || hat.Price > 130 {
			t.Errorf("Chapéu %s tem preço %.2f, deveria estar entre 50 e 130", hat.Nome, hat.Price)
		}
	}
}

func TestGetHats_FilterCategoriaAndPrice(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/hats?categoria=importado&max=150", nil, GetHats)

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	if len(hats) != 1 {
		t.Errorf("Esperava 1 chapéu importado com preço <= 150, recebeu %d", len(hats))
	}

	if len(hats) > 0 && hats[0].Nome != "Chapéu Importado" {
		t.Errorf("Esperava 'Chapéu Importado', recebeu '%s'", hats[0].Nome)
	}
}

func TestGetHats_EmptyResult(t *testing.T) {
	setupTestDB(t)
	// Não faz seed - banco vazio

	rr := executeRequest("GET", "/api/hats", nil, GetHats)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var hats []HatComEstoque
	json.Unmarshal(rr.Body.Bytes(), &hats)

	if len(hats) != 0 {
		t.Errorf("Esperava 0 chapéus, recebeu %d", len(hats))
	}
}

// =============================================================================
// Testes do ListarEstoque
// =============================================================================

func TestListarEstoque_Success(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequest("GET", "/api/estoque", nil, ListarEstoque)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var estoque []struct {
		ID         uint   `json:"id"`
		Nome       string `json:"nome"`
		Quantidade int    `json:"quantidade"`
	}
	json.Unmarshal(rr.Body.Bytes(), &estoque)

	if len(estoque) != 4 {
		t.Errorf("Esperava 4 itens no estoque, recebeu %d", len(estoque))
	}

	// Verificar que os campos retornados são corretos (sem preço)
	for _, item := range estoque {
		if item.Nome == "" {
			t.Error("Nome não deveria ser vazio")
		}
	}
}

func TestListarEstoque_Empty(t *testing.T) {
	setupTestDB(t)

	rr := executeRequest("GET", "/api/estoque", nil, ListarEstoque)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var estoque []interface{}
	json.Unmarshal(rr.Body.Bytes(), &estoque)

	if len(estoque) != 0 {
		t.Errorf("Esperava 0, recebeu %d", len(estoque))
	}
}

// =============================================================================
// Testes do EditarEstoqueHat
// =============================================================================

func TestEditarEstoqueHat_Success(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)

	body := makeJSON(t, EstoquePayload{Quantidade: 50})

	rr := executeRequestWithVars("PATCH", fmt.Sprintf("/api/hats/%d/estoque", hats[0].ID), body, EditarEstoqueHat, map[string]string{
		"id": fmt.Sprintf("%d", hats[0].ID),
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	// Verificar que o estoque foi atualizado no banco
	var hat models.Hat
	config.DB.Where("id = ?", hats[0].ID).First(&hat)
	if hat.Quantidade != 50 {
		t.Errorf("Esperava quantidade 50, recebeu %d", hat.Quantidade)
	}
}

func TestEditarEstoqueHat_InvalidID(t *testing.T) {
	setupTestDB(t)

	body := makeJSON(t, EstoquePayload{Quantidade: 10})

	rr := executeRequestWithVars("PATCH", "/api/hats/abc/estoque", body, EditarEstoqueHat, map[string]string{
		"id": "abc",
	})

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEditarEstoqueHat_NotFound(t *testing.T) {
	setupTestDB(t)

	body := makeJSON(t, EstoquePayload{Quantidade: 10})

	rr := executeRequestWithVars("PATCH", "/api/hats/9999/estoque", body, EditarEstoqueHat, map[string]string{
		"id": "9999",
	})

	if rr.Code != http.StatusNotFound {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusNotFound, rr.Code)
	}
}

func TestEditarEstoqueHat_NegativeQuantity(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)

	body := makeJSON(t, map[string]int{"quantidade": -1})

	rr := executeRequestWithVars("PATCH", fmt.Sprintf("/api/hats/%d/estoque", hats[0].ID), body, EditarEstoqueHat, map[string]string{
		"id": fmt.Sprintf("%d", hats[0].ID),
	})

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEditarEstoqueHat_OverLimit(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)

	body := makeJSON(t, map[string]int{"quantidade": 201})

	rr := executeRequestWithVars("PATCH", fmt.Sprintf("/api/hats/%d/estoque", hats[0].ID), body, EditarEstoqueHat, map[string]string{
		"id": fmt.Sprintf("%d", hats[0].ID),
	})

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEditarEstoqueHat_ZeroQuantity(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)

	body := makeJSON(t, EstoquePayload{Quantidade: 0})

	rr := executeRequestWithVars("PATCH", fmt.Sprintf("/api/hats/%d/estoque", hats[0].ID), body, EditarEstoqueHat, map[string]string{
		"id": fmt.Sprintf("%d", hats[0].ID),
	})

	if rr.Code != http.StatusOK {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var hat models.Hat
	config.DB.Where("id = ?", hats[0].ID).First(&hat)
	if hat.Quantidade != 0 {
		t.Errorf("Esperava quantidade 0, recebeu %d", hat.Quantidade)
	}
}

func TestEditarEstoqueHat_MaxQuantity(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)

	body := makeJSON(t, EstoquePayload{Quantidade: 200})

	rr := executeRequestWithVars("PATCH", fmt.Sprintf("/api/hats/%d/estoque", hats[0].ID), body, EditarEstoqueHat, map[string]string{
		"id": fmt.Sprintf("%d", hats[0].ID),
	})

	if rr.Code != http.StatusOK {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}
}

func TestEditarEstoqueHat_InvalidJSON(t *testing.T) {
	setupTestDB(t)
	seedHats(t)

	rr := executeRequestWithVars("PATCH", "/api/hats/1/estoque", []byte("{bad json"), EditarEstoqueHat, map[string]string{
		"id": "1",
	})

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

// =============================================================================
// Testes do parseFloat
// =============================================================================

func TestParseFloat_Integer(t *testing.T) {
	result, err := parseFloat("100")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if result != 100.0 {
		t.Errorf("Esperava 100.0, recebeu %f", result)
	}
}

func TestParseFloat_WithDot(t *testing.T) {
	result, err := parseFloat("49.90")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if result != 49.90 {
		t.Errorf("Esperava 49.90, recebeu %f", result)
	}
}

func TestParseFloat_WithComma(t *testing.T) {
	result, err := parseFloat("49,90")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if result != 49.90 {
		t.Errorf("Esperava 49.90, recebeu %f", result)
	}
}

func TestParseFloat_Invalid(t *testing.T) {
	_, err := parseFloat("abc")
	if err == nil {
		t.Error("Esperava erro para input inválido")
	}
}

func TestParseFloat_Empty(t *testing.T) {
	_, err := parseFloat("")
	if err == nil {
		t.Error("Esperava erro para string vazia")
	}
}
