# AGENTS.md

This repository is a **hotel management system** with:

- **Go backend** at the repository root (`cmd/server/`)
- **Nuxt/Vue frontend** (`frontend/`)
- **Wails desktop app** (`wails/`)
- **PostgreSQL** as persistence (auto-migrated by GORM)

Use this guide when making changes so backend/frontend behavior stays consistent.

## What the project does

The app manages core hotel operations:

- Authentication and session management
- Rooms, guests, reservations
- Accounting (accounts, expenses, income)
- Parking (lots, spots, vehicles, transactions, stats)

The frontend is a SPA (Nuxt with `ssr: false`) that talks to the Go API under `/api/*`.

## High-level architecture

### Backend (Go)

Entry points and app lifecycle:

- `cmd/server/main.go`: starts server, exposes `seed` command
- `app.go`: wires config, DB, repositories, services, router, graceful shutdown
- `router.go`: registers health and API modules; applies middleware

Layers:

- `internal/models`: GORM models and schema types
- `internal/db`: DB open/migrate/seed logic
- `internal/repository`: DB access interfaces + implementations
- `internal/service`: business services wrapping repositories
- `internal/httpapi`: middleware, helpers, and domain modules

Route modules:

- `_system` (`/healthz`, `/readyz`)
- `auth` (`/api/auth/*`)
- `rooms`, `guests`, `reservation`, `users`, `accounting`, `parking`

### Frontend (Nuxt/Vue)

Primary structure in `frontend/app/`:

- `pages/`: route screens by domain
- `layouts/default.vue`: shell/sidebar/topbar
- `middleware/auth.global.ts`: global auth and role checks
- `stores/auth.ts`: Pinia auth/session state
- `assets/css/main.css`: global styles

Nuxt config:

- `frontend/nuxt.config.ts`
  - `ssr: false`
  - dev proxy for `/api`, `/healthz`, `/readyz` to `http://127.0.0.1:8080`

## Data and API conventions

- IDs are numeric (`uint` in Go), parsed from path params.
- JSON field names are mostly **camelCase** (`json:"firstName"` etc.).
- Backend update handlers accept partial maps and normalize camelCase to snake_case before DB updates.
- Generic CRUD handlers live in `internal/httpapi/resource_helpers.go` and are reused across modules.
- Error responses use `{"error":"..."}` with stable error keys (for example `invalid_id`, `not_found`, `update_failed`).
- Auth uses cookie-based sessions (`SESSION_COOKIE`, default `auth_session`).

## Auth and authorization behavior

- Protected API routes are wrapped with `a.Auth` middleware in `router.go`.
- Public auth endpoints:
  - `POST /api/auth/login`
  - `POST /api/auth/logout` (requires auth middleware in module)
  - `GET /api/auth/me` (requires auth)
- Frontend route access is enforced in `frontend/app/middleware/auth.global.ts`.
- Page-level role checks are declared via `definePageMeta({ requiresRole: [...] })`.

## Permissions

- Defined in `internal/db/seed/permissions.json` as an object: `{ category: { pageKey: { action: "page:action" } } }`
- Page keys are camelCase (e.g., `guestsSettle`, `roomsRack`) to be valid TS object keys
- Values use slashes (e.g., `"guests/settle:read"`)
- Embedded in Go via `//go:embed`; seed package's `init()` generates `frontend/app/utils/permissions.gen.ts`
- Translations in `internal/db/seed/translations.json` use `page:action` as keys (e.g., `"guests/settle:read"`)
- Model uses `Resource`, `Action`, and `CategoryID`; `PermissionCategory` has `Slug` field

## Local development commands

From repo root:

- `just dev` — run backend (`air`) + frontend dev server in parallel
- `just build` — build backend and frontend
- `just seed` — seed reference data (also generates permissions.gen.ts on init)

Direct commands:

- Backend: `go run ./cmd/server`
- Wails dev: `cd wails && wails dev`
- Build Wails app: `cd wails && wails build`

## Environment variables

Backend (`internal/config/config.go`):

- `APP_ADDR` (default `:8080`)
- `DB_PATH` (default `./data/hotel.db`)
- `SESSION_COOKIE`
- `SESSION_TTL_HOURS`
- `REQUEST_TIMEOUT_SECONDS`
- `SHUTDOWN_TIMEOUT_SECONDS`
- `READ_TIMEOUT_SECONDS`, `WRITE_TIMEOUT_SECONDS`, `IDLE_TIMEOUT_SECONDS`
- `SEED_ADMIN_EMAIL`, `SEED_ADMIN_PASSWORD`, `SEED_ADMIN_FIRST_NAME`, `SEED_ADMIN_LAST_NAME`

Frontend (`frontend/nuxt.config.ts` runtime config):

- `BACKEND_URL`
- `NUXT_PUBLIC_HOTEL_NAME`
- optional auth/admin config values used by Nuxt runtime config

## Coding conventions for this repo

### Go backend conventions

- Keep domain logic in `service`, persistence in `repository`, HTTP concerns in `httpapi`.
- Reuse generic CRUD handlers where possible before adding one-off handlers.
- For new resources:
  1. add model in `internal/models`
  2. include it in `models.AllPtr()` / `models.All()` when appropriate
  3. add module handler package under `internal/httpapi/<domain>`
  4. register module in `router.go`
- Prefer consistent response shapes and status codes with existing modules.
- Keep middleware behavior consistent: timeout, panic recovery/logging, JSON content type.

### Frontend conventions

- Domain pages live in `frontend/app/pages/<domain>/...`.
- Use `definePageMeta` for role requirements.
- Use Pinia auth store (`useAuthStore`) for session state and role checks.
- Consume backend via `$fetch('/api/...')`.

## Notes and gotchas

- Wails app embeds `frontend/.output`; production builds must generate frontend output before embedding.
- Wails has its own `go.mod` (`wails/go.mod`) to isolate it from the server which uses `go-fuego`.
- No dedicated test suite is currently present; validate changes by running the app and exercising affected API/UI flows.

## Change checklist (recommended)

When shipping a feature touching both stacks:

1. Implement backend model/repository/service/http handlers.
2. Register routes and verify auth/role requirements.
3. Regenerate TS types with `just gen` if model contracts changed.
4. Implement/update frontend pages and API calls.
5. Run build/dev checks (`just build` and targeted manual verification).
