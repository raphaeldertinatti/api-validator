package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvMongoURI     = "MONGO_URI"
	EnvDBName       = "DB_NAME"
	EnvGeminiAPIKey = "GEMINI_API_KEY"
	EnvPort         = "PORT"
)

type Config struct {
	MongoURI     string
	DBName       string
	GeminiAPIKey string
	Port         string
}

func Load() *Config {
	// Tenta carregar o arquivo .env, mas não falha se não existir (ex: em produção)
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	cfg := &Config{
		MongoURI:     getEnv(EnvMongoURI, "mongodb://localhost:27017"),
		DBName:       getEnv(EnvDBName, "Validator"),
		GeminiAPIKey: getEnv(EnvGeminiAPIKey, "AIzaSyDx6U4iE9oBg0yj-7v3nVfs6tK49Rfsrq4"), // Valor padrão para desenvolvimento, mas deve ser sobrescrito em produção
		Port:         getEnv(EnvPort, "8080"),
	}

	if cfg.GeminiAPIKey == "" {
		log.Fatal("ERRO: A variável GEMINI_API_KEY é obrigatória.")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
