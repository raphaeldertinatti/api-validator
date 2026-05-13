package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Nomes das variáveis de ambiente — fonte única da verdade
const (
	envMongoURI     = "MONGO_URI"
	envDBName       = "DB_NAME"
	envGeminiAPIKey = "GEMINI_API_KEY"
	envPort         = "PORT"
)

// variáveis obrigatórias — validadas no Load()
var requiredEnvVars = []string{
	envMongoURI,
	envDBName,
	envGeminiAPIKey,
}

type Config struct {
	MongoURI     string
	DBName       string
	GeminiAPIKey string
	Port         string
}

// Load carrega a configuração a partir de variáveis de ambiente.
// Em desenvolvimento, lê do arquivo .env se existir.
// Retorna erro se alguma variável obrigatória estiver ausente.
func Load() (*Config, error) {
	loadDotEnv()

	if err := validateRequiredEnvVars(); err != nil {
		return nil, err
	}

	return &Config{
		MongoURI:     os.Getenv(envMongoURI),
		DBName:       os.Getenv(envDBName),
		GeminiAPIKey: os.Getenv(envGeminiAPIKey),
		Port:         getEnvWithDefault(envPort, "8080"),
	}, nil
}

// loadDotEnv tenta carregar o .env mas não falha se não existir —
// em produção as variáveis vêm do ambiente do sistema/container.
func loadDotEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Info("arquivo .env não encontrado, usando variáveis do sistema")
	}
}

// validateRequiredEnvVars verifica todas as variáveis obrigatórias de uma vez
// e retorna um erro consolidado listando todas as que estão faltando.
func validateRequiredEnvVars() error {
	var missing []string

	for _, key := range requiredEnvVars {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("variáveis de ambiente obrigatórias ausentes: %s",
			strings.Join(missing, ", "))
	}

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}
