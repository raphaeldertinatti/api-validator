package service

import (
	"api-validator/internal/domains"
	"context"
	"fmt"
)

type IcmsStRepository interface {
	FindByCode(ctx context.Context, cest string) (*domains.IcmsStDocument, error)
}

type IcmsStService struct {
	repo IcmsStRepository
}

func NewIcmsStService(repo IcmsStRepository) *IcmsStService {
	return &IcmsStService{
		repo: repo,
	}
}

func (s *IcmsStService) ValidateIcmsSt(ctx context.Context, req domains.ValidateRequest, cestResult *domains.CESTValidacaoResponse) (*domains.IcmsStValidacaoResponse, error) {

	icmsStResponse := &domains.IcmsStValidacaoResponse{}
	if cestResult != nil && cestResult.Status == "DEFINIDO" && cestResult.Enquadrado != nil {
		doc, err := s.repo.FindByCode(ctx, cestResult.Enquadrado.CEST)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar ICMS ST: %w", err)
		}
		if doc != nil {
			icmsStResponse.Portaria = doc.Fundamentacao.Portaria
			icmsStResponse.TipoBaseCalculo = doc.TipoBaseCalculo
			if doc.MVAOriginal != nil {
				icmsStResponse.MVA = *doc.MVAOriginal
			} else {
				icmsStResponse.MVA = 0
			}
		} else {
			icmsStResponse.MVA = 0
		}
	} else {
		icmsStResponse = &domains.IcmsStValidacaoResponse{
			MVA: 0,
		}
	}
	return icmsStResponse, nil
}
