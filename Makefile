.PHONY: help dev build test clean docker-up docker-down migrate lint format

# Default target
help:
	@echo "Kodia Framework - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev              - Start all dev servers"
	@echo "  make build            - Build all components"
	@echo ""
	@echo "Testing:"
	@echo "  make test             - Run all tests"
	@echo "  make test-unit        - Run unit tests only"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make test-e2e         - Run E2E API tests only"
	@echo "  make test-short       - Run fast tests (skip slow tests)"
	@echo "  make test-watch       - Run tests in watch mode"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make test-frontend    - Run frontend component tests"
	@echo "  make test-e2e-ui      - Run E2E tests with UI"
	@echo ""
	@echo "Database:"
	@echo "  make migrate          - Run database migrations"
	@echo "  make db-reset         - Reset database to fresh state"
	@echo "  make db-seed          - Run database seeders"
	@echo "  make db-shell         - Open database shell"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint             - Run linter"
	@echo "  make format           - Format code"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up        - Start Docker services"
	@echo "  make docker-down      - Stop Docker services"
	@echo ""
	@echo "Other:"
	@echo "  make clean            - Clean build artifacts"

# Development
dev:
	@echo "Starting Kodia Development Servers..."
	@make -C backend dev & make -C frontend dev

# Build all components
build:
	@echo "Building Kodia framework components..."
	@make -C backend build
	@make -C cli build

# ============================================================================
# Testing Targets
# ============================================================================

# Run all tests
test:
	@echo "Running all tests..."
	@make -C backend test
	@make -C frontend test

# Unit tests only (fast)
test-unit:
	@echo "Running unit tests..."
	@cd backend && go test -v -short ./tests/unit/...

# Integration tests
test-integration:
	@echo "Running integration tests..."
	@cd backend && go test -v ./tests/integration/...

# E2E API tests
test-e2e:
	@echo "Running E2E API tests..."
	@cd backend && go test -v ./tests/e2e/...

# Fast tests (no integration/E2E)
test-short:
	@echo "Running fast tests..."
	@cd backend && go test -v -short ./...

# Watch mode for backend tests
test-watch:
	@echo "Starting test watch mode..."
	@cd backend && air --cmd.tmp_dir tmp --cmd.args='test -v ./...'

# Coverage report
test-coverage:
	@echo "Generating coverage report..."
	@cd backend && go test -v -coverprofile=coverage.out ./...
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: backend/coverage.html"
	@open backend/coverage.html || echo "Please open backend/coverage.html manually"

# Parallel tests (4 workers)
test-parallel:
	@echo "Running tests in parallel..."
	@cd backend && go test -v -parallel 4 ./...

# Frontend component tests
test-frontend:
	@echo "Running frontend component tests..."
	@cd frontend && npm run test

# Frontend tests with coverage
test-frontend-coverage:
	@echo "Running frontend tests with coverage..."
	@cd frontend && npm run test:coverage

# Frontend E2E tests
test-e2e-browser:
	@echo "Running browser E2E tests..."
	@cd frontend && npm run test:e2e

# Frontend E2E with UI
test-e2e-ui:
	@echo "Running E2E tests with UI..."
	@cd frontend && npm run test:e2e -- --ui

# ============================================================================
# Database Targets
# ============================================================================

migrate:
	@echo "Running database migrations..."
	@cd backend && go run cmd/server/main.go migrate

db-reset:
	@echo "Resetting database to fresh state..."
	@docker-compose exec postgres dropdb -U postgres kodia 2>/dev/null || true
	@docker-compose exec postgres createdb -U postgres kodia
	@make migrate
	@echo "Database reset complete"

db-seed:
	@echo "Running database seeders..."
	@cd backend && go run cmd/server/main.go seed

db-shell:
	@echo "Opening database shell..."
	@docker-compose exec postgres psql -U postgres -d kodia

# ============================================================================
# Documentation Targets
# ============================================================================

docs:
	@echo "Generating OpenAPI/Swagger documentation..."
	@command -v swag >/dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@cd backend && swag init -g cmd/server/main.go -o docs
	@echo "✅ OpenAPI documentation generated!"
	@echo "📖 Access Swagger UI at: http://localhost:8080/api/docs"
	@echo "📥 Download spec at: http://localhost:8080/api/docs/swagger.json"

docs-clean:
	@echo "Cleaning generated documentation..."
	@rm -rf backend/docs/swagger.json backend/docs/swagger.yaml backend/docs/docs.go
	@echo "Documentation cleaned"

docs-view:
	@echo "Opening Swagger UI..."
	@open http://localhost:8080/api/docs || echo "Please visit http://localhost:8080/api/docs"

# ============================================================================
# Code Quality Targets
# ============================================================================

lint:
	@echo "Running linter..."
	@cd backend && golangci-lint run
	@cd frontend && npm run lint

lint-fix:
	@echo "Fixing lint issues..."
	@cd backend && golangci-lint run --fix
	@cd frontend && npm run lint:fix

format:
	@echo "Formatting code..."
	@cd backend && go fmt ./...
	@cd frontend && npm run format

# Security scanning
security:
	@echo "Running security checks..."
	@cd backend && gosec ./...

# ============================================================================
# Docker Targets
# ============================================================================

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services ready! PostgreSQL: localhost:5432, Redis: localhost:6379"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs:
	@echo "Following Docker logs..."
	docker-compose logs -f

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

# ============================================================================
# Build Targets
# ============================================================================

backend-build:
	@echo "Building backend binary..."
	@cd backend && go build -o server cmd/server/main.go

backend-build-prod:
	@echo "Building production backend binary..."
	@cd backend && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

frontend-build:
	@echo "Building frontend..."
	@cd frontend && npm run build

cli-build:
	@echo "Building CLI..."
	@cd cli && go build -o kodia kodia/main.go

# ============================================================================
# Utility Targets
# ============================================================================

# Install dependencies
install:
	@echo "Installing dependencies..."
	@cd backend && go mod download
	@cd frontend && npm install
	@cd cli && go mod download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf backend/server backend/coverage.html backend/coverage.out
	@rm -rf backend/tmp cli/kodia frontend/.svelte-kit frontend/build
	@echo "Clean complete"

# Full cleanup (including Docker)
clean-all: clean docker-clean

# Run all checks (lint, test, build)
check: lint test build
	@echo "All checks passed!"

# Development setup
setup:
	@echo "Setting up development environment..."
	@make install
	@make docker-up
	@make migrate
	@echo "Setup complete! Run 'make dev' to start development servers"

# ============================================================================
# Docker Image Building
# ============================================================================

docker-build-backend:
	@echo "Building backend Docker image..."
	docker build -t kodia-backend:latest -f backend/Dockerfile ./backend

docker-build-frontend:
	@echo "Building frontend Docker image..."
	docker build -t kodia-frontend:latest -f frontend/Dockerfile ./frontend

docker-push-backend: docker-build-backend
	@echo "Pushing backend image..."
	docker push kodia-backend:latest

docker-push-frontend: docker-build-frontend
	@echo "Pushing frontend image..."
	docker push kodia-frontend:latest
