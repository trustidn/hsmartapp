package report

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hsmart/app/backend/pkg/cache"
)

type Service struct {
	repo  *Repository
	redis *cache.Redis
}

func NewService(repo *Repository, r *cache.Redis) *Service {
	return &Service{repo: repo, redis: r}
}

type DashboardResult struct {
	Today       *DailySummary     `json:"today"`
	ProductRank []ProductRankItem `json:"product_rank"`
	SalesChart  []SalesChartDay   `json:"sales_chart,omitempty"`
}

func (s *Service) DailySummary(ctx context.Context, tenantID uuid.UUID, date time.Time) (*DailySummary, error) {
	key := cache.SummaryKey(tenantID.String(), date.Format("2006-01-02"))
	if s.redis != nil {
		b, err := s.redis.Client().Get(ctx, key).Bytes()
		if err == nil {
			var ds DailySummary
			if json.Unmarshal(b, &ds) == nil {
				return &ds, nil
			}
		}
	}
	ds, err := s.repo.DailySummary(ctx, tenantID, date)
	if err != nil {
		return nil, err
	}
	if s.redis != nil {
		if b, _ := json.Marshal(ds); len(b) > 0 {
			s.redis.Client().Set(ctx, key, b, 5*time.Minute)
		}
	}
	return ds, nil
}

func (s *Service) ProductRanking(ctx context.Context, tenantID uuid.UUID, date time.Time, limit int) ([]ProductRankItem, error) {
	return s.repo.ProductRanking(ctx, tenantID, date, limit)
}

func (s *Service) Dashboard(ctx context.Context, tenantID uuid.UUID, date time.Time, tzName string) (*DashboardResult, error) {
	today, err := s.DailySummary(ctx, tenantID, date)
	if err != nil {
		return nil, err
	}
	rank, _ := s.repo.ProductRanking(ctx, tenantID, date, 10)
	loc, _ := time.LoadLocation(tzName)
	chart, _ := s.repo.SalesChartThisWeek(ctx, tenantID, loc, tzName)
	return &DashboardResult{Today: today, ProductRank: rank, SalesChart: chart}, nil
}

// DashboardForRange returns dashboard for explicit UTC time range (untuk "hari ini" dengan timezone).
func (s *Service) DashboardForRange(ctx context.Context, tenantID uuid.UUID, start, end time.Time, tzName string) (*DashboardResult, error) {
	today, err := s.repo.DailySummaryForRange(ctx, tenantID, start, end)
	if err != nil {
		return nil, err
	}
	rank, _ := s.repo.ProductRankingForRange(ctx, tenantID, start, end, 10)
	loc, _ := time.LoadLocation(tzName)
	chart, _ := s.repo.SalesChartThisWeek(ctx, tenantID, loc, tzName)
	return &DashboardResult{Today: today, ProductRank: rank, SalesChart: chart}, nil
}

// InvalidateSummary clears cache for a tenant/date (call after new sale/expense).
func (s *Service) InvalidateSummary(ctx context.Context, tenantID uuid.UUID, date time.Time) {
	if s.redis != nil {
		key := cache.SummaryKey(tenantID.String(), date.Format("2006-01-02"))
		s.redis.Client().Del(ctx, key)
	}
}

// DashboardRangeResult for range dashboard (7d, 30d, 12m).
type DashboardRangeResult struct {
	Range       *RangeSummary     `json:"range"`
	ProductRank []ProductRankItem `json:"product_rank"`
	SalesChart  []SalesChartDay   `json:"sales_chart"`
}

func (s *Service) DashboardRange(ctx context.Context, tenantID uuid.UUID, from, to time.Time, tzName string) (*DashboardRangeResult, error) {
	rangeSummary, err := s.repo.RangeSummary(ctx, tenantID, from, to)
	if err != nil {
		return nil, err
	}
	rank, _ := s.repo.RangeProductRanking(ctx, tenantID, from, to, 10)
	// Grafik selalu Minggu-Sabtu minggu ini, tidak terpengaruh filter periode
	loc, _ := time.LoadLocation(tzName)
	chart, _ := s.repo.SalesChartThisWeek(ctx, tenantID, loc, tzName)
	return &DashboardRangeResult{Range: rangeSummary, ProductRank: rank, SalesChart: chart}, nil
}
