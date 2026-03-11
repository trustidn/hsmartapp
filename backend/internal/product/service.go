package product

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hsmart/app/backend/internal/planconfig"
	"github.com/hsmart/app/backend/internal/tenant"
)

var ErrProductLimit = errors.New("product limit reached for plan")

type Service struct {
	repo         *Repository
	tenantRepo   *tenant.Repository
	planConfig   *planconfig.Repository
}

func NewService(repo *Repository, tenantRepo *tenant.Repository, planConfig *planconfig.Repository) *Service {
	return &Service{repo: repo, tenantRepo: tenantRepo, planConfig: planConfig}
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]Product, error) {
	return s.repo.ListByTenant(ctx, tenantID, activeOnly)
}

type CreateInput struct {
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	SortOrder int    `json:"sort_order"`
}

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, input CreateInput) (*Product, error) {
	if input.Name == "" || input.Price < 0 {
		return nil, errors.New("name required and price >= 0")
	}
	t, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	maxProducts := s.planConfig.GetMaxProducts(ctx, t.Plan)
	count, err := s.repo.CountByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if maxProducts >= 0 && count >= maxProducts {
		return nil, ErrProductLimit
	}
	if input.SortOrder == 0 {
		input.SortOrder = count + 1
	}
	return s.repo.Create(ctx, tenantID, input.Name, input.Price, input.SortOrder)
}

func (s *Service) Get(ctx context.Context, id, tenantID uuid.UUID) (*Product, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

type UpdateInput struct {
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	IsActive  *bool  `json:"is_active"`
	SortOrder int    `json:"sort_order"`
}

func (s *Service) Update(ctx context.Context, id, tenantID uuid.UUID, input UpdateInput) error {
	p, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return err
	}
	name := input.Name
	if name == "" {
		name = p.Name
	}
	price := input.Price
	if price < 0 {
		price = p.Price
	}
	active := p.IsActive
	if input.IsActive != nil {
		active = *input.IsActive
	}
	sortOrder := input.SortOrder
	if sortOrder == 0 {
		sortOrder = p.SortOrder
	}
	return s.repo.Update(ctx, id, tenantID, name, price, active, sortOrder)
}

func (s *Service) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	return s.repo.Delete(ctx, id, tenantID)
}
