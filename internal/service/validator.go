// service/validator.go
package service

import (
	"api-validator/internal/domains"
	"context"
)

type ValidatorService struct {
	ncm *NCMService
	// ipi      *IPIService
	// cest     *CESTService
	// pisCofins *PISCOFINSService
	// icms     *ICMSService
}

func NewValidatorService(ncmRepo NCMRepository, apiKey string) *ValidatorService {
	// Inicializa o serviço base de IA
	gemini := NewGeminiService(apiKey)

	// Inicializa os sub-serviços especializados
	ncm := NewNCMService(ncmRepo, gemini)

	return &ValidatorService{
		ncm: ncm,
	}
}

func (v *ValidatorService) Validate(ctx context.Context, req domains.ValidateRequest) (*domains.ValidateResponse, error) {
	// 1. NCM — base de tudo
	ncmResult, err := v.ncm.ValidateNCM(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domains.ValidateResponse{
		NCM: *ncmResult,
	}, nil
}
