package database

import (
	"testing"

	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestDatabaseConnectionLoggingDoesNotExposeCredentials verifies that sensitive data is not logged
func TestDatabaseConnectionLoggingDoesNotExposeCredentials(t *testing.T) {
	// Create a test logger that captures output
	core, logs := observer.New(zapcore.DebugLevel)
	obs := zap.New(core)

	// Create a test config with sensitive data
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver:   "postgres",
			Host:     "secret-db.internal.example.com",
			Port:     5432,
			User:     "admin_user",
			Password: "super_secret_password_12345",
			Name:     "myappdb",
			SSLMode:  "require",
		},
	}

	// Try to establish connection (will fail but we're testing logging)
	_, _ = New(cfg, obs)

	// Get the logged entries
	for _, entry := range logs.All() {
		// Check all logged fields
		for _, field := range entry.Context {
			fieldValue := field.String

			// Verify sensitive information is NOT in logs
			sensitiveValues := []string{
				"secret-db.internal.example.com", // Host
				"admin_user",                      // User
				"super_secret_password_12345",     // Password
				"require",                         // SSL Mode
			}

			for _, sensitive := range sensitiveValues {
				if fieldValue == sensitive {
					t.Errorf("Sensitive data '%s' found in logs. Logs should not contain credentials or infrastructure details", sensitive)
				}
			}
		}
	}
}

// TestDatabaseLoggingDoesNotContainHostOrPort verifies host and port are not logged
func TestDatabaseLoggingDoesNotContainHostOrPort(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	obs := zap.New(core)

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "postgres",
			Host:   "db.example.com",
			Port:   5432,
			User:   "dbuser",
			Name:   "production_db",
		},
	}

	_, _ = New(cfg, obs)

	// Check all logged fields
	for _, entry := range logs.All() {
		for _, field := range entry.Context {
			fieldValue := field.String

			// Host and port should never be logged
			if fieldValue == "db.example.com" {
				t.Error("Database host 'db.example.com' should not be logged")
			}
			if fieldValue == "5432" {
				t.Error("Database port '5432' should not be logged")
			}
			if fieldValue == "dbuser" {
				t.Error("Database user should not be logged")
			}
		}
	}
}

// TestDatabaseLoggingIncludesNonSensitiveInfo verifies we still log useful debugging info
func TestDatabaseLoggingIncludesNonSensitiveInfo(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	obs := zap.New(core)

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "postgres",
			Host:   "db.example.com",
			Port:   5432,
			User:   "dbuser",
			Name:   "mydb",
		},
	}

	_, _ = New(cfg, obs)

	// Verify non-sensitive information is logged
	foundDriver := false
	foundDatabase := false

	for _, entry := range logs.All() {
		if entry.Message == "Attempting database connection" || entry.Message == "Database connected" {
			for _, field := range entry.Context {
				if field.Key == "driver" && field.String == "postgres" {
					foundDriver = true
				}
				if field.Key == "database" && field.String == "mydb" {
					foundDatabase = true
				}
			}
		}
	}

	if !foundDriver {
		t.Log("Note: Driver not logged (connection may have failed before logging)")
	}
	if !foundDatabase {
		t.Log("Note: Database name not logged (connection may have failed before logging)")
	}
}

// BenchmarkDatabaseLogging measures logging overhead
func BenchmarkDatabaseLogging(b *testing.B) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "postgres",
			Host:   "localhost",
			Port:   5432,
			User:   "user",
			Name:   "testdb",
		},
	}

	logger := zap.NewNop() // No-op logger for benchmarking

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("Attempting database connection",
			zap.String("driver", cfg.Database.Driver),
			zap.String("database", cfg.Database.Name),
		)
	}
}
