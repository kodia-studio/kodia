# Testing Infrastructure

Kodia provides a robust testing infrastructure to ensure your application remains stable and reliable. We follow the "Build like a user, code like a pro" philosophy by making integration testing as easy as unit testing.

---

## 1. HTTP Test Server (`NewTestServer`)

Kodia allows you to spin up a fully-booted application instance in milliseconds for integration testing.

```go
func TestMyFeature(t *testing.T) {
    app := kodia.NewApp(cfg, log)
    // Setup dependencies...
    
    ts := app.NewTestServer(t) // Automatically boots and cleans up
    
    resp := tests.JSONRequest(t, http.DefaultClient, "GET", ts.URL+"/api/my-feature", nil)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

---

## 2. Database Test Helper & Reset

When running integration tests, Kodia ensures a clean state by using `TRUNCATE` (much faster than migrations) before each test.

### Automatic Reset
In your test sub-tests, call `td.Reset()` to clear all tables and restart identity sequences.

```go
t.Run("it should create a user", func(t *testing.T) {
    td.Reset()
    // Test logic here...
})
```

---

## 3. Data Factories

Instead of manually creating database records, use the `Factory` pattern to generate domain entities with sensible defaults.

```go
factory := tests.NewFactory(t, td.DB)

// Create a regular user
user := factory.CreateUser()

// Create an admin with custom email
admin := factory.CreateAdmin(func(u *domain.User) {
    u.Email = "custom-admin@kodia.id"
})
```

---

## 4. Contract Testing

Ensure your API responses don't break over time by validating their structure against a contract.

```go
validation.ValidateContract(t, responseBody, map[string]interface{}{
    "id":    "number",
    "name":  "string",
    "email": "string",
    "roles": "array",
})
```

---

## 5. Mock Generation

Kodia uses **Mockery** to automatically generate mocks from your interfaces.

### Generate Mocks
Run the following command to update all mocks in `tests/mocks`:
```bash
make mock
```

### Usage
```go
mockRepo := new(mocks.UserRepository)
mockRepo.On("FindByID", mock.Anything, "123").Return(user, nil)

service := NewUserService(mockRepo)
```

---

## 6. Running Tests

- **Unit Tests**: `make test`
- **Integration Tests (Requires Docker)**: `make test-integration`
- **Short Mode**: `go test -short ./...` (Skips slow integration tests)
