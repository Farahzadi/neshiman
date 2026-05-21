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
- `frontend-admin/` — SolidJS app (admin management UI)
- `frontend-viewer/` — SolidJS app (seat reservation UI)
- `docker-compose.yml` — production orchestration

## Key Decisions

- **Backend**: Hexagonal architecture — keep domain/ports free of infrastructure deps. Adapters (HTTP handlers, DB repos) live outside the core.
- **sqlc**: Generates Go structs from SQL queries. Never edit generated files; modify `.sql` sources in `backend/db/query/` and re-run `sqlc generate`.
- **Frontend**: SolidJS ≠ React. Use signals (`createSignal`), stores (`createStore`), no virtual DOM diffing. JSX compiles to real DOM.
- **Code quality**: ESLint + Prettier enforced via git hooks. Run lint/format before committing.
- **Swagger**: Use `swaggo/swag` for OpenAPI spec generation from Go handler annotations. Serves `/swagger/index.html` in dev.
- **Turborepo**: Frontend monorepo tool for shared packages. Standard layout: `apps/` (admin, viewer) + `packages/` (shared, tsconfig, tailwind-config).

## Backend Commands

Run all from `backend/`:

| Command | Description |
|---------|-------------|
| `go run ./cmd/server` | Start dev server on :8080 |
| `go test ./...` | Run all tests |
| `sqlc generate` | Regenerate `backend/db/sqlc/` from `backend/db/query/` |
| `migrate -path db/migrations -database "$DB_URL" up` | Apply migrations |
| `go mod tidy` | Sync Go module dependencies |

The entrypoint is `backend/cmd/server/main.go` — wires config → DB pool → repositories → services → HTTP router.

## Frontend Commands

Run from each frontend app directory:

| Command | Description |
|---------|-------------|
| `npm run dev` | Dev server (admin: :3000, viewer: :3001) |
| `npm run build` | Production build to `dist/` |
| `npm run test` | Vitest |
| `npm run lint` | ESLint |
| `npm run format` | Prettier |
| `npm run typecheck` | `tsc --noEmit` |

## Agent Notes

- **sqlc generated code**: `backend/db/sqlc/` is DO NOT EDIT. Modify `backend/db/query/*.sql` and run `sqlc generate`.
- **Migrations**: Written in `backend/db/migrations/` using golang-migrate naming convention (`NNNNNN_name.up.sql` / `.down.sql`).
- **Domain pure**: `internal/domain/` imports nothing outside stdlib. No DB, no HTTP.
- **Wiring only in main.go**: `cmd/server/main.go` is the composition root — the only place concrete types cross package boundaries.
- **Seat geometry**: Rectangular (not triangular). Grid-based room layout.
- **Memberships**: Non-overlapping. A user belongs to exactly one team.
- **Reservations**: Per-day, one person per seat. Weekly limit configurable per user by team admin.
- **Cross-team**: Requests require approval from the owning team's admin.

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

### Current state
- `go build` + `go vet` pass across all packages
- Postgres running in Docker, 6 migrations applied, `sqlc generate` done
- Room and Reservation: fully wired end-to-end
- Seat and User: domain + ports + postgres adapters exist, but no services/handlers/routes
- Team and CrossTeamRequest: domain + ports + SQL queries exist, no postgres adapters/services/handlers
- Auth/Role middleware: defined but stubs, not wired to routes
- Both frontends: SolidJS scaffold with rplaceholders for all feature pages, no API layer
- Zero tests anywhere
- Next steps: see `roadmap.md`
