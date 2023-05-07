package ports

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	mock.Mock
}

func (ms *MockStorage) Find(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	args := ms.Called(ctx, filter, result)
	return args.Error(0)
}

func (ms *MockStorage) Insert(ctx context.Context, obj interface{}) error {
	args := ms.Called(ctx, obj)
	return args.Error(0)
}

func (ms *MockStorage) Update(ctx context.Context, id interface{}, obj interface{}) error {
	args := ms.Called(ctx, id, obj)
	return args.Error(0)
}

func TestRepositoryFind(t *testing.T) {
	t.Run("return no error if found", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewPortRepository(StorageTypeInMem, storageMock)

		storageMock.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		_, err := repository.Find(context.Background(), "TC-0001")
		require.NoError(t, err)
		storageMock.AssertNumberOfCalls(t, "Find", 1)
	})

	t.Run("return error if not found", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewPortRepository(StorageTypeInMem, storageMock)

		storageMock.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("mock error"))

		_, err := repository.Find(context.Background(), "TC-0001")
		require.Error(t, err)
		storageMock.AssertNumberOfCalls(t, "Find", 1)
	})
}

func TestCreate(t *testing.T) {
	t.Run("return no error if created", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewPortRepository(StorageTypeInMem, storageMock)

		storageMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

		err := repository.Create(context.Background(), Port{PortCode: "TC-0001"})
		require.NoError(t, err)
		storageMock.AssertNumberOfCalls(t, "Insert", 1)
	})

	t.Run("return error if not created", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewPortRepository(StorageTypeInMem, storageMock)

		storageMock.On("Insert", mock.Anything, mock.Anything).Return(errors.New("not created"))

		err := repository.Create(context.Background(), Port{PortCode: "TC-0001"})
		require.Error(t, err)
		storageMock.AssertNumberOfCalls(t, "Insert", 1)
	})

	t.Run("use in-memory key value if storage is in-memory storage", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewInMemoryRepository(storageMock)
		storageMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

		err := repository.Create(context.Background(), Port{PortCode: "TC-0001"})
		require.NoError(t, err)
		storageMock.AssertNumberOfCalls(t, "Insert", 1)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("return no error if updated", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewPortRepository(StorageTypeInMem, storageMock)

		storageMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := repository.Update(context.Background(), Port{PortCode: "TC-0001"})
		require.NoError(t, err)
		storageMock.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("return no error if updated", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewPortRepository(StorageTypeInMem, storageMock)

		storageMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update error"))

		err := repository.Update(context.Background(), Port{PortCode: "TC-0001"})
		require.Error(t, err)
		storageMock.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("handle port-code if storage is mongo", func(t *testing.T) {
		storageMock := new(MockStorage)
		repository := NewPortRepository(StorageTypeInMem, storageMock)

		storageMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := repository.Update(context.Background(), Port{PortCode: "TC-0001"})
		require.NoError(t, err)
		storageMock.AssertNumberOfCalls(t, "Update", 1)
	})
}

func TestRegisterStrategy(t *testing.T) {
	storageMock := new(MockStorage)

	key := RegisterPortRepositoryStrategy(NewMongoRepository)
	require.Equal(t, key, RepositoryStrategy(3))

	storageMock.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	repository := NewPortRepository(key, storageMock)

	_, err := repository.Find(context.Background(), "TC-0001")
	require.NoError(t, err)
	storageMock.AssertNumberOfCalls(t, "Find", 1)
}
