package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api-validator/internal/domains"
	"api-validator/internal/storage"
)

type NCMRepository struct {
	collection *mongo.Collection
}

func NewNCMRepository(db *storage.MongoDB) *NCMRepository {
	return &NCMRepository{
		collection: db.Collection("base_ncm"),
	}
}

func (r *NCMRepository) FindByCode(ctx context.Context, ncmCode string) (*domains.NCMDocument, error) {
	var result domains.NCMDocument
	err := r.collection.FindOne(ctx, bson.M{"ncm": ncmCode}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // não encontrado não é erro
		}
		return nil, fmt.Errorf("erro ao buscar NCM %s: %w", ncmCode, err)
	}
	return &result, nil
}

func (r *NCMRepository) FindDescricaoCompleta(ctx context.Context, ncmCode string) (string, error) {
	doc, err := r.FindByCode(ctx, ncmCode)
	if err != nil {
		return "", err
	}
	if doc == nil {
		return "", nil
	}
	return doc.Descricao.Completa, nil
}
