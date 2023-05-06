package storage

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Storage interface {
	Find(ctx context.Context, filter map[string]interface{}) (interface{}, error)
	Insert(ctx context.Context, obj interface{}) error
	Update(ctx context.Context, id interface{}, obj interface{}) error
}
