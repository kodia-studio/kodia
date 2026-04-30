.PHONY: run build test test-integration migrate-status migrate-rollback migrate-fresh migrate-reset migration-create swagger mock

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
migrate-status:
	go run cmd/migrate/main.go status

migrate-rollback:
	go run cmd/migrate/main.go down

migrate-fresh:
	go run cmd/migrate/main.go fresh

migrate-reset:
	go run cmd/migrate/main.go reset

# Create a new migration file with timestamp
# Usage: make migration-create NAME=create_posts_table
migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migration-create NAME=create_xxx_table"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	STRUCT_NAME="Migration_$${TIMESTAMP}"; \
	TABLE_NAME=$$(echo "$(NAME)" | sed 's/create_//g' | sed 's/_table//g'); \
	FILE="internal/infrastructure/database/migrations/go/$${TIMESTAMP}_$$(echo $(NAME) | tr '[:upper:]' '[:lower:]').go"; \
	printf 'package migrations\n\nimport "github.com/kodia-studio/kodia/pkg/database"\n\ntype %s struct{}\n\nfunc (m *%s) Up(schema *database.Schema) error {\n\treturn schema.Create("%s", func(table *database.Blueprint) {\n\t\ttable.ID()\n\t\ttable.Timestamps()\n\t})\n}\n\nfunc (m *%s) Down(schema *database.Schema) error {\n\treturn schema.Drop("%s")\n}\n' \
	"$${STRUCT_NAME}" "$${STRUCT_NAME}" "$${TABLE_NAME}" "$${STRUCT_NAME}" "$${TABLE_NAME}" > "$${FILE}"; \
	echo "✓ Created: $${FILE}"; \
	echo "→ Register it in: internal/infrastructure/database/migrations/go/registry.go"

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

# Generate Mocks using Mockery
mock:
	@command -v mockery >/dev/null 2>&1 || (echo "mockery not found, installing..." && go install github.com/vektra/mockery/v2@latest)
	mockery --all --recursive --keeptree --output ./tests/mocks
