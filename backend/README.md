# neshiman/backend

Seat plan application — Go backend.

## Prerequisites

- Go 1.22+
- PostgreSQL 15+
- sqlc (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- golang-migrate (`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)
- Docker (for local Postgres)

## Quick start

```bash
make setup    # from repo root: start postgres, migrate, generate code, install deps
make dev      # start postgres + backend + both frontends
```

## Commands

| Command | Description |
|---------|-------------|
| `go run ./cmd/server` | Start dev server on :8080 |
| `go test ./...` | Run all tests |
| `sqlc generate` | Regenerate `db/sqlc/` from `db/query/` |
| `migrate -path db/migrations -database "$DB_URL" up` | Apply migrations |
| `migrate -path db/migrations -database "$DB_URL" down` | Rollback migrations |
| `go mod tidy` | Sync Go module dependencies |
