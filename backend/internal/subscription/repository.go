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
	// Ambil subscription yang belum kadaluarsa dengan expired_at terjauh, atau yang terakhir jika semua expired
	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, plan, status, started_at::text,
			CASE WHEN expired_at IS NOT NULL THEN to_char(expired_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') ELSE NULL END
		FROM subscriptions WHERE tenant_id = $1
		ORDER BY
			CASE WHEN expired_at > NOW() THEN 0 WHEN expired_at IS NULL THEN 1 ELSE 2 END,
			expired_at DESC NULLS LAST,
			started_at DESC
		LIMIT 1
	`, tenantID).Scan(&s.ID, &s.TenantID, &s.Plan, &s.Status, &s.StartedAt, &expiredAt)
	if err != nil {
		return nil, err
	}
	s.ExpiredAt = expiredAt
	return &s, nil
}

// ListByTenant returns subscription history for a tenant, newest first.
func (r *Repository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Subscription, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, plan, status, started_at::text,
			CASE WHEN expired_at IS NOT NULL THEN to_char(expired_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') ELSE NULL END
		FROM subscriptions WHERE tenant_id = $1
		ORDER BY started_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Subscription
	for rows.Next() {
		var s Subscription
		var expiredAt *string
		if err := rows.Scan(&s.ID, &s.TenantID, &s.Plan, &s.Status, &s.StartedAt, &expiredAt); err != nil {
			return nil, err
		}
		s.ExpiredAt = expiredAt
		list = append(list, s)
	}
	return list, nil
}

// Create inserts a new subscription row. Used for adding/renewing (accumulated).
func (r *Repository) Create(ctx context.Context, tenantID uuid.UUID, plan, status string, expiredAt *string) error {
	var val interface{}
	if expiredAt == nil || *expiredAt == "" {
		val = nil
	} else {
		val = expiredAt
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO subscriptions (tenant_id, plan, status, expired_at)
		VALUES ($1, $2, $3, $4)
	`, tenantID, plan, status, val)
	return err
}

// UpdateLatest updates the most recent subscription for a tenant (admin only).
// Pass empty string for plan/status to leave unchanged.
// expiredAt: nil = don't change, non-nil = set (use empty string to clear/set NULL).
// Valid plans for subscription
var ValidPlans = map[string]bool{"free": true, "premium": true, "premium_1m": true, "premium_3m": true, "premium_6m": true, "premium_1y": true, "platinum": true}

// DeleteLatestPaid menghapus langganan berbayar yang terakhir ditambahkan (by started_at).
// Tidak menghapus row plan free (subscription awal).
func (r *Repository) DeleteLatestPaid(ctx context.Context, tenantID uuid.UUID) (deleted bool, err error) {
	res, err := r.pool.Exec(ctx, `
		DELETE FROM subscriptions
		WHERE tenant_id = $1 AND plan != 'free'
		AND id = (SELECT id FROM subscriptions WHERE tenant_id = $1 AND plan != 'free' ORDER BY started_at DESC LIMIT 1)
	`, tenantID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (r *Repository) UpdateLatest(ctx context.Context, tenantID uuid.UUID, plan, status string, expiredAt *string) error {
	var err error
	if expiredAt == nil {
		_, err = r.pool.Exec(ctx, `
			UPDATE subscriptions SET
				plan = CASE WHEN $2 != '' THEN $2 ELSE plan END,
				status = CASE WHEN $3 != '' THEN $3 ELSE status END
			WHERE id = (SELECT id FROM subscriptions WHERE tenant_id = $1 ORDER BY started_at DESC LIMIT 1)
		`, tenantID, plan, status)
	} else {
		var val interface{}
		if *expiredAt == "" {
			val = nil
		} else {
			val = expiredAt
		}
		_, err = r.pool.Exec(ctx, `
			UPDATE subscriptions SET
				plan = CASE WHEN $2 != '' THEN $2 ELSE plan END,
				status = CASE WHEN $3 != '' THEN $3 ELSE status END,
				expired_at = $4
			WHERE id = (SELECT id FROM subscriptions WHERE tenant_id = $1 ORDER BY started_at DESC LIMIT 1)
		`, tenantID, plan, status, val)
	}
	return err
}
