package main

import (
	"api-validator/internal/config"
	"api-validator/internal/repository"
	"api-validator/internal/server/handler"
	"api-validator/internal/service"
	"api-validator/internal/storage"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
)

func main() {
	// 0. Carregar Configurações (Lê do .env ou variáveis de sistema)
	cfg, err := config.Load()
	if err != nil {
		slog.Error("falha ao carregar configuração", "erro", err)
		os.Exit(1)
	}
	// 1. Inicializar Armazenamento (MongoDB)
	mongoDB, err := storage.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Falha ao conectar ao MongoDB: %v", err)
	}
	defer mongoDB.Disconnect(context.Background())

	// 2. Inicializar Repositórios
	ncmRepo := repository.NewNCMRepository(mongoDB)

	// 3. Inicializar Serviços (Lógica de Negócio)
	validatorService := service.NewValidatorService(ncmRepo, cfg.GeminiAPIKey)

	// 4. Inicializar Handler (Rotas e Injeção de Dependência)
	h := handler.NewHandler(validatorService)

	// 4. Iniciar Servidor HTTP
	router := h.InitRoutes()
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Servidor rodando em http://localhost%s\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
