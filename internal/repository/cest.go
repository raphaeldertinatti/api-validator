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

func (r *CESTRepository) FindByCode(ctx context.Context, ncmCode string) ([]domains.CESTDocument, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"ncms_codigos": ncmCode})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar CEST para NCM %s: %w", ncmCode, err)
	}
	defer cursor.Close(ctx)

	var results []domains.CESTDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("erro ao deserializar resultados CEST: %w", err)
	}

	return results, nil
}
