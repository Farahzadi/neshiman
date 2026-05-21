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
