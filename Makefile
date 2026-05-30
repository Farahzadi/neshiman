DB_URL ?= postgres://postgres:postgres@localhost:5432/neshiman?sslmode=disable
TEST_DB_URL ?= postgres://postgres:postgres@localhost:5432/neshiman_test?sslmode=disable

.PHONY: setup deps db-up db-wait db-migrate db-codegen swagger-gen types-gen dev build test lint clean seed test-db db-test deploy-prod deploy-prod-down deploy-prod-logs deploy-local deploy-local-seed

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

seed:
	cd backend && go run ./cmd/seed/

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

db-test-setup:
	docker compose exec -T postgres psql -U postgres -c "CREATE DATABASE neshiman_test" 2>/dev/null || true
	migrate -path backend/db/migrations -database "$(TEST_DB_URL)" up

test: check-go db-test-setup
	cd backend && DB_URL="$(TEST_DB_URL)" go test ./... -count=1
	pnpm turbo test

lint: check-go
	cd backend && go vet ./...
	pnpm turbo lint

deploy-prod:
	@echo "Starting production stack..."
	docker compose -f docker-compose.prod.yml up -d

deploy-prod-down:
	docker compose -f docker-compose.prod.yml down

deploy-prod-logs:
	docker compose -f docker-compose.prod.yml logs -f

# --- Deploy to local network machine via SSH ---
# Override with: make deploy-local DEPLOY_HOST=user@ip DEPLOY_DIR=/path
DEPLOY_HOST ?= ashpaz@192.168.104.27
DEPLOY_DIR  ?= ~/neshiman
REPO_OWNER  ?= neshiman

BACKEND_IMG = ghcr.io/$(REPO_OWNER)/neshiman-backend:latest
ADMIN_IMG   = ghcr.io/$(REPO_OWNER)/neshiman-admin:latest
VIEWER_IMG  = ghcr.io/$(REPO_OWNER)/neshiman-viewer:latest

deploy-local: check-docker
	@test -f .env || { echo "Error: .env file not found. Copy .env.example to .env and fill in your secrets."; exit 1; }
	@echo "=== Building images locally (linux/amd64) ==="
	docker build --platform linux/amd64 -f backend/Dockerfile -t $(BACKEND_IMG) backend/
	docker build --platform linux/amd64 -f apps/admin/Dockerfile -t $(ADMIN_IMG) .
	docker build --platform linux/amd64 -f apps/viewer/Dockerfile -t $(VIEWER_IMG) .
	@echo "=== Saving images to tar ==="
	docker save $(BACKEND_IMG) $(ADMIN_IMG) $(VIEWER_IMG) | gzip > /tmp/neshiman-images.tar.gz
	@echo "=== Transferring to $(DEPLOY_HOST):$(DEPLOY_DIR) ==="
	ssh $(DEPLOY_HOST) "mkdir -p $(DEPLOY_DIR)"
	scp /tmp/neshiman-images.tar.gz $(DEPLOY_HOST):$(DEPLOY_DIR)/
	scp docker-compose.prod.yml $(DEPLOY_HOST):$(DEPLOY_DIR)/
	scp .env $(DEPLOY_HOST):$(DEPLOY_DIR)/
	@echo "=== Loading and starting on remote ==="
	ssh $(DEPLOY_HOST) "cd $(DEPLOY_DIR) && docker load < neshiman-images.tar.gz && docker compose -f docker-compose.prod.yml up -d && rm neshiman-images.tar.gz"
	@rm -f /tmp/neshiman-images.tar.gz
	@echo "=== Deployed to $(DEPLOY_HOST):$(DEPLOY_DIR) ==="
	@echo "⚠  Run 'make deploy-local-seed' to seed the database with initial data."

deploy-local-seed:
	@echo "=== Seeding database on $(DEPLOY_HOST) ==="
	@echo "⚠  This will DESTROY all existing data and recreate from scratch!"
	ssh $(DEPLOY_HOST) "cd $(DEPLOY_DIR) && docker compose -f docker-compose.prod.yml exec backend /seed"
	@echo "=== Done. Default password for all users is: password ==="

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
