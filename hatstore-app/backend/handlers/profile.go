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

// GetProfile retorna o perfil do usuário logado (requer new_hatstore ativo)
func GetProfile(w http.ResponseWriter, r *http.Request) {
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

	var user models.User
	if err := config.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("[ERROR] GetProfile: usuário não encontrado: %s", email)
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
