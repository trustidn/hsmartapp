package middleware

import (
	"net/http"
	"sync"
	"time"
)

// Simple in-memory rate limiter for API (per IP).
// For production at scale, use Redis-based limiter.
const (
	rateLimitRequests = 100
	rateLimitWindow   = time.Minute
)

type rateEntry struct {
	count int
	start time.Time
}

var (
	rateMu   sync.Mutex
	rateMap  = make(map[string]*rateEntry)
	rateStop = make(chan struct{})
)

func init() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-rateStop:
				return
			case <-ticker.C:
				rateMu.Lock()
				for k, e := range rateMap {
					if time.Since(e.start) > rateLimitWindow {
						delete(rateMap, k)
					}
				}
				rateMu.Unlock()
			}
		}
	}()
}

// RateLimit returns middleware that limits requests per IP.
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			key = xff
		}
		rateMu.Lock()
		e, ok := rateMap[key]
		if !ok || time.Since(e.start) > rateLimitWindow {
			e = &rateEntry{count: 0, start: time.Now()}
			rateMap[key] = e
		}
		e.count++
		count := e.count
		rateMu.Unlock()
		if count > rateLimitRequests {
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
