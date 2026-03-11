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
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, phone, plan, status, created_at::text,
			COALESCE(receipt_footer, ''), COALESCE(default_payment, 'cash'), COALESCE(whatsapp_number, ''),
			COALESCE(logo_url, '')
		FROM tenants WHERE id = $1
	`, id).Scan(&t.ID, &t.Name, &t.Phone, &t.Plan, &t.Status, &t.CreatedAt, &t.ReceiptFooter, &t.DefaultPayment, &t.WhatsAppNumber, &t.LogoURL)
	if err != nil {
		return nil, err
	}
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
