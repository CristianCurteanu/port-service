package ports

import (
	"context"

	"github.com/CristianCurteanu/koken-api/internal/infra/storage"
	"golang.org/x/sync/errgroup"
)

type PortService interface {
	GetByPortCode(ctx context.Context, code string) (Port, error)
	CreateOrUpdate(ctx context.Context, port Port) error
	CreateOrUpdateMany(ctx context.Context, ports []Port) error
}

type portsService struct {
	repo PortRepository
}

func NewPortService(repo PortRepository) PortService {
	return &portsService{repo}
}

func (ps *portsService) GetByPortCode(ctx context.Context, code string) (Port, error) {
	return ps.repo.Find(ctx, code)
}

func (ps *portsService) CreateOrUpdate(ctx context.Context, port Port) error {
	if _, err := ps.repo.Find(ctx, port.PortCode); err == storage.ErrNotFound {
		return ps.repo.Create(ctx, port)
	}
	return ps.repo.Update(ctx, port)
}

func (ps *portsService) CreateOrUpdateMany(ctx context.Context, ports []Port) error {
	g := new(errgroup.Group)

	for _, port := range ports {
		g.Go(ps.wrapCreateOrUpdate(ctx, port))
	}

	return g.Wait()
}

func (ps *portsService) wrapCreateOrUpdate(ctx context.Context, port Port) func() error {
	return func() error {
		return ps.CreateOrUpdate(ctx, port)
	}
}
