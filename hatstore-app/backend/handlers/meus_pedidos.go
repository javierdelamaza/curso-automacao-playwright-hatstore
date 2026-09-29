package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"hat-store-training/backend/config"
	"hat-store-training/backend/models"

	"github.com/golang-jwt/jwt/v5"
)

// MeusPedidos retorna os pedidos do usuário logado (requer new_hatstore ativo)
func MeusPedidos(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWTKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	email := claims.Subject

	var pedidos []models.Pedido
	result := config.GetDB().Preload("Itens").Where("email = ?", email).Order("created_at DESC").Find(&pedidos)
	if result.Error != nil {
		log.Printf("[ERROR] MeusPedidos: erro ao buscar pedidos para %s: %v", email, result.Error)
		http.Error(w, "Erro ao buscar pedidos", http.StatusInternalServerError)
		return
	}

	if len(pedidos) == 0 {
		pedidos = []models.Pedido{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}
