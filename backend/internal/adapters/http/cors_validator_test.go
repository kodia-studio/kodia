package http

import (
	"testing"

	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
)

// TestValidateCORSConfigWildcardWithCredentials verifies that wildcard origins with credentials are rejected
func TestValidateCORSConfigWildcardWithCredentials(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
		},
	}

	err := ValidateCORSConfig(cfg, logger)
	if err == nil {
		t.Error("Expected error for wildcard origin with credentials, but got none")
	}

	if err != nil {
		errMsg := err.Error()
		if !contains(errMsg, "wildcard") || !contains(errMsg, "credentials") {
			t.Errorf("Error message should mention wildcard and credentials: %v", err)
		}
	}
}

// TestValidateCORSConfigSpecificOrigins verifies that specific origins are allowed
func TestValidateCORSConfigSpecificOrigins(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name          string
		origins       []string
		shouldError   bool
		description   string
	}{
		{
			name:        "single https origin",
			origins:     []string{"https://example.com"},
			shouldError: false,
			description: "Single specific HTTPS origin should be allowed",
		},
		{
			name:        "multiple https origins",
			origins:     []string{"https://app.example.com", "https://admin.example.com"},
			shouldError: false,
			description: "Multiple specific HTTPS origins should be allowed",
		},
		{
			name:        "localhost for development",
			origins:     []string{"http://localhost:3000"},
			shouldError: false,
			description: "Localhost should be allowed",
		},
		{
			name:        "mixed http and https",
			origins:     []string{"http://localhost:3000", "https://example.com"},
			shouldError: false,
			description: "Mix of localhost and HTTPS should be allowed",
		},
		{
			name:        "origin without scheme",
			origins:     []string{"example.com"},
			shouldError: true,
			description: "Origin without scheme should be rejected",
		},
		{
			name:        "invalid origin format",
			origins:     []string{"not-a-url"},
			shouldError: true,
			description: "Invalid origin format should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				CORS: config.CORSConfig{
					AllowedOrigins: tt.origins,
				},
			}

			err := ValidateCORSConfig(cfg, logger)
			hasError := err != nil

			if hasError != tt.shouldError {
				if tt.shouldError {
					t.Errorf("%s: Expected error but got none", tt.description)
				} else {
					t.Errorf("%s: Expected no error but got: %v", tt.description, err)
				}
			}
		})
	}
}

// TestValidateCORSConfigEmptyConfig verifies that empty config is allowed
func TestValidateCORSConfigEmptyConfig(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name      string
		cfg       *config.Config
		wantError bool
	}{
		{
			name: "empty CORS config",
			cfg: &config.Config{
				CORS: config.CORSConfig{},
			},
			wantError: false,
		},
		{
			name: "empty origins list",
			cfg: &config.Config{
				CORS: config.CORSConfig{
					AllowedOrigins: []string{},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCORSConfig(tt.cfg, logger)
			if (err != nil) != tt.wantError {
				if tt.wantError {
					t.Errorf("Expected error but got none")
				} else {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestCORSOriginValidationAttackScenarios verifies protection against common CORS attacks
func TestCORSOriginValidationAttackScenarios(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name        string
		origins     []string
		shouldError bool
		description string
	}{
		{
			name:        "wildcard origin attack",
			origins:     []string{"*"},
			shouldError: true,
			description: "Wildcard should be rejected to prevent credential theft",
		},
		{
			name:        "prefix wildcard attack",
			origins:     []string{"*.example.com"},
			shouldError: true,
			description: "Prefix wildcards are not proper absolute URLs",
		},
		{
			name:        "subdomain enumeration prevention",
			origins:     []string{"https://*.example.com"},
			shouldError: true,
			description: "Wildcard subdomains should be rejected",
		},
		{
			name:        "credential theft prevention",
			origins:     []string{"https://attacker.com"},
			shouldError: false,
			description: "Attacker origin is syntactically valid but won't be attacked by us",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				CORS: config.CORSConfig{
					AllowedOrigins: tt.origins,
				},
			}

			err := ValidateCORSConfig(cfg, logger)
			hasError := err != nil

			if hasError != tt.shouldError {
				if tt.shouldError {
					t.Errorf("%s: Expected error but got none", tt.description)
				} else {
					t.Errorf("%s: Expected no error but got: %v", tt.description, err)
				}
			}
		})
	}
}

// TestGetCORSConfigReturnsValidConfig verifies that CORS config is properly formatted
func TestGetCORSConfigReturnsValidConfig(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"https://example.com"},
		},
	}

	corsConfig := GetCORSConfig(cfg)

	if corsConfig["AllowOrigins"] == nil {
		t.Error("CORS config should have AllowOrigins")
	}

	if corsConfig["AllowMethods"] == nil {
		t.Error("CORS config should have AllowMethods")
	}

	if corsConfig["AllowHeaders"] == nil {
		t.Error("CORS config should have AllowHeaders")
	}

	if corsConfig["AllowCredentials"] != true {
		t.Error("CORS config should have AllowCredentials=true")
	}
}

// Helper function
func contains(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// BenchmarkCORSValidation measures validation performance
func BenchmarkCORSValidation(b *testing.B) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{
				"https://app.example.com",
				"https://api.example.com",
				"https://admin.example.com",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateCORSConfig(cfg, logger)
	}
}
