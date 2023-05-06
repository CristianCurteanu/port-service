package ports

import (
	"context"
	"log"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
	"github.com/CristianCurteanu/koken-api/internal/infra/storage/database"
	"github.com/CristianCurteanu/koken-api/internal/infra/storage/inmemory"
	"go.mongodb.org/mongo-driver/bson"
)

type PortRepository interface {
	Find(ctx context.Context, code string) (Port, error)
	Create(ctx context.Context, port Port) error
	Update(ctx context.Context, port Port) error
}

type portsRepository struct {
	storage storage.Storage
}

func NewPortRepositories(storage storage.Storage) PortRepository {
	return &portsRepository{storage}
}

func (pr *portsRepository) Find(ctx context.Context, code string) (port Port, err error) {
	err = pr.storage.Find(ctx, bson.M{"port_code": code}, &port)
	if err != nil {
		log.Println("error while looking up for element, err:", err)
		return
	}

	return
}

func (pr *portsRepository) Create(ctx context.Context, port Port) error {
	_, isInMem := pr.storage.(*inmemory.InMemoryStorage)
	if isInMem {
		return pr.storage.Insert(ctx, inmemory.KeyValue{
			Key:   port.PortCode,
			Value: port,
		})
	}
	return pr.storage.Insert(ctx, port)
}

func (pr *portsRepository) Update(ctx context.Context, port Port) error {
	_, isMongo := pr.storage.(*database.MongoDB)
	if isMongo {
		filter := bson.M{
			"port_code": port.PortCode,
		}
		return pr.storage.Update(ctx, filter, bson.M{"$set": port.AsBson()})
	}

	return pr.storage.Update(ctx, port.PortCode, port)
}
