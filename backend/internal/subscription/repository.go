package subscription

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

type Subscription struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	StartedAt string    `json:"started_at"`
	ExpiredAt *string   `json:"expired_at,omitempty"`
}

func (r *Repository) GetByTenant(ctx context.Context, tenantID uuid.UUID) (*Subscription, error) {
	var s Subscription
	var expiredAt *string
	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, plan, status, started_at::text, expired_at::text
		FROM subscriptions WHERE tenant_id = $1 ORDER BY started_at DESC LIMIT 1
	`, tenantID).Scan(&s.ID, &s.TenantID, &s.Plan, &s.Status, &s.StartedAt, &expiredAt)
	if err != nil {
		return nil, err
	}
	s.ExpiredAt = expiredAt
	return &s, nil
}
