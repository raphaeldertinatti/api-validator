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
	// Gera todos os prefixos possíveis a partir do NCM completo (8 dígitos)
	// Ex: "85442000" → ["85442000", "8544200", "854420", "85442", "8544", "85"]
	candidates := buildNCMCandidates(ncm)

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

// buildNCMCandidates gera o NCM completo + todos os prefixos significativos.
// Só inclui prefixos com comprimento que realmente existe na collection:
// 8 (completo), 7, 6, 5, 4, 2 dígitos — nunca 3 (não existe na CEST).
func buildNCMCandidates(ncm string) []string {
	validLengths := []int{8, 7, 6, 5, 4, 2}
	seen := make(map[string]struct{})
	candidates := []string{}

	for _, l := range validLengths {
		if len(ncm) >= l {
			prefix := ncm[:l]
			if _, exists := seen[prefix]; !exists {
				seen[prefix] = struct{}{}
				candidates = append(candidates, prefix)
			}
		}
	}

	return candidates
}
