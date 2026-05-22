package service

import (
	"api-validator/internal/domains"
	"context"
	"fmt"
)

type PISCOFINSRepository interface {
	FindByCode(ctx context.Context, ncm string) (*domains.PISCOFINSDocument, error)
}

type PISCOFINSService struct {
	repo PISCOFINSRepository
}

func NewPISCOFINSService(repo PISCOFINSRepository) *PISCOFINSService {
	return &PISCOFINSService{
		repo: repo,
	}
}

func (s *PISCOFINSService) ValidatePISCOFINS(ctx context.Context, req domains.ValidateRequest) (*domains.PISCOFINSValidacaoResponse, error) {
	doc, err := s.repo.FindByCode(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar PIS/COFINS: %w", err)
	}

	// Mapeia os CSTs genéricos para os campos específicos de PIS e COFINS
	return &domains.PISCOFINSValidacaoResponse{
		Regime:           doc.Regime,
		Descricao:        doc.Descricao,
		CSTPisEntrada:    doc.CST.Entrada,
		CSTPisSaida:      doc.CST.Saida,
		CSTCofinsEntrada: doc.CST.Entrada,
		CSTCofinsSaida:   doc.CST.Saida,
	}, nil
}
