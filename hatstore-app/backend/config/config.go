package config

import (
	"log"
	"os"
)

// JWTKey é a chave secreta usada para assinar tokens JWT (obrigatória via env)
var JWTKey []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Em ambiente de teste, usa uma chave placeholder
		if os.Getenv("GO_TEST") == "1" || isTestBinary() {
			secret = "test-secret-key"
		} else {
			log.Fatal("[FATAL] JWT_SECRET environment variable not set")
		}
	}
	JWTKey = []byte(secret)
}

// isTestBinary detecta se o binário atual é um teste Go
func isTestBinary() bool {
	for _, arg := range os.Args {
		if len(arg) >= 5 && arg[:5] == "-test" {
			return true
		}
	}
	return false
}
