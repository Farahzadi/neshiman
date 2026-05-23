# AGENTS.md

## Project

Seat plan application. Admins define rooms with rectangular seats arranged on a grid. Seats have labels, can rotate/move, and are owned by teams. Users (viewers) reserve seats per day with a weekly limit (e.g., 2 days/week). Team admins configure reservation capacity for their members. Role hierarchy: superadmin > team admin > viewer. Cross-team seat requests require approval from the owning team's admin.

## Quick Start

```bash
make setup    # First time: start postgres, migrate, sqlc generate, install deps
make dev      # Start all services (postgres + backend + both frontends)
```

## Monorepo Layout

- `backend/` — Go service, hexagonal architecture, PostgreSQL, sqlc
- `apps/admin/` — SolidJS app (admin management UI) [@neshiman/admin]
- `apps/viewer/` — SolidJS app (seat reservation UI) [@neshiman/viewer]
- `packages/api-types/` — Auto-generated TypeScript types from backend Swagger spec [@neshiman/api-types]
- `packages/api-client/` — TanStack Query hooks + typed fetch wrapper [@neshiman/api-client]
- `packages/tsconfig/` — Shared TypeScript base configs [@neshiman/tsconfig]
- `packages/tailwind-config/` — Shared TailwindCSS preset [@neshiman/tailwind-config]
- `docker-compose.yml` — production orchestration
- `turbo.json` — Turborepo pipeline config
- `pnpm-workspace.yaml` — pnpm workspace definition

## Key Decisions

- **Backend**: Hexagonal architecture — keep domain/ports free of infrastructure deps. Adapters (HTTP handlers, DB repos) live outside the core.
- **sqlc**: Generates Go structs from SQL queries. Never edit generated files; modify `.sql` sources in `backend/db/query/` and re-run `sqlc generate`.
- **Frontend**: SolidJS ≠ React. Use signals (`createSignal`), stores (`createStore`), no virtual DOM diffing. JSX compiles to real DOM.
- **Code quality**: ESLint + Prettier enforced via git hooks. Run lint/format before committing.
- **Swagger**: Use `swaggo/swag` for Swagger 2.0 spec generation from Go handler annotations. Serves `/swagger/index.html` in dev.
- **Turborepo + pnpm**: Frontend monorepo via pnpm workspaces + Turborepo. Root `package.json` has `packageManager: pnpm@9.15.0`. All frontend commands run at root via `pnpm turbo <task>`. Lockfile: `pnpm-lock.yaml`.
- **openapi-typescript**: Generates TypeScript types from `backend/docs/swagger.json` into `packages/api-types/src/v1.d.ts`. Regenerate with `make types-gen` or `pnpm generate:types`.
- **Shared packages**: `@neshiman/api-types` (types), `@neshiman/tsconfig` (base tsconfig), `@neshiman/tailwind-config` (Tailwind preset). Imported via `workspace:*` protocol.

## Backend Commands

Run all from `backend/`:

| Command | Description |
|---------|-------------|
| `go run ./cmd/server` | Start dev server on :8080 |
| `go test ./...` | Run all tests |
| `go test -run Integration ./...` | Run only integration tests |
| `sqlc generate` | Regenerate `backend/db/sqlc/` from `backend/db/query/` |
| `migrate -path db/migrations -database "$DB_URL" up` | Apply migrations |
| `swag init -g ./cmd/server/main.go --output ./docs` | Regenerate swagger docs from handler annotations |
| `go run ./cmd/seed` | Seed database with sample data |
| `go mod tidy` | Sync Go module dependencies |

The entrypoint is `backend/cmd/server/main.go` — wires config → DB pool → repositories → services → HTTP router.

## Frontend Commands

Run from repo root (or via `pnpm --filter <package>`):

| Command | Description |
|---------|-------------|
| `pnpm turbo dev` | Start both dev servers (admin: :3000, viewer: :3001) |
| `pnpm turbo build` | Build all apps to `dist/` |
| `pnpm turbo test` | Run Vitest in all apps |
| `pnpm turbo lint` | ESLint across all apps |
| `pnpm turbo typecheck` | `tsc --noEmit` across all apps |
| `pnpm --filter @neshiman/api-types generate` | Regenerate API types from swagger.json |
| `pnpm --filter <app> dev` | Dev server for a single app |
| `pnpm --filter <app> <script>` | Run any script in a specific app |

## Agent Notes

- **sqlc generated code**: `backend/db/sqlc/` is DO NOT EDIT. Modify `backend/db/query/*.sql` and run `sqlc generate`.
- **Migrations**: Written in `backend/db/migrations/` using golang-migrate naming convention (`NNNNNN_name.up.sql` / `.down.sql`).
- **Domain pure**: `internal/domain/` imports nothing outside stdlib. No DB, no HTTP.
- **Wiring only in main.go**: `cmd/server/main.go` is the composition root — the only place concrete types cross package boundaries.
- **Swagger**: Annotate new handlers with `@Summary`, `@Tags`, `@Param`, `@Success`, `@Router` comments. Regenerate with `make swagger-gen` or `swag init` from `backend/`.
- **openapi-typescript**: After swagger-gen, run `make types-gen` (or `pnpm generate:types`) to regenerate TypeScript types from the updated swagger.json. Generated file `packages/api-types/src/v1.d.ts` is committed.
- **Turborepo pipeline**: `build` depends on `^build` so `api-types` compiles before apps that depend on it. `dev` is persistent (no caching).
- **pnpm**: Enable via `corepack enable`. Workspaces are defined in `pnpm-workspace.yaml`. Use `workspace:*` protocol for local deps.
- **Seat geometry**: Rectangular (not triangular). Grid-based room layout.
- **Memberships**: Non-overlapping. A user belongs to exactly one team.
- **Reservations**: Per-day, one person per seat. Weekly limit configurable per user by team admin.
- **Cross-team**: Requests require approval from the owning team's admin.
- **Auth**: JWT via `golang-jwt/jwt/v5` (HS256, 24h expiry). `POST /auth/login` returns `{ token, user }`. Role hierarchy: superadmin > team_admin > viewer. `AuthMiddleware` validates and sets context; `RequireRole(minimumRole)` enforces hierarchy. Migration 7 adds `password_hash` column. Seed users all have password "password".
- **Both apps**: Store JWT in `localStorage` key `neshiman_token`, user info in `neshiman_user`. Login pages at `/login`. Route guards redirect unauthenticated users. Admin app additionally restricts to superadmin/team_admin roles.

## Session Progress (May 2026)

### Built
- **Backend** (`backend/`): Full hexagonal Go service with chi router, pgx/v5, sqlc, golang-migrate
  - Domain: Room, Seat, Team, User, Reservation, CrossTeamRequest + value objects + typed errors
  - Ports: 6 repository interfaces + TxManager
  - Application: RoomService, ReservationService (weekly limit enforcement)
  - HTTP: Chi router, Room + Reservation handlers, CORS/RequestID/Auth/Role middleware, DTOs
  - Postgres: Pool, TxManager, Room/Seat/User/Reservation repositories
  - DB: 6 migration pairs + 6 sqlc query files + Mermaid ER diagram
- **Frontend admin** (`frontend-admin/`): SolidJS + Vite + TailwindCSS + TypeScript scaffold
  - Pages: Dashboard, Rooms, Seats, Teams (placeholder)
- **Frontend viewer** (`frontend-viewer/`): SolidJS + Vite + TailwindCSS + TypeScript scaffold
  - Pages: Dashboard, Reservations (placeholder)
- **Tooling**: Root Makefile, AGENTS.md, SPEC.md (updated), .gitignore, .editorconfig
- **Frontend monorepo**: Turborepo + pnpm workspaces, 3 shared packages (`api-types`, `tsconfig`, `tailwind-config`)

### Current state
- `go build` + `go vet` pass across all packages
- Postgres running in Docker, 6 migrations applied, `sqlc generate` done
- Room, Reservation, Team, Seat, User, CrossTeamRequest: all fully wired end-to-end
- Auth/Role middleware: wired to `/api/v1` routes, stub identity until Phase 6
- Swagger: `/swagger/index.html` serves browsable API docs (22 paths documented, incl. bulk sync)
- Phase 2 complete: Turborepo + pnpm workspaces. Frontend apps moved to `apps/`, shared packages in `packages/`. TypeScript types auto-generated from backend DTOs via openapi-typescript. Shared tsconfig + tailwind preset.
- **Phase 3 complete**: Seed script (`make seed`), test DB setup in `make test`, 144 passing tests (28 domain + 62 service + 54 integration)
- **Phase 3.5 complete**: `packages/api-client/` with TanStack Query hooks for all 21 endpoints, typed fetch wrapper with `setAuthHeader()`, query key factory, `QueryClientProvider` + `SolidQueryDevtools` wired into both apps
- **Phase 4 complete**: Full admin UI with sidebar layout, Rooms CRUD, Teams with inline user management, Users page with team filter and weekly limit editing, Cross-Team Requests approval queue
- **Phase 4.5 complete**: Room Visual Editor — canvas-based pan/zoom room editor with click-to-place and drag-to-move seats, sparse grid support (only placed seats rendered), Select/Place modes, properties panel (label/team/rotation), dirty state tracking, bulk save. Backend `PUT /rooms/{id}/seats` endpoint with atomic transaction. 8 new tests (144 total). Old `Seats.tsx` replaced, route changed to `/rooms/:id/edit`.
- **Phase 5 complete**: Viewer app — Dashboard with onboarding/overview, Rooms card grid, RoomDetail with visual seat map and reservation flow, Reservations single-day list with cancel, Cross-Team Requests tab-filtered list. Auth bypass: localStorage user ID sent as Bearer token, backend middleware treats it as identity. Backend fix: added `GET /reservations` with `date`/`user_id` query params. 63 application tests passing.
- **Phase 6 complete**: JWT authentication — Login endpoint (`POST /auth/login`) with bcrypt password verification, JWT signing (HS256, 24h expiry), auth middleware validates tokens and extracts user ID + role, role middleware enforces hierarchy (superadmin > team_admin > viewer) on sensitive routes. Migration adds `password_hash` column. Seed users default password: "password". Login pages for both admin and viewer, route guards redirect to `/login` when unauthenticated.
- **Phase 7 — Soft Delete + Polish**: Migration `000008_add_deleted_at` adds `deleted_at TIMESTAMPTZ` to all 6 tables. All SELECTs filter `deleted_at IS NULL`, DELETEs become soft-delete UPDATEs. Domain entities have `DeletedAt *time.Time`. Room deletion cascades to seats, returns seat count. HTTP handlers return `200 { "deleted": true }` instead of 204. Admin UI has `ConfirmModal` component replacing `confirm()` dialogs. Toast notification system replaces `alert()` calls. Login uses username+password (migration `000009` makes email nullable, adds name UNIQUE). Superadmin protected from deletion. Weekly limit capped at 5. GitHub Actions workflow builds & pushes 3 images (backend, admin, viewer) to `ghcr.io` on push to main. 147 tests passing.
- **Phase 8 — Deployment**: `GET /health` endpoint (unauthenticated). Auto-run DB migrations on startup via `golang-migrate` library in `main.go`. Production CORS middleware reads `CORS_ORIGINS` env var (comma-separated origins, falls back to `*` in dev). `docker-compose.prod.yml` uses pre-built ghcr.io images with required secrets. CI pipeline (`ci.yml`) runs tests + lint + typecheck + build, then Docker push on main. `README.md` with architecture, setup guide, deployment docs. `.env.example` for required vars. Backend Go version bumped to 1.24. `Makefile` has `deploy-prod` targets.
- Next steps: None — all phases complete
