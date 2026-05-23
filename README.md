# Neshiman

Seat plan management system — admins define rooms with grid-based seats, users reserve seats per day with a configurable weekly limit.

## Architecture

```
┌─────────────┐  ┌─────────────┐
│  Admin UI   │  │ Viewer UI   │
│  :3000      │  │ :3001       │
│  SolidJS    │  │ SolidJS     │
└──────┬──────┘  └──────┬──────┘
       │  /api/*        │ /api/*
       └──────┬─────────┘
              │
       ┌──────▼──────┐
       │  Backend    │
       │  :8080      │
       │  Go + chi   │
       └──────┬──────┘
              │
       ┌──────▼──────┐
       │  Postgres   │
       │  :5432      │
       └─────────────┘
```

### Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go, chi router, pgx/v5, sqlc, golang-migrate |
| Database | PostgreSQL 15 |
| Admin UI | SolidJS, Vite, TailwindCSS, TanStack Query |
| Viewer UI | SolidJS, Vite, TailwindCSS, TanStack Query |
| API docs | Swagger 2.0 (swaggo/swag) |
| Auth | JWT (HS256, 24h expiry), bcrypt |
| Container | Docker, GitHub Container Registry |

## Quick Start (Development)

### Prerequisites

- Go 1.24+
- Docker (for Postgres)
- pnpm 9.15+ (enable via `corepack enable`)
- golang-migrate (`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)

### Setup

```bash
make setup    # Start postgres, apply migrations, sqlc generate, install frontend deps
make dev      # Start postgres + backend + both frontends
```

This starts:
- **Backend** — http://localhost:8080 (Swagger UI at /swagger)
- **Admin UI** — http://localhost:3000
- **Viewer UI** — http://localhost:3001

### Seed Data

```bash
make seed
```

Creates sample rooms, teams, users, seats, and reservations. All seeded users have password `password`.

### Tests

```bash
make test     # Backend + frontend tests
```

## Monorepo Layout

```
neshiman/
├── backend/              # Go backend (hexagonal architecture)
│   ├── cmd/server/       # Entrypoint — wires everything
│   ├── cmd/seed/         # Database seeder
│   ├── db/
│   │   ├── migrations/   # SQL migration files
│   │   ├── query/        # sqlc query definitions
│   │   └── sqlc/         # Generated code (DO NOT EDIT)
│   ├── internal/
│   │   ├── domain/       # Core domain entities + value objects
│   │   ├── ports/        # Repository interfaces
│   │   ├── application/  # Service layer
│   │   └── adapters/
│   │       ├── http/     # Handlers, middleware, router
│   │       └── postgres/ # Repository implementations
│   └── docs/             # Swagger spec (generated)
├── apps/
│   ├── admin/            # Admin management UI (SolidJS)
│   └── viewer/           # Seat reservation UI (SolidJS)
├── packages/
│   ├── api-client/       # TanStack Query hooks + fetch wrapper
│   ├── api-types/        # Generated TypeScript types from Swagger
│   ├── tsconfig/         # Shared TypeScript configs
│   └── tailwind-config/  # Shared TailwindCSS preset
├── docker-compose.yml    # Development compose
├── docker-compose.prod.yml # Production compose
└── .github/workflows/    # CI/CD
```

## Deployment

### Production

1. Copy `.env.example` to `.env` and fill in secrets:
   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

2. Start the stack:
   ```bash
   make deploy-prod
   ```

This uses pre-built images from `ghcr.io/<repo>/neshiman-*`. Images are built automatically on push to `main` via GitHub Actions.

### Required environment variables

| Variable | Description |
|----------|-------------|
| `DB_PASSWORD` | Postgres password |
| `JWT_SECRET` | Secret for JWT signing (generate a long random string) |
| `CORS_ORIGINS` | Comma-separated allowed origins (e.g. `https://admin.example.com,https://viewer.example.com`) |

## API

Browse the interactive Swagger documentation at `/swagger/index.html` when the backend is running.

All authenticated endpoints require a JWT via the `Authorization: Bearer <token>` header. Obtain a token via `POST /api/v1/auth/login` with `{"username": "...", "password": "..."}`.

### Role Hierarchy

```
superadmin > team_admin > viewer
```

- `superadmin`: Full access to everything
- `team_admin`: Manage their team's seats, users, and reservations
- `viewer`: Reserve seats, view their own reservations
