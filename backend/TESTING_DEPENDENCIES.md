# Testing Dependencies for Go

Add these dependencies to your `go.mod` for comprehensive testing support:

## Installation

```bash
cd backend

# Testing & Mocking
go get -u github.com/stretchr/testify
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/require
go get -u github.com/stretchr/testify/mock

# Test Containers (PostgreSQL, Redis)
go get -u github.com/testcontainers/testcontainers-go
go get -u github.com/testcontainers/testcontainers-go/wait

# Coverage
go get -u golang.org/x/tools/cmd/cover
```

## Dependencies

```go
// Required test dependencies
require (
    github.com/stretchr/testify v1.8.4      // assertions & mocking
    github.com/testcontainers/testcontainers-go v0.26.0  // database testing
    github.com/testcontainers/testcontainers-go/wait v0.26.0
)
```

## Frontend Testing Dependencies

Add to `frontend/package.json`:

```json
{
  "devDependencies": {
    "vitest": "^1.0.0",
    "@testing-library/svelte": "^4.0.0",
    "@testing-library/user-event": "^14.0.0",
    "@playwright/test": "^1.40.0"
  }
}
```

## Installation Commands

```bash
# Go testing
cd backend && go get -u github.com/stretchr/testify
cd backend && go get -u github.com/testcontainers/testcontainers-go

# Frontend testing
cd frontend && npm install --save-dev vitest @testing-library/svelte @playwright/test
```

## Verification

```bash
# Backend
cd backend && go test ./tests/unit/...

# Frontend
cd frontend && npm run test
```

All dependencies will be available in your test code!
