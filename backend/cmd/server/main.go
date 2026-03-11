package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hsmart/app/backend/internal/admin"
	"github.com/hsmart/app/backend/internal/adminauth"
	"github.com/hsmart/app/backend/internal/auth"
	"github.com/hsmart/app/backend/internal/expense"
	"github.com/hsmart/app/backend/internal/planconfig"
	"github.com/hsmart/app/backend/internal/plans"
	"github.com/hsmart/app/backend/internal/product"
	"github.com/hsmart/app/backend/internal/report"
	"github.com/hsmart/app/backend/internal/saasconfig"
	"github.com/hsmart/app/backend/internal/sales"
	"github.com/hsmart/app/backend/internal/suborder"
	"github.com/hsmart/app/backend/internal/subscription"
	"github.com/hsmart/app/backend/internal/tenant"
	"github.com/hsmart/app/backend/internal/upload"
	"github.com/hsmart/app/backend/pkg/cache"
	"github.com/hsmart/app/backend/pkg/database"
	"github.com/hsmart/app/backend/pkg/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // load .env from current dir (ignore error if missing)
	validateProductionEnv()
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

// validateProductionEnv warns when critical production vars use defaults.
func validateProductionEnv() {
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" || jwt == "change-me-in-production" || jwt == "rahasia-dev" {
		log.Printf("[WARN] JWT_SECRET menggunakan nilai default - HARUS diubah di production!")
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
	planConfigRepo := planconfig.NewRepository(pool)
	subRepo := subscription.NewRepository(pool)
	subSvc := subscription.NewService(subRepo)
	subHandler := subscription.NewHandler(subSvc)

	productRepo := product.NewRepository(pool)
	productSvc := product.NewService(productRepo, tenantRepo, planConfigRepo, subRepo)
	productHandler := product.NewHandler(productSvc)

	salesRepo := sales.NewRepository(pool)
	salesSvc := sales.NewService(salesRepo)
	salesHandler := sales.NewHandler(salesSvc)

	expenseRepo := expense.NewRepository(pool)
	expenseSvc := expense.NewService(expenseRepo)
	expenseHandler := expense.NewHandler(expenseSvc)

	reportRepo := report.NewRepository(pool)
	reportSvc := report.NewService(reportRepo, redisCache, subRepo, planConfigRepo)
	reportHandler := report.NewHandler(reportSvc)

	suborderRepo := suborder.NewRepository(pool)
	suborderSvc := suborder.NewService(suborderRepo)
	suborderHandler := suborder.NewHandler(suborderSvc, planConfigRepo)
	plansHandler := plans.NewHandler(planConfigRepo)
	saasConfigRepo := saasconfig.NewRepository(pool)
	saasConfigHandler := saasconfig.NewHandler(saasConfigRepo)

	mux := http.NewServeMux()

	// Admin auth (superadmin, no tenant)
	adminAuthRepo := adminauth.NewRepository(pool)
	adminAuthSvc := adminauth.NewService(adminAuthRepo, []byte(jwtSecret))
	adminAuthHandler := adminauth.NewHandler(adminAuthSvc)
	adminGuard := middleware.AdminGuard([]byte(jwtSecret))
	adminHandler := admin.NewHandler(tenantRepo, subRepo, planConfigRepo, salesRepo, suborderRepo, authRepo)

	authWrap := middleware.Auth([]byte(jwtSecret))
	protect := func(h http.Handler) http.Handler {
		return middleware.Tenant(authWrap(h))
	}

	mux.HandleFunc("POST /api/admin/auth/login", adminAuthHandler.Login)
	mux.Handle("GET /api/admin/me", adminGuard(http.HandlerFunc(adminAuthHandler.Me)))
	mux.Handle("PATCH /api/admin/profile", adminGuard(http.HandlerFunc(adminAuthHandler.UpdateProfile)))
	mux.Handle("GET /api/admin/dashboard/stats", adminGuard(http.HandlerFunc(adminHandler.DashboardStats)))
	mux.Handle("GET /api/admin/tenants", adminGuard(http.HandlerFunc(adminHandler.ListTenants)))
	mux.Handle("GET /api/admin/tenants/get", adminGuard(http.HandlerFunc(adminHandler.GetTenant)))
	mux.Handle("PATCH /api/admin/tenants/status", adminGuard(http.HandlerFunc(adminHandler.UpdateTenantStatus)))
	mux.Handle("PATCH /api/admin/tenants/reset-password", adminGuard(http.HandlerFunc(adminHandler.ResetTenantPassword)))
	mux.Handle("PATCH /api/admin/tenants/subscription", adminGuard(http.HandlerFunc(adminHandler.UpdateTenantSubscription)))
	mux.Handle("POST /api/admin/tenants/subscription/revoke", adminGuard(http.HandlerFunc(adminHandler.RevokeTenantSubscription)))
	mux.Handle("GET /api/admin/plans", adminGuard(http.HandlerFunc(adminHandler.ListPlanConfig)))
	mux.Handle("PATCH /api/admin/plans", adminGuard(http.HandlerFunc(adminHandler.UpdatePlanConfig)))
	mux.Handle("DELETE /api/admin/plans", adminGuard(http.HandlerFunc(adminHandler.DeletePlan)))
	mux.Handle("POST /api/admin/plans/restore", adminGuard(http.HandlerFunc(adminHandler.RestorePlan)))
	mux.Handle("GET /api/admin/subscription-orders", adminGuard(http.HandlerFunc(adminHandler.ListSubscriptionOrders)))
	mux.Handle("POST /api/admin/subscription-orders/approve", adminGuard(http.HandlerFunc(adminHandler.ApproveSubscriptionOrder)))
	mux.Handle("POST /api/admin/subscription-orders/reject", adminGuard(http.HandlerFunc(adminHandler.RejectSubscriptionOrder)))
	mux.Handle("GET /api/admin/saas-settings", adminGuard(http.HandlerFunc(saasConfigHandler.Get)))
	mux.Handle("PATCH /api/admin/saas-settings", adminGuard(http.HandlerFunc(saasConfigHandler.Update)))

	uploadDir := getEnv("UPLOAD_DIR", "./uploads")
	baseURL := getEnv("BASE_URL", "") // kosong = relative path (/uploads/...) untuk proxy dev & production
	if err := os.MkdirAll(filepath.Join(uploadDir, "logos"), 0755); err != nil {
		log.Printf("[main] mkdir uploads/logos: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(uploadDir, "payment-proofs"), 0755); err != nil {
		log.Printf("[main] mkdir uploads/payment-proofs: %v", err)
	}
	uploadHandler := upload.NewHandler(uploadDir, baseURL)
	mux.Handle("POST /api/admin/upload/logo", adminGuard(http.HandlerFunc(uploadHandler.UploadLogo)))
	mux.Handle("POST /api/upload/payment-proof", protect(http.HandlerFunc(uploadHandler.UploadPaymentProof)))
	mux.Handle("/uploads/", upload.ServeUploads(uploadDir))

	// Public
	mux.HandleFunc("GET /api/public/branding", saasConfigHandler.GetPublic)
	mux.HandleFunc("GET /api/public/manifest", saasConfigHandler.GetManifest)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/register", tenantHandler.Register)

	// Protected: require Tenant header + JWT (Tenant first, then Auth)
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

	mux.Handle("GET /api/plans", protect(http.HandlerFunc(plansHandler.ListActive)))
	mux.Handle("GET /api/subscription", protect(http.HandlerFunc(subHandler.Get)))
	mux.Handle("GET /api/subscription/history", protect(http.HandlerFunc(subHandler.ListHistory)))
	mux.Handle("POST /api/subscription/orders", protect(http.HandlerFunc(suborderHandler.CreateOrder)))
	mux.Handle("PATCH /api/subscription/orders/payment-proof", protect(http.HandlerFunc(suborderHandler.SetPaymentProof)))
	mux.Handle("GET /api/subscription/orders", protect(http.HandlerFunc(suborderHandler.ListMyOrders)))
	mux.Handle("GET /api/saas-settings", protect(http.HandlerFunc(saasConfigHandler.Get)))

	mux.Handle("GET /api/tenant/settings", protect(http.HandlerFunc(tenantHandler.GetSettings)))
	mux.Handle("PUT /api/tenant/settings", protect(http.HandlerFunc(tenantHandler.UpdateSettings)))

	handler := middleware.Logging(middleware.CORS(middleware.RateLimit(mux)))

	log.Printf("HSmart API listening on %s", addr)
	return http.ListenAndServe(addr, handler)
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
