package main

import (
	"fmt"
	"hat-store-training/backend/config"
	_ "hat-store-training/backend/docs" // Importa os docs gerados pelo swag
	"hat-store-training/backend/routes"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Hat Store API
// @version 1.0
// @description API da loja de chapéus Hat Store.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Token JWT no formato: Bearer {token}

func main() {
	// 0. Inicializa o banco de dados PostgreSQL
	if err := config.InitDatabase(); err != nil {
		log.Fatal("[FATAL] Falha ao inicializar banco de dados:", err)
	}

	// 0.1 Carrega configuração de auditoria do banco para memória
	config.LoadAuditSetting()

	// 0.2 Inicia sincronização periódica da flag de auditoria (a cada 30s)
	config.StartAuditSettingSyncJob()

	// 0.3 Carrega feature flags do banco para memória
	config.LoadFeatureFlags()

	// 0.4 Inicia sincronização periódica das feature flags (a cada 10min)
	config.StartFeatureFlagsSyncJob()

	// 0.5 Inicia buffer de audit logs (batch insert a cada 5s)
	config.StartAuditBufferJob()

	// 0.6 Inicia job de limpeza de logs de auditoria
	go startAuditLogCleanupJob()

	// 1. Inicializa o nosso roteador gorilla/mux que já conhece TODAS as rotas (/api, /auth)
	r := routes.InitializeRoutes()

	// 2. Registra a rota da documentação Swagger DIRETAMENTE no roteador gorilla/mux
	// O PathPrefix garante que ele capture todos os subcaminhos como /swagger/index.html
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 3. Define a rota "catch-all" para servir os arquivos do frontend.
	// O gorilla/mux é inteligente: ele primeiro tenta casar as rotas mais específicas (/api, /auth, /swagger)
	// Se nenhuma casar, ele usa esta regra genérica.
	// O noCacheMiddleware adiciona headers para evitar cache do browser nos arquivos estáticos.
	fs := http.FileServer(http.Dir("./frontend"))
	r.PathPrefix("/").Handler(noCacheMiddleware(featureFlagInjector(fs)))

	// 4. Inicia o servidor e diz a ele para usar o roteador 'r' para TODAS as requisições.
	// A MUDANÇA MAIS IMPORTANTE: Passamos 'r' em vez de 'nil' para ListenAndServe.
	log.Println("Servidor iniciado em http://0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", platformGuard(r)); err != nil {
		log.Fatal(err)
	}
}

// noCacheMiddleware adiciona headers Cache-Control para que o browser sempre revalide
// os arquivos estáticos do frontend após cada deploy.
func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// featureFlagInjector injeta as feature flags no HTML como uma tag <script> antes de </head>.
// Isso permite que o frontend acesse as flags via window.FEATURE_FLAGS sem precisar de um endpoint.
func featureFlagInjector(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || strings.HasSuffix(path, ".html") {
			filePath := "./frontend" + path
			if path == "/" {
				filePath = "./frontend/index.html"
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				fs.ServeHTTP(w, r)
				return
			}

			platformEnabled := config.IsPlatformEnabled()
			var flagsScript string
			if !platformEnabled {
				flagsScript = fmt.Sprintf(
					`<script>window.FEATURE_FLAGS={filters_enabled:%t,platform_enabled:false};`+
						`document.addEventListener('DOMContentLoaded',function(){`+
						`document.body.style.margin='0';document.body.style.padding='0';document.body.style.background='#181818';`+
						`document.body.innerHTML='<div style="background:#181818;color:#c9a94b;display:flex;justify-content:center;align-items:center;min-height:100vh;font-family:Segoe UI,Arial,sans-serif">`+
						`<div style="text-align:center;padding:2rem">`+
						`<div style="font-size:4rem;margin-bottom:1.5rem">🎩</div>`+
						`<h1 style="font-size:2.5rem;margin-bottom:1rem">Em manutenção</h1>`+
						`<p style="font-size:1.2rem;color:#aaa;margin-bottom:.5rem">A Hat Store está temporariamente fora do ar.</p>`+
						`<p style="font-size:1.2rem;color:#aaa">Voltamos em breve!</p>`+
						`</div></div>';});</script>`,
					config.IsFiltersEnabled(),
				)
			} else {
				flagsScript = fmt.Sprintf(
					`<script>window.FEATURE_FLAGS={filters_enabled:%t,platform_enabled:true,new_hatstore:%t};</script>`,
					config.IsFiltersEnabled(),
					config.IsNewHatstoreEnabled(),
				)
			}

			modified := strings.Replace(string(content), "</head>", flagsScript+"\n</head>", 1)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(modified))
			return
		}

		fs.ServeHTTP(w, r)
	})
}

// platformGuard bloqueia rotas de API e autenticação quando platform_enabled=false,
// retornando 503 com JSON. Páginas HTML passam pelo featureFlagInjector que mostra manutenção no frontend.
func platformGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.IsPlatformEnabled() {
			path := r.URL.Path
			if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/auth") || strings.HasPrefix(path, "/swagger") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error":"Site em manutenção"}`))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// startAuditLogCleanupJob inicia um job que limpa logs de auditoria a cada 24 horas
func startAuditLogCleanupJob() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Executar limpeza uma vez na inicialização
	config.CleanOldAuditLogs()

	for range ticker.C {
		config.CleanOldAuditLogs()
	}
}
