package tenant

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

type Tenant struct {
	ID              uuid.UUID
	Name            string
	Phone           string
	Plan            string
	Status          string
	CreatedAt       string
	ReceiptFooter   string
	DefaultPayment  string
	WhatsAppNumber  string
	LogoURL         string
}

func (r *Repository) Create(ctx context.Context, name, phone, passwordHash string) (*Tenant, uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var tenantID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO tenants (name, phone, plan, status)
		VALUES ($1, $2, 'free', 'active')
		RETURNING id
	`, name, phone).Scan(&tenantID)
	if err != nil {
		return nil, uuid.Nil, err
	}

	ownerID := uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, tenant_id, name, email, password_hash, role)
		VALUES ($1, $2, $3, NULL, $4, 'owner')
	`, ownerID, tenantID, name, passwordHash)
	if err != nil {
		return nil, uuid.Nil, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO subscriptions (tenant_id, plan, status)
		VALUES ($1, 'free', 'active')
	`, tenantID)
	if err != nil {
		return nil, uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, uuid.Nil, err
	}

	t := &Tenant{ID: tenantID, Name: name, Phone: phone, Plan: "free", Status: "active"}
	return t, ownerID, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	var t Tenant
	// Only use base columns (migration 001) for compatibility when 002/003 not run
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, phone, plan, status, created_at::text
		FROM tenants WHERE id = $1
	`, id).Scan(&t.ID, &t.Name, &t.Phone, &t.Plan, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	t.DefaultPayment = "cash"
	return &t, nil
}

func (r *Repository) UpdateSettings(ctx context.Context, id uuid.UUID, name, receiptFooter, defaultPayment, whatsappNumber string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE tenants SET
			name = COALESCE(NULLIF($2, ''), name),
			receipt_footer = COALESCE($3, receipt_footer),
			default_payment = COALESCE(NULLIF($4, ''), default_payment),
			whatsapp_number = COALESCE($5, whatsapp_number)
		WHERE id = $1
	`, id, name, receiptFooter, defaultPayment, whatsappNumber)
	return err
}

func (r *Repository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM tenants WHERE phone = $1`, phone).Scan(&n)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ListResult for admin pagination
type ListResult struct {
	Tenants []Tenant
	Total   int64
}

func (r *Repository) List(ctx context.Context, limit, offset int, search string) (*ListResult, error) {
	searchCond := "1=1"
	args := []interface{}{limit, offset}
	if search != "" {
		searchCond = "(name ILIKE $3 OR phone ILIKE $3)"
		args = append(args, "%"+search+"%")
	}
	query := `SELECT id, name, phone, plan, status, created_at::text
		FROM tenants WHERE ` + searchCond + ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tenants []Tenant
	for rows.Next() {
		var t Tenant
		err := rows.Scan(&t.ID, &t.Name, &t.Phone, &t.Plan, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	var total int64
	if search != "" {
		err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tenants WHERE (name ILIKE $1 OR phone ILIKE $1)`, "%"+search+"%").Scan(&total)
	} else {
		err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tenants`).Scan(&total)
	}
	if err != nil {
		return nil, err
	}
	return &ListResult{Tenants: tenants, Total: total}, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE tenants SET status = $2 WHERE id = $1`, id, status)
	return err
}

func (r *Repository) UpdatePlan(ctx context.Context, id uuid.UUID, plan string) error {
	_, err := r.pool.Exec(ctx, `UPDATE tenants SET plan = $2 WHERE id = $1`, id, plan)
	return err
}

func (r *Repository) CountActive(ctx context.Context) (int64, error) {
	var n int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tenants WHERE status = 'active'`).Scan(&n)
	return n, err
}

func (r *Repository) CountByMonth(ctx context.Context, months int) ([]struct{ Month string; Count int64 }, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT to_char(created_at, 'YYYY-MM') as month, COUNT(*)::bigint
		FROM tenants
		WHERE created_at >= NOW() - ($1::text || ' months')::interval
		GROUP BY to_char(created_at, 'YYYY-MM')
		ORDER BY month
	`, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct{ Month string; Count int64 }
	for rows.Next() {
		var m string
		var c int64
		if err := rows.Scan(&m, &c); err != nil {
			return nil, err
		}
		out = append(out, struct{ Month string; Count int64 }{m, c})
	}
	return out, rows.Err()
}
