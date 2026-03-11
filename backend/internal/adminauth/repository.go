package adminauth

import (
	"context"
	"strconv"
	"strings"

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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*SuperadminRow, error) {
	var row SuperadminRow
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, name
		FROM superadmins WHERE id = $1
	`, id).Scan(&row.ID, &row.Email, &row.PasswordHash, &row.Name)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, name, email, passwordHash *string) error {
	if name != nil || email != nil || passwordHash != nil {
		// Build dynamic update - only non-nil fields
		updates := []string{}
		args := []interface{}{}
		argNum := 1
		if name != nil {
			updates = append(updates, "name = $"+strconv.Itoa(argNum))
			args = append(args, *name)
			argNum++
		}
		if email != nil {
			updates = append(updates, "email = $"+strconv.Itoa(argNum))
			args = append(args, *email)
			argNum++
		}
		if passwordHash != nil {
			updates = append(updates, "password_hash = $"+strconv.Itoa(argNum))
			args = append(args, *passwordHash)
			argNum++
		}
		args = append(args, id)
		q := `UPDATE superadmins SET ` + strings.Join(updates, ", ") + ` WHERE id = $` + strconv.Itoa(argNum)
		_, err := r.pool.Exec(ctx, q, args...)
		return err
	}
	return nil
}
