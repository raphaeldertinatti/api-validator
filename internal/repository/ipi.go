package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api-validator/internal/domains"
	"api-validator/internal/storage"
)

type IPIRepository struct {
	collection *mongo.Collection
}

func NewIPRepository(db *storage.MongoDB) *IPIRepository {
	return &IPIRepository{
		collection: db.Collection("ipi_tipi"),
	}
}

func (r *IPIRepository) FindByCode(ctx context.Context, ipiCode string) (*domains.IPIDocument, error) {
	var result domains.IPIDocument
	err := r.collection.FindOne(ctx, bson.M{"ncm": ipiCode}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // não encontrado não é erro
		}
		return nil, fmt.Errorf("erro ao buscar IPI %s: %w", ipiCode, err)
	}
	return &result, nil
}
