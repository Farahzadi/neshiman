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

- [x] **JWT token generation** — Login endpoint (`POST /auth/login`) returns signed JWT (HS256, 24h expiry) with `sub` (user ID) and `role` claims
- [x] **Auth middleware** — Validate JWT, extract user ID and role into request context
- [x] **Role middleware** — Hierarchy-based enforcement (superadmin > team_admin > viewer), wired to protect sensitive routes (user create/delete requires team_admin+)
- [x] **Password support** — New migration adds `password_hash` column, bcrypt hashing for seed users (default: "password"), `POST /users/{id}/password` endpoint
- [x] **Frontend auth** — Login pages for both admin and viewer apps, JWT stored in localStorage, `setAuthHeader()` sends Bearer token on all requests
- [x] **Route guards** — Protected route wrappers redirect to `/login` when unauthenticated, admin app checks for superadmin/team_admin role
- [x] **Dependencies** — `github.com/golang-jwt/jwt/v5` for JWT, `golang.org/x/crypto` (already indirect) for bcrypt

## Phase 7 — Soft Delete & Polish

- [x] **Soft delete migration** — `000008_add_deleted_at` adds `deleted_at TIMESTAMPTZ` to all 6 tables
- [x] **Soft delete queries** — All SELECTs filter `deleted_at IS NULL`, DELETEs become `UPDATE ... SET deleted_at = now()`
- [x] **Domain entities** — `DeletedAt *time.Time` on all 6 entities
- [x] **Cascade soft delete** — Room deletion cascades to seats, returns seat count (`seats_deleted`)
- [x] **HTTP responses** — DELETE handlers return `200 { "deleted": true, "seats_deleted": N }` instead of 204
- [x] **`ConfirmModal` component** — Replaces native `confirm()` dialogs in admin UI with styled modal
- [x] **Toast notification system** — Replaces `alert()` calls with `showToast()` (success/error/info, auto-dismiss)
- [x] **Login by username** — Migration `000009_make_email_optional` drops email NOT NULL, adds UNIQUE on name; `POST /auth/login` accepts `{"username", "password"}`
- [x] **Superadmin deletion protection** — `UserService.DeleteUser` returns `ErrCannotDeleteSuperAdmin`; delete button hidden in UI
- [x] **Weekly limit cap** — `MaxWeeklyLimit = 5` enforced in service layer; form inputs capped at `max="5"`
- [x] **Docker production build** — Multi-stage Dockerfiles for backend (Go), admin (pnpm → nginx), viewer (pnpm → nginx)
- [x] **CI pipeline** — GitHub Actions (`docker.yml`) builds & pushes all 3 images to `ghcr.io` on push to main

## Phase 8 — Deployment

- [x] **Health check endpoint** — `GET /health` returns `{"status":"ok"}` (unauthenticated)
- [x] **Auto-run migrations** — Backend runs `golang-migrate` on startup via `file://db/migrations`
- [x] **Production CORS** — `CORS_ORIGINS` env var with origin-validating middleware (falls back to `*` in dev)
- [x] **Production docker-compose** — `docker-compose.prod.yml` uses pre-built ghcr.io images, requires `DB_PASSWORD`, `JWT_SECRET`, `CORS_ORIGINS`
- [x] **CI pipeline** — `.github/workflows/ci.yml` runs backend tests (with Postgres service), frontend typecheck/lint/build, then Docker push on main
- [x] **README** — Architecture diagram, quick start, monorepo layout, deployment guide, API reference
- [x] **`.env.example`** — Template for required production env vars

## Phase 9 — Permanent Seat Assignment

Allow superadmins to permanently assign a seat to a user. The seat is always reserved to that user — no daily reservation needed, doesn't count toward weekly limit. One seat per user.

- [x] **Migration** — Add `assigned_user_id UUID REFERENCES users(id)` column to `seats` table with partial unique index (`WHERE assigned_user_id IS NOT NULL AND deleted_at IS NULL`)
- [x] **Domain** — Add `AssignedUserID *uuid.UUID` to `Seat`; add `ErrSeatPermanentlyAssigned`, `ErrUserAlreadyAssigned` errors
- [x] **sqlc** — Add `assigned_user_id` to all seat queries; add `GetSeatByAssignedUser`
- [x] **Seat repository** — Add `GetByAssignedUser(ctx, userID)` to interface + postgres implementation
- [x] **Seat service** — Add `AssignUser()` (superadmin only, one-per-user check) and `UnassignUser()` methods
- [x] **Reservation service** — Block all reservations on permanently assigned seats (`ErrSeatPermanentlyAssigned`)
- [x] **HTTP handlers** — Add `POST /seats/{id}/assign` and `DELETE /seats/{id}/assign` routes (superadmin only)
- [x] **DTOs** — Add `AssignedUserID`, `AssignedUserName` to `SeatResponse`; add `AssignSeatRequest`
- [x] **API client** — Add `useAssignSeat()`, `useUnassignSeat()` hooks
- [x] **Admin Room Editor** — Properties panel shows assigned user dropdown + unassign button (superadmin only)
- [x] **Viewer RoomDetail** — Show permanently assigned seats as locked with user name; no click action
- [x] **Viewer WeeklyCalendar** — Same treatment: permanently assigned cells show user name, locked

## Phase 10 — Viewer Calendar Enhancements & Cross-Team Management

Improve the viewer weekly calendar and cross-team requests page with today highlight, team admin reservation delegation, and role-based access.

- [x] **Calendar today highlight** — Add light blue background to today's column header and data cells in weekly calendar view
- [x] **Team admin reserve for others** — Show user dropdown in confirm modal for team admins to reserve seats on behalf of team members via `POST /reservations/admin`
- [x] **Hide Requests from viewers** — Conditionally show "Requests" nav item only for team_admin/superadmin in viewer Layout
- [x] **Team-filtered requests** — For team admins, filter cross-team requests list to only show requests where the seat belongs to their team
- [x] **Approve/reject for team admins** — Add approve/reject buttons to viewer CrossTeamRequests page for pending requests to the admin's team
