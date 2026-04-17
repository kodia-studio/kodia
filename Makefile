.PHONY: dev build test clean docker-up docker-down

# Development
dev:
	@echo "Starting Kodia Development Servers..."
	@make -C backend dev & make -C frontend dev

# Build all components
build:
	@echo "Building Kodia framework components..."
	@make -C backend build
	@make -C kodia-cli build

# Test all components
test:
	@make -C backend test

# Docker management
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Cleanup
clean:
	@rm -rf backend/bin kodia-cli/bin frontend/.svelte-kit
