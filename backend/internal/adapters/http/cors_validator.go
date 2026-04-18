package http

import (
	"fmt"
	"strings"

	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
)

// ValidateCORSConfig validates CORS configuration for security issues.
// Returns an error if insecure configuration is detected.
func ValidateCORSConfig(cfg *config.Config, log *zap.Logger) error {
	if len(cfg.CORS.AllowedOrigins) == 0 {
		// No CORS configuration is fine
		return nil
	}

	corsOrigins := cfg.CORS.AllowedOrigins

	// Check for wildcard origin with credentials enabled
	hasWildcard := false
	for _, origin := range corsOrigins {
		// Check for any wildcard (either exact "*" or wildcard in subdomain like "*.example.com")
		if origin == "*" || strings.Contains(origin, "*") {
			hasWildcard = true
			break
		}
	}

	// Security check: Wildcard origins with AllowCredentials=true is insecure
	// Credentials allow cookies and authorization headers to be sent
	// If wildcard origins are allowed, any website can make requests with user credentials
	if hasWildcard {
		return fmt.Errorf(
			"insecure CORS configuration: cannot use wildcard origin '*' with credentials. "+
				"Wildcard origins allow ANY website to access your API with user credentials. "+
				"Solution: Use specific allowed origins instead (e.g., https://example.com, https://app.example.com)",
		)
	}

	// Validate that all origins are absolute URLs (have scheme)
	for _, origin := range corsOrigins {
		if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
			return fmt.Errorf(
				"invalid CORS origin '%s': origins must be absolute URLs with scheme (http:// or https://)",
				origin,
			)
		}
	}

	// Production security warnings
	if cfg.IsProduction() {
		// Warn if HTTP origins are used in production
		for _, origin := range corsOrigins {
			if strings.HasPrefix(origin, "http://") && !strings.Contains(origin, "localhost") {
				log.Warn(
					"Insecure CORS origin in production",
					zap.String("origin", origin),
					zap.String("recommendation", "Use https:// origins in production for security"),
				)
			}
		}

		// Warn if localhost is allowed in production
		for _, origin := range corsOrigins {
			if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
				log.Warn(
					"Localhost origin allowed in production",
					zap.String("origin", origin),
					zap.String("recommendation", "Remove localhost origins from production configuration"),
				)
			}
		}
	}

	return nil
}

// GetCORSConfig returns the CORS configuration for the Gin CORS middleware.
// This function should only be called after ValidateCORSConfig has passed.
func GetCORSConfig(cfg *config.Config) map[string]interface{} {
	origins := cfg.CORS.AllowedOrigins
	if origins == nil {
		origins = []string{}
	}

	return map[string]interface{}{
		"AllowOrigins":     origins,
		"AllowMethods":     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		"AllowHeaders":     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		"ExposeHeaders":    []string{"Content-Length"},
		"AllowCredentials": true,
	}
}
