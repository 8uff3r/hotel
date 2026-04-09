# Hotel Management - Nuxt Frontend + Go Backend

The app now runs with:

- Nuxt frontend/UI in this repository (`app/`)
- Go backend API in `backend/`

The old Nuxt server routes were moved to `legacy_nuxt_server/` and replaced by proxying `/api/**` to Go.

## Architecture

- **Frontend:** Nuxt 4 + Nuxt UI
- **Backend:** Go HTTP API + SQLite (WAL mode)
- **Auth:** Cookie-based sessions (`auth_session`)

## Go backend reliability defaults

- SQLite WAL mode, foreign keys, busy timeout
- Connection pool and request timeout controls
- Panic recovery and structured logs
- Graceful shutdown on SIGINT/SIGTERM
- Health endpoints: `/healthz`, `/readyz`
- Migration tracking table: `schema_migrations`
- Expired session cleanup worker

## Run

1) Start backend:

```bash
cd backend
go mod tidy
go run ./cmd/server
```

2) Start Nuxt frontend (project root):

```bash
bun install
bun run dev
```

Nuxt proxies API traffic to Go via `routeRules`.

## Environment variables

Frontend / Nuxt:

- `BACKEND_URL` (default `http://127.0.0.1:8080`)
- `NUXT_PUBLIC_HOTEL_NAME`

Backend:

- `APP_ADDR` (default `:8080`)
- `DB_PATH` (default `./data/hotel.db`)
- `SESSION_COOKIE`
- `SESSION_TTL_HOURS`
- `REQUEST_TIMEOUT_SECONDS`
- `SHUTDOWN_TIMEOUT_SECONDS`
- `SEED_ADMIN_EMAIL`
- `SEED_ADMIN_PASSWORD`
- `SEED_ADMIN_FIRST_NAME`
- `SEED_ADMIN_LAST_NAME`

Default seeded admin: `admin@hotel.local` / `admin123`
