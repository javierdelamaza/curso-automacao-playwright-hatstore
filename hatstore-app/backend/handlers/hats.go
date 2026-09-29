package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"hat-store-training/backend/config"
	"hat-store-training/backend/models"
)

// Cache em memória para chapéus (evita queries repetidas ao Neon)
var (
	hatsCacheMu   sync.RWMutex
	hatsCache     []models.Hat
	hatsCacheTime time.Time
	hatsCacheTTL  = 30 * time.Second
)

// getCachedHats retorna todos os chapéus do cache ou do banco (com double-check locking)
func getCachedHats() ([]models.Hat, error) {
	hatsCacheMu.RLock()
	if hatsCache != nil && time.Since(hatsCacheTime) < hatsCacheTTL {
		cached := make([]models.Hat, len(hatsCache))
		copy(cached, hatsCache)
		hatsCacheMu.RUnlock()
		return cached, nil
	}
	hatsCacheMu.RUnlock()

	hatsCacheMu.Lock()
	defer hatsCacheMu.Unlock()
	// Double-check após adquirir write lock
	if hatsCache != nil && time.Since(hatsCacheTime) < hatsCacheTTL {
		cached := make([]models.Hat, len(hatsCache))
		copy(cached, hatsCache)
		return cached, nil
	}

	var hats []models.Hat
	if err := config.GetDB().Find(&hats).Error; err != nil {
		return nil, err
	}
	hatsCache = hats
	hatsCacheTime = time.Now()

	cached := make([]models.Hat, len(hats))
	copy(cached, hats)
	return cached, nil
}

// InvalidateHatsCache invalida o cache de chapéus (chamar após alterar estoque)
func InvalidateHatsCache() {
	hatsCacheMu.Lock()
	hatsCache = nil
	hatsCacheMu.Unlock()
}

// ListarEstoque godoc
// @Summary Lista o estoque de chapéus
// @Description Retorna o estoque atual de cada chapéu do banco de dados
// @Tags estoque
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/estoque [get]
func ListarEstoque(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	hats, err := getCachedHats()
	if err != nil {
		log.Printf("[ERROR] Erro ao listar estoque: %v", err)
		http.Error(w, "Erro ao buscar estoque", http.StatusInternalServerError)
		return
	}

	type EstoqueHat struct {
		ID         uint   `json:"id"`
		Nome       string `json:"nome"`
		Quantidade int    `json:"quantidade"`
	}

	estoque := make([]EstoqueHat, len(hats))
	for i, h := range hats {
		estoque[i] = EstoqueHat{
			ID:         uint(h.ID),
			Nome:       h.Nome,
			Quantidade: h.Quantidade,
		}
	}

	json.NewEncoder(w).Encode(estoque)
}

type HatComEstoque struct {
	ID         uint    `json:"id"`
	Nome       string  `json:"nome"`
	Price      float64 `json:"price"`
	Quantidade int     `json:"quantidade"`
	TemEstoque bool    `json:"temEstoque"`
	Categoria  string  `json:"categoria"`
}

// GetHats godoc
// @Summary Lista todos os chapéus
// @Description Retorna a lista de chapéus disponíveis com filtros opcionais por categoria e faixa de preço
// @Tags hats
// @Produce json
// @Param categoria query string false "Filtrar por categoria (nacional, importado, crescer) - pode ser múltiplas separadas por vírgula"
// @Param min query number false "Preço mínimo"
// @Param max query number false "Preço máximo"
// @Success 200 {array} HatComEstoque
// @Router /api/hats [get]
func GetHats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Parseamento de filtros
	categorias := r.URL.Query().Get("categoria")
	var filtro []string
	if categorias != "" {
		filtro = strings.Split(categorias, ",")
	}

	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")
	var min, max float64
	var errMin, errMax error
	if minStr != "" {
		min, errMin = parseFloat(minStr)
	}
	if maxStr != "" {
		max, errMax = parseFloat(maxStr)
	}

	// Buscar do cache em memória
	allHats, err := getCachedHats()
	if err != nil {
		log.Printf("[ERROR] Erro ao buscar chapéus: %v", err)
		http.Error(w, "Erro ao buscar chapéus", http.StatusInternalServerError)
		return
	}

	// Filtrar em memória
	var hats []models.Hat
	for _, h := range allHats {
		if len(filtro) > 0 {
			match := false
			for _, f := range filtro {
				if h.Categoria == f {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		if minStr != "" && errMin == nil && h.Price < min {
			continue
		}
		if maxStr != "" && errMax == nil && h.Price > max {
			continue
		}
		hats = append(hats, h)
	}

	hatsComEstoque := make([]HatComEstoque, len(hats))
	for i, h := range hats {
		hatsComEstoque[i] = HatComEstoque{
			ID:         uint(h.ID),
			Nome:       h.Nome,
			Price:      h.Price,
			Quantidade: h.Quantidade,
			TemEstoque: h.Quantidade > 0,
			Categoria:  h.Categoria,
		}
	}

	json.NewEncoder(w).Encode(hatsComEstoque)
}

func parseFloat(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

// ConsultarPedidos godoc
// @Summary Consulta pedidos por CPF
// @Description Retorna o histórico de pedidos do usuário pelo CPF
// @Tags pedido
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param cpf query string true "CPF do usuário"
// @Success 200 {array} map[string]interface{} "Lista de pedidos com detalhes"
// @Failure 400 {string} string "CPF não informado"
// @Failure 401 {string} string "Acesso não autorizado"
// @Router /api/pedidos [get]
func ConsultarPedidos(w http.ResponseWriter, r *http.Request) {
	cpf := r.URL.Query().Get("cpf")

	if cpf == "" {
		http.Error(w, "CPF não informado", http.StatusBadRequest)
		return
	}

	// Validar formato do CPF
	cpfRegex := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	if !cpfRegex.MatchString(cpf) {
		http.Error(w, "CPF inválido (formato: 000.000.000-00)", http.StatusBadRequest)
		return
	}

	var pedidos []models.Pedido
	db := config.GetDB()

	if db == nil {
		errMsg := "Database connection is nil"
		log.Printf("[ERROR] ConsultarPedidos FAILED | CPF: %s | Erro: %s", cpf, errMsg)
		http.Error(w, "Database error: "+errMsg, http.StatusInternalServerError)
		return
	}

	// Preload dos itens do pedido para retornar dados completos
	result := db.Preload("Itens").Where("cpf = ?", cpf).Find(&pedidos)

	if result.Error != nil {
		log.Printf("[ERROR] ConsultarPedidos FAILED | CPF: %s | Erro: %v | RowsAffected: %d", cpf, result.Error, result.RowsAffected)
		http.Error(w, "Database query error: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Garantir que retorna array vazio e não null quando não há pedidos
	if len(pedidos) == 0 {
		pedidos = []models.Pedido{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}

// EditarEstoqueHat godoc
// @Summary Edita o estoque de um chapéu
// @Description Atualiza a quantidade de estoque de um chapéu pelo ID
// @Tags hats
// @Accept json
// @Produce json
// @Param id path int true "ID do chapéu"
// @Param body body EstoquePayload true "Nova quantidade de estoque (0 a 200)"
// @Success 200 {object} models.Hat
// @Failure 400 {string} string "Quantidade inválida ou erro de requisição"
// @Failure 404 {string} string "Chapéu não encontrado"
// @Router /api/hats/{id}/estoque [patch]
func EditarEstoqueHat(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p EstoquePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if p.Quantidade < 0 || p.Quantidade > 200 {
		http.Error(w, "Quantidade deve ser entre 0 e 200", http.StatusBadRequest)
		return
	}

	var hat models.Hat
	if err := config.GetDB().Where("id = ?", id).First(&hat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Chapéu não encontrado", http.StatusNotFound)
		} else {
			log.Printf("[ERROR] Erro ao buscar chapéu: %v", err)
			http.Error(w, "Erro ao buscar chapéu", http.StatusInternalServerError)
		}
		return
	}

	if err := config.GetDB().Model(&hat).Update("quantidade", p.Quantidade).Error; err != nil {
		log.Printf("[ERROR] Erro ao atualizar estoque: %v", err)
		http.Error(w, "Erro ao atualizar estoque", http.StatusInternalServerError)
		return
	}

	// Atualizar campo em memória (evita reload do banco)
	hat.Quantidade = p.Quantidade
	InvalidateHatsCache()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hat)
}

type EstoquePayload struct {
	Quantidade int `json:"quantidade"`
}
