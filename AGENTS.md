# AGENTS.md

This repository is a **hotel management system** with:

- **Go backend** at the repository root (`cmd/server/`)
- **Nuxt/Vue frontend** (`frontend/`)
- **Wails desktop app** (`wails/`)
- **PostgreSQL** as persistence (auto-migrated by GORM)

Use this guide when making changes so backend/frontend behavior stays consistent.

## What the project does

The app manages core hotel operations:

- **Authentication and session management** (cookie-based sessions)
- **Rooms** (types, floors, amenities, status tracking)
- **Guests** (profile, companions, status tracking)
- **Reservations** (future bookings with status workflow)
- **Stays / Reception** (actual check-ins with invoice generation)
- **Accounting** (accounts, expenses, income, payment methods)
- **Parking** (lots, spots, vehicles, transactions)
- **Restaurant** (inventory, bills, meal transactions)
- **Travel Agencies** (contracting parties)
- **Services** (billable service catalog)
- **Sana Integration** (Iranian hotel regulatory system)
- **Permissions** (role-based access with templates)

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
- `rooms` (`/api/rooms/*`, `/api/rooms/rack`, `/api/rooms/statuses`, etc.)
- `guests` (`/api/guests/*`, `/api/guests/{id}/settle`, `/api/guests/{id}/checkout`, `/api/guests/{id}/with-stay`)
- `stays` (`/api/stays/*`, `/api/stays/{id}/check-in`, `/api/stays/{id}/change-room`, `/api/stays/{id}/change-duration`)
- `reservation` (`/api/reservation/*`, `/api/reservation/statuses`, `/api/reservation/check-availability`)
- `users` (`/api/users/*`)
- `admins` (`/api/admins/*`)
- `accounting` (`/api/accounting/*`)
- `parking` (`/api/parking/*`)
- `restaurant` (`/api/restaurant/*`)
- `services` (`/api/services/*`)
- `travel-agencies` (`/api/travel-agencies/*`)
- `hotels` (`/api/hotels/*`, `/api/hotels/{id}/settings`)
- `permissions` (`/api/permissions/*`, `/api/permissions/templates`, `/api/permissions/user/{id}`)
- `dashboard` (`/api/dashboard/*`)
- `common` (`/api/common/*`)
- `sana` (`/api/sana/*`)

### Frontend (Nuxt/Vue)

Primary structure in `frontend/app/`:

- `pages/`: route screens by domain
  - `login.vue`, `index.vue` (dashboard)
  - `guests/`, `rooms/`, `reservations/`, `stays/`
  - `accounting/`, `parking/`, `restaurant/`
  - `users/`, `admins/`, `permissions/`
  - `hotels/`, `travel-agencies/`, `sana/`
  - `settings/`, `profile/`
- `layouts/default.vue`: shell/sidebar/topbar
- `middleware/auth.global.ts`: global auth and permission checks
- `stores/auth.ts`: Pinia auth/session state
- `assets/css/main.css`: global styles
- `utils/permissions.gen.ts`: auto-generated permissions enum
- `utils/client/`: auto-generated hey-api SDK

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
- Data is retrieved by utilizing **hey-api** generated SDK functions and **pinia-colada** (`useQuery`) for the fetching layer and async state management.
- Status values (rooms, guests, stays, reservations, parking) are **DB-driven** via `TranslateBase` tables with `slug`, `label`, and `colorHex`. The backend resolves translations based on the `Accept-Language` header.

## Room Status System

Room statuses are **static and DB-driven** (single source of truth):

- `available` — room is free for assignment
- `occupied` — guest is currently staying
- `reserved` — future reservation holds the room
- `cleaning` — post-checkout, being prepared
- `under_repair` — temporarily out of service

Status is stored in `rooms.status_id` (foreign key to `room_statuses`). Updates happen explicitly in handlers (e.g., check-in sets `occupied`, checkout sets `cleaning`). No dynamic computation.

## Guest Status System

Guest statuses are **DB-driven** with auto-update via `AfterSave` hooks:

- `waiting` — has an active reservation or pending stay
- `resident` — has an active stay with `resident` status
- `checked_out` — all stays completed, no active stay
- `cancelled` — stays or reservations cancelled/expired

Guest status updates automatically when stays or reservations change. Priority: `resident` > `waiting` > `checked_out` > `cancelled`.

## Stay and Invoice Flow

1. **Guest creation** (`POST /api/guests/with-stay`) creates guest + stay + invoice + marks room `occupied`
2. **Check-in** (`POST /api/stays/{id}/check-in`) marks stay `resident`, room `occupied`, generates invoice items
3. **During stay** — add services (`POST /api/stays/{id}/services`), parking, restaurant charges
4. **Payment** — go to guest settlement page (`/guests/{id}/settle`), pay via `POST /api/guests/{id}/settle` (applies to invoice first)
5. **Checkout** — once balance is 0, `POST /api/guests/{id}/checkout` marks stay `checked_out`, room `cleaning`, creates `GuestCheckout` record

## Checkout and Settlement

The **only** checkout endpoint is `POST /api/guests/{id}/checkout`. It:

- Validates all balances are 0 (invoice + parking + restaurant)
- Creates a `GuestCheckout` record (audit trail)
- Marks stay `checked_out`, room `cleaning`
- Settles all parking transactions and restaurant bills
- Links all charges to the checkout record via `checkout_id`

There is no stay-level checkout or payment. All payment goes through the guest settlement page.

## Auth and authorization behavior

- Protected API routes are wrapped with `a.Auth` middleware in `router.go`.
- Public auth endpoints:
  - `POST /api/auth/login`
  - `POST /api/auth/logout` (requires auth middleware in module)
  - `GET /api/auth/me` (requires auth)
- Frontend route access is enforced in `frontend/app/middleware/auth.global.ts`.
- Page-level permission checks are declared via `definePageMeta({ requiresPermission: PERMISSIONS.guests.guests.read })`.

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
- `just gen-perms` — regenerate permissions TypeScript enum

Direct commands:

- Backend: `go run ./cmd/server`
- Seed: `go run ./cmd/server seed`
- Wails dev: `cd wails && wails dev`
- Build Wails app: `cd wails && wails build`

## Environment variables

Backend (`internal/config/config.go`):

- `APP_ADDR` (default `:8080`)
- `DB_PATH` (default `postgres://hotel_user:hotel_password@localhost:5432/hotel_db?sslmode=disable`)
- `SESSION_COOKIE` (default `auth_session`)
- `HOTEL_COOKIE` (default `hotel_id`)
- `SESSION_TTL_HOURS`
- `REQUEST_TIMEOUT_SECONDS`
- `SHUTDOWN_TIMEOUT_SECONDS`
- `READ_TIMEOUT_SECONDS`, `WRITE_TIMEOUT_SECONDS`, `IDLE_TIMEOUT_SECONDS`
- `SEED_ADMIN_EMAIL`, `SEED_ADMIN_PASSWORD`, `SEED_ADMIN_FIRST_NAME`, `SEED_ADMIN_LAST_NAME`
- `SEED_HOTEL_NAME`, `SEED_HOTEL_ADDRESS`, `SEED_HOTEL_PHONE`, `SEED_HOTEL_EMAIL`
- `SANA_KELID_VAHAED`, `SANA_KELID_PEIMANKAR`, `SANA_CODE_VAHAED`

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
  2. include it in `models.AllForDB()` / `models.AllForTypeGen()` when appropriate
  3. add module handler package under `internal/httpapi/<domain>`
  4. register module in `router.go`
- Prefer consistent response shapes and status codes with existing modules.
- Keep middleware behavior consistent: timeout, panic recovery/logging, JSON content type.
- Status enums (slugs) use typed constants (e.g., `RoomStatusSlug`, `GuestStatusSlug`, `StayStatusSlug`).
- For DB-driven statuses with translations, use `TranslateBase` embedded in the model (e.g., `RoomStatus`, `GuestStatus`).
- Use `BeforeSave` hooks to resolve slug strings to `StatusID` foreign keys.
- Use `AfterSave` hooks for cascading status updates (e.g., stay changes trigger guest status recalculation).

### Frontend conventions

- Domain pages live in `frontend/app/pages/<domain>/...`.
- Use `definePageMeta` for permission requirements.
- Use Pinia auth store (`useAuthStore`) for session state and role checks.
- Consume backend via the **auto-generated hey-api SDK** (`getApi*`, `postApi*`, etc.) from `~/utils/client`.
- Use `useQuery` from `@pinia/colada` for async data fetching.
- Use `useI18n()` for translations. Backend-driven status labels are auto-translated via the API and should not have frontend i18n keys.
- Use `UBadge` with `colorHex` from the backend for status display.

## Notes and gotchas

- Wails app embeds `frontend/.output`; production builds must generate frontend output before embedding.
- Wails has its own `go.mod` (`wails/go.mod`) to isolate it from the server which uses `go-fuego`.
- No dedicated test suite is currently present; validate changes by running the app and exercising affected API/UI flows.
- The `doc/openapi.json` is auto-generated by fuego when the server starts. The frontend SDK (`frontend/app/utils/client`) is auto-generated from this file by `@hey-api/nuxt`.
- **Always** run the server briefly to regenerate `openapi.json` after adding/modifying backend routes, then regenerate the frontend SDK via `bun run typecheck` or `bun run build`.
- When adding DB-driven statuses, also add seed data in `internal/db/seed/seed.go` and translations in `internal/db/seed/translations.json`.
- The `guests` list page has a settlement action (credit card icon) that navigates to `/guests/{id}/settle`.
- The `stays` detail page no longer has its own checkout or payment flow. It links to the guest settlement page.
- Guest checkout is **only** possible when all balances (invoice + parking + restaurant) are zero.

## Change checklist (recommended)

When shipping a feature touching both stacks:

1. Implement backend model/repository/service/http handlers.
2. Register routes and verify auth/role requirements.
3. Run the server to regenerate `doc/openapi.json`.
4. Regenerate frontend SDK via `bun run typecheck` or `bun run build`.
5. Implement/update frontend pages and API calls.
6. Run build/dev checks (`just build` and targeted manual verification).
7. Update `AGENTS.md` if the change affects architectural conventions or adds new patterns.
