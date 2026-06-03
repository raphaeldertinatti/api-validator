package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api-validator/internal/domains"
	"api-validator/internal/storage"
)

type AliquotaRepository struct {
	collection *mongo.Collection
}

func NewAliquotaRepository(db *storage.MongoDB) *AliquotaRepository {
	return &AliquotaRepository{
		collection: db.Collection("icms_aliquota"),
	}
}

func (r *AliquotaRepository) FindByCode(ctx context.Context, ncm string) ([]domains.AliquotaDocument, error) {
	var allResults []domains.AliquotaDocument
	seen := make(map[string]struct{})

	// Etapa 1: busca por NCM (completo + prefixos 7,6,5,4)
	lengths := []int{8, 7, 6, 5, 4}
	candidates := domains.BuildNCMCandidates(ncm, lengths)

	cursor, err := r.collection.Find(ctx, bson.M{
		"uf":           "SP",
		"ativo":        true,
		"ncms_codigos": bson.M{"$in": candidates},
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar aliquotas por NCM: %w", err)
	}
	defer cursor.Close(ctx)

	var byNCM []domains.AliquotaDocument
	if err := cursor.All(ctx, &byNCM); err != nil {
		return nil, fmt.Errorf("erro ao decodificar aliquotas por NCM: %w", err)
	}

	for _, doc := range byNCM {
		key := doc.DocKey()
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			allResults = append(allResults, doc)
		}
	}

	// Etapa 2: SEMPRE busca pelo capítulo também.
	if len(ncm) >= 2 {
		capitulo := ncm[:2]

		cursor2, err := r.collection.Find(ctx, bson.M{
			"uf":            "SP",
			"ativo":         true,
			"capitulos_ncm": capitulo,
		})
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar aliquotas por capítulo: %w", err)
		}
		defer cursor2.Close(ctx)

		var byCapitulo []domains.AliquotaDocument
		if err := cursor2.All(ctx, &byCapitulo); err != nil {
			return nil, fmt.Errorf("erro ao decodificar aliquotas por capítulo: %w", err)
		}

		for _, doc := range byCapitulo {
			key := doc.DocKey()
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				allResults = append(allResults, doc)
			}
		}
	}

	return allResults, nil
}
