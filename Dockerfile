# =============================================================================
# HSmart Production Dockerfile - Multi-stage build
# =============================================================================
# Stage 1: Frontend build (Vue 3 + Vite)
# =============================================================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Install all dependencies (devDependencies needed for build)
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci 2>/dev/null || npm install

COPY frontend/ ./
RUN npm run build

# =============================================================================
# Stage 2: Backend build (Go)
# =============================================================================
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app/backend

# Download modules
COPY backend/go.mod backend/go.sum* ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /server \
    ./cmd/server

# =============================================================================
# Stage 3: Runtime (minimal)
# =============================================================================
FROM alpine:3.19

# Security: run as non-root user
RUN adduser -D -u 1000 appuser

# Minimal packages: ca-certificates for HTTPS clients, wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Copy Go binary
COPY --from=backend-builder /server /app/server

# Copy frontend static files
COPY --from=frontend-builder /app/frontend/dist /app/static

# Create writable uploads directory (persisted via volume)
RUN mkdir -p /app/uploads/logos /app/uploads/payment-proofs \
    && chown -R appuser:appuser /app/uploads

# Switch to non-root user
USER appuser

EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O- http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/app/server"]
