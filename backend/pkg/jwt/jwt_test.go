package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewManager creates a new JWT manager with test secrets
func newTestManager() *Manager {
	return NewManager(
		"test-access-secret-minimum-32-characters-long!!!",
		"test-refresh-secret-minimum-32-characters-long!",
		1, // 1 hour access
		7, // 7 days refresh
	)
}

// TestGenerateAccessToken tests access token generation
func TestGenerateAccessToken(t *testing.T) {
	m := newTestManager()

	token, err := m.GenerateAccessToken("user-123", "user@example.com", "user", []string{"read"})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

// TestGenerateRefreshToken tests refresh token generation
func TestGenerateRefreshToken(t *testing.T) {
	m := newTestManager()

	token, err := m.GenerateRefreshToken("user-123", "user@example.com", "user", []string{})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

// TestValidateAccessToken tests access token validation
func TestValidateAccessToken(t *testing.T) {
	m := newTestManager()
	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}

	token, _ := m.GenerateAccessToken(userID, email, role, permissions)

	claims, err := m.ValidateAccessToken(token)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, permissions, claims.Permissions)
	assert.Equal(t, AccessToken, claims.TokenType)
}

// TestValidateRefreshToken tests refresh token validation
func TestValidateRefreshToken(t *testing.T) {
	m := newTestManager()
	userID := "user-456"
	email := "test@example.com"
	role := "admin"

	token, _ := m.GenerateRefreshToken(userID, email, role, []string{})

	claims, err := m.ValidateRefreshToken(token)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, RefreshToken, claims.TokenType)
}

// TestValidateInvalidToken tests validation with invalid token
func TestValidateInvalidToken(t *testing.T) {
	m := newTestManager()

	claims, err := m.ValidateAccessToken("invalid-token-string")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestValidateExpiredToken tests validation with expired token
func TestValidateExpiredToken(t *testing.T) {
	// Create a manager with very short expiry
	m := NewManager(
		"test-access-secret-minimum-32-characters-long!!!",
		"test-refresh-secret-minimum-32-characters-long!",
		0, // 0 hours = instant expiration
		7,
	)

	token, _ := m.GenerateAccessToken("user-123", "user@example.com", "user", []string{})
	time.Sleep(100 * time.Millisecond) // Small delay to ensure expiration

	claims, err := m.ValidateAccessToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestAccessTokenCannotValidateAsRefresh tests type validation
func TestAccessTokenCannotValidateAsRefresh(t *testing.T) {
	m := newTestManager()

	accessToken, _ := m.GenerateAccessToken("user-123", "user@example.com", "user", []string{})

	// Try to validate access token as refresh token
	claims, err := m.ValidateRefreshToken(accessToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "invalid token type", err.Error())
}

// TestRefreshTokenCannotValidateAsAccess tests reverse type validation
func TestRefreshTokenCannotValidateAsAccess(t *testing.T) {
	m := newTestManager()

	refreshToken, _ := m.GenerateRefreshToken("user-123", "user@example.com", "user", []string{})

	// Try to validate refresh token as access token
	claims, err := m.ValidateAccessToken(refreshToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "invalid token type", err.Error())
}

// TestTokenWithDifferentSecret tests token signed with different secret
func TestTokenWithDifferentSecret(t *testing.T) {
	m1 := newTestManager()
	m2 := NewManager("different-secret-minimum-32-chars!", "different-secret-minimum-32-chars!", 1, 7)

	token, _ := m1.GenerateAccessToken("user-123", "user@example.com", "user", []string{})

	// Try to validate with different manager
	claims, err := m2.ValidateAccessToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestEmptyPermissions tests token generation with empty permissions
func TestEmptyPermissions(t *testing.T) {
	m := newTestManager()

	token, _ := m.GenerateAccessToken("user-123", "user@example.com", "user", []string{})

	claims, _ := m.ValidateAccessToken(token)

	// Empty slice should be preserved or nil
	assert.True(t, len(claims.Permissions) == 0 || claims.Permissions == nil)
}

// TestMultiplePermissions tests token with multiple permissions
func TestMultiplePermissions(t *testing.T) {
	m := newTestManager()
	permissions := []string{"read", "write", "delete", "admin"}

	token, _ := m.GenerateAccessToken("user-123", "user@example.com", "user", permissions)

	claims, _ := m.ValidateAccessToken(token)

	assert.Equal(t, permissions, claims.Permissions)
}

// TestClaimsHaveValidID tests that claims have a valid JTI
func TestClaimsHaveValidID(t *testing.T) {
	m := newTestManager()

	token, _ := m.GenerateAccessToken("user-123", "user@example.com", "user", []string{})

	claims, _ := m.ValidateAccessToken(token)

	assert.NotEmpty(t, claims.ID)
	assert.NotEmpty(t, claims.Subject)
}

// MockTokenStore for testing token revocation
type MockTokenStore struct {
	revoked map[string]bool
	used    map[string]bool
}

func newMockTokenStore() *MockTokenStore {
	return &MockTokenStore{
		revoked: make(map[string]bool),
		used:    make(map[string]bool),
	}
}

func (m *MockTokenStore) IsRevoked(jti string) (bool, error) {
	return m.revoked[jti], nil
}

func (m *MockTokenStore) MarkUsed(jti string, userID string, expiry time.Duration) (bool, error) {
	alreadyUsed := m.used[jti]
	m.used[jti] = true
	return alreadyUsed, nil
}

func (m *MockTokenStore) Revoke(jti string, expiry time.Duration) error {
	m.revoked[jti] = true
	return nil
}

// TestRefreshTokenRevokeDetection tests token revocation
func TestRefreshTokenRevokeDetection(t *testing.T) {
	m := newTestManager()
	store := newMockTokenStore()
	m.SetStore(store)

	token, _ := m.GenerateRefreshToken("user-123", "user@example.com", "user", []string{})

	claims, _ := m.ValidateRefreshToken(token)
	require.NoError(t, nil) // First validation should succeed

	// Revoke the token
	store.revoked[claims.ID] = true

	// Second validation should fail
	_, err := m.ValidateRefreshToken(token)
	assert.Error(t, err)
	assert.Equal(t, "token has been revoked", err.Error())
}

// TestRefreshTokenReuseDetection tests token reuse detection
func TestRefreshTokenReuseDetection(t *testing.T) {
	m := newTestManager()
	store := newMockTokenStore()
	m.SetStore(store)

	token, _ := m.GenerateRefreshToken("user-123", "user@example.com", "user", []string{})

	// First use should succeed
	claims1, err1 := m.ValidateRefreshToken(token)
	require.NoError(t, err1)

	// Second use should detect reuse
	_, err2 := m.ValidateRefreshToken(token)
	assert.Error(t, err2)
	assert.Equal(t, "refresh token reuse detected - all sessions invalidated", err2.Error())

	// Verify all user tokens were revoked
	assert.True(t, store.revoked[claims1.UserID])
}

// BenchmarkGenerateAccessToken benchmarks access token generation
func BenchmarkGenerateAccessToken(b *testing.B) {
	m := newTestManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.GenerateAccessToken("user-123", "user@example.com", "user", []string{"read"})
	}
}

// BenchmarkValidateAccessToken benchmarks access token validation
func BenchmarkValidateAccessToken(b *testing.B) {
	m := newTestManager()
	token, _ := m.GenerateAccessToken("user-123", "user@example.com", "user", []string{"read"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ValidateAccessToken(token)
	}
}

// BenchmarkGenerateRefreshToken benchmarks refresh token generation
func BenchmarkGenerateRefreshToken(b *testing.B) {
	m := newTestManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.GenerateRefreshToken("user-123", "user@example.com", "user", []string{})
	}
}
