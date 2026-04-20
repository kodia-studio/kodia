package integration

import (
	"context"
	"testing"
	"github.com/kodia-studio/kodia/pkg/pagination"

	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/domain"
	tests "github.com/kodia-studio/kodia/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserRepositorySave tests creating a new user
func TestUserRepositorySave(t *testing.T) {
	tests.SkipIfShort(t)

	// Setup
	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	user := &domain.User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     "user",
	}

	// Act
	err := repo.Create(ctx, user)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// Verify in database
	retrieved, err := repo.FindByEmail(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", retrieved.Email)
	assert.Equal(t, "user", retrieved.Role)
}

// TestUserRepositoryFindByID tests finding user by ID
func TestUserRepositoryFindByID(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test user
	user := tests.CreateTestUser(t, testDB.DB, "findme@example.com")

	// Act
	found, err := repo.FindByID(ctx, user.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, "findme@example.com", found.Email)
}

// TestUserRepositoryFindByIDNotFound tests finding non-existent user
func TestUserRepositoryFindByIDNotFound(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Act
	found, err := repo.FindByID(ctx, "nonexistent-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRepositoryFindByEmail tests finding user by email
func TestUserRepositoryFindByEmail(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test user
	tests.CreateTestUser(t, testDB.DB, "email@example.com")

	// Act
	user, err := repo.FindByEmail(ctx, "email@example.com")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "email@example.com", user.Email)
}

// TestUserRepositoryFindAll tests finding all users with pagination
func TestUserRepositoryFindAll(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Create multiple test users
	tests.CreateTestUser(t, testDB.DB, "user1@example.com")
	tests.CreateTestUser(t, testDB.DB, "user2@example.com")
	tests.CreateTestUser(t, testDB.DB, "user3@example.com")

	// Act
	users, _, err := repo.FindAll(ctx, &pagination.Params{Page: 1, PerPage: 10})

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 3)
}

// TestUserRepositoryUpdate tests updating user
func TestUserRepositoryUpdate(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test user
	user := tests.CreateTestUser(t, testDB.DB, "update@example.com")

	// Act - Update user
	user.Role = "admin"
	err := repo.Update(ctx, user)

	// Assert
	require.NoError(t, err)

	// Verify update
	updated, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.UserRole("admin"), updated.Role)
}

// TestUserRepositoryDelete tests deleting user
func TestUserRepositoryDelete(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test user
	user := tests.CreateTestUser(t, testDB.DB, "delete@example.com")

	// Act
	err := repo.Delete(ctx, user.ID)

	// Assert
	require.NoError(t, err)

	// Verify deletion
	found, err := repo.FindByID(ctx, user.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRepositoryDuplicateEmail tests that duplicate emails are rejected
func TestUserRepositoryDuplicateEmail(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Create first user
	user1 := &domain.User{
		Email:    "duplicate@example.com",
		Password: "hash1",
		Role:     "user",
	}
	err := repo.Create(ctx, user1)
	require.NoError(t, err)

	// Try to create user with same email
	user2 := &domain.User{
		Email:    "duplicate@example.com",
		Password: "hash2",
		Role:     "user",
	}
	err = repo.Create(ctx, user2)

	// Assert - Should error on duplicate
	assert.Error(t, err)
}

// TestUserRepositoryConcurrentWrites tests concurrent write operations
func TestUserRepositoryConcurrentWrites(t *testing.T) {
	tests.SkipIfShort(t)

	testDB := tests.NewTestDatabase(t)
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Create users concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			user := &domain.User{
				Email:    "concurrent" + string(rune(index)) + "@example.com",
				Password: "hash",
				Role:     "user",
			}
			repo.Create(ctx, user)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all users were created
	users, _, err := repo.FindAll(ctx, &pagination.Params{Page: 1, PerPage: 100})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 10)
}

// BenchmarkUserRepositorySave benchmarks user creation
func BenchmarkUserRepositorySave(b *testing.B) {
	testDB := tests.NewTestDatabase(&testing.T{})
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &domain.User{
			Email:    "bench" + string(rune(i)) + "@example.com",
			Password: "hash",
			Role:     "user",
		}
		repo.Create(ctx, user)
	}
}

// BenchmarkUserRepositoryFindByID benchmarks finding user by ID
func BenchmarkUserRepositoryFindByID(b *testing.B) {
	testDB := tests.NewTestDatabase(&testing.T{})
	defer testDB.Cleanup()

	repo := postgres.NewUserRepository(testDB.DB)
	ctx := context.Background()

	user := tests.CreateTestUser(nil, testDB.DB, "benchmark@example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindByID(ctx, user.ID)
	}
}
