// Package config provides application configuration loading for Kodia Framework.
// Supports loading from environment variables and config files via Viper.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App           AppConfig           `mapstructure:"app"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Redis         RedisConfig         `mapstructure:"redis"`
	JWT           JWTConfig           `mapstructure:"jwt"`
	CORS          CORSConfig          `mapstructure:"cors"`
	Storage       StorageConfig       `mapstructure:"storage"`
	Mail          MailConfig          `mapstructure:"mail"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Notification  NotificationConfig  `mapstructure:"notification"`
}

// AppConfig holds general application settings.
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Env         string `mapstructure:"env"`
	Port        int    `mapstructure:"port"`
	Debug       bool   `mapstructure:"debug"`
	BaseURL     string `mapstructure:"base_url"`
	FrontendURL string `mapstructure:"frontend_url"`
	Locale      string `mapstructure:"locale"`
}

// DatabaseConfig holds database connection settings.
// Driver can be "postgres" or "mysql".
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Name            string        `mapstructure:"name"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	TimeZone        string        `mapstructure:"timezone"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DSN returns the database connection string based on the configured driver.
func (d DatabaseConfig) DSN() string {
	switch strings.ToLower(d.Driver) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.User, d.Password, d.Host, d.Port, d.Name)
	case "sqlite", "sqlite3":
		return d.Name
	default: // postgres
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.TimeZone)
	}
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Addr returns the Redis address in host:port format.
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// JWTConfig holds JWT signing settings.
type JWTConfig struct {
	AccessSecret      string `mapstructure:"access_secret"`
	RefreshSecret     string `mapstructure:"refresh_secret"`
	AccessExpiryHours int    `mapstructure:"access_expiry_hours"`
	RefreshExpiryDays int    `mapstructure:"refresh_expiry_days"`
}

// CORSConfig holds CORS allowed origins.
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// StorageConfig holds file storage settings.
type StorageConfig struct {
	Provider  string `mapstructure:"provider"` // local, s3
	LocalDir  string `mapstructure:"local_dir"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Endpoint  string `mapstructure:"endpoint"`
	AccessID  string `mapstructure:"access_id"`
	SecretKey string `mapstructure:"secret_key"`
	PublicURL string `mapstructure:"public_url"`
}

// MailConfig holds email service settings.
type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	FromAddr string `mapstructure:"from_addr"`
	FromName string `mapstructure:"from_name"`
}

// ObservabilityConfig holds monitoring and telemetry settings.
type ObservabilityConfig struct {
	TracingEnabled  bool    `mapstructure:"tracing_enabled"`
	MetricsEnabled  bool    `mapstructure:"metrics_enabled"`
	SentryDSN       string  `mapstructure:"sentry_dsn"`
	PrometheusPort  int     `mapstructure:"prometheus_port"`
	SamplingRate    float64 `mapstructure:"sampling_rate"`
	OTLPEndpoint    string  `mapstructure:"otlp_endpoint"` // e.g. localhost:4317
	ServiceName     string  `mapstructure:"service_name"`
}

// NotificationConfig holds settings for all notification channels.
type NotificationConfig struct {
	// Twilio SMS
	TwilioAccountSID string `mapstructure:"twilio_account_sid"`
	TwilioAuthToken  string `mapstructure:"twilio_auth_token"`
	TwilioFromNumber string `mapstructure:"twilio_from_number"`

	// Slack
	SlackWebhookURL string `mapstructure:"slack_webhook_url"`

	// Firebase Cloud Messaging
	FCMServerKey string `mapstructure:"fcm_server_key"`
}

// Load reads the application configuration from environment variables and/or config.yaml.
// Environment variables take precedence over file values.
// Example: APP_PORT=8080 overrides app.port in the file.
// In production, JWT secrets must be at least 32 characters long.
func Load() (*Config, error) {
	// Try to load .env file from current directory or parent directory
	// This ensures it works if running from root or from backend/ folder
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	v := viper.New()

	// Defaults
	v.SetDefault("app.name", "Kodia App")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.debug", true)
	v.SetDefault("app.base_url", "http://localhost:8080")
	v.SetDefault("app.frontend_url", "http://localhost:3000")
	v.SetDefault("app.locale", "en")

	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.name", "kodia.sqlite")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.timezone", "UTC")
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.access_secret", "")
	v.SetDefault("jwt.refresh_secret", "")
	v.SetDefault("jwt.access_expiry_hours", 1)
	v.SetDefault("jwt.refresh_expiry_days", 30)

	v.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})

	v.SetDefault("storage.provider", "local")
	v.SetDefault("storage.local_dir", "./uploads")
	v.SetDefault("storage.bucket", "kodia-bucket")
	v.SetDefault("storage.region", "us-east-1")

	v.SetDefault("mail.host", "localhost")
	v.SetDefault("mail.port", 1025)
	v.SetDefault("mail.user", "")
	v.SetDefault("mail.password", "")
	v.SetDefault("mail.from_addr", "no-reply@kodia.studio")
	v.SetDefault("mail.from_name", "Kodia App")
	
	v.SetDefault("observability.tracing_enabled", false)
	v.SetDefault("observability.metrics_enabled", true)
	v.SetDefault("observability.prometheus_port", 9090)
	v.SetDefault("observability.sampling_rate", 1.0)
	v.SetDefault("observability.otlp_endpoint", "localhost:4317")
	v.SetDefault("observability.service_name", "kodia-api")

	v.SetDefault("notification.twilio_account_sid", "")
	v.SetDefault("notification.twilio_auth_token", "")
	v.SetDefault("notification.twilio_from_number", "")
	v.SetDefault("notification.slack_webhook_url", "")
	v.SetDefault("notification.fcm_server_key", "")

	// Config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("APP")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found — OK, use env vars and defaults
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate JWT secrets
	if err := cfg.ValidateJWTSecrets(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// IsProduction returns true when the app is running in production mode.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.App.Env, "production")
}

// IsDevelopment returns true when the app is running in development mode.
func (c *Config) IsDevelopment() bool {
	return strings.EqualFold(c.App.Env, "development")
}

// ValidateJWTSecrets ensures JWT secrets meet security requirements.
// In production, secrets must be at least 32 characters long.
// In development, if secrets are empty, a warning is issued but execution continues.
func (c *Config) ValidateJWTSecrets() error {
	const minSecretLength = 32

	if c.JWT.AccessSecret == "" {
		if c.IsProduction() {
			return fmt.Errorf("JWT_ACCESS_SECRET must be set in production (minimum 32 characters)")
		}
		// Development mode with empty secret - this is acceptable but warn
		return nil
	}

	if c.JWT.RefreshSecret == "" {
		if c.IsProduction() {
			return fmt.Errorf("JWT_REFRESH_SECRET must be set in production (minimum 32 characters)")
		}
		return nil
	}

	// Check secret lengths in production
	if c.IsProduction() {
		if len(c.JWT.AccessSecret) < minSecretLength {
			return fmt.Errorf("JWT_ACCESS_SECRET must be at least %d characters long in production (got %d)", minSecretLength, len(c.JWT.AccessSecret))
		}
		if len(c.JWT.RefreshSecret) < minSecretLength {
			return fmt.Errorf("JWT_REFRESH_SECRET must be at least %d characters long in production (got %d)", minSecretLength, len(c.JWT.RefreshSecret))
		}
	}

	return nil
}
