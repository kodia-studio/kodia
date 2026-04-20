# 🧪 Testing Guide

Comprehensive guide to testing in Kodia Framework with examples and best practices.

**Table of Contents:**
- [Testing Pyramid](#testing-pyramid)
- [Backend Testing](#backend-testing)
- [Frontend Testing](#frontend-testing)
- [Test Fixtures & Mocks](#test-fixtures--mocks)
- [Running Tests](#running-tests)
- [Coverage Reporting](#coverage-reporting)
- [Best Practices](#best-practices)
- [CI/CD Integration](#cicd-integration)

---

## Testing Pyramid

The testing pyramid shows the recommended distribution of tests:

```
        △
       /|\
      / | \
     /  |  \
    /   |   \
   / E2E|    \       1-3 tests per feature
  /_____|_____\      (slow, expensive, comprehensive)
  \           /
   \    API   /      3-5 tests per feature
    \  Tests /       (moderate speed & cost)
     \     /
      \   /
       \ /
   ______▼______
   |            |     10-20 tests per feature
   | Unit Tests |     (fast, cheap, specific)
   |____________|
```

**Distribution:**
- 70% Unit tests (fast, isolated)
- 20% Integration tests (database involved)
- 10% E2E tests (complete flows)

---

## Backend Testing

### Unit Tests

Test individual functions in isolation with mocks.

**Structure:**
```
backend/tests/unit/
├── services/
│   ├── auth_service_test.go
│   ├── user_service_test.go
│   └── post_service_test.go
├── handlers/
│   ├── auth_handler_test.go
│   └── user_handler_test.go
└── validators/
    └── auth_validator_test.go
```

**Example:**
```go
func TestAuthServiceRegister(t *testing.T) {
    // Arrange - Setup test data
    mockUserRepo := new(MockUserRepository)
    mockUserRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    authService := services.NewAuthService(mockUserRepo, mockLogger)
    user := &domain.User{Email: "test@example.com"}
    
    // Act - Execute function
    err := authService.RegisterUser(context.Background(), user)
    
    // Assert - Verify results
    assert.NoError(t, err)
    mockUserRepo.AssertCalled(t, "Save", mock.Anything, mock.Anything)
}
```

**Benefits:**
- ✅ Fast execution
- ✅ No external dependencies
- ✅ Easy to debug
- ✅ Cheap to run

---

### Integration Tests

Test with real database and dependencies.

**Structure:**
```
backend/tests/integration/
├── user_repository_test.go
├── post_repository_test.go
└── auth_flow_test.go
```

**Example:**
```go
func TestUserRepositorySave(t *testing.T) {
    // Setup real database
    testDB := tests.NewTestDatabase(t)
    defer testDB.Cleanup()
    
    repo := postgres.NewUserRepository(testDB.DB)
    ctx := context.Background()
    
    user := &domain.User{
        Email:        "test@example.com",
        PasswordHash: "hash",
        Role:         "user",
    }
    
    // Act - Save to real database
    err := repo.Save(ctx, user)
    
    // Assert
    require.NoError(t, err)
    
    // Verify in database
    retrieved, _ := repo.FindByEmail(ctx, "test@example.com")
    assert.Equal(t, "test@example.com", retrieved.Email)
}
```

**Benefits:**
- ✅ Tests real database behavior
- ✅ Catches ORM issues
- ✅ Validates data persistence
- ✅ More realistic scenarios

---

### E2E API Tests

Test complete HTTP request/response flows.

**Structure:**
```
backend/tests/e2e/
├── auth_e2e_test.go
├── users_e2e_test.go
└── posts_e2e_test.go
```

**Example:**
```go
func TestAuthLoginFlow(t *testing.T) {
    // Setup
    testDB := tests.NewTestDatabase(t)
    defer testDB.Cleanup()
    
    tests.CreateTestUser(t, testDB.DB, "user@example.com")
    router := setupTestRouter(t, testDB)
    
    // Create HTTP request
    loginReq := dto.LoginRequest{
        Email:    "user@example.com",
        Password: "password",
    }
    
    body, _ := json.Marshal(loginReq)
    req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    // Act - Execute request
    router.ServeHTTP(w, req)
    
    // Assert - Verify HTTP response
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response struct {
        AccessToken string `json:"access_token"`
    }
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotEmpty(t, response.AccessToken)
}
```

**Benefits:**
- ✅ Tests complete flows
- ✅ Validates HTTP contracts
- ✅ Catches integration issues
- ✅ Mimics real usage

---

## Frontend Testing

### Component Tests

Test individual Svelte components.

**Setup:**
```bash
npm install --save-dev vitest @testing-library/svelte @testing-library/user-event
```

**Example:**
```typescript
// src/lib/components/LoginForm.test.ts
import { render, screen } from '@testing-library/svelte'
import userEvent from '@testing-library/user-event'
import LoginForm from './LoginForm.svelte'

describe('LoginForm', () => {
    it('renders email and password inputs', () => {
        render(LoginForm)
        
        expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
        expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    })
    
    it('calls onSubmit when form is submitted', async () => {
        const user = userEvent.setup()
        const handleSubmit = vi.fn()
        
        render(LoginForm, { props: { onSubmit: handleSubmit } })
        
        await user.type(screen.getByLabelText(/email/i), 'test@example.com')
        await user.type(screen.getByLabelText(/password/i), 'password')
        await user.click(screen.getByText(/login/i))
        
        expect(handleSubmit).toHaveBeenCalled()
    })
})
```

---

### E2E Browser Tests

Test complete user flows in browser.

**Setup:**
```bash
npm install --save-dev @playwright/test
```

**Example:**
```typescript
// tests/e2e/auth.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
    test('user can register and login', async ({ page }) => {
        // Visit registration page
        await page.goto('/register')
        
        // Fill form
        await page.fill('input[name="email"]', 'newuser@example.com')
        await page.fill('input[name="password"]', 'SecurePassword123!')
        
        // Submit
        await page.click('button[type="submit"]')
        
        // Wait for redirect
        await expect(page).toHaveURL('/login')
        
        // Now login
        await page.fill('input[name="email"]', 'newuser@example.com')
        await page.fill('input[name="password"]', 'SecurePassword123!')
        await page.click('button[type="submit"]')
        
        // Verify logged in
        await expect(page).toHaveURL('/dashboard')
    })
})
```

---

## Test Fixtures & Mocks

### Using Test Helpers

```go
// Create test database with automatic cleanup
testDB := tests.NewTestDatabase(t)
defer testDB.Cleanup()

// Reset database between tests
testDB.Reset()

// Create test user
user := tests.CreateTestUser(t, testDB.DB, "test@example.com")

// Create test cache
testCache := tests.NewTestCache(t)
defer testCache.Cleanup()
```

### Mocking Repositories

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Usage
mockRepo := new(MockUserRepository)
mockRepo.On("FindByID", mock.Anything, "123").Return(user, nil)
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
```

---

## Running Tests

### Backend Tests

```bash
# Run all tests
make test

# Run specific test file
go test ./tests/unit/services/auth_service_test.go

# Run specific test function
go test -run TestAuthServiceRegister ./tests/unit/services/

# Run with verbose output
go test -v ./...

# Run short tests only (skip integration tests)
go test -short ./...

# Run tests in parallel
go test -parallel 4 ./...

# Run with timeout
go test -timeout 30s ./...

# Watch mode (requires air)
air
```

### Frontend Tests

```bash
# Run all tests
npm run test

# Run specific test file
npm run test -- LoginForm.test.ts

# Watch mode
npm run test:watch

# E2E tests
npm run test:e2e

# E2E with UI
npm run test:e2e -- --ui

# E2E headed (see browser)
npm run test:e2e -- --headed
```

---

## Coverage Reporting

### Backend Coverage

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
open coverage.html

# Check coverage percentage
go test -cover ./...

# Coverage for specific package
go test -cover ./internal/core/services/

# CI coverage (required >80%)
go test -cover ./... | grep coverage
```

### Frontend Coverage

```bash
# Generate coverage report
npm run test:coverage

# View coverage report
open coverage/index.html
```

---

## Best Practices

### ✅ DO:

```go
// 1. Use descriptive test names
func TestAuthServiceCreateUserWithValidEmail(t *testing.T) {
    // ✅ Clear what is being tested
}

// 2. Follow Arrange-Act-Assert pattern
func TestUserRepository(t *testing.T) {
    // Arrange - Setup
    repo := setupRepository()
    user := createTestUser()
    
    // Act - Execute
    err := repo.Save(user)
    
    // Assert - Verify
    assert.NoError(t, err)
}

// 3. Test behaviors, not implementation
func TestLoginSucceedsWithCorrectPassword(t *testing.T) {
    // What should happen: login succeeds
    // Not: how repository is called
}

// 4. Use table-driven tests for multiple scenarios
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name  string
        email string
        valid bool
    }{
        {"valid email", "test@example.com", true},
        {"missing @", "testexample.com", false},
        {"empty", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ValidateEmail(tt.email)
            assert.Equal(t, tt.valid, result)
        })
    }
}

// 5. Test error cases
func TestAuthServiceInvalidPassword(t *testing.T) {
    // Test what happens with wrong password
    err := authService.Login(ctx, email, "wrong_password")
    assert.Error(t, err)
}
```

### ❌ DON'T:

```go
// ❌ Don't use vague test names
func TestAuth(t *testing.T) {  // What is being tested?
}

// ❌ Don't have multiple assertions in one test
func TestUser(t *testing.T) {
    user := CreateUser()
    assert.NotNil(t, user)
    assert.Equal(t, user.Email, "test@example.com")
    assert.Equal(t, user.Role, "user")
    assert.Equal(t, user.Status, "active")
    // 4 assertions = 4 things to break
}

// ❌ Don't test implementation details
func TestRepositoryCalls(t *testing.T) {
    mockRepo.On("Query", mock.Anything).Return(user, nil)
    // ✅ Instead: Test behavior (user is retrieved correctly)
}

// ❌ Don't use sleeps for timing
func TestNotification(t *testing.T) {
    triggerNotification()
    time.Sleep(1 * time.Second)  // Flaky! Sometimes fails on slow systems
    assert.True(t, notificationSent)
}

// ❌ Don't leave test data in database
func TestUser(t *testing.T) {
    db.Create(&user)  // Cleanup?
    // ✅ Instead: Use testcontainers or defer cleanup
}
```

---

## Test Organization

### By Layer

```
tests/
├── unit/                    # Single component tests
│   ├── services/
│   ├── handlers/
│   └── validators/
├── integration/             # Database tests
│   ├── user_repository_test.go
│   ├── post_repository_test.go
│   └── cache_test.go
└── e2e/                     # Complete flows
    ├── auth_e2e_test.go
    ├── posts_e2e_test.go
    └── user_e2e_test.go
```

### By Feature

```
tests/
├── auth/
│   ├── unit_test.go
│   ├── integration_test.go
│   └── e2e_test.go
├── posts/
│   ├── unit_test.go
│   ├── integration_test.go
│   └── e2e_test.go
└── users/
    ├── unit_test.go
    ├── integration_test.go
    └── e2e_test.go
```

---

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
      redis:
        image: redis:7
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test -v ./...
      - run: go test -cover ./... | grep coverage
      - uses: codecov/codecov-action@v3

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm ci
      - run: npm run test
      - run: npm run test:e2e
```

---

## Coverage Goals

```
Minimum coverage targets:

Backend:    ✅ 80% overall
├── services:    90%
├── handlers:    85%
├── repository:  75%
└── domain:      100%

Frontend:   ✅ 70% overall
├── components:  75%
├── stores:      85%
└── utils:       80%
```

---

## Test Performance

### Optimization Tips

1. **Use short flag to skip slow tests**
   ```bash
   go test -short ./...  # Skip integration tests
   ```

2. **Run tests in parallel**
   ```bash
   go test -parallel 4 ./...
   ```

3. **Cache dependencies**
   ```bash
   go mod download  # Download once, reuse
   ```

4. **Use test containers efficiently**
   ```go
   // Reuse database for multiple tests
   testDB := tests.NewTestDatabase(t)
   defer testDB.Cleanup()
   
   // Reset between tests
   testDB.Reset()
   ```

5. **Profile slow tests**
   ```bash
   go test -cpuprofile=cpu.prof ./...
   go tool pprof cpu.prof
   ```

---

## Troubleshooting

### Test Fails Intermittently (Flaky)

```go
// ❌ Avoid timing-based assertions
time.Sleep(1 * time.Second)

// ✅ Use proper synchronization
select {
case <-done:
    // Test complete
case <-time.After(5*time.Second):
    t.Fatal("timeout")
}
```

### Database Tests Conflict

```go
// ✅ Isolate each test
func TestOne(t *testing.T) {
    db := tests.NewTestDatabase(t)
    defer db.Cleanup()
    // Each test gets fresh database
}
```

### Tests Too Slow

```bash
# Run only fast tests
go test -short ./...

# Profile to find bottlenecks
go test -cpuprofile=prof.out ./...
go tool pprof prof.out
```

---

## Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Testify Assertions](https://github.com/stretchr/testify)
- [Vitest Documentation](https://vitest.dev/)
- [Playwright Documentation](https://playwright.dev/)
- [Testcontainers](https://testcontainers.com/)

---

**Happy testing! 🧪**
