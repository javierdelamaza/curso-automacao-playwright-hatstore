package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hat-store-training/backend/config"
	"hat-store-training/backend/models"
)

func TestMeusPedidos_Success(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "comprador@test.com", "senha123")

	// Criar pedidos para o usuário
	config.DB.Create(&models.Pedido{
		Nome: "Comprador", CPF: "123.456.789-00", Email: "comprador@test.com",
		Telefone: "(11) 91234-5678", Endereco: "Rua Teste", Numero: "100",
		Bairro: "Centro", CEP: "01000-000", Cidade: "São Paulo", UF: "SP",
		Pagamento: "pix", Total: 99.90,
		Itens: []models.PedidoItem{
			{HatID: 1, Nome: "Boné Nacional", Price: 49.90, Quantidade: 2},
		},
	})
	config.DB.Create(&models.Pedido{
		Nome: "Comprador", CPF: "123.456.789-00", Email: "comprador@test.com",
		Telefone: "(11) 91234-5678", Endereco: "Rua Teste", Numero: "100",
		Bairro: "Centro", CEP: "01000-000", Cidade: "São Paulo", UF: "SP",
		Pagamento: "boleto", Total: 129.90, Cupom: "HATOFF",
		Itens: []models.PedidoItem{
			{HatID: 2, Nome: "Chapéu Importado", Price: 129.90, Quantidade: 1},
		},
	})

	token := generateTestToken(t, "comprador@test.com")

	req := httptest.NewRequest("GET", "/api/meus-pedidos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	MeusPedidos(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var pedidos []models.Pedido
	if err := json.Unmarshal(rr.Body.Bytes(), &pedidos); err != nil {
		t.Fatalf("Erro ao parsear resposta: %v", err)
	}

	if len(pedidos) != 2 {
		t.Errorf("Esperava 2 pedidos, recebeu %d", len(pedidos))
	}

	// Primeiro pedido deve ser o mais recente (ORDER BY created_at DESC)
	if len(pedidos) > 0 && len(pedidos[0].Itens) == 0 {
		t.Error("Itens do pedido não foram carregados (Preload)")
	}
}

func TestMeusPedidos_Empty(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "semcompras@test.com", "senha123")

	token := generateTestToken(t, "semcompras@test.com")

	req := httptest.NewRequest("GET", "/api/meus-pedidos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	MeusPedidos(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var pedidos []models.Pedido
	if err := json.Unmarshal(rr.Body.Bytes(), &pedidos); err != nil {
		t.Fatalf("Erro ao parsear resposta: %v", err)
	}

	if len(pedidos) != 0 {
		t.Errorf("Esperava 0 pedidos, recebeu %d", len(pedidos))
	}
}

func TestMeusPedidos_InvalidToken(t *testing.T) {
	setupTestDB(t)
	config.JWTKey = []byte("test-secret-key")

	req := httptest.NewRequest("GET", "/api/meus-pedidos", nil)
	req.Header.Set("Authorization", "Bearer token-invalido")
	rr := httptest.NewRecorder()
	MeusPedidos(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Esperava status %d, recebeu %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestMeusPedidos_DoesNotReturnOtherUsersOrders(t *testing.T) {
	setupTestDB(t)
	seedUser(t, "user1@test.com", "senha123")
	seedUser(t, "user2@test.com", "senha123")

	// Pedido do user1
	config.DB.Create(&models.Pedido{
		Nome: "User1", CPF: "111.111.111-11", Email: "user1@test.com",
		Telefone: "(11) 91234-5678", Endereco: "Rua A", Numero: "1",
		Bairro: "B1", CEP: "01000-000", Cidade: "SP", UF: "SP",
		Pagamento: "pix", Total: 50.00,
		Itens: []models.PedidoItem{{HatID: 1, Nome: "Boné", Price: 50.00, Quantidade: 1}},
	})

	// Pedido do user2
	config.DB.Create(&models.Pedido{
		Nome: "User2", CPF: "222.222.222-22", Email: "user2@test.com",
		Telefone: "(11) 91234-5678", Endereco: "Rua B", Numero: "2",
		Bairro: "B2", CEP: "02000-000", Cidade: "SP", UF: "SP",
		Pagamento: "boleto", Total: 75.00,
		Itens: []models.PedidoItem{{HatID: 2, Nome: "Chapéu", Price: 75.00, Quantidade: 1}},
	})

	// Buscar como user1 — deve ver apenas seu pedido
	token := generateTestToken(t, "user1@test.com")
	req := httptest.NewRequest("GET", "/api/meus-pedidos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	MeusPedidos(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rr.Code)
	}

	var pedidos []models.Pedido
	if err := json.Unmarshal(rr.Body.Bytes(), &pedidos); err != nil {
		t.Fatalf("Erro ao parsear resposta: %v", err)
	}

	if len(pedidos) != 1 {
		t.Errorf("Esperava 1 pedido para user1, recebeu %d", len(pedidos))
	}

	if len(pedidos) > 0 && pedidos[0].Email != "user1@test.com" {
		t.Errorf("Pedido retornado não pertence ao user1: %s", pedidos[0].Email)
	}
}
