package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDatabase manages test database lifecycle
type TestDatabase struct {
	DB        *gorm.DB
	Container testcontainers.Container
	t         *testing.T
}

// TestCache manages test Redis cache
type TestCache struct {
	Client    *redis.Client
	Container testcontainers.Container
	t         *testing.T
}

// NewTestDatabase creates a test PostgreSQL database using testcontainers
func NewTestDatabase(t *testing.T) *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "kodia_test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Get connection string
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), "testuser", "testpass", "kodia_test",
	)

	// Connect to database
	db, err := gorm.Open(gorm_postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Run migrations
	runMigrations(t, db)

	return &TestDatabase{
		DB:        db,
		Container: container,
		t:         t,
	}
}

// Cleanup stops the database container
func (td *TestDatabase) Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := td.Container.Terminate(ctx); err != nil {
		td.t.Logf("failed to stop database container: %v", err)
	}
}

// Reset clears all data from database using TRUNCATE for speed
func (td *TestDatabase) Reset() {
	// List tables to truncate
	tables := []string{"users", "refresh_tokens", "audit_logs"}
	for _, table := range tables {
		if td.DB.Migrator().HasTable(table) {
			if err := td.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
				td.t.Logf("failed to truncate table %s: %v", table, err)
			}
		}
	}
}

// NewTestCache creates a test Redis cache using testcontainers
func NewTestCache(t *testing.T) *TestCache {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create Redis container
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	// Get connection string
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "6379")

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("failed to connect to redis: %v", err)
	}

	return &TestCache{
		Client:    client,
		Container: container,
		t:         t,
	}
}

// Cleanup stops the Redis container
func (tc *TestCache) Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tc.Client.Close()

	if err := tc.Container.Terminate(ctx); err != nil {
		tc.t.Logf("failed to stop redis container: %v", err)
	}
}

// Flush clears all data from Redis
func (tc *TestCache) Flush() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tc.Client.FlushDB(ctx).Err(); err != nil {
		tc.t.Logf("failed to flush redis: %v", err)
	}
}

// NewTestLogger creates a test logger
func NewTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

// NewTestConfig creates a test configuration
func NewTestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Name:        "kodia-test",
			Env:         "test",
			Port:        8080,
			Debug:       true,
		},
		JWT: config.JWTConfig{
			AccessSecret:       "test-secret-key-minimum-32-characters-long!!!",
			RefreshSecret:      "test-refresh-secret-minimum-32-characters-long!",
			AccessExpiryHours:  24,
			RefreshExpiryDays:  7,
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:5173"},
		},
	}
}

// runMigrations runs all database migrations
func runMigrations(t *testing.T, db *gorm.DB) {
	// Auto migrate all domain models
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
		&auditEntry{}, // Internal for audit logging tests
	); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
}

type auditEntry struct {
	ID        uint      `gorm:"primaryKey"`
	Timestamp time.Time `gorm:"index"`
	ActorID   string    `gorm:"index"`
	Actor     string
	Action    string `gorm:"index"`
	Resource  string `gorm:"index"`
	Details   string `gorm:"type:text"`
	Status    string `gorm:"index"`
	IP        string
	UserAgent string
}

func (auditEntry) TableName() string {
	return "audit_logs"
}

// --- Request Helpers ---

// JSONRequest makes an HTTP request with JSON body and returns the response
func JSONRequest(t *testing.T, client *http.Client, method, url string, body interface{}) *http.Response {
	t.Helper()

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	return resp
}

// ParseJSON parses a JSON response body into a struct
func ParseJSON(t *testing.T, resp *http.Response, dest interface{}) {
	t.Helper()
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
}

// SkipIfShort skips test if running with -short flag
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
}

// SkipCI skips test if running in CI environment
func SkipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping test in CI environment")
	}
}
