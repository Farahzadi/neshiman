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

- [x] **Turborepo setup** — Root `turbo.json`, root `package.json` with workspaces
- [x] **Create `packages/api-types/`** — Auto-generated TypeScript types from backend Swagger spec via openapi-typescript
- [x] **Create `packages/tsconfig/`** — Shared TypeScript base config
- [x] **Create `packages/tailwind-config/`** — Shared TailwindCSS preset
- [x] **Migrate both apps into workspace** — `apps/admin/`, `apps/viewer/` (renamed from `frontend-admin/`, `frontend-viewer/`)
- [x] **Import from `@neshiman/api-types`** in both frontend apps
- [x] **Define API types** — Auto-generated from `backend/docs/swagger.json` via openapi-typescript v5

## Phase 3 — Seed Data & Test Fixtures

- [x] **Seed script** — `backend/cmd/seed/` populates DB with sample rooms, teams, users (all roles), seats, reservations, cross-team requests
- [x] **Backend test helpers** — `TestMain` with pool setup, `truncate()`, manual mock implementations for all 6 repository interfaces
- [x] **Backend unit tests** — 28 domain tests, 59 application service tests (87 total)
- [x] **Backend integration tests** — 49 repository tests against real Postgres (136 total across all packages)

## Phase 3.5 — API Client Library

Build a shared TanStack Query-based API client package used by both frontends.

- [x] **Create `packages/api-client/`** — Workspace package with typed fetch wrapper
- [x] **`client.ts`** — Fetch wrapper with `setAuthHeader()` config, error extraction, base URL
- [x] **`query-keys.ts`** — Query key factory for all entities
- [x] **Hook files** — 6 entity files covering all 20 API endpoints
  - Rooms (5 hooks: list, detail, create, update, delete)
  - Seats (6 hooks: list-by-room, detail, create, delete, move, rotate)
  - Teams (4 hooks: list, detail, create, delete)
  - Users (6 hooks: list-by-team, detail, by-email, create, delete, weekly-limit)
  - Reservations (3 hooks: by-date, by-user-date, create)
  - Cross-team requests (6 hooks: by-status, pending-by-team, detail, create, approve, reject)
- [x] **Providers** — `QueryClientProvider` + `SolidQueryDevtools` added to both apps
- [x] **Both apps updated** — `App.tsx` wraps router in Providers, lazy-loads devtools in dev

## Phase 4 — Frontend Admin

Build the admin management UI.

- [x] **Dashboard** — Summary stats (room count, team count)
- [x] **Rooms page** — CRUD table, create/edit modal form, grid dimensions, delete confirmation
- [x] **Seats page** — Visual grid editor (replaced by Phase 4.5 Room Editor)
- [x] **Teams page** — CRUD list with expandable inline user management (create/delete users, edit weekly limits)
- [x] **Users page** — Team filter dropdown, user CRUD table, inline weekly limit editing
- [x] **Cross-team requests** — Cross-team requests** — Status tab filter (pending/approved/rejected), table with approve/reject actions, user names displayed instead of truncated UUIDs
- [x] **Layout** — Sidebar navigation with active state highlighting

## Phase 4.5 — Room Visual Editor

Replace the dense CSS-grid seat editor with a canvas-based room editor supporting sparse grids, pan/zoom navigation, and a realistic room view.

- [x] **Backend bulk seats endpoint** — `PUT /api/v1/rooms/{id}/seats` (create/update/delete in one atomic call)
- [x] **Canvas room editor** — Pan/zoom canvas with grid lines, seats as positioned elements, click-to-place, drag-to-move
- [x] **Toolbar + properties panel** — Mode switching (select/place), team/rotation/label editing, delete
- [x] **Dirty state + save** — Track local edits, bulk save on changes
- [x] **Route change** — `/rooms/:id/edit` replaces `/rooms/:id/seats`
- [x] **Swagger + types regeneration** — Regenerate API docs and TypeScript types
- [x] **Backend tests** — 3 service tests + 3 integration tests for bulk sync (144 total)

## Phase 5 — Frontend Viewer

Build the seat reservation UI.

- [x] **Room grid view** — Visual seat map with availability indicators per date
- [x] **Reservation flow** — Select date → select seat → confirm
- [x] **My reservations** — Single-day view of own reservations, cancel action
- [x] **Cross-team requests** — Request a seat in another team's area, tab-filtered list
- [x] **Dashboard** — Welcome/onboarding, stat cards, quick action links, profile info
- [x] **Auth bypass** — Viewer sends stored user ID as Bearer token; backend middleware treats it as user identity

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
