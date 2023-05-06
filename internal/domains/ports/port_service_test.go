package ports_test

import (
	"context"
	"testing"

	"github.com/CristianCurteanu/koken-api/internal/domains/ports"
	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPortRepo struct {
	mock.Mock
}

func (mpr *MockPortRepo) Find(ctx context.Context, code string) (ports.Port, error) {
	args := mpr.Called(ctx, code)
	return args.Get(0).(ports.Port), args.Error(1)
}

func (mpr *MockPortRepo) Create(ctx context.Context, port ports.Port) error {
	// panic("not implemented") // TODO:
	args := mpr.Called(ctx, port)
	return args.Error(0)
}

func (mpr *MockPortRepo) Update(ctx context.Context, port ports.Port) error {
	args := mpr.Called(ctx, port)
	return args.Error(0)
}

func TestCreateOrUpdateMany(t *testing.T) {
	t.Run("create if not found", func(t *testing.T) {
		mockRepo := new(MockPortRepo)
		service := ports.NewPortService(mockRepo)
		p := ports.Port{PortCode: "TPC-00001"}

		ctx := context.Background()
		mockRepo.On("Find", mock.Anything, mock.Anything).Return(ports.Port{}, storage.ErrNotFound)
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		pts := []ports.Port{p}
		err := service.CreateOrUpdateMany(ctx, pts)
		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "Find", len(pts))
		mockRepo.AssertNumberOfCalls(t, "Create", len(pts))
	})

	t.Run("update if found", func(t *testing.T) {
		mockRepo := new(MockPortRepo)
		service := ports.NewPortService(mockRepo)
		p := ports.Port{PortCode: "TPC-00001"}

		ctx := context.Background()
		mockRepo.On("Find", mock.Anything, mock.Anything).Return(ports.Port{PortCode: "TPC-00001"}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		pts := []ports.Port{p}
		err := service.CreateOrUpdateMany(ctx, pts)
		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "Find", len(pts))
		mockRepo.AssertNumberOfCalls(t, "Update", len(pts))
	})
}

func TestGetByCode(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockRepo := new(MockPortRepo)
		service := ports.NewPortService(mockRepo)
		p := ports.Port{PortCode: "TPC-00001"}

		ctx := context.Background()
		mockRepo.On("Find", mock.Anything, p.PortCode).Return(p, nil)

		port, err := service.GetByPortCode(ctx, p.PortCode)
		require.NoError(t, err)
		require.Equal(t, port.PortCode, p.PortCode)

		mockRepo.AssertNumberOfCalls(t, "Find", 1)
	})
}
