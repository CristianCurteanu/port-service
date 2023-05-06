package database

import (
	"context"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	client     *mongo.Client
	dbName     string
	collection *mongo.Collection
}

func NewMongoDB(ctx context.Context, url, dbName, collectionName string) (storage.Storage, error) {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &mongoDB{
		client:     client,
		dbName:     dbName,
		collection: client.Database(dbName).Collection(collectionName),
	}, nil
}

func (m *mongoDB) Find(ctx context.Context, filter map[string]interface{}) (interface{}, error) {
	var result interface{}

	err := m.collection.FindOne(ctx, filter).Decode(result)
	return result, err
}

func (m *mongoDB) Insert(ctx context.Context, document interface{}) error {
	_, err := m.collection.InsertOne(ctx, document)
	return err
}

func (m *mongoDB) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}
