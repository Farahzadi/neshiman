DB_URL ?= postgres://postgres:postgres@localhost:5432/neshiman?sslmode=disable

.PHONY: setup deps db-up db-wait db-migrate db-codegen swagger-gen types-gen dev build test lint clean

check-docker:
	$(call check_tool,docker)

check-go:
	$(call check_tool,go)

check-pnpm:
	$(call check_tool,pnpm)

define check_tool
	@if ! command -v $(1) >/dev/null 2>&1; then \
		echo "Error: $(1) is not installed."; \
		case "$(1)" in \
			migrate) echo "  Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest" ;; \
			sqlc)    echo "  Install: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest" ;; \
			pnpm)   echo "  Install: enable via corepack, or brew install pnpm" ;; \
			*)       echo "  Install using your package manager." ;; \
		esac; \
		exit 1; \
	fi
endef

setup: check-docker check-go db-up db-wait db-migrate db-codegen deps
	@echo "Setup complete."

deps:
	$(call check_tool,go)
	pnpm install

db-up:
	$(call check_tool,docker)
	docker compose up -d postgres

db-wait:
	@echo "Waiting for postgres..."
	@for i in $$(seq 1 30); do \
		if docker compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then \
			echo "Postgres is ready."; \
			exit 0; \
		fi; \
		printf "."; \
		sleep 1; \
	done; \
	echo ""; \
	echo "Timed out waiting for postgres."; \
	exit 1

db-migrate:
	$(call check_tool,migrate)
	migrate -path backend/db/migrations -database "$(DB_URL)" up

swagger-gen:
	cd backend && swag init -g ./cmd/server/main.go --output ./docs

types-gen: swagger-gen
	pnpm generate:types

db-migrate-down:
	$(call check_tool,migrate)
	migrate -path backend/db/migrations -database "$(DB_URL)" down $(filter-out $@,$(MAKECMDGOALS))

db-codegen:
	$(call check_tool,sqlc)
	cd backend && sqlc generate

dev: check-docker check-go db-up db-wait db-migrate
	@trap 'kill 0 2>/dev/null; exit' INT TERM; \
	echo "Starting backend..."; \
	cd backend && go run ./cmd/server & \
	echo "Starting frontend apps via Turborepo..."; \
	pnpm turbo dev & \
	echo ""; \
	echo "--- Services ---"; \
	echo "  Backend:       http://localhost:8080"; \
	echo "  Admin UI:      http://localhost:3000"; \
	echo "  Viewer UI:     http://localhost:3001"; \
	echo "Press Ctrl+C to stop all services."; \
	echo "-----------------"; \
	wait

build: check-go
	cd backend && go build -o server ./cmd/server
	pnpm turbo build

test: check-go
	cd backend && go test ./...
	pnpm turbo test

lint: check-go
	cd backend && go vet ./...
	pnpm turbo lint

clean:
	rm -f backend/server
	rm -rf apps/admin/dist
	rm -rf apps/viewer/dist
	rm -rf packages/*/dist
	rm -rf packages/api-types/src/v1.d.ts
	rm -rf node_modules
	rm -rf apps/admin/node_modules
	rm -rf apps/viewer/node_modules
	rm -rf packages/*/node_modules
