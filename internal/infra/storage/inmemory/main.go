package inmemory

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
)

type InMemoryStorage struct {
	store map[string]interface{}
	mx    *sync.RWMutex
}

func NewInMemoryStorage() storage.Storage {
	return &InMemoryStorage{
		store: make(map[string]interface{}),
		mx:    &sync.RWMutex{},
	}
}

func (im *InMemoryStorage) Find(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	key, keyFound := filter["port_code"]
	if !keyFound {
		return errors.New("no `port_code` key set to filter for in memory lookup")
	}

	im.mx.Lock()
	res, found := im.store[key.(string)]
	im.mx.Unlock()

	if !found {
		return storage.ErrNotFound
	}
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr {
		return errors.New("result should be a pointer")
	}
	resultValue.Elem().Set(reflect.ValueOf(res))

	return nil
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func (im *InMemoryStorage) Insert(ctx context.Context, obj interface{}) error {
	insertable, ok := obj.(KeyValue)
	if !ok {
		return errors.New("not a KeyValue type added")
	}

	im.mx.Lock()
	defer im.mx.Unlock()

	im.store[insertable.Key] = insertable.Value

	return nil
}

func (im *InMemoryStorage) Update(ctx context.Context, id interface{}, obj interface{}) error {
	key, isString := id.(string)
	if !isString {
		return errors.New("the key given to store is not a `string`")
	}

	im.mx.Lock()
	defer im.mx.Unlock()

	im.store[key] = obj

	return nil
}
