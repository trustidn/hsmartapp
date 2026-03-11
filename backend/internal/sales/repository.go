package sales

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

type Sale struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Total         int64      `json:"total"`
	PaymentMethod string     `json:"payment_method"`
	CreatedAt     time.Time  `json:"created_at"`
	Items         []SaleItem `json:"items,omitempty"`
}

type SaleItem struct {
	ID          uuid.UUID `json:"id"`
	SaleID      uuid.UUID `json:"sale_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	Qty         int       `json:"qty"`
	Price       int64     `json:"price"`
	Subtotal    int64     `json:"subtotal"`
}

type CreateSaleInput struct {
	Total         int64
	PaymentMethod string
	Items         []CreateSaleItemInput
}

type CreateSaleItemInput struct {
	ProductID uuid.UUID
	Qty       int
	Price     int64
	Subtotal  int64
}

func (r *Repository) Create(ctx context.Context, tenantID uuid.UUID, input CreateSaleInput) (*Sale, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var saleID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO sales (tenant_id, total, payment_method)
		VALUES ($1, $2, $3)
		RETURNING id
	`, tenantID, input.Total, input.PaymentMethod).Scan(&saleID)
	if err != nil {
		return nil, err
	}

	for _, it := range input.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO sale_items (sale_id, product_id, qty, price, subtotal)
			VALUES ($1, $2, $3, $4, $5)
		`, saleID, it.ProductID, it.Qty, it.Price, it.Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Fetch full sale with items
	return r.GetByID(ctx, saleID, tenantID)
}

func (r *Repository) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*Sale, error) {
	var s Sale
	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, total, payment_method, created_at
		FROM sales WHERE id = $1 AND tenant_id = $2
	`, id, tenantID).Scan(&s.ID, &s.TenantID, &s.Total, &s.PaymentMethod, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT si.id, si.sale_id, si.product_id, p.name, si.qty, si.price, si.subtotal
		FROM sale_items si
		LEFT JOIN products p ON p.id = si.product_id
		WHERE si.sale_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var it SaleItem
		if err := rows.Scan(&it.ID, &it.SaleID, &it.ProductID, &it.ProductName, &it.Qty, &it.Price, &it.Subtotal); err != nil {
			return nil, err
		}
		s.Items = append(s.Items, it)
	}
	return &s, rows.Err()
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID uuid.UUID, from, to *time.Time, limit, offset int) ([]Sale, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q := `SELECT id, tenant_id, total, payment_method, created_at FROM sales WHERE tenant_id = $1`
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
	var list []Sale
	for rows.Next() {
		var s Sale
		if err := rows.Scan(&s.ID, &s.TenantID, &s.Total, &s.PaymentMethod, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// TenantSaleStats for admin: count and last sale per tenant
type TenantSaleStats struct {
	TenantID    uuid.UUID
	Count       int64
	LastSaleAt  *string
}

func (r *Repository) GetStatsByTenantIDs(ctx context.Context, tenantIDs []uuid.UUID) (map[uuid.UUID]TenantSaleStats, error) {
	if len(tenantIDs) == 0 {
		return map[uuid.UUID]TenantSaleStats{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, COUNT(*), MAX(created_at)::text
		FROM sales WHERE tenant_id = ANY($1)
		GROUP BY tenant_id
	`, tenantIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[uuid.UUID]TenantSaleStats)
	for rows.Next() {
		var s TenantSaleStats
		var lastAt *string
		if err := rows.Scan(&s.TenantID, &s.Count, &lastAt); err != nil {
			return nil, err
		}
		s.LastSaleAt = lastAt
		result[s.TenantID] = s
	}
	return result, rows.Err()
}
