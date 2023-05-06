package ports

import (
	"context"
	"errors"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/multierr"
)

var (
	ErrPortNotFound = errors.New("port not found")
)

type PortRepository interface {
	Find(ctx context.Context, code string) (Port, error)
	Create(ctx context.Context, port Port) error
	Update(ctx context.Context, port Port) error
}

type portsRepository struct {
	mongoStorage storage.Storage
	inMemory     storage.Storage
}

func NewPortRepositories(inMemStorage, mongoStorage storage.Storage) PortRepository {
	return &portsRepository{mongoStorage, inMemStorage}
}

func (pr *portsRepository) Find(ctx context.Context, code string) (port Port, err error) {
	var portRes interface{}

	portRes, err = pr.inMemory.Find(ctx, map[string]interface{}{"port_code": code})
	if errors.Is(err, ErrPortNotFound) && pr.mongoStorage != nil {
		portRes, err = pr.mongoStorage.Find(ctx, map[string]interface{}{"port_code": code})
		if err != nil {
			return
		}
	} else if err != nil {
		return
	}
	port = portRes.(Port)
	return
}

func (pr *portsRepository) Create(ctx context.Context, port Port) error {
	var err, storageErr error

	storageErr = pr.inMemory.Insert(ctx, port)
	if storageErr != nil {
		multierr.Append(err, storageErr)
	}

	if pr.mongoStorage != nil {
		storageErr = pr.mongoStorage.Insert(ctx, port)
		if err != nil {
			multierr.Append(err, storageErr)
		}
	}

	return err
}

func (pr *portsRepository) Update(ctx context.Context, port Port) error {
	var err, storageErr error

	storageErr = pr.inMemory.Update(ctx, port.PortCode, port)
	if storageErr != nil {
		multierr.Append(err, storageErr)
	}

	if pr.mongoStorage != nil {
		filter := bson.M{
			"port_code": port.PortCode,
		}
		storageErr = pr.mongoStorage.Update(ctx, filter, bson.M{"$set": port.AsBson()})
		if err != nil {
			multierr.Append(err, storageErr)
		}
	}

	return err
}
