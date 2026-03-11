package suborder

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

type Order struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	PlanSlug         string     `json:"plan_slug"`
	AmountRupiah     int64      `json:"amount_rupiah"`
	Status           string     `json:"status"`
	PaymentNote      string     `json:"payment_note,omitempty"`
	PaymentProofURL  string     `json:"payment_proof_url,omitempty"`
	CreatedAt        string     `json:"created_at"`
	ApprovedAt       *string    `json:"approved_at,omitempty"`
	ApprovedBy       *uuid.UUID `json:"approved_by,omitempty"`
	RejectionReason  string     `json:"rejection_reason,omitempty"`
}

func (r *Repository) Create(ctx context.Context, tenantID uuid.UUID, planSlug string, amountRupiah int64, paymentNote string) (*Order, error) {
	var o Order
	var paymentNoteOut *string
	var approvedAt *string
	var approvedBy *uuid.UUID
	var rejectionReason *string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO subscription_orders (tenant_id, plan_slug, amount_rupiah, payment_note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tenant_id, plan_slug, amount_rupiah, status, payment_note, COALESCE(payment_proof_url,''), created_at::text, approved_at::text, approved_by, rejection_reason
	`, tenantID, planSlug, amountRupiah, nullIfEmpty(paymentNote)).Scan(
		&o.ID, &o.TenantID, &o.PlanSlug, &o.AmountRupiah, &o.Status, &paymentNoteOut, &o.PaymentProofURL, &o.CreatedAt,
		&approvedAt, &approvedBy, &rejectionReason,
	)
	if err != nil {
		return nil, err
	}
	if paymentNoteOut != nil {
		o.PaymentNote = *paymentNoteOut
	}
	o.ApprovedAt = approvedAt
	o.ApprovedBy = approvedBy
	if rejectionReason != nil {
		o.RejectionReason = *rejectionReason
	}
	return &o, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	var o Order
	var paymentNote *string
	var approvedAt *string
	var approvedBy *uuid.UUID
	var rejectionReason *string
	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, plan_slug, amount_rupiah, status, payment_note, COALESCE(payment_proof_url,''), created_at::text, approved_at::text, approved_by, rejection_reason
		FROM subscription_orders WHERE id = $1
	`, id).Scan(
		&o.ID, &o.TenantID, &o.PlanSlug, &o.AmountRupiah, &o.Status, &paymentNote, &o.PaymentProofURL, &o.CreatedAt,
		&approvedAt, &approvedBy, &rejectionReason,
	)
	if err != nil {
		return nil, err
	}
	if paymentNote != nil {
		o.PaymentNote = *paymentNote
	}
	o.ApprovedAt = approvedAt
	o.ApprovedBy = approvedBy
	if rejectionReason != nil {
		o.RejectionReason = *rejectionReason
	}
	return &o, nil
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Order, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, plan_slug, amount_rupiah, status, payment_note, COALESCE(payment_proof_url,''), created_at::text, approved_at::text, approved_by, rejection_reason
		FROM subscription_orders WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrderRows(rows)
}

func (r *Repository) ListPending(ctx context.Context) ([]Order, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, plan_slug, amount_rupiah, status, payment_note, COALESCE(payment_proof_url,''), created_at::text, approved_at::text, approved_by, rejection_reason
		FROM subscription_orders WHERE status IN ('pending','paid') ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrderRows(rows)
}

func (r *Repository) ListAll(ctx context.Context, limit, offset int, status string) ([]Order, int, error) {
	var total int
	if status != "" {
		err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM subscription_orders WHERE status = $1`, status).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM subscription_orders`).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}
	args := []interface{}{limit, offset}
	q := `SELECT id, tenant_id, plan_slug, amount_rupiah, status, payment_note, COALESCE(payment_proof_url,''), created_at::text, approved_at::text, approved_by, rejection_reason
		FROM subscription_orders`
	if status != "" {
		q += ` WHERE status = $3`
		args = append(args, status)
	}
	q += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	list, err := scanOrderRows(rows)
	return list, total, err
}

func (r *Repository) Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
	var by interface{}
	if approvedBy != uuid.Nil {
		by = approvedBy
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE subscription_orders SET status = 'approved', approved_at = NOW(), approved_by = $2, rejection_reason = NULL
		WHERE id = $1 AND status IN ('pending', 'paid')
	`, id, by)
	return err
}

func (r *Repository) Reject(ctx context.Context, id uuid.UUID, reason string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE subscription_orders SET status = 'rejected', rejection_reason = $2
		WHERE id = $1 AND status IN ('pending', 'paid')
	`, id, nullIfEmpty(reason))
	return err
}

func (r *Repository) SetPaymentProof(ctx context.Context, orderID, tenantID uuid.UUID, proofURL string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE subscription_orders SET status = 'paid', payment_proof_url = $3
		WHERE id = $1 AND tenant_id = $2 AND status = 'pending'
	`, orderID, tenantID, proofURL)
	return err
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func scanOrderRows(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
}) ([]Order, error) {
	var list []Order
	for rows.Next() {
		var o Order
		var paymentNote *string
		var approvedAt *string
		var approvedBy *uuid.UUID
		var rejectionReason *string
		if err := rows.Scan(&o.ID, &o.TenantID, &o.PlanSlug, &o.AmountRupiah, &o.Status, &paymentNote, &o.PaymentProofURL, &o.CreatedAt,
			&approvedAt, &approvedBy, &rejectionReason); err != nil {
			return nil, err
		}
		if paymentNote != nil {
			o.PaymentNote = *paymentNote
		}
		o.ApprovedAt = approvedAt
		o.ApprovedBy = approvedBy
		if rejectionReason != nil {
			o.RejectionReason = *rejectionReason
		}
		list = append(list, o)
	}
	return list, nil
}
