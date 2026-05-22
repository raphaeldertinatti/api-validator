package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api-validator/internal/domains"
	"api-validator/internal/storage"
)

type PISCOFINSRepository struct {
	collection *mongo.Collection
}

func NewPISCOFINSRepository(db *storage.MongoDB) *PISCOFINSRepository {
	return &PISCOFINSRepository{
		collection: db.Collection("pis_cofins"),
	}
}

func (r *PISCOFINSRepository) FindByCode(ctx context.Context, ncm string) (*domains.PISCOFINSDocument, error) {
	// Tenta NCM completo primeiro
	var result domains.PISCOFINSDocument
	err := r.collection.FindOne(ctx, bson.M{
		"ncm":   ncm,
		"ativo": true,
	}).Decode(&result)

	if err == nil {
		return &result, nil // achou exato, retorna
	}
	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("erro ao buscar PIS/COFINS para NCM %s: %w", ncm, err)
	}

	// Tenta prefixos do mais específico para o mais genérico
	// NCM tem 8 dígitos, tentamos prefixos de 7, 6, 5, 4, 2
	validLengths := []int{7, 6, 5, 4, 2}
	for _, l := range validLengths {
		if len(ncm) < l {
			continue
		}
		prefix := ncm[:l]
		err := r.collection.FindOne(ctx, bson.M{
			"ncm":   prefix,
			"ativo": true,
		}).Decode(&result)

		if err == nil {
			return &result, nil // primeiro prefixo que bater, retorna
		}
		if err != mongo.ErrNoDocuments {
			return nil, fmt.Errorf("erro ao buscar PIS/COFINS por prefixo %s: %w", prefix, err)
		}
	}

	// Nenhum encontrado — fallback 50/01 (Tributado Normal)
	return &domains.PISCOFINSDocument{
		NCM:    ncm,
		Regime: "tributado_normal",
		CST: domains.PISCOFINSCST{
			Entrada: "50",
			Saida:   "01",
		},
		Ativo: true,
	}, nil
}
