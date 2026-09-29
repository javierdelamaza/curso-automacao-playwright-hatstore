package routes

import (
	"hat-store-training/backend/config"
	"hat-store-training/backend/handlers"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// Middleware para logar cada request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log padronizado: método, rota, IP, usuário (se disponível)
		user := r.Header.Get("X-User-Email")
		ip := config.GetClientIP(r)
		if user != "" {
			log.Printf("[INFO][%s][%s] %s %s", ip, user, r.Method, r.URL.Path)
		} else {
			log.Printf("[INFO][%s] %s %s", ip, r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

// authMiddleware verifica o token JWT
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("[ERRO] Autenticação falhou: Cabeçalho de autorização ausente em %s %s", r.Method, r.URL.Path)
			http.Error(w, "Cabeçalho de autorização ausente", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey, nil
		})

		if err != nil || !token.Valid {
			log.Printf("[ERRO] Token inválido: %v | Rota: %s %s", err, r.Method, r.URL.Path)
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func InitializeRoutes() *mux.Router {
	// 1. Roteador Principal
	router := mux.NewRouter()
	// Aplica o logging em TODAS as requisições
	router.Use(loggingMiddleware)

	// 2. Sub-roteador para rotas de AUTENTICAÇÃO (PÚBLICAS)
	// Caminho: /auth/...
	// Middlewares: Nenhum (além do logging)
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	authRouter.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// 3. Sub-roteador para rotas da API públicas
	apiPublic := router.PathPrefix("/api").Subrouter()
	apiPublic.HandleFunc("/hats", handlers.GetHats).Methods("GET")
	apiPublic.HandleFunc("/estoque", handlers.ListarEstoque).Methods("GET")
	apiPublic.HandleFunc("/hats/{id}/estoque", handlers.EditarEstoqueHat).Methods("PATCH")
	apiPublic.HandleFunc("/frete", handlers.ConsultarFrete).Methods("GET")

	// 4. Sub-roteador para rotas da API protegidas
	apiProtected := router.PathPrefix("/api").Subrouter()
	apiProtected.Use(authMiddleware)
	apiProtected.HandleFunc("/me", handlers.GetProfile).Methods("GET")
	apiProtected.HandleFunc("/pedido", handlers.CriarPedido).Methods("POST")
	apiProtected.HandleFunc("/pedidos", handlers.ConsultarPedidos).Methods("GET")
	apiProtected.HandleFunc("/meus-pedidos", handlers.MeusPedidos).Methods("GET")

	return router
}
