package subscription

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByTenant(ctx context.Context, tenantID uuid.UUID) (*Subscription, error) {
	return s.repo.GetByTenant(ctx, tenantID)
}

func (s *Service) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Subscription, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}
