package planconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

type PlanConfig struct {
	PlanSlug       string
	Name           string
	DurationMonths int
	DurationDays   int
	MaxProducts    int
	ReportDays     int
	PriceRupiah    int64
	SortOrder      int
	IsActive       bool
}

func (r *Repository) GetAll(ctx context.Context) ([]PlanConfig, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT plan_slug, name, duration_months, COALESCE(duration_days, 0), max_products, report_days,
			COALESCE(price_rupiah, 0), COALESCE(sort_order, 999), COALESCE(is_active, true)
		FROM plan_config ORDER BY COALESCE(sort_order, 999), plan_slug
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PlanConfig
	for rows.Next() {
		var p PlanConfig
		if err := rows.Scan(&p.PlanSlug, &p.Name, &p.DurationMonths, &p.DurationDays, &p.MaxProducts, &p.ReportDays, &p.PriceRupiah, &p.SortOrder, &p.IsActive); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *Repository) GetByPlan(ctx context.Context, plan string) (*PlanConfig, error) {
	var p PlanConfig
	err := r.pool.QueryRow(ctx, `
		SELECT plan_slug, name, duration_months, COALESCE(duration_days, 0), max_products, report_days,
			COALESCE(price_rupiah, 0), COALESCE(sort_order, 999), COALESCE(is_active, true)
		FROM plan_config WHERE plan_slug = $1
	`, plan).Scan(&p.PlanSlug, &p.Name, &p.DurationMonths, &p.DurationDays, &p.MaxProducts, &p.ReportDays, &p.PriceRupiah, &p.SortOrder, &p.IsActive)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Default limits when plan_config not found or table doesn't exist
var DefaultMaxProducts = map[string]int{"free": 10, "premium": -1, "premium_1m": -1, "premium_3m": -1, "premium_6m": -1, "premium_1y": -1, "platinum": -1}
var DefaultReportDays = map[string]int{"free": 7, "premium": 30, "premium_1m": 30, "premium_3m": 90, "premium_6m": 180, "premium_1y": 365, "platinum": 365}

func (r *Repository) GetMaxProducts(ctx context.Context, plan string) int {
	p, err := r.GetByPlan(ctx, plan)
	if err != nil {
		if n, ok := DefaultMaxProducts[plan]; ok {
			return n
		}
		return 10
	}
	return p.MaxProducts
}

func (r *Repository) GetReportDays(ctx context.Context, plan string) int {
	p, err := r.GetByPlan(ctx, plan)
	if err != nil {
		if n, ok := DefaultReportDays[plan]; ok {
			return n
		}
		return 7
	}
	return p.ReportDays
}

func (r *Repository) Update(ctx context.Context, planSlug string, maxProducts, reportDays, durationDays, sortOrder *int, priceRupiah *int64, isActive *bool) error {
	updates := []string{}
	args := []interface{}{planSlug}
	argIdx := 2
	if sortOrder != nil {
		updates = append(updates, fmt.Sprintf("sort_order = $%d", argIdx))
		args = append(args, *sortOrder)
		argIdx++
	}
	if durationDays != nil {
		updates = append(updates, fmt.Sprintf("duration_days = $%d", argIdx))
		args = append(args, *durationDays)
		argIdx++
	}
	if maxProducts != nil {
		updates = append(updates, fmt.Sprintf("max_products = $%d", argIdx))
		args = append(args, *maxProducts)
		argIdx++
	}
	if reportDays != nil {
		updates = append(updates, fmt.Sprintf("report_days = $%d", argIdx))
		args = append(args, *reportDays)
		argIdx++
	}
	if priceRupiah != nil {
		updates = append(updates, fmt.Sprintf("price_rupiah = $%d", argIdx))
		args = append(args, *priceRupiah)
		argIdx++
	}
	if isActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *isActive)
		argIdx++
	}
	if len(updates) == 0 {
		return nil
	}
	updates = append(updates, "updated_at = NOW()")
	_, err := r.pool.Exec(ctx, `UPDATE plan_config SET `+strings.Join(updates, ", ")+` WHERE plan_slug = $1`, args...)
	return err
}

// Delete soft-deletes (sets is_active = false). Plan 'free' cannot be deleted.
func (r *Repository) Delete(ctx context.Context, planSlug string) error {
	if planSlug == "free" {
		return fmt.Errorf("plan free tidak dapat dinonaktifkan")
	}
	f := false
	return r.Update(ctx, planSlug, nil, nil, nil, nil, nil, &f)
}

// Restore re-activates a plan (sets is_active = true).
func (r *Repository) Restore(ctx context.Context, planSlug string) error {
	t := true
	return r.Update(ctx, planSlug, nil, nil, nil, nil, nil, &t)
}
