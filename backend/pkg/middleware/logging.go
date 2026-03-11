package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// responseRecorder wraps http.ResponseWriter to capture status code and bytes written.
type responseRecorder struct {
	http.ResponseWriter
	status int
	written int64
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

// requestID extracts or generates a simple request identifier for tracing.
// In production you might use X-Request-ID header or UUID.
func requestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return ""
}

// clientIP returns the client IP, considering X-Forwarded-For when behind proxy.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can be "client, proxy1, proxy2" - take first (original client)
		if idx := strings.Index(xff, ","); idx > 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	return r.RemoteAddr
}

// Logging returns middleware that logs every request centrally:
// - method, path, status, duration, client IP
// - For 4xx/5xx: logs at error level with status for quick monitoring
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		dur := time.Since(start)
		status := rec.status
		if status == 0 {
			status = http.StatusOK
		}
		ip := clientIP(r)
		reqID := requestID(r)
		msg := ""
		if reqID != "" {
			msg = " req_id=" + reqID
		}
		logLine := "method=" + r.Method + " path=" + r.URL.Path + " status=" + strconv.Itoa(status) + " duration_ms=" + strconv.Itoa(int(dur.Milliseconds())) + " ip=" + ip + msg
		if status >= 400 {
			log.Printf("[ERROR] %s", logLine)
		} else {
			log.Printf("[INFO] %s", logLine)
		}
	})
}
