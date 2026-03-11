# HSmart Production Deployment Guide

## Architecture Overview

```
                    ┌─────────────────────────────────────────┐
                    │     Nginx Proxy Manager (external)      │
                    │     hsmart.app, www.hsmart.app          │
                    │     SSL termination                     │
                    └─────────────────┬───────────────────────┘
                                      │ HTTP (port 80/443)
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  Server: 10.10.10.50 (Ubuntu 24.04, 2 CPU, 4 GB RAM)                       │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │  PostgreSQL (host, not Docker)                                       │  │
│  │  Port: 5432                                                           │  │
│  │  DB: hsmartapp_saas                                                   │  │
│  └──────────────────────────────────────▲──────────────────────────────┘  │
│                                           │ host.docker.internal:5432       │
│  ┌───────────────────────────────────────┼──────────────────────────────┐  │
│  │  Docker Container: hsmart-app                                          │  │
│  │  - Go API (:8080)                                                      │  │
│  │  - Vue static (built-in)                                                │  │
│  │  - /uploads volume                                                     │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Step 1: Repository Analysis

### Backend
- **Language:** Go 1.22
- **Framework:** Standard library `net/http`
- **Entry:** `backend/cmd/server/main.go`
- **Build:** `go build -ldflags="-s -w" -o server ./cmd/server`
- **Port:** 8080 (HTTP_ADDR)
- **Dependencies:** pgx, go-redis, jwt, bcrypt, uuid

### Frontend
- **Framework:** Vue 3 + Vite 6
- **Build:** `npm run build` → `frontend/dist/`
- **PWA:** vite-plugin-pwa (service worker, offline)

### Environment
- **DATABASE_URL:** PostgreSQL DSN (required)
- **JWT_SECRET:** Required in production
- **STATIC_DIR:** /app/static (in container)
- **UPLOAD_DIR:** /app/uploads (persisted volume)
- **REDIS_ADDR:** Optional (report cache)

### Migrations
- SQL files in `backend/migrations/` (001–011)
- Run manually with `psql` on host database

---

## Step 2: Production Requirements

| Requirement | Solution |
|-------------|----------|
| Frontend static | Served by Go from STATIC_DIR |
| API | /api/* on port 8080 |
| Uploads | /uploads/*, volume-backed |
| Health | GET /api/health |
| Database | Host PostgreSQL via host.docker.internal |

---

## Step 3: Pre-Deployment Checklist

### 1. Host PostgreSQL Setup

```bash
# Create database and user
sudo -u postgres psql -c "CREATE USER hsmartapp_user WITH PASSWORD 'P4ssword90901';"
sudo -u postgres psql -c "CREATE DATABASE hsmartapp_saas OWNER hsmartapp_user;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE hsmartapp_saas TO hsmartapp_user;"
sudo -u postgres psql -c "ALTER DATABASE hsmartapp_saas SET timezone TO 'UTC';"

# PostgreSQL must accept connections from Docker:
# postgresql.conf: listen_addresses = 'localhost' or '*' 
# pg_hba.conf: 
#   host all hsmartapp_user 127.0.0.1/32 md5
#   host all hsmartapp_user 172.17.0.0/16 md5    # Docker bridge
#   host all hsmartapp_user 172.18.0.0/16 md5    # Docker compose default
```

### 2. Run Migrations

```bash
cd /path/to/hsmartapp
for f in backend/migrations/*.up.sql; do
  psql -h localhost -U hsmartapp_user -d hsmartapp_saas -f "$f"
done
```

### 3. Seed Superadmin (first time only)

```bash
cd backend
DATABASE_URL="postgres://hsmartapp_user:P4ssword90901@localhost:5432/hsmartapp_saas?sslmode=disable" \
ADMIN_EMAIL=admin@hsmart.app \
ADMIN_PASSWORD="your-secure-password" \
ADMIN_NAME="Super Admin" \
go run ./cmd/seed-admin
```

---

## Step 4: Deployment Commands

### Initial Deployment

```bash
git clone https://github.com/trustidn/hsmartapp.git
cd hsmartapp

# Configure environment
cp .env.example .env
# Edit .env: set JWT_SECRET (openssl rand -hex 32)

# Build and run
docker compose -f docker-compose.production.yml build
docker compose -f docker-compose.production.yml up -d
```

### Check Status

```bash
docker compose -f docker-compose.production.yml ps
docker compose -f docker-compose.production.yml logs -f hsmart
curl http://localhost:8080/api/health
```

### Update Application

```bash
git pull
docker compose -f docker-compose.production.yml down
docker compose -f docker-compose.production.yml build --no-cache
docker compose -f docker-compose.production.yml up -d
```

### View Logs

```bash
docker logs hsmart-app -f
docker logs hsmart-app --tail 100
```

---

## Step 5: Nginx Proxy Manager Configuration

1. Add **Proxy Host** for `hsmart.app` and `www.hsmart.app`
2. **Forward Hostname/IP:** `10.10.10.50` (or `host.docker.internal` if NPM is on same host)
3. **Forward Port:** `8080`
4. **Scheme:** `http`
5. **SSL:** Request new certificate (Let's Encrypt)
6. **Advanced:** Custom locations if needed (default proxy all)

---

## Step 6: Security Hardening

| Measure | Implementation |
|---------|----------------|
| Non-root user | Container runs as `appuser` (UID 1000) |
| Minimal image | Alpine 3.19, only ca-certificates + wget |
| No secrets in image | All config via environment |
| Resource limits | 768MB RAM, 1 CPU |
| Log rotation | json-file, max 10MB × 3 files |

---

## Step 7: Performance Tuning

| Setting | Value | Rationale |
|---------|-------|-----------|
| Memory limit | 768M | Leaves ~3GB for OS, PostgreSQL |
| CPU limit | 1 | 2-core server |
| Go binary | -ldflags="-s -w" | Strips debug info |
| Health interval | 30s | Reduces overhead |

---

## Step 8: Troubleshooting

### Container cannot reach database

```bash
# On Linux, host.docker.internal needs extra_hosts (already in compose)
# Verify PostgreSQL accepts connections from Docker:
psql -h 172.17.0.1 -U hsmartapp_user -d hsmartapp_saas -c "SELECT 1"
```

### Health check failing

```bash
docker exec hsmart-app wget -q -O- http://127.0.0.1:8080/api/health
```

### Uploads not persisting

```bash
docker volume inspect hsmartapp_hsmart_uploads
```

---

## File Reference

| File | Purpose |
|------|---------|
| `Dockerfile` | Multi-stage production build |
| `docker-compose.production.yml` | Production compose |
| `.env.example` | Environment template |
| `.dockerignore` | Build context exclusions |
