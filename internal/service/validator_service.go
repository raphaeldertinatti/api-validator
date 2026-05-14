package service

import (
	"api-validator/internal/domains"
	"context"
	"fmt"
)

type NCMRepository interface {
	FindDescricaoCompleta(ctx context.Context, ncmCode string) (string, error)
}

type ValidatorService struct {
	ncmRepo      NCMRepository
	geminiClient *GeminiClient
}

func NewValidatorService(ncmRepo NCMRepository, apiKey string) *ValidatorService {
	return &ValidatorService{
		ncmRepo:      ncmRepo,
		geminiClient: NewGeminiClient(apiKey),
	}
}

func (s *ValidatorService) ValidateProduct(ctx context.Context, req domains.ValidateRequest) (*domains.ValidateResponse, error) {
	// 1. Buscar descrição do NCM no MongoDB via repositório
	ncmDescricao, err := s.ncmRepo.FindDescricaoCompleta(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCM: %w", err)
	}

	if ncmDescricao == "" {
		ncmDescricao = "NCM não encontrado na base de dados."
	}

	// 2. Validar com Gemini usando o novo cliente REST
	validacaoLLM, status, err := s.geminiClient.ValidateNCM(ctx, req.Descricao, ncmDescricao)
	if err != nil {
		return nil, fmt.Errorf("erro na validação com Gemini: %w", err)
	}

	return &domains.ValidateResponse{
		NCMDescricao:    ncmDescricao,
		ValidacaoLLM:    validacaoLLM,
		StatusValidacao: status,
	}, nil
}
