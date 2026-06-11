// service/validator.go
package service

import (
	"api-validator/internal/domains"
	"context"
)

type ValidatorService struct {
	ncm         *NCMService
	ipi         *IPIService
	cest        *CESTService
	piscofins   *PISCOFINSService
	isencao     *IsencaoService
	diferimento *DiferimentoService
	aliquota    *AliquotaService
	reducao     *ReducaoService
}

func NewValidatorService(
	ncmRepo NCMRepository,
	ipiRepo IPIRepository,
	cestRepo CESTRepository,
	piscofinsRepo PISCOFINSRepository,
	isencaoRepo IsencaoRepository,
	diferimentoRepo DiferimentoRepository,
	aliquotaRepo AliquotaRepository,
	reducaoRepo ReducaoRepository,
	apiKey string,
) *ValidatorService {
	// Inicializa o serviço base de IA
	gemini := NewGeminiService(apiKey)

	// Inicializa os sub-serviços especializados
	ncm := NewNCMService(ncmRepo, gemini)
	ipi := NewIPIService(ncmRepo, ipiRepo, gemini)
	cest := NewCESTService(cestRepo, gemini)
	piscofins := NewPISCOFINSService(piscofinsRepo)
	isencao := NewIsencaoService(ncmRepo, isencaoRepo, gemini)
	diferimento := NewDiferimentoService(ncmRepo, diferimentoRepo, gemini)
	aliquota := NewAliquotaService(ncmRepo, aliquotaRepo, gemini)
	reducao := NewReducaoService(ncmRepo, reducaoRepo, gemini) // Passa nil por enquanto, o MongoDB será injetado depois

	return &ValidatorService{
		ncm:         ncm,
		ipi:         ipi,
		cest:        cest,
		piscofins:   piscofins,
		isencao:     isencao,
		diferimento: diferimento,
		aliquota:    aliquota,
		reducao:     reducao,
	}
}

func (v *ValidatorService) Validate(ctx context.Context, req domains.ValidateRequest) (*domains.ValidateResponse, error) {
	// 1. NCM — base de tudo
	ncmResult, err := v.ncm.ValidateNCM(ctx, req)
	if err != nil {
		return nil, err
	}
	if !ncmResult.RetornoLLM.Compativel {
		// Se a NCM não for compatível, retorna apenas o NCM e omite o resto
		return &domains.ValidateResponse{
			NCM: ncmResult,
		}, nil
	}

	// 2. IPI — depende do NCM
	ipiResult, err := v.ipi.ValidateIPI(ctx, req)
	if err != nil {
		return nil, err
	}

	// 3. CEST — depende do NCM
	cestResult, err := v.cest.ValidateCEST(ctx, req)
	if err != nil {
		return nil, err
	}

	// 4. PIS/COFINS — depende do NCM
	piscofinsResult, err := v.piscofins.ValidatePISCOFINS(ctx, req)
	if err != nil {
		return nil, err
	}

	// 5. Isenção — depende do NCM
	isencaoResult, err := v.isencao.ValidateIsencao(ctx, req)
	if err != nil {
		return nil, err
	}

	// 6. Diferimento — depende do NCM
	diferimentoResult, err := v.diferimento.ValidateDiferimento(ctx, req)
	if err != nil {
		return nil, err
	}

	// 7. Se for Isento retorna a response e não valida mais nada, isenção prevalece sobre os outros impostos.
	if isencaoResult.Isento {
		return &domains.ValidateResponse{
			NCM:         ncmResult,
			IPI:         ipiResult,
			CEST:        cestResult,
			PISCOFINS:   piscofinsResult,
			Isencao:     isencaoResult,
			Diferimento: diferimentoResult,
		}, nil
	} else if diferimentoResult.Diferido { // 8. Se for Diferido retorna a response e não valida mais nada. Diferimento prevalece sobre os outros impostos (exceto isenção).
		return &domains.ValidateResponse{
			NCM:         ncmResult,
			IPI:         ipiResult,
			CEST:        cestResult,
			PISCOFINS:   piscofinsResult,
			Isencao:     isencaoResult,
			Diferimento: diferimentoResult,
		}, nil
	}

	// 9. Aliquota — depende do NCM
	aliquotaResult, err := v.aliquota.ValidateAliquota(ctx, req)
	if err != nil {
		return nil, err
	}

	// 10. Redução — depende do NCM e da Alíquota
	reducaoResult, err := v.reducao.ValidateReducao(*aliquotaResult, ctx, req)
	if err != nil {
		return nil, err
	}

	return &domains.ValidateResponse{
		NCM:         ncmResult,
		IPI:         ipiResult,
		CEST:        cestResult,
		PISCOFINS:   piscofinsResult,
		Isencao:     isencaoResult,
		Diferimento: diferimentoResult,
		Aliquota:    aliquotaResult,
		Reducao:     reducaoResult,
	}, nil
}
