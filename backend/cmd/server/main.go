package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/hsmart/app/backend/internal/auth"
	"github.com/hsmart/app/backend/internal/expense"
	"github.com/hsmart/app/backend/internal/product"
	"github.com/hsmart/app/backend/internal/report"
	"github.com/hsmart/app/backend/internal/sales"
	"github.com/hsmart/app/backend/internal/subscription"
	"github.com/hsmart/app/backend/internal/tenant"
	"github.com/hsmart/app/backend/pkg/cache"
	"github.com/hsmart/app/backend/pkg/database"
	"github.com/hsmart/app/backend/pkg/middleware"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	dsn := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/hsmart?sslmode=disable")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPass := getEnv("REDIS_PASSWORD", "")
	jwtSecret := getEnv("JWT_SECRET", "change-me-in-production")
	addr := getEnv("HTTP_ADDR", ":8080")

	pool, err := database.NewPostgres(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	var redisCache *cache.Redis
	if redisAddr != "" {
		redisCache, _ = cache.NewRedis(redisAddr, redisPass, 0)
	}

	// Auth (no tenant required)
	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, []byte(jwtSecret))
	authHandler := auth.NewHandler(authSvc)

	// Tenant (registration - no auth)
	tenantRepo := tenant.NewRepository(pool)
	tenantSvc := tenant.NewService(tenantRepo)
	tenantHandler := tenant.NewHandler(tenantSvc)

	// Protected modules (tenant + auth)
	productRepo := product.NewRepository(pool)
	productSvc := product.NewService(productRepo)
	productHandler := product.NewHandler(productSvc)

	salesRepo := sales.NewRepository(pool)
	salesSvc := sales.NewService(salesRepo)
	salesHandler := sales.NewHandler(salesSvc)

	expenseRepo := expense.NewRepository(pool)
	expenseSvc := expense.NewService(expenseRepo)
	expenseHandler := expense.NewHandler(expenseSvc)

	reportRepo := report.NewRepository(pool)
	reportSvc := report.NewService(reportRepo, redisCache)
	reportHandler := report.NewHandler(reportSvc)

	subRepo := subscription.NewRepository(pool)
	subSvc := subscription.NewService(subRepo)
	subHandler := subscription.NewHandler(subSvc)

	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/register", tenantHandler.Register)

	// Protected: require Tenant header + JWT (Tenant first, then Auth)
	authWrap := middleware.Auth([]byte(jwtSecret))
	protect := func(h http.Handler) http.Handler {
		return middleware.Tenant(authWrap(h))
	}
	mux.Handle("GET /api/products", protect(http.HandlerFunc(productHandler.List)))
	mux.Handle("POST /api/products", protect(http.HandlerFunc(productHandler.Create)))
	mux.Handle("GET /api/products/get", protect(http.HandlerFunc(productHandler.Get)))
	mux.Handle("PUT /api/products", protect(http.HandlerFunc(productHandler.Update)))
	mux.Handle("DELETE /api/products", protect(http.HandlerFunc(productHandler.Delete)))

	mux.Handle("POST /api/sales", protect(http.HandlerFunc(salesHandler.Create)))
	mux.Handle("GET /api/sales", protect(http.HandlerFunc(salesHandler.List)))
	mux.Handle("GET /api/sales/get", protect(http.HandlerFunc(salesHandler.Get)))

	mux.Handle("POST /api/expenses", protect(http.HandlerFunc(expenseHandler.Create)))
	mux.Handle("GET /api/expenses", protect(http.HandlerFunc(expenseHandler.List)))
	mux.Handle("DELETE /api/expenses", protect(http.HandlerFunc(expenseHandler.Delete)))

	mux.Handle("GET /api/report/daily", protect(http.HandlerFunc(reportHandler.DailySummary)))
	mux.Handle("GET /api/report/ranking", protect(http.HandlerFunc(reportHandler.ProductRanking)))
	mux.Handle("GET /api/report/dashboard", protect(http.HandlerFunc(reportHandler.Dashboard)))

	mux.Handle("GET /api/subscription", protect(http.HandlerFunc(subHandler.Get)))

	mux.Handle("GET /api/tenant/settings", protect(http.HandlerFunc(tenantHandler.GetSettings)))
	mux.Handle("PUT /api/tenant/settings", protect(http.HandlerFunc(tenantHandler.UpdateSettings)))

	handler := middleware.CORS(middleware.RateLimit(mux))

	log.Printf("HSmart API listening on %s", addr)
	return http.ListenAndServe(addr, handler)
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
