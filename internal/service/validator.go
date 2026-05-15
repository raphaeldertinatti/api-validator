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

	// // 2. IPI — direto da TIPI
	// ipiResult, err := v.ipi.Validate(ctx, ncmCode)
	// if err != nil {
	// 	return nil, err
	// }

	// // 3. CEST — candidatos para ST
	// cestResult, err := v.cest.Identify(ctx, ncmCode, prodDesc)
	// if err != nil {
	// 	return nil, err
	// }

	// // 4. PIS/COFINS — CST entrada/saída
	// pisCofinsResult, err := v.pisCofins.Validate(ctx, ncmCode)
	// if err != nil {
	// 	return nil, err
	// }

	// // 5. ICMS — sequência: isenção → ST → redução → alíquota
	// icmsResult, err := v.icms.Validate(ctx, ncmCode, prodDesc, cestResult)
	// if err != nil {
	// 	return nil, err
	// }

	// return &domains.ValidationResult{
	// 	NCM:       *ncmResult,
	// 	IPI:       *ipiResult,
	// 	CEST:      *cestResult,
	// 	PISCofins: *pisCofinsResult,
	// 	ICMS:      *icmsResult,
	// }, nil

	return &domains.ValidateResponse{
		NCMDescricao:    ncmResult.NCMDescricao,
		ValidacaoLLM:    ncmResult.ValidacaoLLM,
		StatusValidacao: ncmResult.StatusValidacao,
	}, nil
}
