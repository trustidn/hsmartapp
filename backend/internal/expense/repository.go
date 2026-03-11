package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

type Expense struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	Amount    int64     `json:"amount"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Repository) Create(ctx context.Context, tenantID uuid.UUID, name string, amount int64, note string) (*Expense, error) {
	var e Expense
	err := r.pool.QueryRow(ctx, `
		INSERT INTO expenses (tenant_id, name, amount, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tenant_id, name, amount, note, created_at
	`, tenantID, name, amount, note).Scan(&e.ID, &e.TenantID, &e.Name, &e.Amount, &e.Note, &e.CreatedAt)
	return &e, err
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID uuid.UUID, from, to *time.Time, limit, offset int) ([]Expense, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q := `SELECT id, tenant_id, name, amount, note, created_at FROM expenses WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	idx := 2
	if from != nil {
		q += ` AND created_at >= $` + fmt.Sprintf("%d", idx)
		args = append(args, *from)
		idx++
	}
	if to != nil {
		q += ` AND created_at <= $` + fmt.Sprintf("%d", idx)
		args = append(args, *to)
		idx++
	}
	q += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", idx) + ` OFFSET $` + fmt.Sprintf("%d", idx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.TenantID, &e.Name, &e.Amount, &e.Note, &e.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *Repository) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM expenses WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}
