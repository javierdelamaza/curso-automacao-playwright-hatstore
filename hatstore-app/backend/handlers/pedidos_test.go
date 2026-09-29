package handlers

import (
	"encoding/json"
	"hat-store-training/backend/config"
	"hat-store-training/backend/models"
	"net/http"
	"os"
	"testing"
)

// =============================================================================
// Testes do CriarPedido
// =============================================================================

func TestCriarPedido_Success(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	// Desabilitar verificação de ambiente de desenvolvimento
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	// Verificar que o pedido foi criado
	var pedido models.Pedido
	if err := config.DB.First(&pedido).Error; err != nil {
		t.Fatalf("Pedido não foi criado: %v", err)
	}

	if pedido.Nome != "João da Silva" {
		t.Errorf("Nome esperado 'João da Silva', recebeu '%s'", pedido.Nome)
	}
	if pedido.CPF != "123.456.789-00" {
		t.Errorf("CPF esperado '123.456.789-00', recebeu '%s'", pedido.CPF)
	}
	if pedido.Pagamento != "pix" {
		t.Errorf("Pagamento esperado 'pix', recebeu '%s'", pedido.Pagamento)
	}

	// Verificar que o estoque foi descontado
	var hat models.Hat
	config.DB.Where("id = ?", hats[0].ID).First(&hat)
	if hat.Quantidade != 9 {
		t.Errorf("Estoque esperado 9 (era 10 - 1), recebeu %d", hat.Quantidade)
	}
}

func TestCriarPedido_CalculatesTotal(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "false")

	payload := validPedidoPayload(uint(hats[0].ID))
	// 2 unidades a 49.90
	payload["itens"] = []map[string]interface{}{
		{
			"id":         hats[0].ID,
			"nome":       "Boné Nacional",
			"price":      49.90,
			"quantidade": 2,
		},
	}
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	var pedido models.Pedido
	config.DB.First(&pedido)

	expectedTotal := 49.90 * 2
	if pedido.Total != expectedTotal {
		t.Errorf("Total esperado %.2f, recebeu %.2f", expectedTotal, pedido.Total)
	}
}

func TestCriarPedido_EmptyCart(t *testing.T) {
	setupTestDB(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(1)
	payload["itens"] = []map[string]interface{}{}
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCriarPedido_InvalidJSON(t *testing.T) {
	setupTestDB(t)
	os.Setenv("ENVIRONMENT", "")

	rr := executeRequest("POST", "/api/pedido", []byte("{invalid}"), CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingNome(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["nome"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_NomeMinLength(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["nome"] = "AB" // menos de 3 caracteres
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d para nome curto, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingCPF(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cpf"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_InvalidCPF(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cpf"] = "12345678900" // sem formatação
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_InvalidEmail(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["email"] = "email-invalido"
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingEmail(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["email"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_InvalidTelefone(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["telefone"] = "11999990000" // sem formatação
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_InvalidCEP(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cep"] = "01234567" // sem hífen
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_InvalidUF(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["uf"] = "SPA" // 3 letras
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_InvalidPagamento(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["pagamento"] = "cartao"
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_PagamentoBoleto(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["pagamento"] = "boleto"
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Errorf("Esperava status %d para pagamento boleto, recebeu %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
}

func TestCriarPedido_MissingEndereco(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["endereco"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_EnderecoMinLength(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["endereco"] = "AB" // menos de 3
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_BairroMinLength(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["bairro"] = "AB"
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_CidadeMinLength(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cidade"] = "SP" // menos de 3
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_InsufficientStock(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	// Pedir mais do que o estoque (10) permite
	payload["itens"] = []map[string]interface{}{
		{
			"id":         hats[0].ID,
			"nome":       "Boné Nacional",
			"price":      49.90,
			"quantidade": 999,
		},
	}
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusConflict {
		t.Errorf("Esperava status %d (estoque insuficiente), recebeu %d: %s", http.StatusConflict, rr.Code, rr.Body.String())
	}
}

func TestCriarPedido_HatNotFound(t *testing.T) {
	setupTestDB(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(9999) // ID inexistente
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Esperava status %d (hat não encontrado), recebeu %d: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestCriarPedido_HATOFF_DiscountApplied(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "false")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cupom"] = "HATOFF"
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	var pedido models.Pedido
	config.DB.First(&pedido)

	// Preço original: 49.90 * 1 = 49.90, com 25% desconto = 37.425
	expectedTotal := 49.90 * 0.75
	if pedido.Total != expectedTotal {
		t.Errorf("Total esperado %.2f (com HATOFF), recebeu %.2f", expectedTotal, pedido.Total)
	}

	if pedido.Cupom != "HATOFF" {
		t.Errorf("Cupom esperado 'HATOFF', recebeu '%s'", pedido.Cupom)
	}

	// Verificar que o uso do cupom foi registrado
	var cupomUsado models.CupomHatOffUsado
	if err := config.DB.Where("cpf = ?", "123.456.789-00").First(&cupomUsado).Error; err != nil {
		t.Errorf("Uso do cupom HATOFF não foi registrado: %v", err)
	}
}

func TestCriarPedido_HATOFF_AlreadyUsed(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	// Registrar uso prévio do cupom
	config.DB.Create(&models.CupomHatOffUsado{CPF: "123.456.789-00"})

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cupom"] = "HATOFF"
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Esperava status %d (cupom já usado), recebeu %d: %s", http.StatusForbidden, rr.Code, rr.Body.String())
	}
}

func TestCriarPedido_MultipleItems(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "false")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["itens"] = []map[string]interface{}{
		{
			"id":         hats[0].ID,
			"nome":       "Boné Nacional",
			"price":      49.90,
			"quantidade": 2,
		},
		{
			"id":         hats[1].ID,
			"nome":       "Chapéu Importado",
			"price":      129.90,
			"quantidade": 1,
		},
	}
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	// Verificar total
	var pedido models.Pedido
	config.DB.Preload("Itens").First(&pedido)

	expectedTotal := (49.90 * 2) + (129.90 * 1)
	if pedido.Total != expectedTotal {
		t.Errorf("Total esperado %.2f, recebeu %.2f", expectedTotal, pedido.Total)
	}

	if len(pedido.Itens) != 2 {
		t.Errorf("Esperava 2 itens, recebeu %d", len(pedido.Itens))
	}

	// Verificar estoque descontado
	var hat0, hat1 models.Hat
	config.DB.Where("id = ?", hats[0].ID).First(&hat0)
	config.DB.Where("id = ?", hats[1].ID).First(&hat1)

	if hat0.Quantidade != 8 {
		t.Errorf("Estoque hat0 esperado 8 (10-2), recebeu %d", hat0.Quantidade)
	}
	if hat1.Quantidade != 4 {
		t.Errorf("Estoque hat1 esperado 4 (5-1), recebeu %d", hat1.Quantidade)
	}
}

func TestCriarPedido_DevEnvWithoutUUID(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "development")
	os.Setenv("DEV_UUID", "test-uuid-123")
	defer func() {
		os.Setenv("ENVIRONMENT", "")
		os.Setenv("DEV_UUID", "")
	}()

	payload := validPedidoPayload(uint(hats[0].ID))
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Esperava status %d sem UUID em dev, recebeu %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestCriarPedido_MissingTelefone(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["telefone"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingNumero(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["numero"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingBairro(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["bairro"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingCEP(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cep"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingCidade(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cidade"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCriarPedido_MissingUF(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["uf"] = ""
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

// =============================================================================
// Testes do ConsultarPedidos
// =============================================================================

func TestConsultarPedidos_Success(t *testing.T) {
	setupTestDB(t)
	os.Setenv("ENVIRONMENT", "")

	// Criar um pedido no banco diretamente
	pedido := models.Pedido{
		Nome:      "Maria",
		CPF:       "111.222.333-44",
		Email:     "maria@email.com",
		Telefone:  "(11) 99999-0000",
		Endereco:  "Rua do Teste",
		Numero:    "100",
		Bairro:    "Centro",
		CEP:       "01000-000",
		Cidade:    "São Paulo",
		UF:        "SP",
		Pagamento: "pix",
		Total:     99.90,
	}
	config.DB.Create(&pedido)

	// Criar itens do pedido
	item := models.PedidoItem{
		PedidoID:   pedido.ID,
		HatID:      1,
		Nome:       "Boné Teste",
		Price:      99.90,
		Quantidade: 1,
	}
	config.DB.Create(&item)

	rr := executeRequest("GET", "/api/pedidos?cpf=111.222.333-44", nil, ConsultarPedidos)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var pedidos []models.Pedido
	json.Unmarshal(rr.Body.Bytes(), &pedidos)

	if len(pedidos) != 1 {
		t.Errorf("Esperava 1 pedido, recebeu %d", len(pedidos))
	}

	if len(pedidos) > 0 && pedidos[0].CPF != "111.222.333-44" {
		t.Errorf("CPF esperado '111.222.333-44', recebeu '%s'", pedidos[0].CPF)
	}
}

func TestConsultarPedidos_MissingCPF(t *testing.T) {
	setupTestDB(t)

	rr := executeRequest("GET", "/api/pedidos", nil, ConsultarPedidos)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestConsultarPedidos_InvalidCPF(t *testing.T) {
	setupTestDB(t)

	rr := executeRequest("GET", "/api/pedidos?cpf=12345678900", nil, ConsultarPedidos)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d para CPF sem formatação, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestConsultarPedidos_EmptyResult(t *testing.T) {
	setupTestDB(t)

	rr := executeRequest("GET", "/api/pedidos?cpf=000.000.000-00", nil, ConsultarPedidos)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var pedidos []models.Pedido
	json.Unmarshal(rr.Body.Bytes(), &pedidos)

	if len(pedidos) != 0 {
		t.Errorf("Esperava 0 pedidos, recebeu %d", len(pedidos))
	}
}

func TestConsultarPedidos_MultiplePedidos(t *testing.T) {
	setupTestDB(t)

	cpf := "111.222.333-44"
	for i := 0; i < 3; i++ {
		config.DB.Create(&models.Pedido{
			Nome: "Maria", CPF: cpf, Email: "m@e.com",
			Telefone: "(11) 99999-0000", Endereco: "Rua X",
			Numero: "100", Bairro: "Centro", CEP: "01000-000",
			Cidade: "SP", UF: "SP", Pagamento: "pix", Total: 10.0,
		})
	}

	rr := executeRequest("GET", "/api/pedidos?cpf=111.222.333-44", nil, ConsultarPedidos)

	var pedidos []models.Pedido
	json.Unmarshal(rr.Body.Bytes(), &pedidos)

	if len(pedidos) != 3 {
		t.Errorf("Esperava 3 pedidos, recebeu %d", len(pedidos))
	}
}

// =============================================================================
// Testes do PedidoJaExiste e CPFUsouHatOff (utilitários)
// =============================================================================

func TestPedidoJaExiste_True(t *testing.T) {
	setupTestDB(t)

	config.DB.Create(&models.Pedido{
		Nome: "Teste", CPF: "111.222.333-44", Email: "t@e.com",
		Telefone: "(11) 99999-0000", Endereco: "Rua T",
		Numero: "1", Bairro: "Centro", CEP: "01000-000",
		Cidade: "SP", UF: "SP", Pagamento: "pix", Total: 10.0,
	})

	if !PedidoJaExiste("111.222.333-44") {
		t.Error("Esperava true, pedido já existe")
	}
}

func TestPedidoJaExiste_False(t *testing.T) {
	setupTestDB(t)

	if PedidoJaExiste("000.000.000-00") {
		t.Error("Esperava false, nenhum pedido existe")
	}
}

func TestCPFUsouHatOff_True(t *testing.T) {
	setupTestDB(t)

	config.DB.Create(&models.CupomHatOffUsado{CPF: "111.222.333-44"})

	if !CPFUsouHatOff("111.222.333-44") {
		t.Error("Esperava true, CPF já usou HATOFF")
	}
}

func TestCPFUsouHatOff_False(t *testing.T) {
	setupTestDB(t)

	if CPFUsouHatOff("000.000.000-00") {
		t.Error("Esperava false, CPF não usou HATOFF")
	}
}

// =============================================================================
// Testes de Frete na criação de pedido
// =============================================================================

func TestCriarPedido_FreteCalculadoQuandoHabilitado(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "true")

	payload := validPedidoPayload(uint(hats[0].ID))
	// CEP 01234-567 = Sudeste, 1 item → frete = 12.50
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	var pedido models.Pedido
	config.DB.First(&pedido)

	expectedFrete := 12.50
	if pedido.Frete != expectedFrete {
		t.Errorf("Frete esperado %.2f, recebeu %.2f", expectedFrete, pedido.Frete)
	}

	expectedTotal := 49.90 + expectedFrete // itens + frete
	if pedido.Total != expectedTotal {
		t.Errorf("Total esperado %.2f, recebeu %.2f", expectedTotal, pedido.Total)
	}
}

func TestCriarPedido_FreteZeroQuandoDesabilitado(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "false")

	payload := validPedidoPayload(uint(hats[0].ID))
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	var pedido models.Pedido
	config.DB.First(&pedido)

	if pedido.Frete != 0 {
		t.Errorf("Frete esperado 0 (desabilitado), recebeu %.2f", pedido.Frete)
	}

	if pedido.Total != 49.90 {
		t.Errorf("Total esperado 49.90 (sem frete), recebeu %.2f", pedido.Total)
	}
}

func TestCriarPedido_FreteMultiplosItens(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "true")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["itens"] = []map[string]interface{}{
		{
			"id":         hats[0].ID,
			"nome":       "Boné Nacional",
			"price":      49.90,
			"quantidade": 2,
		},
		{
			"id":         hats[1].ID,
			"nome":       "Chapéu Importado",
			"price":      129.90,
			"quantidade": 1,
		},
	}
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	var pedido models.Pedido
	config.DB.First(&pedido)

	// 3 itens no total, Sudeste: 12.50 + 12.50*0.30 + 12.50*0.20 = 18.75
	expectedFrete := 18.75
	if pedido.Frete != expectedFrete {
		t.Errorf("Frete esperado %.2f, recebeu %.2f", expectedFrete, pedido.Frete)
	}

	expectedTotal := (49.90*2 + 129.90*1) + expectedFrete
	if pedido.Total != expectedTotal {
		t.Errorf("Total esperado %.2f, recebeu %.2f", expectedTotal, pedido.Total)
	}
}

func TestCriarPedido_FreteComCupomHATOFF(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "true")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cupom"] = "HATOFF"
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	var pedido models.Pedido
	config.DB.First(&pedido)

	// Desconto HATOFF aplica nos itens, frete é somado depois
	// Itens: 49.90 * 0.75 = 37.425, Frete: 12.50
	expectedFrete := 12.50
	expectedTotal := (49.90 * 0.75) + expectedFrete
	if pedido.Frete != expectedFrete {
		t.Errorf("Frete esperado %.2f, recebeu %.2f", expectedFrete, pedido.Frete)
	}
	if pedido.Total != expectedTotal {
		t.Errorf("Total esperado %.4f, recebeu %.4f", expectedTotal, pedido.Total)
	}
}

func TestCriarPedido_FreteCEPNaoAtendido(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "true")

	payload := validPedidoPayload(uint(hats[0].ID))
	payload["cep"] = "00000-000" // CEP fora de qualquer faixa
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperava status %d (CEP não atendido), recebeu %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCriarPedido_FreteIgnoradoSemNewHatstore(t *testing.T) {
	setupTestDB(t)
	hats := seedHats(t)
	os.Setenv("ENVIRONMENT", "")
	setNewHatstoreFlag(t, "false")

	payload := validPedidoPayload(uint(hats[0].ID))
	body := makeJSON(t, payload)

	rr := executeRequest("POST", "/api/pedido", body, CriarPedido)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Esperava status %d: %s", http.StatusCreated, rr.Body.String())
	}

	var pedido models.Pedido
	config.DB.First(&pedido)

	if pedido.Frete != 0 {
		t.Errorf("Frete esperado 0 (new_hatstore desabilitado), recebeu %.2f", pedido.Frete)
	}

	if pedido.Total != 49.90 {
		t.Errorf("Total esperado 49.90 (sem frete), recebeu %.2f", pedido.Total)
	}
}
