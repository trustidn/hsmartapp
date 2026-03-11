package suborder

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

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, planSlug string, amountRupiah int64, paymentNote string) (*Order, error) {
	return s.repo.Create(ctx, tenantID, planSlug, amountRupiah, paymentNote)
}

func (s *Service) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Order, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *Service) ListPending(ctx context.Context) ([]Order, error) {
	return s.repo.ListPending(ctx)
}

func (s *Service) ListAll(ctx context.Context, limit, offset int, status string) ([]Order, int, error) {
	return s.repo.ListAll(ctx, limit, offset, status)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
	return s.repo.Approve(ctx, id, approvedBy)
}

func (s *Service) Reject(ctx context.Context, id uuid.UUID, reason string) error {
	return s.repo.Reject(ctx, id, reason)
}

func (s *Service) SetPaymentProof(ctx context.Context, orderID, tenantID uuid.UUID, proofURL string) error {
	return s.repo.SetPaymentProof(ctx, orderID, tenantID, proofURL)
}
