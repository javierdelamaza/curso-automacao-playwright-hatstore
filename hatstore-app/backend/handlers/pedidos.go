// pedidos.go - Handlers auxiliares para pedidos (integração menor com rotas principais)
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"hat-store-training/backend/config"
	"hat-store-training/backend/models"
)

// Utilitários para validação e verificação de pedidos via banco de dados

// PedidoJaExiste verifica se um CPF já tem pedidos no banco
func PedidoJaExiste(cpf string) bool {
	var count int64
	config.GetDB().Model(&models.Pedido{}).Where("cpf = ?", cpf).Count(&count)
	return count > 0
}

// CPFUsouHatOff verifica se um CPF já utilizou o cupom HATOFF
func CPFUsouHatOff(cpf string) bool {
	var count int64
	config.GetDB().Model(&models.CupomHatOffUsado{}).Where("cpf = ?", cpf).Count(&count)
	return count > 0
}

// CriarPedido godoc
// @Summary Cria um novo pedido
// @Description Cria um novo pedido no banco de dados com validação completa, desconto de estoque e registro de auditoria
// @Tags pedido
// @Accept json
// @Produce json
// @Param X-Dev-UUID header string false "UUID de desenvolvimento"
// @Param pedido body models.Pedido true "Dados do pedido"
// @Success 201 {object} models.Pedido "Pedido criado com sucesso"
// @Failure 400 {string} string "Dados inválidos ou carrinho vazio"
// @Failure 409 {string} string "Estoque insuficiente ou cupom já utilizado"
// @Router /api/pedido [post]
// @Security Bearer
func CriarPedido(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENVIRONMENT") == "development" {
		uuid := r.Header.Get("X-Dev-UUID")
		expected := os.Getenv("DEV_UUID")
		if uuid == "" || uuid != expected {
			http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
			return
		}
	}

	var pedidoPayload struct {
		Nome      string `json:"nome"`
		CPF       string `json:"cpf"`
		Email     string `json:"email"`
		Telefone  string `json:"telefone"`
		Endereco  string `json:"endereco"`
		Numero    string `json:"numero"`
		Bairro    string `json:"bairro"`
		CEP       string `json:"cep"`
		Cidade    string `json:"cidade"`
		UF        string `json:"uf"`
		Pagamento string `json:"pagamento"`
		Itens     []struct {
			ID         uint    `json:"id"`
			Nome       string  `json:"nome"`
			Price      float64 `json:"price"`
			Quantidade int     `json:"quantidade"`
		} `json:"itens"`
		Cupom string `json:"cupom"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&pedidoPayload); err != nil {
		log.Printf("[ERROR] CriarPedido: Dados inválidos | Erro: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validar carrinho vazio IMEDIATAMENTE
	if len(pedidoPayload.Itens) == 0 {
		http.Error(w, "Seu carrinho está vazio. Adicione produtos antes de finalizar o pedido.", http.StatusBadRequest)
		return
	}

	// Validações - Campos vazios
	if strings.TrimSpace(pedidoPayload.Nome) == "" {
		http.Error(w, "Nome é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.CPF) == "" {
		http.Error(w, "CPF é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.Email) == "" {
		http.Error(w, "Email é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.Telefone) == "" {
		http.Error(w, "Telefone é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.Endereco) == "" {
		http.Error(w, "Endereço é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.Numero) == "" {
		http.Error(w, "Número é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.Bairro) == "" {
		http.Error(w, "Bairro é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.CEP) == "" {
		http.Error(w, "CEP é obrigatório", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.Cidade) == "" {
		http.Error(w, "Cidade é obrigatória", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(pedidoPayload.UF) == "" {
		http.Error(w, "UF é obrigatório", http.StatusBadRequest)
		return
	}

	// Validações de formato
	if len(pedidoPayload.Nome) < 3 {
		http.Error(w, "Nome deve ter pelo menos 3 caracteres", http.StatusBadRequest)
		return
	}

	cpfRegex := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	if !cpfRegex.MatchString(pedidoPayload.CPF) {
		http.Error(w, "CPF inválido (formato: 000.000.000-00)", http.StatusBadRequest)
		return
	}

	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !emailRegex.MatchString(pedidoPayload.Email) {
		http.Error(w, "Email inválido", http.StatusBadRequest)
		return
	}

	telRegex := regexp.MustCompile(`^\(\d{2}\) \d{5}-\d{4}$`)
	if !telRegex.MatchString(pedidoPayload.Telefone) {
		http.Error(w, "Telefone inválido (formato: (00) 00000-0000)", http.StatusBadRequest)
		return
	}

	cepRegex := regexp.MustCompile(`^\d{5}-\d{3}$`)
	if !cepRegex.MatchString(pedidoPayload.CEP) {
		http.Error(w, "CEP inválido (formato: 00000-000)", http.StatusBadRequest)
		return
	}

	ufRegex := regexp.MustCompile(`^[A-Za-z]{2}$`)
	if !ufRegex.MatchString(pedidoPayload.UF) {
		http.Error(w, "UF inválido (deve ter 2 letras)", http.StatusBadRequest)
		return
	}

	if pedidoPayload.Pagamento != "pix" && pedidoPayload.Pagamento != "boleto" {
		http.Error(w, "Forma de pagamento inválida (pix ou boleto)", http.StatusBadRequest)
		return
	}

	// Validações adicionais para tamanho mínimo
	if len(pedidoPayload.Endereco) < 3 {
		http.Error(w, "Endereço deve ter pelo menos 3 caracteres", http.StatusBadRequest)
		return
	}

	if len(pedidoPayload.Bairro) < 3 {
		http.Error(w, "Bairro deve ter pelo menos 3 caracteres", http.StatusBadRequest)
		return
	}

	if len(pedidoPayload.Cidade) < 3 {
		http.Error(w, "Cidade deve ter pelo menos 3 caracteres", http.StatusBadRequest)
		return
	}

	// Iniciar transação para garantir atomicidade
	tx := config.GetDB().Begin()
	if tx.Error != nil {
		log.Printf("[ERROR] Erro ao iniciar transação: %v", tx.Error)
		http.Error(w, "Erro ao processar pedido", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback automático se não houver Commit

	// Verificar uso do cupom HATOFF dentro da transação
	if pedidoPayload.Cupom == "HATOFF" {
		var cupomCount int64
		tx.Model(&models.CupomHatOffUsado{}).Where("cpf = ?", pedidoPayload.CPF).Count(&cupomCount)
		if cupomCount > 0 {
			http.Error(w, "Cupom HATOFF já utilizado por este CPF", http.StatusForbidden)
			return
		}
	}

	// Verificar e descontar estoque dentro da transação
	total := 0.0
	for _, item := range pedidoPayload.Itens {
		var hat models.Hat
		// Lock pessimista para evitar race condition no estoque
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", item.ID).First(&hat).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Chapéu não encontrado: "+item.Nome, http.StatusNotFound)
			} else {
				http.Error(w, "Erro ao buscar chapéu", http.StatusInternalServerError)
			}
			return
		}

		if hat.Quantidade < item.Quantidade {
			http.Error(w, "Quantidade insuficiente do chapéu: "+hat.Nome, http.StatusConflict)
			return
		}

		total += item.Price * float64(item.Quantidade)

		if err := tx.Model(&hat).Update("quantidade", gorm.Expr("quantidade - ?", item.Quantidade)).Error; err != nil {
			log.Printf("[ERROR] Erro ao atualizar estoque: %v", err)
			http.Error(w, "Erro ao processar pedido", http.StatusInternalServerError)
			return
		}
	}

	// Aplicar desconto HATOFF se for primeira compra
	if pedidoPayload.Cupom == "HATOFF" {
		if config.IsFixHatoffDiscountEnabled() {
			total = total * 0.80 // 20% de desconto (corrigido)
		} else {
			total = total * 0.75 // 25% de desconto (bug)
		}
	}

	// Calcular frete apenas no modelo novo (new_hatstore)
	valorFrete := 0.0
	if config.IsNewHatstoreEnabled() {
		var totalQuantidade int
		for _, item := range pedidoPayload.Itens {
			totalQuantidade += item.Quantidade
		}

		cepLimpo := regexp.MustCompile(`\D`).ReplaceAllString(pedidoPayload.CEP, "")
		cepNumerico, _ := strconv.Atoi(cepLimpo)

		_, _, frete, encontrou := calcularFrete(cepNumerico, totalQuantidade)
		if !encontrou {
			http.Error(w, "CEP não atendido pela nossa logística", http.StatusBadRequest)
			return
		}
		valorFrete = frete
	}

	total += valorFrete

	// Criar pedido
	pedido := models.Pedido{
		Nome:      pedidoPayload.Nome,
		CPF:       pedidoPayload.CPF,
		Email:     pedidoPayload.Email,
		Telefone:  pedidoPayload.Telefone,
		Endereco:  pedidoPayload.Endereco,
		Numero:    pedidoPayload.Numero,
		Bairro:    pedidoPayload.Bairro,
		CEP:       pedidoPayload.CEP,
		Cidade:    pedidoPayload.Cidade,
		UF:        pedidoPayload.UF,
		Pagamento: pedidoPayload.Pagamento,
		Total:     total,
		Frete:     valorFrete,
		Cupom:     pedidoPayload.Cupom,
	}

	if err := tx.Create(&pedido).Error; err != nil {
		log.Printf("[ERROR] Erro ao criar pedido: %v", err)
		http.Error(w, "Erro ao registrar pedido", http.StatusInternalServerError)
		return
	}

	// Criar itens do pedido (batch insert: 1 query em vez de N)
	var pedidoItems []models.PedidoItem
	for _, item := range pedidoPayload.Itens {
		pedidoItems = append(pedidoItems, models.PedidoItem{
			PedidoID:   pedido.ID,
			HatID:      int(item.ID),
			Nome:       item.Nome,
			Price:      item.Price,
			Quantidade: item.Quantidade,
		})
	}
	if err := tx.Create(&pedidoItems).Error; err != nil {
		log.Printf("[ERROR] Erro ao criar itens do pedido: %v", err)
		http.Error(w, "Erro ao registrar itens", http.StatusInternalServerError)
		return
	}

	// Registrar uso do cupom HATOFF se aplicável
	if pedidoPayload.Cupom == "HATOFF" {
		cupomUso := models.CupomHatOffUsado{
			CPF: pedidoPayload.CPF,
		}
		if err := tx.Create(&cupomUso).Error; err != nil {
			log.Printf("[ERROR] Erro ao registrar uso do cupom: %v", err)
			http.Error(w, "Erro ao registrar cupom", http.StatusInternalServerError)
			return
		}
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] Erro ao confirmar transação: %v", err)
		http.Error(w, "Erro ao finalizar pedido", http.StatusInternalServerError)
		return
	}

	// Invalidar cache de chapéus (estoque alterado)
	InvalidateHatsCache()

	ip := config.GetClientIP(r)

	// Registrar auditoria
	config.LogAuditForRequest(r, "pedido_criado", pedidoPayload.Email, map[string]interface{}{
		"pedido_id": pedido.ID,
		"cpf":       pedido.CPF,
		"total":     pedido.Total,
		"itens":     len(pedidoPayload.Itens),
		"cupom":     pedidoPayload.Cupom,
		"pagamento": pedidoPayload.Pagamento,
	}, ip)

	log.Printf("[INFO][%s] Pedido criado | Nome: %s | CPF: %s | Total: %.2f | Pagamento: %s | Itens: %d", ip, pedido.Nome, pedido.CPF, pedido.Total, pedido.Pagamento, len(pedidoPayload.Itens))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pedido)
}
