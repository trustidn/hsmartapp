package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AdminClaims for superadmin JWT (no tenant_id).
type AdminClaims struct {
	jwt.RegisteredClaims
	AdminID string `json:"aid"`
	Role    string `json:"role"`
}

const AdminIDKey contextKey = "admin_id"

// AdminGuard validates Bearer JWT and ensures role is superadmin.
// Does NOT require X-Tenant-ID.
func AdminGuard(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			var claims AdminClaims
			token, err := jwt.ParseWithClaims(tokenStr, &claims, func(*jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			if claims.Role != "superadmin" {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			aid, _ := uuid.Parse(claims.AdminID)
			ctx := context.WithValue(r.Context(), AdminIDKey, aid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAdminID returns admin UUID from context (AdminGuard only).
func GetAdminID(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(AdminIDKey)
	if v == nil {
		return uuid.Nil, false
	}
	return v.(uuid.UUID), true
}
