package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/dto"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	kodia_http "github.com/kodia-studio/kodia/internal/adapters/http"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	tests "github.com/kodia-studio/kodia/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthRegisterFlow tests complete registration flow
func TestAuthRegisterFlow(t *testing.T) {
	tests.SkipIfShort(t)

	// Setup
	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	router := setupTestRouter(t, testDB)

	// Create request
	registerReq := dto.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "SecurePassword123!",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "newuser@example.com", response.Email)
}

// TestAuthRegisterDuplicateEmail tests registering with duplicate email
func TestAuthRegisterDuplicateEmail(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	// Create existing user
	factory := tests.NewFactory(t, testDB.DB)
	factory.CreateUser(func(u *domain.User) {
		u.Email = "existing@example.com"
	})

	router := setupTestRouter(t, testDB)

	// Try to register with same email
	registerReq := dto.RegisterRequest{
		Email:    "existing@example.com",
		Password: "SecurePassword123!",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
}

// TestAuthLoginFlow tests complete login flow
func TestAuthLoginFlow(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	// Create test user
	factory := tests.NewFactory(t, testDB.DB)
	user := factory.CreateUser(func(u *domain.User) {
		u.Email = "login@example.com"
	})

	router := setupTestRouter(t, testDB)

	// Create login request
	loginReq := dto.LoginRequest{
		Email:    user.Email,
		Password: "password123", // The factory uses "password123" by default
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
}

// TestAuthLoginInvalidCredentials tests login with wrong password
func TestAuthLoginInvalidCredentials(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	factory := tests.NewFactory(t, testDB.DB)
	factory.CreateUser(func(u *domain.User) {
		u.Email = "user@example.com"
	})

	router := setupTestRouter(t, testDB)

	loginReq := dto.LoginRequest{
		Email:    "user@example.com",
		Password: "wrong_password",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthLoginUserNotFound tests login with non-existent email
func TestAuthLoginUserNotFound(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	router := setupTestRouter(t, testDB)

	loginReq := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthProtectedEndpoint tests accessing protected endpoint with token
func TestAuthProtectedEndpoint(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	factory := tests.NewFactory(t, testDB.DB)
	user := factory.CreateUser(func(u *domain.User) {
		u.Email = "protected@example.com"
	})
	
	jwtManager := jwt.NewManager(
		"test-secret-key-minimum-32-characters-long!!!",
		"test-refresh-secret-minimum-32-characters-long!",
		24,
		7,
	)

	// Generate token
	token, _ := jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role), []string{})

	router := setupTestRouter(t, testDB)

	// Access protected endpoint
	req := httptest.NewRequest("GET", "/api/users/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAuthProtectedEndpointNoToken tests accessing protected endpoint without token
func TestAuthProtectedEndpointNoToken(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	router := setupTestRouter(t, testDB)

	// Access protected endpoint without token
	req := httptest.NewRequest("GET", "/api/users/me", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthProtectedEndpointInvalidToken tests with invalid token
func TestAuthProtectedEndpointInvalidToken(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	router := setupTestRouter(t, testDB)

	// Access with invalid token
	req := httptest.NewRequest("GET", "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthRefreshToken tests token refresh flow
func TestAuthRefreshToken(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	factory := tests.NewFactory(t, testDB.DB)
	user := factory.CreateUser(func(u *domain.User) {
		u.Email = "refresh@example.com"
	})
	
	jwtManager := jwt.NewManager(
		"test-secret-key-minimum-32-characters-long!!!",
		"test-refresh-secret-minimum-32-characters-long!",
		24,
		7,
	)

	// Generate initial tokens
	refreshToken, _ := jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role), []string{})

	router := setupTestRouter(t, testDB)

	// Request new access token
	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: refreshToken,
	}

	body, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		AccessToken string `json:"access_token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
}

// BenchmarkAuthLogin benchmarks login performance
func BenchmarkAuthLogin(b *testing.B) {
	testDB := tests.NewTestDatabase(&testing.T{})
	defer testDB.Cleanup()

	factory := tests.NewFactory(nil, testDB.DB)
	factory.CreateUser(func(u *domain.User) {
		u.Email = "bench@example.com"
	})
	
	router := setupTestRouter(nil, testDB)

	loginReq := dto.LoginRequest{
		Email:    "bench@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// setupTestRouter creates a test router with all dependencies
func setupTestRouter(_ *testing.T, testDB *tests.TestDatabase) *gin.Engine {
	config := tests.NewTestConfig()
	jwtManager := jwt.NewManager(
		config.JWT.AccessSecret,
		config.JWT.RefreshSecret,
		config.JWT.AccessExpiryHours,
		config.JWT.RefreshExpiryDays,
	)
	logger := tests.NewTestLogger()

	// Create repositories
	userRepo := postgres.NewUserRepository(testDB.DB)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(testDB.DB)

	// Create services
	mockCache := new(MockCacheProvider)
	mockMailer := new(MockMailer)
	authService := services.NewAuthService(
		userRepo,
		refreshTokenRepo,
		jwtManager,
		mockCache,
		mockMailer,
		"http://localhost:8080",
		"http://localhost:3000",
		logger,
	)
	userService := services.NewUserService(userRepo, logger)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, nil, logger)
	userHandler := handlers.NewUserHandler(userService, nil, logger)
	graphqlHandler := handlers.NewGraphQLHandler(authService, userService, logger)
	healthHandler := handlers.NewHealthHandler(testDB.DB, nil, logger)
	pulseHandler := handlers.NewPulseHandler(nil, logger) // Mock or nil for E2E if not needed

	// Create router
	// want (*Config, *Logger, *jwt.Manager, *AuthHandler, *UserHandler, *redis.Client, *websocket.Handler, *GraphQLHandler, *observability.Manager, *PulseHandler, *HealthHandler)
	router := kodia_http.NewRouter(config, logger, jwtManager, authHandler, userHandler, nil, nil, graphqlHandler, nil, pulseHandler, healthHandler)
	return router.Setup()
}

// MockCacheProvider is a simple mock for E2E tests
type MockCacheProvider struct{}
func (m *MockCacheProvider) Get(ctx context.Context, key string, dest interface{}) error { return nil }
func (m *MockCacheProvider) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error { return nil }
func (m *MockCacheProvider) Delete(ctx context.Context, key string) error { return nil }
func (m *MockCacheProvider) Remember(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error { return nil }

// MockMailer is a simple mock for E2E tests
type MockMailer struct{}
func (m *MockMailer) Send(ctx context.Context, to []string, subject string, body string) error { return nil }
func (m *MockMailer) SendHTML(ctx context.Context, to []string, subject string, htmlBody string) error { return nil }
func (m *MockMailer) SendWithTemplate(ctx context.Context, to []string, subject string, templatePath string, data interface{}) error { return nil }
