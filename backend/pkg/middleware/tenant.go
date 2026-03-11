package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	TenantIDKey    contextKey = "tenant_id"
	UserIDKey      contextKey = "user_id"
	UserTenantIDKey contextKey = "user_tenant_id"
)

// Tenant extracts X-Tenant-ID, validates UUID, and attaches to context.
// Returns 400 if header missing or invalid.
func Tenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.Header.Get("X-Tenant-ID")
		if idStr == "" {
			http.Error(w, `{"error":"X-Tenant-ID required"}`, http.StatusBadRequest)
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid X-Tenant-ID"}`, http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), TenantIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTenantID returns tenant UUID from context. Panics if not set (use after Tenant middleware).
func GetTenantID(ctx context.Context) uuid.UUID {
	return ctx.Value(TenantIDKey).(uuid.UUID)
}
