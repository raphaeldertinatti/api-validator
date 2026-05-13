package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DBName string
}

func NewMongoDB(uri string, dbName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return &MongoDB{Client: client, DBName: dbName}, nil
}

func (m *MongoDB) Disconnect() error {
	return m.Client.Disconnect(context.Background())
}

func (m *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	return m.Client.Database(m.DBName).Collection(collectionName)
}

func (m *MongoDB) FindNCM(ctx context.Context, ncmCode string) (string, error) {
	collection := m.GetCollection("base_ncm")
	var result struct {
		Descricao struct {
			Completa string `bson:"completa"`
		} `bson:"descricao"`
	}
	err := collection.FindOne(ctx, map[string]string{"ncm": ncmCode}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	return result.Descricao.Completa, nil
}
