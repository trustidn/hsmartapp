package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims for JWT.
type Claims struct {
	jwt.RegisteredClaims
	UserID   string `json:"uid"`
	TenantID string `json:"tid"`
	Role     string `json:"role"`
}

// Auth validates Bearer JWT and sets user/tenant in context.
func Auth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			var claims Claims
			token, err := jwt.ParseWithClaims(tokenStr, &claims, func(*jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			uid, _ := uuid.Parse(claims.UserID)
			tid, _ := uuid.Parse(claims.TenantID)
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, uid)
			ctx = context.WithValue(ctx, UserTenantIDKey, tid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID returns user UUID from context.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(UserIDKey)
	if v == nil {
		return uuid.Nil, false
	}
	return v.(uuid.UUID), true
}

// GetUserTenantID returns tenant UUID from auth context.
func GetUserTenantID(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(UserTenantIDKey)
	if v == nil {
		return uuid.Nil, false
	}
	return v.(uuid.UUID), true
}
