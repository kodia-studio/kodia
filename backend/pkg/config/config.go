// Package config provides application configuration loading for Kodia Framework.
// Supports loading from environment variables and config files via Viper.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

// AppConfig holds general application settings.
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Env         string `mapstructure:"env"`
	Port        int    `mapstructure:"port"`
	Debug       bool   `mapstructure:"debug"`
	BaseURL     string `mapstructure:"base_url"`
	FrontendURL string `mapstructure:"frontend_url"`
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

// Load reads the application configuration from environment variables and/or config.yaml.
// Environment variables take precedence over file values.
// Example: APP_PORT=8080 overrides app.port in the file.
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("app.name", "Kodia App")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.debug", true)
	v.SetDefault("app.base_url", "http://localhost:8080")
	v.SetDefault("app.frontend_url", "http://localhost:3000")

	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.timezone", "UTC")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "30m")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.access_expiry_hours", 1)
	v.SetDefault("jwt.refresh_expiry_days", 30)

	v.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})

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
