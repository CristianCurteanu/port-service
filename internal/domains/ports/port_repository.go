package ports

import (
	"context"
	"log"
	"sort"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
	"github.com/CristianCurteanu/koken-api/internal/infra/storage/inmemory"
	"go.mongodb.org/mongo-driver/bson"
)

type PortRepository interface {
	Find(ctx context.Context, code string) (Port, error)
	Create(ctx context.Context, port Port) error
	Update(ctx context.Context, port Port) error
}

var storageStrategies map[RepositoryStrategy]func(storage.Storage) PortRepository

func init() {
	storageStrategies = make(map[RepositoryStrategy]func(storage.Storage) PortRepository)
	storageStrategies[StorageTypeInMem] = NewInMemoryRepository
	storageStrategies[StorageTypeMongoDB] = NewMongoRepository
}

// RegisterPortRepositoryStrategy adds a new strategy for port repository, to handle different storage mechanisms
func RegisterPortRepositoryStrategy(constructor func(storage.Storage) PortRepository) RepositoryStrategy {
	keys := make([]RepositoryStrategy, 0, len(storageStrategies))
	for key := range storageStrategies {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	newKey := keys[len(keys)-1] + 1
	storageStrategies[newKey] = constructor

	return newKey
}

type portsRepository struct {
	// This is supposed to have an underlying layer of PortRepository implementation, for specific storage type
	repositoryStrategy PortRepository
}

type RepositoryStrategy int

const (
	storageTypeUndefined RepositoryStrategy = iota
	StorageTypeInMem
	StorageTypeMongoDB
)

func NewPortRepository(key RepositoryStrategy, st storage.Storage) PortRepository {
	storageStrategy := storageStrategies[key](st)
	return &portsRepository{storageStrategy}
}

func (pr *portsRepository) Find(ctx context.Context, code string) (port Port, err error) {
	return pr.repositoryStrategy.Find(ctx, code)
}

func (pr *portsRepository) Create(ctx context.Context, port Port) error {
	return pr.repositoryStrategy.Create(ctx, port) //.Insert(ctx, obj)
}

func (pr *portsRepository) Update(ctx context.Context, port Port) error {
	return pr.repositoryStrategy.Update(ctx, port)
}

/*
inMemoryRepository is a repository strategy, that is created to handle
in memory data access layer, and to use this kind of storage type structs for insert
*/
type inMemoryRepository struct {
	store storage.Storage
}

func NewInMemoryRepository(st storage.Storage) PortRepository {
	return &inMemoryRepository{st}
}

func (imr *inMemoryRepository) Find(ctx context.Context, code string) (port Port, err error) {
	err = imr.store.Find(ctx, bson.M{"port_code": code}, &port)
	if err != nil {
		log.Println("error while looking up for element, err:", err)
		return
	}
	return
}

func (pr *inMemoryRepository) Create(ctx context.Context, port Port) error {
	return pr.store.Insert(ctx, inmemory.KeyValue{
		Key:   port.PortCode,
		Value: port,
	})
}
func (pr *inMemoryRepository) Update(ctx context.Context, port Port) error {
	return pr.store.Update(ctx, port.PortCode, port)
}

/*
mongoRepository is a repository strategy, that is created to handle
MongoDB data access layer, and to use this kind of storage type structs for update, like bson.M{"$set": ...}
*/
type mongoRepository struct {
	store storage.Storage
}

func NewMongoRepository(st storage.Storage) PortRepository {
	return &mongoRepository{st}
}

func (imr *mongoRepository) Find(ctx context.Context, code string) (port Port, err error) {
	err = imr.store.Find(ctx, bson.M{"port_code": code}, &port)
	if err != nil {
		log.Println("error while looking up for element, err:", err)
		return
	}
	return
}

func (pr *mongoRepository) Create(ctx context.Context, port Port) error {
	return pr.store.Insert(ctx, port)
}

func (pr *mongoRepository) Update(ctx context.Context, port Port) error {
	return pr.store.Update(ctx, bson.M{"port_code": port.PortCode}, bson.M{"$set": port.AsBson()})
}
