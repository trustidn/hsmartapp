package product

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

type Product struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	ImageURL  string    `json:"image_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	SortOrder int       `json:"sort_order"`
	CreatedAt string    `json:"created_at,omitempty"`
}

func (r *Repository) Create(ctx context.Context, tenantID uuid.UUID, name string, price int64, imageURL string, sortOrder int) (*Product, error) {
	var p Product
	err := r.pool.QueryRow(ctx, `
		INSERT INTO products (tenant_id, name, price, image_url, is_active, sort_order)
		VALUES ($1, $2, $3, NULLIF($4, ''), true, $5)
		RETURNING id, tenant_id, name, price, COALESCE(image_url, ''), is_active, sort_order, created_at::text
	`, tenantID, name, price, imageURL, sortOrder).Scan(&p.ID, &p.TenantID, &p.Name, &p.Price, &p.ImageURL, &p.IsActive, &p.SortOrder, &p.CreatedAt)
	return &p, err
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]Product, error) {
	q := `SELECT id, tenant_id, name, price, COALESCE(image_url, ''), is_active, sort_order, created_at::text FROM products WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	if activeOnly {
		q += ` AND is_active = true`
	}
	q += ` ORDER BY sort_order ASC, name ASC`

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Price, &p.ImageURL, &p.IsActive, &p.SortOrder, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*Product, error) {
	var p Product
	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, name, price, COALESCE(image_url, ''), is_active, sort_order, created_at::text
		FROM products WHERE id = $1 AND tenant_id = $2
	`, id, tenantID).Scan(&p.ID, &p.TenantID, &p.Name, &p.Price, &p.ImageURL, &p.IsActive, &p.SortOrder, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Update(ctx context.Context, id, tenantID uuid.UUID, name string, price int64, imageURL *string, isActive bool, sortOrder int) error {
	if imageURL != nil {
		_, err := r.pool.Exec(ctx, `
			UPDATE products SET name = $1, price = $2, image_url = NULLIF($3, ''), is_active = $4, sort_order = $5
			WHERE id = $6 AND tenant_id = $7
		`, name, price, *imageURL, isActive, sortOrder, id, tenantID)
		return err
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE products SET name = $1, price = $2, is_active = $3, sort_order = $4
		WHERE id = $5 AND tenant_id = $6
	`, name, price, isActive, sortOrder, id, tenantID)
	return err
}

func (r *Repository) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

func (r *Repository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE tenant_id = $1`, tenantID).Scan(&n)
	return n, err
}
