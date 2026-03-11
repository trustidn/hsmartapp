package product

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrProductLimit = errors.New("product limit reached for free plan")

const FreePlanProductLimit = 10

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
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
	// Free plan limit
	count, err := s.repo.CountByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if count >= FreePlanProductLimit {
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
