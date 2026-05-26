package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api-validator/internal/domains"
	"api-validator/internal/storage"
)

type DiferimentoRepository struct {
	collection *mongo.Collection
}

func NewDiferimentoRepository(db *storage.MongoDB) *DiferimentoRepository {
	return &DiferimentoRepository{
		collection: db.Collection("icms_diferimento"),
	}
}

func (r *DiferimentoRepository) FindByCode(ctx context.Context, ncm string) ([]domains.DiferimentoDocument, error) {
	if len(ncm) < 2 {
		return nil, nil
	}

	capitulo := ncm[:2]
	var allResults []domains.DiferimentoDocument
	seen := make(map[string]struct{})

	cursor, err := r.collection.Find(ctx, bson.M{
		"uf":            "SP",
		"ativo":         true,
		"capitulos_ncm": capitulo,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar diferimentos por capítulo: %w", err)
	}
	defer cursor.Close(ctx)

	var results []domains.DiferimentoDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("erro ao decodificar diferimentos: %w", err)
	}

	for _, doc := range results {
		key := doc.DocKey()
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			allResults = append(allResults, doc)
		}
	}

	return allResults, nil
}
