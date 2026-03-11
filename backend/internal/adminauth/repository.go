package adminauth

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

type SuperadminRow struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*SuperadminRow, error) {
	var row SuperadminRow
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, name
		FROM superadmins
		WHERE email = $1
		LIMIT 1
	`, email).Scan(&row.ID, &row.Email, &row.PasswordHash, &row.Name)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
