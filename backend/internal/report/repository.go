package report

import (
	"context"
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

type DailySummary struct {
	Date         time.Time `json:"date"`
	SalesTotal   int64     `json:"sales_total"`
	ExpenseTotal int64     `json:"expense_total"`
	Profit       int64     `json:"profit"`
	Transactions int       `json:"transactions"`
}

type ProductRankItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Qty         int    `json:"qty"`
	Total       int64  `json:"total"`
}

func (r *Repository) DailySummary(ctx context.Context, tenantID uuid.UUID, date time.Time) (*DailySummary, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var salesTotal, expenseTotal int64
	var transactions int
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total), 0), COUNT(*)
		FROM sales
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at < $3
	`, tenantID, dayStart, dayEnd).Scan(&salesTotal, &transactions)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at < $3
	`, tenantID, dayStart, dayEnd).Scan(&expenseTotal)
	if err != nil {
		return nil, err
	}
	return &DailySummary{
		Date:         dayStart,
		SalesTotal:   salesTotal,
		ExpenseTotal: expenseTotal,
		Profit:       salesTotal - expenseTotal,
		Transactions: transactions,
	}, nil
}

// DailySummaryForRange returns summary for explicit start/end (UTC).
func (r *Repository) DailySummaryForRange(ctx context.Context, tenantID uuid.UUID, start, end time.Time) (*DailySummary, error) {
	var salesTotal, expenseTotal int64
	var transactions int
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total), 0), COUNT(*)
		FROM sales
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3
	`, tenantID, start, end).Scan(&salesTotal, &transactions)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3
	`, tenantID, start, end).Scan(&expenseTotal)
	if err != nil {
		return nil, err
	}
	return &DailySummary{
		Date:         start,
		SalesTotal:   salesTotal,
		ExpenseTotal: expenseTotal,
		Profit:       salesTotal - expenseTotal,
		Transactions: transactions,
	}, nil
}

func (r *Repository) ProductRankingForRange(ctx context.Context, tenantID uuid.UUID, start, end time.Time, limit int) ([]ProductRankItem, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			COALESCE(si.product_id::text, 'cust-' || COALESCE(si.product_name, '')) as product_id,
			COALESCE(p.name, si.product_name, 'Item') as product_name,
			SUM(si.qty)::int as qty,
			SUM(si.subtotal) as total
		FROM sale_items si
		JOIN sales s ON s.id = si.sale_id AND s.tenant_id = $1
		LEFT JOIN products p ON p.id = si.product_id
		WHERE s.created_at >= $2 AND s.created_at <= $3
		GROUP BY COALESCE(si.product_id::text, 'cust-' || COALESCE(si.product_name, '')), COALESCE(p.name, si.product_name, 'Item')
		ORDER BY qty DESC, total DESC
		LIMIT $4
	`, tenantID, start, end, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ProductRankItem
	for rows.Next() {
		var it ProductRankItem
		if err := rows.Scan(&it.ProductID, &it.ProductName, &it.Qty, &it.Total); err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	return list, rows.Err()
}

func (r *Repository) ProductRanking(ctx context.Context, tenantID uuid.UUID, date time.Time, limit int) ([]ProductRankItem, error) {
	if limit <= 0 {
		limit = 10
	}
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	return r.ProductRankingForRange(ctx, tenantID, dayStart, dayEnd, limit)
}

// SalesChartDay for chart (last N days or monthly)
type SalesChartDay struct {
	Date  string `json:"date"`
	Total int64  `json:"total"`
}

func (r *Repository) RangeSummary(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*RangeSummary, error) {
	dayStart := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	dayEnd := time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, to.Location()).Add(time.Second)

	var salesTotal, expenseTotal int64
	var transactions int
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total), 0), COUNT(*)
		FROM sales
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3
	`, tenantID, dayStart, dayEnd).Scan(&salesTotal, &transactions)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3
	`, tenantID, dayStart, dayEnd).Scan(&expenseTotal)
	if err != nil {
		return nil, err
	}
	return &RangeSummary{
		From:         dayStart.Format("2006-01-02"),
		To:           to.Format("2006-01-02"),
		SalesTotal:   salesTotal,
		ExpenseTotal: expenseTotal,
		Profit:       salesTotal - expenseTotal,
		Transactions: transactions,
	}, nil
}

func (r *Repository) RangeProductRanking(ctx context.Context, tenantID uuid.UUID, from, to time.Time, limit int) ([]ProductRankItem, error) {
	if limit <= 0 {
		limit = 10
	}
	dayStart := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	dayEnd := time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, to.Location()).Add(time.Second)

	rows, err := r.pool.Query(ctx, `
		SELECT
			COALESCE(si.product_id::text, 'cust-' || COALESCE(si.product_name, '')) as product_id,
			COALESCE(p.name, si.product_name, 'Item') as product_name,
			SUM(si.qty)::int as qty,
			SUM(si.subtotal) as total
		FROM sale_items si
		JOIN sales s ON s.id = si.sale_id AND s.tenant_id = $1
		LEFT JOIN products p ON p.id = si.product_id
		WHERE s.created_at >= $2 AND s.created_at <= $3
		GROUP BY COALESCE(si.product_id::text, 'cust-' || COALESCE(si.product_name, '')), COALESCE(p.name, si.product_name, 'Item')
		ORDER BY qty DESC, total DESC
		LIMIT $4
	`, tenantID, dayStart, dayEnd, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ProductRankItem
	for rows.Next() {
		var it ProductRankItem
		if err := rows.Scan(&it.ProductID, &it.ProductName, &it.Qty, &it.Total); err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	return list, rows.Err()
}

type RangeSummary struct {
	From         string `json:"from"`
	To           string `json:"to"`
	SalesTotal   int64  `json:"sales_total"`
	ExpenseTotal int64  `json:"expense_total"`
	Profit       int64  `json:"profit"`
	Transactions int    `json:"transactions"`
}

// SalesChartThisWeek returns daily sales for current week (Minggu-Sabtu) in given timezone.
func (r *Repository) SalesChartThisWeek(ctx context.Context, tenantID uuid.UUID, loc *time.Location, tzName string) ([]SalesChartDay, error) {
	if loc == nil {
		loc = time.UTC
	}
	if tzName == "" {
		tzName = "Asia/Jakarta"
	}
	now := time.Now().In(loc)
	weekday := now.Weekday() // 0=Sun, 6=Sat
	daysSinceSunday := int(weekday)
	sunday := now.AddDate(0, 0, -daysSinceSunday)
	saturday := sunday.AddDate(0, 0, 6)
	dayStart := time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, loc).UTC()
	dayEnd := time.Date(saturday.Year(), saturday.Month(), saturday.Day(), 23, 59, 59, 999999999, loc).UTC()

	rows, err := r.pool.Query(ctx, `
		SELECT (created_at AT TIME ZONE $4)::date::text as d, COALESCE(SUM(total), 0)::bigint as t
		FROM sales
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY (created_at AT TIME ZONE $4)::date
		ORDER BY d ASC
	`, tenantID, dayStart, dayEnd, tzName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []SalesChartDay
	for rows.Next() {
		var it SalesChartDay
		if err := rows.Scan(&it.Date, &it.Total); err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Ensure all 7 days (Sun-Sat) exist with 0 if missing
	dateToTotal := make(map[string]int64)
	for _, it := range list {
		dateToTotal[it.Date] = it.Total
	}
	result := make([]SalesChartDay, 0, 7)
	for i := 0; i < 7; i++ {
		d := sunday.AddDate(0, 0, i).Format("2006-01-02")
		result = append(result, SalesChartDay{Date: d, Total: dateToTotal[d]})
	}
	return result, nil
}

func (r *Repository) SalesChart(ctx context.Context, tenantID uuid.UUID, days int) ([]SalesChartDay, error) {
	if days <= 0 {
		days = 7
	}
	from := time.Now().AddDate(0, 0, -days)
	rows, err := r.pool.Query(ctx, `
		SELECT DATE(created_at)::text as d, COALESCE(SUM(total), 0)::bigint as t
		FROM sales
		WHERE tenant_id = $1 AND created_at >= $2
		GROUP BY DATE(created_at)
		ORDER BY d ASC
	`, tenantID, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []SalesChartDay
	for rows.Next() {
		var it SalesChartDay
		if err := rows.Scan(&it.Date, &it.Total); err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	return list, rows.Err()
}

func (r *Repository) SalesChartMonthly(ctx context.Context, tenantID uuid.UUID, months int) ([]SalesChartDay, error) {
	if months <= 0 {
		months = 12
	}
	from := time.Now().AddDate(0, -months, 0)
	rows, err := r.pool.Query(ctx, `
		SELECT TO_CHAR(created_at, 'YYYY-MM') as d, COALESCE(SUM(total), 0)::bigint as t
		FROM sales
		WHERE tenant_id = $1 AND created_at >= $2
		GROUP BY TO_CHAR(created_at, 'YYYY-MM')
		ORDER BY d ASC
	`, tenantID, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []SalesChartDay
	for rows.Next() {
		var it SalesChartDay
		if err := rows.Scan(&it.Date, &it.Total); err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	return list, rows.Err()
}
