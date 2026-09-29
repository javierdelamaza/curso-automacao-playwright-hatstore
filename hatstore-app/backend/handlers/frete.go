package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"hat-store-training/backend/config"
)

// Tabela de frete por faixa de CEP (região)
// Cada região tem um custo base por unidade
var tabelaFrete = []struct {
	Regiao       string
	CepInicio    int
	CepFim       int
	CustoPorItem float64
}{
	{"Sudeste", 1000000, 39999999, 12.50},
	{"Sul", 80000000, 99999999, 15.00},
	{"Nordeste", 40000000, 65999999, 22.00},
	{"Centro-Oeste", 70000000, 76999999, 20.00},
	{"Norte", 66000000, 69999999, 28.00},
	// Faixas complementares do Centro-Oeste
	{"Centro-Oeste", 77000000, 77999999, 20.00},
}

// FreteResponse é a resposta da consulta de frete
type FreteResponse struct {
	CEP          string  `json:"cep"`
	Regiao       string  `json:"regiao"`
	CustoPorItem float64 `json:"custo_por_item"`
	Quantidade   int     `json:"quantidade"`
	ValorFrete   float64 `json:"valor_frete"`
}

// Percentual adicional por unidade extra: 2ª=30%, 3ª=20%, 4ª em diante=10%
func percentualUnidade(n int) float64 {
	switch n {
	case 2:
		return 0.30
	case 3:
		return 0.20
	default:
		return 0.10
	}
}

// calcularFrete calcula o frete baseado no CEP e quantidade.
// A 1ª unidade paga o custo base. Cada unidade adicional soma um percentual do custo base.
func calcularFrete(cepNumerico int, quantidade int) (string, float64, float64, bool) {
	for _, faixa := range tabelaFrete {
		if cepNumerico >= faixa.CepInicio && cepNumerico <= faixa.CepFim {
			valor := faixa.CustoPorItem
			for i := 2; i <= quantidade; i++ {
				valor += faixa.CustoPorItem * percentualUnidade(i)
			}
			valor = math.Round(valor*100) / 100
			return faixa.Regiao, faixa.CustoPorItem, valor, true
		}
	}
	return "", 0, 0, false
}

// ConsultarFrete calcula o frete por CEP (requer new_hatstore ativo)
func ConsultarFrete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Verificar feature flag
	if !config.IsNewHatstoreEnabled() {
		http.Error(w, "Funcionalidade de frete desabilitada", http.StatusServiceUnavailable)
		return
	}

	cep := r.URL.Query().Get("cep")
	quantidadeStr := r.URL.Query().Get("quantidade")

	// Validar CEP
	if cep == "" {
		http.Error(w, "CEP é obrigatório", http.StatusBadRequest)
		return
	}

	cepRegex := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	if !cepRegex.MatchString(cep) {
		http.Error(w, "CEP inválido (formato: 00000-000)", http.StatusBadRequest)
		return
	}

	// Validar quantidade
	if quantidadeStr == "" {
		http.Error(w, "Quantidade é obrigatória", http.StatusBadRequest)
		return
	}

	quantidade, err := strconv.Atoi(quantidadeStr)
	if err != nil || quantidade < 1 {
		http.Error(w, "Quantidade deve ser um número inteiro maior que zero", http.StatusBadRequest)
		return
	}

	// Converter CEP para numérico (remove o hífen)
	cepLimpo := regexp.MustCompile(`\D`).ReplaceAllString(cep, "")
	cepNumerico, _ := strconv.Atoi(cepLimpo)

	regiao, custoPorItem, valorFrete, encontrou := calcularFrete(cepNumerico, quantidade)
	if !encontrou {
		log.Printf("[WARN] CEP não atendido: %s", cep)
		http.Error(w, "CEP não atendido pela nossa logística", http.StatusNotFound)
		return
	}

	resp := FreteResponse{
		CEP:          cep,
		Regiao:       regiao,
		CustoPorItem: custoPorItem,
		Quantidade:   quantidade,
		ValorFrete:   valorFrete,
	}

	json.NewEncoder(w).Encode(resp)
}
