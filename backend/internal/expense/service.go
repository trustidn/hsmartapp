package expense

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateExpenseRequest struct {
	Name   string `json:"name"`
	Amount int64  `json:"amount"`
	Note   string `json:"note"`
}

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, req CreateExpenseRequest) (*Expense, error) {
	if req.Amount < 0 {
		req.Amount = 0
	}
	return s.repo.Create(ctx, tenantID, req.Name, req.Amount, req.Note)
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, from, to *time.Time, limit, offset int) ([]Expense, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByTenant(ctx, tenantID, from, to, limit, offset)
}

func (s *Service) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	return s.repo.Delete(ctx, id, tenantID)
}
