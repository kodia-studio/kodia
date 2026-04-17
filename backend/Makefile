.PHONY: run build test migrate-up migrate-down swagger

# Run the application using air (hot reload)
dev:
	air

# Run the application normally
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Database migrations
migrate-up:
	migrate -path internal/infrastructure/database/migrations/sql -database "$(DB_URL)" up

migrate-down:
	migrate -path internal/infrastructure/database/migrations/sql -database "$(DB_URL)" down

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs
