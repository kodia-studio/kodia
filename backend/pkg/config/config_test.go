package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabaseConfigDSN_Postgres tests DSN generation for PostgreSQL
func TestDatabaseConfigDSN_Postgres(t *testing.T) {
	cfg := DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}

	dsn := cfg.DSN()

	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=testuser")
	assert.Contains(t, dsn, "password=testpass")
	assert.Contains(t, dsn, "dbname=testdb")
	assert.Contains(t, dsn, "sslmode=disable")
}

// TestDatabaseConfigDSN_MySQL tests DSN generation for MySQL
func TestDatabaseConfigDSN_MySQL(t *testing.T) {
	cfg := DatabaseConfig{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "mysqluser",
		Password: "mysqlpass",
		Name:     "mysqldb",
	}

	dsn := cfg.DSN()

	assert.Contains(t, dsn, "mysqluser")
	assert.Contains(t, dsn, "mysqlpass")
	assert.Contains(t, dsn, "localhost:3306")
	assert.Contains(t, dsn, "mysqldb")
	assert.Contains(t, dsn, "charset=utf8mb4")
}

// TestDatabaseConfigDSN_SQLite tests DSN generation for SQLite
func TestDatabaseConfigDSN_SQLite(t *testing.T) {
	cfg := DatabaseConfig{
		Driver: "sqlite",
		Name:   "test.db",
	}

	dsn := cfg.DSN()

	assert.Equal(t, "test.db", dsn)
}

// TestDatabaseConfigDSN_DefaultsToPostgres tests that unknown driver defaults to Postgres
func TestDatabaseConfigDSN_DefaultsToPostgres(t *testing.T) {
	cfg := DatabaseConfig{
		Driver:   "unknown-db",
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		Name:     "db",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}

	dsn := cfg.DSN()

	assert.Contains(t, dsn, "host=localhost")
}

// TestRedisConfigAddr tests Redis address generation
func TestRedisConfigAddr(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{"Standard Redis", "localhost", 6379, "localhost:6379"},
		{"Remote Redis", "redis.example.com", 6379, "redis.example.com:6379"},
		{"Custom port", "localhost", 9999, "localhost:9999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := RedisConfig{
				Host: tt.host,
				Port: tt.port,
			}

			assert.Equal(t, tt.expected, cfg.Addr())
		})
	}
}

// TestAppConfigDefaults tests app config values
func TestAppConfigDefaults(t *testing.T) {
	cfg := AppConfig{
		Name:        "test-app",
		Env:         "test",
		Port:        8080,
		Debug:       true,
		BaseURL:     "http://localhost:8080",
		FrontendURL: "http://localhost:3000",
		Locale:      "en",
	}

	assert.Equal(t, "test-app", cfg.Name)
	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, 8080, cfg.Port)
	assert.True(t, cfg.Debug)
}

// TestJWTConfigValues tests JWT config
func TestJWTConfigValues(t *testing.T) {
	cfg := JWTConfig{
		AccessSecret:      "access-secret-minimum-32-characters-long!!!",
		RefreshSecret:     "refresh-secret-minimum-32-characters-long!",
		AccessExpiryHours: 24,
		RefreshExpiryDays: 7,
	}

	assert.NotEmpty(t, cfg.AccessSecret)
	assert.NotEmpty(t, cfg.RefreshSecret)
	assert.True(t, len(cfg.AccessSecret) >= 32)
	assert.True(t, len(cfg.RefreshSecret) >= 32)
	assert.Equal(t, 24, cfg.AccessExpiryHours)
	assert.Equal(t, 7, cfg.RefreshExpiryDays)
}

// TestCORSConfigAllowedOrigins tests CORS configuration
func TestCORSConfigAllowedOrigins(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://example.com",
		},
	}

	assert.Len(t, cfg.AllowedOrigins, 2)
	assert.Contains(t, cfg.AllowedOrigins, "http://localhost:3000")
	assert.Contains(t, cfg.AllowedOrigins, "https://example.com")
}

// TestStorageConfigLocal tests local storage config
func TestStorageConfigLocal(t *testing.T) {
	cfg := StorageConfig{
		Provider: "local",
		LocalDir: "/tmp/uploads",
		PublicURL: "http://localhost:8080/uploads",
	}

	assert.Equal(t, "local", cfg.Provider)
	assert.Equal(t, "/tmp/uploads", cfg.LocalDir)
	assert.Equal(t, "http://localhost:8080/uploads", cfg.PublicURL)
}

// TestStorageConfigS3 tests S3 storage config
func TestStorageConfigS3(t *testing.T) {
	cfg := StorageConfig{
		Provider:  "s3",
		Bucket:    "my-bucket",
		Region:    "us-east-1",
		AccessID:  "AKIAIOSFODNN7EXAMPLE",
		SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		PublicURL: "https://s3.amazonaws.com/my-bucket",
	}

	assert.Equal(t, "s3", cfg.Provider)
	assert.Equal(t, "my-bucket", cfg.Bucket)
	assert.Equal(t, "us-east-1", cfg.Region)
}

// TestMailConfigSMTP tests mail config
func TestMailConfigSMTP(t *testing.T) {
	cfg := MailConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		User:     "user@example.com",
		Password: "password",
		FromAddr: "noreply@example.com",
		FromName: "My App",
	}

	assert.Equal(t, "smtp.gmail.com", cfg.Host)
	assert.Equal(t, 587, cfg.Port)
	assert.Equal(t, "noreply@example.com", cfg.FromAddr)
	assert.Equal(t, "My App", cfg.FromName)
}

// TestObservabilityConfigSentry tests Sentry config
func TestObservabilityConfigSentry(t *testing.T) {
	cfg := ObservabilityConfig{
		SentryDSN:      "https://key@sentry.io/project-id",
		TracingEnabled: true,
		MetricsEnabled: true,
		SamplingRate:   1.0,
		ServiceName:    "kodia-api",
	}

	assert.NotEmpty(t, cfg.SentryDSN)
	assert.True(t, cfg.TracingEnabled)
	assert.True(t, cfg.MetricsEnabled)
	assert.Equal(t, 1.0, cfg.SamplingRate)
	assert.Equal(t, "kodia-api", cfg.ServiceName)
}

// TestNotificationConfigChannels tests notification channels config
func TestNotificationConfigChannels(t *testing.T) {
	cfg := NotificationConfig{
		TwilioAccountSID: "AC123456789",
		TwilioAuthToken:  "token123",
		TwilioFromNumber: "+1234567890",
		SlackWebhookURL:  "https://hooks.slack.com/services/ABC/123/xyz",
		FCMServerKey:     "key123",
	}

	assert.NotEmpty(t, cfg.TwilioAccountSID)
	assert.NotEmpty(t, cfg.SlackWebhookURL)
	assert.NotEmpty(t, cfg.FCMServerKey)
}

// TestDatabaseConfigVariations tests various database configurations
func TestDatabaseConfigVariations(t *testing.T) {
	testCases := []struct {
		name     string
		config   DatabaseConfig
		contains []string
	}{
		{
			name: "Postgres with SSL",
			config: DatabaseConfig{
				Driver:   "postgres",
				Host:     "db.example.com",
				Port:     5432,
				User:     "produser",
				Password: "prodpass",
				Name:     "proddb",
				SSLMode:  "require",
				TimeZone: "UTC",
			},
			contains: []string{"sslmode=require", "host=db.example.com"},
		},
		{
			name: "MySQL with charset",
			config: DatabaseConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "root",
				Name:     "kodia",
			},
			contains: []string{"charset=utf8mb4", "parseTime=True"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dsn := tc.config.DSN()
			for _, expected := range tc.contains {
				assert.Contains(t, dsn, expected)
			}
		})
	}
}

// TestLoad tests configuration loading (mocked environment)
func TestLoad(t *testing.T) {
	// Set minimal required env vars
	os.Setenv("APP_NAME", "test-kodia")
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "8080")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_PORT")
	}()

	cfg, err := Load()

	// Should load without error
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

// TestLoadWithDefaults tests that Load sets reasonable defaults
func TestLoadWithDefaults(t *testing.T) {
	cfg, err := Load()

	require.NoError(t, err)
	assert.NotNil(t, cfg.App)
	assert.NotNil(t, cfg.Database)
	assert.NotNil(t, cfg.JWT)
	assert.NotNil(t, cfg.CORS)
}

// BenchmarkDatabaseDSN benchmarks DSN generation
func BenchmarkDatabaseDSN(b *testing.B) {
	cfg := DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		Name:     "db",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg.DSN()
	}
}

// BenchmarkRedisAddr benchmarks Redis address generation
func BenchmarkRedisAddr(b *testing.B) {
	cfg := RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg.Addr()
	}
}
