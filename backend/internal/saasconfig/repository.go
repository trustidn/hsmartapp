package saasconfig

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

type Config struct {
	ID               string
	AppName          string
	LogoURL          string
	AdminContact     string
	BankName         string
	BankAccountNumber string
	BankAccountName  string
}

func (r *Repository) Get(ctx context.Context) (*Config, error) {
	var c Config
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, COALESCE(app_name, 'HSmart'), COALESCE(logo_url, ''),
			COALESCE(admin_contact, ''), COALESCE(bank_name, ''), COALESCE(bank_account_number, ''), COALESCE(bank_account_name, '')
		FROM saas_config LIMIT 1
	`).Scan(&c.ID, &c.AppName, &c.LogoURL, &c.AdminContact, &c.BankName, &c.BankAccountNumber, &c.BankAccountName)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Update(ctx context.Context, appName, logoURL, adminContact, bankName, bankAccountNumber, bankAccountName *string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE saas_config SET
			app_name = COALESCE($1, app_name),
			logo_url = COALESCE($2, logo_url),
			admin_contact = COALESCE($3, admin_contact),
			bank_name = COALESCE($4, bank_name),
			bank_account_number = COALESCE($5, bank_account_number),
			bank_account_name = COALESCE($6, bank_account_name),
			updated_at = NOW()
	`, appName, logoURL, adminContact, bankName, bankAccountNumber, bankAccountName)
	return err
}
