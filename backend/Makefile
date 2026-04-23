.PHONY: run build test test-integration migrate-up migrate-down swagger mock

# Run the application using air (hot reload)
dev:
	air

# Run the application normally
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run unit tests
test:
	go test -v ./pkg/... ./internal/...

# Run integration tests (using testcontainers)
test-integration:
	go test -v ./tests/...

# Database migrations
migrate-up:
	migrate -path internal/infrastructure/database/migrations/sql -database "$(DB_URL)" up

migrate-down:
	migrate -path internal/infrastructure/database/migrations/sql -database "$(DB_URL)" down

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

# Generate Mocks using Mockery
mock:
	@command -v mockery >/dev/null 2>&1 || (echo "mockery not found, installing..." && go install github.com/vektra/mockery/v2@latest)
	mockery --all --recursive --keeptree --output ./tests/mocks
