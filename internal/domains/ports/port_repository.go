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
	storage     storage.Storage
	storageType int
}

const (
	storageTypeUndefined = iota
	storageTypeInMem
	storageTypeMongoDB
)

func NewPortRepositories(storage storage.Storage) PortRepository {
	var storageType int
	_, isInMem := storage.(*inmemory.InMemoryStorage)
	if isInMem {
		storageType = storageTypeInMem
	}
	_, isMongo := storage.(*database.MongoDB)
	if isMongo {
		storageType = storageTypeMongoDB
	}

	return &portsRepository{storage, storageType}
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
	var obj interface{}

	if pr.storageType == storageTypeInMem {
		obj = inmemory.KeyValue{
			Key:   port.PortCode,
			Value: port,
		}
	} else {
		obj = port
	}
	return pr.storage.Insert(ctx, obj)
}

func (pr *portsRepository) Update(ctx context.Context, port Port) error {
	var filter, obj interface{}

	if pr.storageType == storageTypeMongoDB {
		filter = bson.M{"port_code": port.PortCode}
		obj = bson.M{"$set": port.AsBson()}
	} else {
		filter = port.PortCode
		obj = port
	}

	return pr.storage.Update(ctx, filter, obj)
}
