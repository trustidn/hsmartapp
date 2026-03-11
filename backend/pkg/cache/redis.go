package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis wraps go-redis client.
type Redis struct {
	client *redis.Client
}

// NewRedis creates a new Redis client.
func NewRedis(addr, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &Redis{client: client}, nil
}

// Client returns the underlying redis client.
func (r *Redis) Client() *redis.Client {
	return r.client
}

// SummaryKey returns cache key for daily summary: summary:{tenant_id}:{date}
func SummaryKey(tenantID, date string) string {
	return "summary:" + tenantID + ":" + date
}

// DashboardKey returns cache key for dashboard: dashboard:{tenant_id}
func DashboardKey(tenantID string) string {
	return "dashboard:" + tenantID
}
