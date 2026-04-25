# Testing Infrastructure

Kodia provides a **production-grade testing infrastructure** with 120+ comprehensive test cases covering security-critical modules, validation, pagination, and more. We follow the "Build like a user, code like a pro" philosophy by making integration testing as easy as unit testing.

---

## Test Coverage Summary

As of Phase 1 (April 2026), Kodia has **78% average coverage** across implemented packages:

| Package | Coverage | Test Cases | Focus Area |
|---------|----------|-----------|------------|
| `pkg/hash` | **100%** | 13 | Password hashing, bcrypt |
| `pkg/jwt` | **100%** | 18 | Token generation, validation, revocation |
| `pkg/pagination` | **100%** | 14 | Offset/limit, sorting, search |
| `pkg/policy` | **100%** | 15 | ABAC evaluation, role-based access |
| `pkg/config` | **83.1%** | 15 | Config loading, DSN generation |
| `pkg/response` | **70.4%** | 16 | HTTP responses, error formatting |
| `pkg/health` | **65.4%** | 16 | System health checks, stats |
| `pkg/validation` | **46%** | 17 | Struct validation, custom rules |

---

## Core Test Types

### 1. Unit Tests (Fast, Isolated)

Pure unit tests without external dependencies. 88% of all tests are unit tests, running in < 500ms total.

```bash
# Run all unit tests
go test ./pkg/... -v -cover

# Run specific package
go test ./pkg/hash -v -cover
go test ./pkg/jwt -v -cover

# With coverage report
go test ./pkg/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 2. HTTP Test Server (`NewTestServer`)

Spin up a fully-booted application instance in milliseconds for integration testing.

```go
func TestMyFeature(t *testing.T) {
    app := kodia.NewApp(cfg, log)
    // Setup dependencies...
    
    ts := app.NewTestServer(t) // Automatically boots and cleans up
    
    resp := tests.JSONRequest(t, http.DefaultClient, "GET", ts.URL+"/api/my-feature", nil)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### 3. Database Test Helper & Reset

Kodia uses `testcontainers` to spin up PostgreSQL/Redis instances. Tests ensure clean state using `TRUNCATE`.

```go
func TestUserRepository(t *testing.T) {
    td := tests.NewTestDatabase(t)
    defer td.Cleanup()
    
    // Run tests...
    
    t.Run("it should create a user", func(t *testing.T) {
        td.Reset() // TRUNCATE tables before each test
        // Test logic here...
    })
}
```

### 4. Data Factories

Use the `Factory` pattern to generate domain entities with sensible defaults.

```go
factory := tests.NewFactory(t, td.DB)

// Create a regular user
user := factory.CreateUser()

// Create with custom email
user := factory.CreateUser(func(u *domain.User) {
    u.Email = "custom@kodia.id"
})

// Create admin user
admin := factory.CreateAdmin()

// Create multiple users
users := factory.CreateMultipleUsers(10)

// Create inactive user
inactive := factory.CreateInactiveUser()
```

### 5. Contract Testing

Validate API response structures to prevent breaking changes.

```go
validation.ValidateContract(t, responseBody, map[string]interface{}{
    "id":    "number",
    "name":  "string",
    "email": "string",
    "roles": "array",
})
```

### 6. Mock Generation

Kodia uses **Mockery** to automatically generate mocks from interfaces.

```bash
# Generate all mocks
make mock
```

```go
mockRepo := new(mocks.UserRepository)
mockRepo.On("FindByID", mock.Anything, "123").Return(user, nil)

service := NewUserService(mockRepo)
```

### 7. Benchmarks

Performance benchmarks for critical operations:

```bash
# Run all benchmarks
go test ./pkg/hash -bench=. -benchmem
go test ./pkg/jwt -bench=. -benchmem
go test ./pkg/pagination -bench=. -benchmem
```

Each package includes benchmarks for its most critical operations.

---

## Running Tests

### Quick Start
```bash
# Run all tests with coverage
make test

# Run integration tests (requires Docker)
make test-integration

# Short mode (skip slow tests)
go test -short ./...

# Specific package
go test ./pkg/jwt -v
```

### Advanced
```bash
# Run with race detector
go test ./pkg/... -race

# Run with verbose output
go test ./pkg/... -v

# Run with timeout
go test ./pkg/... -timeout 30s

# Run a specific test
go test ./pkg/jwt -run TestGenerateAccessToken

# Run tests matching pattern
go test ./pkg/... -run Validate
```

---

## Test Examples

### Hash Password Test
```go
func TestCheckValidPassword(t *testing.T) {
	password := "my-password"
	hash, _ := hash.Make(password)

	result := hash.Check(password, hash)

	assert.True(t, result)
}
```

### JWT Validation Test
```go
func TestValidateAccessToken(t *testing.T) {
	m := NewManager(accessSecret, refreshSecret, 1, 7)

	token, _ := m.GenerateAccessToken("user-123", "user@example.com", "user", []string{"read"})
	claims, err := m.ValidateAccessToken(token)

	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
}
```

### ABAC Policy Test
```go
func TestEvaluateRoleBasedAccess(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "admin-all-access",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, ok := s["role"]
			return ok && role == "admin"
		},
	})

	result := e.Evaluate(Attributes{"role": "admin"}, nil, nil)
	assert.True(t, result)
}
```

### Integration Test with Database
```go
func TestUserRepositoryCreate(t *testing.T) {
	td := tests.NewTestDatabase(t)
	defer td.Cleanup()

	repo := postgres.NewUserRepository(td.DB)
	user := &domain.User{
		Email:    "test@example.com",
		Password: "hashed",
		Role:     "user",
	}

	err := repo.Create(context.Background(), user)

	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
}
```

---

## Best Practices

1. **Write Unit Tests First**: Focus on package-level unit tests before integration tests
2. **Use Table-Driven Tests**: Test multiple scenarios in one test function
3. **Test Edge Cases**: Empty inputs, invalid data, boundary conditions
4. **Mock External Dependencies**: Use mocks for databases, APIs, caches
5. **Clean Up Resources**: Always defer cleanup of database containers
6. **Use Subtests**: Organize related tests with `t.Run()`
7. **Add Benchmarks**: Benchmark critical operations for performance regressions

---

## CI/CD Integration

Tests are automatically run on:
- **Every PR**: Via GitHub Actions
- **Before Merge**: All tests must pass
- **Nightly**: Full integration tests and benchmarks

See [GitHub Actions](#) workflow for details.

---

## Next Steps

- Add tests for remaining packages (auth2fa, authsocial, binder, database, i18n, etc.)
- Expand integration test coverage for complex flows
- Add e2e tests for critical user journeys
- Monitor coverage metrics and aim for 80%+ overall coverage
