package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
)

type inMemoryStorage struct {
	store map[string]interface{}
	mx    sync.RWMutex
}

func NewInMemoryStorage() storage.Storage {
	return &inMemoryStorage{}
}

func (im *inMemoryStorage) Find(ctx context.Context, filter map[string]interface{}) (interface{}, error) {
	key, keyFound := filter["port_code"]
	if !keyFound {
		return nil, errors.New("no `port_code` key set to filter for in memory lookup")
	}

	im.mx.Lock()
	defer im.mx.Unlock()
	result, found := im.store[key.(string)]
	if !found {
		return nil, storage.ErrNotFound
	}
	return result, nil
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func (im *inMemoryStorage) Insert(ctx context.Context, obj interface{}) error {
	insertable, ok := obj.(KeyValue)
	if !ok {
		return errors.New("not a KeyValue type added")
	}

	im.mx.Lock()
	defer im.mx.Unlock()

	im.store[insertable.Key] = insertable.Value

	return nil
}

func (im *inMemoryStorage) Update(ctx context.Context, id interface{}, obj interface{}) error {
	key, isString := id.(string)
	if !isString {
		return errors.New("the key given to store is not a `string`")
	}

	im.mx.Lock()
	defer im.mx.Unlock()

	im.store[key] = obj

	return nil
}
