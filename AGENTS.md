# AGENTS.md

## Project

Seat plan application. Admins define rooms with rectangular seats arranged on a grid. Seats have labels, can rotate/move, and are owned by teams. Users (viewers) reserve seats per day with a weekly limit (e.g., 2 days/week). Team admins configure reservation capacity for their members. Role hierarchy: superadmin > team admin > viewer. Cross-team seat requests require approval from the owning team's admin.

## Monorepo Layout

- `backend/` — Go service, hexagonal architecture, PostgreSQL, sqlc
- `frontend-admin/` — SolidJS app (admin management UI)
- `frontend-viewer/` — SolidJS app (seat reservation UI)
- `docker-compose.yml` — production orchestration

## Key Decisions

- **Backend**: Hexagonal architecture — keep domain/ports free of infrastructure deps. Adapters (HTTP handlers, DB repos) live outside the core.
- **sqlc**: Generates Go structs from SQL queries. Never edit generated files; modify `.sql` sources and re-run `sqlc generate`.
- **Frontend**: SolidJS ≠ React. Use signals (`createSignal`), stores (`createStore`), no virtual DOM diffing. JSX compiles to real DOM.
- **Code quality**: ESLint + Prettier enforced via git hooks. Run lint/format before committing.

## Agent Notes

- Seat geometry is rectangular (not triangular). Grid-based room layout.
- Team memberships are non-overlapping. A user belongs to exactly one team.
- Seat reservations are scoped to a single day, one person per seat.
- Users have a weekly reservation limit. Team admins can adjust this per member.
- Cross-team seat requests require approval from the owning team's admin.
