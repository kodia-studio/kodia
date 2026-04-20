package services

import (
	"context"
	"testing"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/hash"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserRepository is a mock implementation of the user repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, params *pagination.Params) ([]*domain.User, int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// MockRefreshTokenRepository is a mock implementation
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestAuthServiceRegister tests user registration
func TestAuthServiceRegister(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshRepo := new(MockRefreshTokenRepository)
	logger := zap.NewNop()
	jwtManager := jwt.NewManager("access-secret-32-chars-long-at-least", "refresh-secret-32-chars-long-at-least", 1, 7)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo, jwtManager, logger)

	input := ports.RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "test@example.com" && u.Name == "Test User"
	})).Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()

	// Act
	resp, err := authService.Register(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "test@example.com", resp.User.Email)
	mockUserRepo.AssertExpectations(t)
}

// TestAuthServiceLogin tests user login
func TestAuthServiceLogin(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshRepo := new(MockRefreshTokenRepository)
	logger := zap.NewNop()
	jwtManager := jwt.NewManager("access-secret-32-chars-long-at-least", "refresh-secret-32-chars-long-at-least", 1, 7)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo, jwtManager, logger)

	password := "password123"
	hashedPassword, _ := hash.Make(password)

	user := &domain.User{
		ID:       "user-123",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     domain.RoleUser,
		IsActive: true,
	}

	input := ports.LoginInput{
		Email:    "test@example.com",
		Password: password,
	}

	mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()

	// Act
	resp, err := authService.Login(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	mockUserRepo.AssertExpectations(t)
}

// TestAuthServiceLogout tests user logout
func TestAuthServiceLogout(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshRepo := new(MockRefreshTokenRepository)
	logger := zap.NewNop()
	jwtManager := jwt.NewManager("access-secret-32-chars-long-at-least", "refresh-secret-32-chars-long-at-least", 1, 7)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo, jwtManager, logger)

	token := "some-refresh-token"
	mockRefreshRepo.On("RevokeByToken", mock.Anything, token).Return(nil)

	ctx := context.Background()

	// Act
	err := authService.Logout(ctx, token)

	// Assert
	assert.NoError(t, err)
	mockRefreshRepo.AssertExpectations(t)
}
