# HSmart

**HSmart** is a simple POS and business tracking application for micro merchants (UMKM): street food vendors, coffee stalls, small kiosks, drink sellers, gorengan sellers, and roadside merchants.

- **Domain:** https://hsmart.app  
- **Stack:** Go (backend), Vue 3 + Vite (frontend), PostgreSQL, Redis, PWA with offline-first (IndexedDB/Dexie).

## Features

- **POS:** Tap product → add to cart → Pay. Target &lt; 3 seconds per transaction.
- **Daily summary:** Sales, expenses, profit, transaction count.
- **Product ranking:** Best-selling products today.
- **Quick expense:** One-tap expense recording.
- **Offline-first:** Works without internet; syncs when back online.
- **PWA:** Installable on Android, iOS, Windows, macOS, Linux.

## Development lokal (tanpa Docker)

Semua dijalankan langsung di mesin Anda. Docker tidak dipakai untuk development.

### Prerequisites

- **Go 1.22+**
- **Node 18+**
- **PostgreSQL 16** (jalan di local)
- **Redis 7** (opsional; kalau tidak di-set, fitur cache report dinonaktifkan)

### 1. Database

Buat database dan jalankan migrasi (sekali saja):

```bash
createdb hsmart
psql -d hsmart -f backend/migrations/001_init.up.sql
```

Sesuaikan user/password jika berbeda (default: `postgres`/`postgres`).

### 2. Backend

```bash
cd backend
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/hsmart?sslmode=disable"
export JWT_SECRET="rahasia-dev"

# Opsional: cache report (tanpa ini pun backend jalan)
# export REDIS_ADDR=localhost:6379

go run ./cmd/server
```

API: **http://localhost:8080**

### 3. Frontend

Di terminal lain:

```bash
cd frontend
npm install
npm run dev
```

App: **http://localhost:5173** — request `/api` di-proxy ke backend.

### 4. Cek

Buka http://localhost:5173 → Daftar (HP + password) → Login → Tambah Produk → POS → tap produk → BAYAR.

### 5. Admin Superadmin (opsional)

```bash
# Jalankan migration superadmins
psql -d hsmart -f backend/migrations/004_superadmins.up.sql
psql -d hsmart -f backend/migrations/005_plan_config.up.sql
psql -d hsmart -f backend/migrations/006_plan_price.up.sql
psql -d hsmart -f backend/migrations/007_plan_duration_days.up.sql
psql -d hsmart -f backend/migrations/008_plan_platinum_and_order.up.sql
psql -d hsmart -f backend/migrations/009_plan_is_active.up.sql
psql -d hsmart -f backend/migrations/010_subscription_orders.up.sql
psql -d hsmart -f backend/migrations/011_saas_config.up.sql

# Seed superadmin pertama (default: admin@hsmart.app / admin123)
cd backend && go run ./cmd/seed-admin

# Custom email/password:
ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=secret go run ./cmd/seed-admin
```

Buka http://localhost:5173/admin/login → Login dengan email & password superadmin.

## Project Structure

```
backend/
  cmd/server/          # Entrypoint
  internal/            # Auth, tenant, product, sales, expense, report, subscription
  pkg/                 # database, cache, middleware
  migrations/         # PostgreSQL migrations
frontend/
  src/
    views/             # Login, Register, POS, Dashboard, Expenses, Products
    stores/            # Pinia: auth, products, pos, sync
    lib/               # api.js, db.js (Dexie IndexedDB)
```

## API Overview

- **Public:** `POST /api/auth/login`, `POST /api/register`, `POST /api/admin/auth/login` (superadmin)
- **Protected (Header: `Authorization: Bearer <token>`, `X-Tenant-ID: <uuid>`):**
  - Products: `GET/POST/PUT/DELETE /api/products`
  - Sales: `POST /api/sales`, `GET /api/sales`, `GET /api/sales/get?id=`
  - Expenses: `POST /api/expenses`, `GET /api/expenses`
  - Report: `GET /api/report/daily`, `GET /api/report/ranking`, `GET /api/report/dashboard`
  - Subscription: `GET /api/subscription`
- **Admin (Header: `Authorization: Bearer <admin_token>` — tidak pakai X-Tenant-ID):**
  - `GET /api/admin/me` — verifikasi token superadmin
  - `GET /api/admin/tenants?limit=&offset=&search=` — daftar tenant (paginated, search by nama/HP)
  - `GET /api/admin/tenants/get?id=` — detail tenant + subscription
  - `PATCH /api/admin/tenants/status` — ubah status tenant (body: `{id, status}` — status: active, suspended, inactive)
  - `PATCH /api/admin/tenants/subscription` — ubah plan & expiry (body: `{id, plan?, expired_at?}` — plan: free, premium_1m, premium_3m, premium_6m, premium_1y)
  - `GET /api/admin/plans` — daftar konfigurasi plan (max_products, report_days)
  - `PATCH /api/admin/plans` — ubah konfigurasi plan (body: `{plan_slug, duration_days?, max_products?, report_days?, price_rupiah?}` — duration_days: hari langganan, 0 = unlimited)

## PWA & Offline

- **vite-plugin-pwa:** Service worker, manifest (`name: HSmart`, `theme_color: #16a34a`), installable app.
- **Offline:** Sales and expenses are stored in IndexedDB (Dexie) when offline and synced when online.

## Production Build (tanpa Docker)

```bash
# Frontend
cd frontend && npm run build
# Output: dist/ — serve dengan Nginx atau static host.

# Backend
cd backend && go build -o bin/server ./cmd/server
```

## Production dengan Docker

Docker dipakai hanya untuk production/deployment, bukan untuk development.

- **Backend:** `backend/Dockerfile` — build image Go API.
- **Infra:** `docker-compose.yml` — PostgreSQL, Redis, dan service backend untuk production.

Contoh deploy:

```bash
# Build & jalankan (production)
docker compose up -d postgres redis
# Migrasi sekali: psql -h localhost -U postgres -d hsmart -f backend/migrations/001_init.up.sql
docker compose up -d backend
# Frontend: build (npm run build) lalu serve dist/ lewat Nginx
```

Atau deploy backend sebagai image, PostgreSQL/Redis di managed service; frontend di CDN/static host.

## Deployment (Nginx → Vue PWA + Go API)

- Nginx: serve `frontend/dist` untuk `/`, proxy `/api` ke backend Go.
- Backend: listen `:8080` (atau di belakang Nginx).
- Set `DATABASE_URL`, `JWT_SECRET`, `REDIS_ADDR`, `BASE_URL`, `UPLOAD_DIR` di environment production. Lihat `backend/.env.example`.

## Backup

Backup database dan folder uploads secara berkala:

```bash
# Backup semua (db + uploads)
DATABASE_URL=postgres://user:pass@host:5432/hsmart ./scripts/backup-all.sh

# Hanya database
DATABASE_URL=postgres://... ./scripts/backup-db.sh

# Hanya uploads
UPLOAD_DIR=/path/to/uploads ./scripts/backup-uploads.sh
```

Output default: `./backups/db/` dan `./backups/uploads/`. Untuk cron harian jam 02:00:

```
0 2 * * * cd /path/to/hsmartapp && DATABASE_URL=postgres://... ./scripts/backup-all.sh
```

## Monitoring & Logging

- **Request logging:** Tiap request dicatat (method, path, status, duration, IP). Error (4xx/5xx) dilog dengan prefix `[ERROR]` untuk monitoring.
- **Rate limiting:** Sudah aktif (100 req/menit per IP). Untuk production skala besar, pertimbangkan Redis-based limiter.
- **Log output:** stdout (gunakan systemd/docker log driver untuk pengumpulan terpusat).

## Multi-tenant

All business data is scoped by `tenant_id`. Every request must send `X-Tenant-ID` (and valid JWT). Registration creates a tenant and an owner user; login returns `tenant_id` for subsequent requests.

## Plans (konfigurasi oleh superadmin)

- **Free:** Max produk & hari laporan diatur di `/admin/plans` (default: 10 produk, 7 hari).
- **Premium 1/3/6/12 Bulan:** Max produk & hari laporan diatur per plan (default: unlimited produk, 30–365 hari).

---

**Principle:** HSmart stays super simple, fast, mobile-first, offline-first, and easier than writing sales in a notebook.
