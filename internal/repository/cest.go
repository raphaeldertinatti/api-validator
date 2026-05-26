package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api-validator/internal/domains"
	"api-validator/internal/storage"
)

type CESTRepository struct {
	collection *mongo.Collection
}

func NewCESTRepository(db *storage.MongoDB) *CESTRepository {
	return &CESTRepository{
		collection: db.Collection("cest"),
	}
}

func (r *CESTRepository) FindByCode(ctx context.Context, ncm string) ([]domains.CESTDocument, error) {
	// 8 (completo), 7, 6, 5, 4, 2 dígitos — nunca 3 (não existe na CEST)
	lengths := []int{8, 7, 6, 5, 4, 2}
	candidates := domains.BuildNCMCandidates(ncm, lengths)

	filter := bson.M{
		"ncms_codigos": bson.M{"$in": candidates},
		"ativo":        true,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar CESTs para NCM %s: %w", ncm, err)
	}
	defer cursor.Close(ctx)

	var results []domains.CESTDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("erro ao decodificar CESTs: %w", err)
	}

	return results, nil
}
