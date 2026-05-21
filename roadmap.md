# Roadmap

Seat plan application — neshiman.

## Phase 1 — Backend Completion

Finish wiring all core entities end-to-end (domain → port → postgres → service → handler → route).

- [x] **Team postgres adapter** — Implement `TeamRepository` in `internal/adapters/postgres/team_repository.go`
- [x] **Team service** — `internal/application/team_service.go` (CRUD)
- [x] **Team handler + routes** — `internal/adapters/http/handlers/team.go` + routes in `router.go`
- [x] **Seat service** — `internal/application/seat_service.go` (CRUD, list-by-room, move, rotate)
- [x] **Seat handler + routes** — `internal/adapters/http/handlers/seat.go` + routes in `router.go`
- [x] **User service** — `internal/application/user_service.go` (CRUD, weekly limit updates)
- [x] **User handler + routes** — `internal/adapters/http/handlers/user.go` + routes in `router.go`
- [x] **CrossTeamRequest postgres adapter** — `internal/adapters/postgres/cross_team_request_repository.go`
- [x] **CrossTeamRequest service** — `internal/application/cross_team_request_service.go` (submit, list, approve, reject)
- [x] **CrossTeamRequest handler + routes** — `internal/adapters/http/handlers/cross_team_request.go` + routes in `router.go`
- [x] **Wire Auth + Role middleware** — Apply to routes; keep stub identity for now (real auth deferred)

## Phase 1.5 — Swagger / OpenAPI

Add API documentation with a browsable Swagger UI in development.

- [x] **Integrate swaggo/swag** — Generate OpenAPI spec from Go handler annotations
- [x] **Swagger UI endpoint** — Serve `/swagger/index.html` in dev mode (chi route)
- [x] **Annotate all handlers** — Add `@Summary`, `@Tags`, `@Param`, `@Success`, `@Router` comments to each handler
- [x] **DTO schemas** — Ensure request/response types are documented for Swagger output

## Phase 2 — Shared Frontend Package

Extract common code into a shared package used by both frontends, using Turborepo.

- [ ] **Turborepo setup** — Root `turbo.json`, root `package.json` with workspaces
- [ ] **Create `packages/shared/`** — API client, TypeScript types matching backend DTOs, fetch wrapper, shared utilities
- [ ] **Create `packages/tsconfig/`** — Shared TypeScript config (optional)
- [ ] **Create `packages/tailwind-config/`** — Shared Tailwind preset (optional)
- [ ] **Migrate both apps into workspace** — `apps/admin/`, `apps/viewer/` (rename from `frontend-admin/`, `frontend-viewer/`)
- [ ] **Import from `@neshiman/shared`** in both frontend apps
- [ ] **Define API types** matching backend DTOs (`dto/request.go`, `dto/response.go`)

## Phase 3 — Seed Data & Test Fixtures

- [ ] **Seed script** — Populate DB with sample rooms, teams, users, seats, reservations
- [ ] **Backend test helpers** — Test DB setup/teardown utilities
- [ ] **Backend unit tests** — Domain logic, application services
- [ ] **Backend integration tests** — Repository tests against real Postgres

## Phase 4 — Frontend Admin

Build the admin management UI.

- [ ] **Rooms page** — CRUD table, create/edit form, grid dimensions
- [ ] **Seats page** — Visual grid editor for room layout, drag to place seats, assign teams
- [ ] **Teams page** — CRUD list, member management
- [ ] **Users page** — List by team, edit weekly limits
- [ ] **Cross-team requests** — Approval queue with approve/reject actions
- [ ] **Dashboard** — Summary stats (rooms, seats, active reservations)

## Phase 5 — Frontend Viewer

Build the seat reservation UI.

- [ ] **Room grid view** — Visual seat map with availability indicators per date
- [ ] **Reservation flow** — Select date → select seat → confirm
- [ ] **My reservations** — Weekly view of own reservations, cancel action
- [ ] **Cross-team requests** — Request a seat in another team's area
- [ ] **Dashboard** — Upcoming reservations, weekly usage

## Phase 6 — Authentication

Replace stub auth with real JWT-based authentication.

- [ ] **JWT token generation** — Login endpoint, token signing
- [ ] **Auth middleware** — Validate JWT, extract user identity
- [ ] **Role middleware** — Enforce superadmin / team admin / viewer permissions
- [ ] **Frontend auth** — Login page, token storage, attach to API requests
- [ ] **Route guards** — Protect admin routes, redirect unauthenticated users

## Phase 7 — Polish & Deployment

- [ ] **Frontend tests** — Component tests with Vitest
- [ ] **Lint + format pass** — Both frontend apps, ESLint + Prettier
- [ ] **Docker production build** — Verify multi-stage builds, docker-compose production profile
- [ ] **README update** — Setup instructions, architecture overview
- [ ] **CI pipeline** — GitHub Actions (lint, test, build)
