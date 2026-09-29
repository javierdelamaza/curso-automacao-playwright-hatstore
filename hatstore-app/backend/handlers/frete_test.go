package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

// TestConsultarFrete_Success testa consulta com CEP e quantidade válidos (Sudeste)
func TestConsultarFrete_Success(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=01001-000&quantidade=2", nil, ConsultarFrete)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d: %s", rr.Code, rr.Body.String())
	}

	var resp FreteResponse
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Regiao != "Sudeste" {
		t.Errorf("Esperava região 'Sudeste', recebeu '%s'", resp.Regiao)
	}
	if resp.CustoPorItem != 12.50 {
		t.Errorf("Esperava custo por item 12.50, recebeu %.2f", resp.CustoPorItem)
	}
	if resp.Quantidade != 2 {
		t.Errorf("Esperava quantidade 2, recebeu %d", resp.Quantidade)
	}
	// 12.50 + 12.50*0.30 = 16.25
	if resp.ValorFrete != 16.25 {
		t.Errorf("Esperava valor frete 16.25, recebeu %.2f", resp.ValorFrete)
	}
}

// TestConsultarFrete_Sul testa CEP da região Sul
func TestConsultarFrete_Sul(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=80000-000&quantidade=3", nil, ConsultarFrete)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d: %s", rr.Code, rr.Body.String())
	}

	var resp FreteResponse
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Regiao != "Sul" {
		t.Errorf("Esperava região 'Sul', recebeu '%s'", resp.Regiao)
	}
	// 15.00 + 15.00*0.30 + 15.00*0.20 = 22.50
	if resp.ValorFrete != 22.50 {
		t.Errorf("Esperava valor frete 22.50, recebeu %.2f", resp.ValorFrete)
	}
}

// TestConsultarFrete_Nordeste testa CEP da região Nordeste
func TestConsultarFrete_Nordeste(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=40000-000&quantidade=1", nil, ConsultarFrete)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d: %s", rr.Code, rr.Body.String())
	}

	var resp FreteResponse
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Regiao != "Nordeste" {
		t.Errorf("Esperava região 'Nordeste', recebeu '%s'", resp.Regiao)
	}
	if resp.ValorFrete != 22.00 {
		t.Errorf("Esperava valor frete 22.00, recebeu %.2f", resp.ValorFrete)
	}
}

// TestConsultarFrete_Norte testa CEP da região Norte
func TestConsultarFrete_Norte(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=66000-000&quantidade=2", nil, ConsultarFrete)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d: %s", rr.Code, rr.Body.String())
	}

	var resp FreteResponse
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Regiao != "Norte" {
		t.Errorf("Esperava região 'Norte', recebeu '%s'", resp.Regiao)
	}
	// 28.00 + 28.00*0.30 = 36.40
	if resp.ValorFrete != 36.40 {
		t.Errorf("Esperava valor frete 36.40, recebeu %.2f", resp.ValorFrete)
	}
}

// TestConsultarFrete_CentroOeste testa CEP da região Centro-Oeste
func TestConsultarFrete_CentroOeste(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=70000-000&quantidade=4", nil, ConsultarFrete)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d: %s", rr.Code, rr.Body.String())
	}

	var resp FreteResponse
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Regiao != "Centro-Oeste" {
		t.Errorf("Esperava região 'Centro-Oeste', recebeu '%s'", resp.Regiao)
	}
	// 20.00 + 20.00*0.30 + 20.00*0.20 + 20.00*0.10 = 32.00
	if resp.ValorFrete != 32.00 {
		t.Errorf("Esperava valor frete 32.00, recebeu %.2f", resp.ValorFrete)
	}
}

// TestConsultarFrete_ProgressivoQuantidade testa cálculo progressivo por quantidade
func TestConsultarFrete_ProgressivoQuantidade(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	// Sudeste: custo base 12.50
	// 1 item: 12.50
	// 2 itens: 12.50 + 3.75 = 16.25
	// 3 itens: 12.50 + 3.75 + 2.50 = 18.75
	// 4 itens: 12.50 + 3.75 + 2.50 + 1.25 = 20.00
	// 5 itens: 12.50 + 3.75 + 2.50 + 1.25 + 1.25 = 21.25
	// 6 itens: 12.50 + 3.75 + 2.50 + 1.25 + 1.25 + 1.25 = 22.50
	cases := []struct {
		qtd      int
		esperado float64
	}{
		{1, 12.50},
		{2, 16.25},
		{3, 18.75},
		{4, 20.00},
		{5, 21.25},
		{6, 22.50},
	}

	for _, tc := range cases {
		rr := executeRequest("GET", "/api/frete?cep=01001-000&quantidade="+strconv.Itoa(tc.qtd), nil, ConsultarFrete)
		var resp FreteResponse
		json.NewDecoder(rr.Body).Decode(&resp)
		if resp.ValorFrete != tc.esperado {
			t.Errorf("Quantidade %d: esperava frete %.2f, recebeu %.2f", tc.qtd, tc.esperado, resp.ValorFrete)
		}
	}
}

// TestConsultarFrete_CEPSemHifen testa CEP sem hífen
func TestConsultarFrete_CEPSemHifen(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=01001000&quantidade=1", nil, ConsultarFrete)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d: %s", rr.Code, rr.Body.String())
	}
}

// TestConsultarFrete_CEPInvalido testa CEP com formato inválido
func TestConsultarFrete_CEPInvalido(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=123&quantidade=1", nil, ConsultarFrete)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status 400, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_CEPVazio testa CEP vazio
func TestConsultarFrete_CEPVazio(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=&quantidade=1", nil, ConsultarFrete)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status 400, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_QuantidadeInvalida testa quantidade não numérica
func TestConsultarFrete_QuantidadeInvalida(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=01001-000&quantidade=abc", nil, ConsultarFrete)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status 400, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_QuantidadeZero testa quantidade zero
func TestConsultarFrete_QuantidadeZero(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=01001-000&quantidade=0", nil, ConsultarFrete)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status 400, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_QuantidadeNegativa testa quantidade negativa
func TestConsultarFrete_QuantidadeNegativa(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=01001-000&quantidade=-1", nil, ConsultarFrete)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status 400, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_QuantidadeVazia testa sem parâmetro quantidade
func TestConsultarFrete_QuantidadeVazia(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=01001-000", nil, ConsultarFrete)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status 400, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_CEPNaoAtendido testa CEP fora de todas as faixas
func TestConsultarFrete_CEPNaoAtendido(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "true")

	rr := executeRequest("GET", "/api/frete?cep=00100-000&quantidade=1", nil, ConsultarFrete)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Esperava status 404, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_FeatureFlagDesabilitada testa com flag desligada
func TestConsultarFrete_FeatureFlagDesabilitada(t *testing.T) {
	setupTestDB(t)
	setNewHatstoreFlag(t, "false")

	rr := executeRequest("GET", "/api/frete?cep=01001-000&quantidade=1", nil, ConsultarFrete)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Esperava status 503, recebeu %d", rr.Code)
	}
}

// TestConsultarFrete_FeatureFlagHabilitadaDepois testa que habilitar a flag libera o endpoint
func TestConsultarFrete_FeatureFlagHabilitadaDepois(t *testing.T) {
	setupTestDB(t)

	// Começa desabilitada
	setNewHatstoreFlag(t, "false")
	rr := executeRequest("GET", "/api/frete?cep=01001-000&quantidade=1", nil, ConsultarFrete)
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Esperava status 503, recebeu %d", rr.Code)
	}

	// Habilita
	setNewHatstoreFlag(t, "true")
	rr = executeRequest("GET", "/api/frete?cep=01001-000&quantidade=1", nil, ConsultarFrete)
	if rr.Code != http.StatusOK {
		t.Errorf("Esperava status 200 após habilitar flag, recebeu %d", rr.Code)
	}
}
