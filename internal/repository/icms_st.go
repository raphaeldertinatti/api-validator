package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api-validator/internal/domains"
	"api-validator/internal/storage"
)

type IcmsStRepository struct {
	collection *mongo.Collection
}

func NewIcmsStRepository(db *storage.MongoDB) *IcmsStRepository {
	return &IcmsStRepository{
		collection: db.Collection("icms_st"),
	}
}

func (r *IcmsStRepository) FindByCode(ctx context.Context, cest string) (*domains.IcmsStDocument, error) {
	var result domains.IcmsStDocument
	err := r.collection.FindOne(ctx, bson.M{
		"cest":  cest,
		"ativo": true,
	}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar ST para CEST %s: %w", cest, err)
	}

	return &result, nil
}
