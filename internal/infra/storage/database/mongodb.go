package database

import (
	"context"
	"log"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
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

	return &MongoDB{
		client:     client,
		dbName:     dbName,
		collection: client.Database(dbName).Collection(collectionName),
	}, nil
}

func (m *MongoDB) Find(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	err := m.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		err = storage.ErrNotFound
	}
	return err
}

func (m *MongoDB) Insert(ctx context.Context, document interface{}) error {
	log.Println("mongo driver insert...")
	_, err := m.collection.InsertOne(ctx, document)
	return err
}

func (m *MongoDB) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}
