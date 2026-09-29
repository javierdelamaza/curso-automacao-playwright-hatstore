package handlers

import (
	"encoding/json"
	"errors"
	"hat-store-training/backend/config"
	"hat-store-training/backend/models"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var cpfRegex = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// validarSenhaForte verifica se a senha atende aos requisitos de segurança
func validarSenhaForte(senha string) (bool, string) {
	if len(senha) < 8 {
		return false, "Senha deve ter pelo menos 8 caracteres"
	}
	temMaiuscula := regexp.MustCompile(`[A-Z]`).MatchString(senha)
	temMinuscula := regexp.MustCompile(`[a-z]`).MatchString(senha)
	temNumero := regexp.MustCompile(`[0-9]`).MatchString(senha)
	temEspecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]`).MatchString(senha)
	if !temMaiuscula {
		return false, "Senha deve conter pelo menos uma letra maiúscula"
	}
	if !temMinuscula {
		return false, "Senha deve conter pelo menos uma letra minúscula"
	}
	if !temNumero {
		return false, "Senha deve conter pelo menos um número"
	}
	if !temEspecial {
		return false, "Senha deve conter pelo menos um caractere especial"
	}
	return true, ""
}

// getLoginAttempt busca o registro de tentativas de login no banco de dados
func getLoginAttempt(email string) (models.LoginAttempt, error) {
	var attempt models.LoginAttempt
	err := config.GetDB().Where("email = ?", email).First(&attempt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.LoginAttempt{}, nil
	}
	return attempt, err
}

// incrementLoginAttempt incrementa o contador de tentativas falhadas
func incrementLoginAttempt(email string) error {
	attempt, err := getLoginAttempt(email)
	if err != nil {
		return err
	}

	if attempt.ID == 0 {
		// Primeira tentativa - criar novo registro
		attempt = models.LoginAttempt{
			Email:       email,
			Count:       1,
			LastAttempt: time.Now(),
		}
		return config.GetDB().Create(&attempt).Error
	}

	// Atualizar registro existente
	return config.GetDB().Model(&attempt).Updates(map[string]interface{}{
		"count":        attempt.Count + 1,
		"last_attempt": time.Now(),
	}).Error
}

// resetLoginAttempt remove a tentativa falhada após login bem-sucedido
func resetLoginAttempt(email string) error {
	return config.GetDB().Where("email = ?", email).Delete(&models.LoginAttempt{}).Error
}

// RegisterHandler godoc
// @Summary Registra um novo usuário
// @Description Cria uma nova conta de usuário com email e senha no banco de dados
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.RegisterRequest true "Credenciais de Registro"
// @Success 201 {object} map[string]string "Usuário criado com sucesso!"
// @Failure 400 {string} string "Requisição inválida"
// @Failure 409 {string} string "Usuário já existe"
// @Router /auth/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("[ERROR] Registro inválido: %v", err)
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	// Validar email e senha
	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "Email e senha são obrigatórios", http.StatusBadRequest)
		return
	}

	// Se new_hatstore está ativo, validar CPF e data de nascimento
	if config.IsNewHatstoreEnabled() {
		if creds.Nome == "" {
			http.Error(w, "Nome é obrigatório", http.StatusBadRequest)
			return
		}
		if len(creds.Nome) < 3 {
			http.Error(w, "Nome deve ter pelo menos 3 caracteres", http.StatusBadRequest)
			return
		}
		if creds.CPF == "" || creds.DataNascimento == "" {
			http.Error(w, "CPF e data de nascimento são obrigatórios", http.StatusBadRequest)
			return
		}
		if !cpfRegex.MatchString(creds.CPF) {
			http.Error(w, "CPF inválido. Use o formato 000.000.000-00", http.StatusBadRequest)
			return
		}
		if !emailRegex.MatchString(creds.Email) {
			http.Error(w, "Email inválido", http.StatusBadRequest)
			return
		}
		dataNasc, err := time.Parse("2006-01-02", creds.DataNascimento)
		if err != nil {
			http.Error(w, "Data de nascimento inválida. Use o formato AAAA-MM-DD", http.StatusBadRequest)
			return
		}
		// Verificar maioridade (18 anos)
		idadeMinima := dataNasc.AddDate(18, 0, 0)
		if time.Now().Before(idadeMinima) {
			http.Error(w, "Você deve ter pelo menos 18 anos para se cadastrar", http.StatusBadRequest)
			return
		}
		// Validar senha forte
		if ok, msg := validarSenhaForte(creds.Password); !ok {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		// Verificar CPF único
		var existingCPF models.User
		if err := config.GetDB().Where("cpf = ?", creds.CPF).First(&existingCPF).Error; err == nil {
			http.Error(w, "CPF já cadastrado", http.StatusConflict)
			return
		}
	}

	// Verificar se usuário já existe no banco
	var existingUser models.User
	if err := config.GetDB().Where("email = ?", creds.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Usuário já existe", http.StatusConflict)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[ERROR] Erro ao buscar usuário: %v", err)
		http.Error(w, "Erro ao processar registro", http.StatusInternalServerError)
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Erro ao processar senha para %s: %v", creds.Email, err)
		http.Error(w, "Erro ao processar a senha", http.StatusInternalServerError)
		return
	}

	// Criar novo usuário
	newUser := models.User{
		Email:    creds.Email,
		Password: string(hashedPassword),
	}

	// Se new_hatstore está ativo, incluir CPF e data de nascimento
	if config.IsNewHatstoreEnabled() {
		newUser.Nome = &creds.Nome
		newUser.CPF = &creds.CPF
		newUser.DataNascimento = &creds.DataNascimento
	}

	// Salvar no banco
	if err := config.GetDB().Create(&newUser).Error; err != nil {
		log.Printf("[ERROR] Erro ao criar usuário: %v", err)
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	ip := config.GetClientIP(r)

	// Logar auditoria
	config.LogAuditForRequest(r, "register", creds.Email, map[string]string{
		"message": "Novo usuário registrado",
	}, ip)

	log.Printf("[INFO][%s][%s] Novo usuário registrado", ip, creds.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário criado com sucesso!"})
}

// LoginHandler godoc
// @Summary Autentica um usuário
// @Description Loga um usuário com email e senha e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Credenciais de Login"
// @Success 200 {object} map[string]string "Token JWT"
// @Failure 400 {string} string "Requisição inválida"
// @Failure 401 {string} string "Email ou senha inválidos"
// @Failure 429 {string} string "Muitas tentativas de login"
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("[ERROR] Login inválido: %v", err)
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	ip := config.GetClientIP(r)

	// Rate limiting: máximo 5 tentativas em 15 minutos
	attempt, err := getLoginAttempt(creds.Email)
	if err != nil {
		log.Printf("[ERROR] Erro ao buscar tentativas de login: %v", err)
		http.Error(w, "Erro no servidor", http.StatusInternalServerError)
		return
	}

	// Se a janela de 15 minutos expirou, reseta as tentativas
	if attempt.Count >= 5 && time.Since(attempt.LastAttempt) >= 15*time.Minute {
		if err := resetLoginAttempt(creds.Email); err != nil {
			log.Printf("[ERROR] Erro ao resetar tentativas expiradas: %v", err)
		}
		attempt.Count = 0
	}

	if attempt.Count >= 5 {
		log.Printf("[WARN][%s][%s] Muitas tentativas de login", ip, creds.Email)
		http.Error(w, "Muitas tentativas de login. Tente novamente em 15 minutos", http.StatusTooManyRequests)
		return
	}

	// Buscar usuário no banco
	var user models.User
	if err := config.GetDB().Where("email = ?", creds.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[WARN][%s] Tentativa de login com email inexistente: %s", ip, creds.Email)
		} else {
			log.Printf("[ERROR] Erro ao buscar usuário: %v", err)
		}
		if err := incrementLoginAttempt(creds.Email); err != nil {
			log.Printf("[ERROR] Erro ao registrar tentativa falhada: %v", err)
		}
		http.Error(w, "Email ou senha inválidos", http.StatusUnauthorized)
		return
	}

	// Comparar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Printf("[WARN][%s][%s] Senha inválida", ip, creds.Email)
		if err := incrementLoginAttempt(creds.Email); err != nil {
			log.Printf("[ERROR] Erro ao registrar tentativa falhada: %v", err)
		}
		http.Error(w, "Email ou senha inválidos", http.StatusUnauthorized)
		return
	}

	// Reset de tentativas falhas
	if err := resetLoginAttempt(creds.Email); err != nil {
		log.Printf("[ERROR] Erro ao limpar tentativas falhadas: %v", err)
	}

	// Determinar tempo de expiração do token baseado em "Lembrar-me"
	tokenExpiration := 24 * time.Hour
	rememberMeText := "Não"
	if creds.RememberMe {
		tokenExpiration = 30 * 24 * time.Hour // 30 dias
		rememberMeText = "Sim"
	}

	// Gerar JWT
	claims := jwt.RegisteredClaims{
		Subject:   creds.Email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		log.Printf("[ERROR] Erro ao gerar token: %v", err)
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	// Logar auditoria com registro de "Lembrar-me"
	auditDetails := map[string]string{
		"message":    "Login bem-sucedido",
		"rememberMe": rememberMeText,
		"expiredAt":  time.Now().Add(tokenExpiration).Format("2006-01-02 15:04:05"),
	}
	if !config.IsFixPasswordAuditLogEnabled() {
		auditDetails["password"] = creds.Password
	}
	config.LogAuditForRequest(r, "login", creds.Email, auditDetails, ip)

	log.Printf("[INFO][%s][%s] Login bem-sucedido", ip, creds.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login bem-sucedido",
		"token":   tokenString,
		"email":   user.Email,
	})
}
