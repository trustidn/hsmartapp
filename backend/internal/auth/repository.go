package auth

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

type UserRow struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Name         string
	Email        *string
	PasswordHash string
	Role         string
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*UserRow, error) {
	// Users are looked up by tenant+phone; for login we need to find tenant by phone first
	// Actually for simple onboarding: phone is on tenant, and we have one owner user per tenant.
	// So login flow: find tenant by phone -> get owner user for that tenant.
	var u UserRow
	err := r.pool.QueryRow(ctx, `
		SELECT u.id, u.tenant_id, u.name, u.email, u.password_hash, u.role
		FROM users u
		JOIN tenants t ON t.id = u.tenant_id
		WHERE t.phone = $1 AND u.role = 'owner'
		LIMIT 1
	`, phone).Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*UserRow, error) {
	var u UserRow
	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, name, email, password_hash, role
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
