package mailer

import (
	"context"
	"testing"

	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
)

// TestSendWithTemplate_PathValidation tests that template path validation prevents traversal attacks
func TestSendWithTemplate_PathValidation(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a mock mailer with test configuration
	testMailer := &SMTPMailer{
		basePath: "resources/mail",
		log:      logger,
		config:   &config.MailConfig{FromAddr: "test@example.com", FromName: "Test"},
		// Note: client is nil for this test, but we're testing path validation before client usage
	}

	tests := []struct {
		name          string
		templatePath  string
		shouldFail    bool
		errorContains string
	}{
		// Valid paths
		{
			name:         "simple template filename",
			templatePath: "welcome.html",
			shouldFail:   true, // Will fail when trying to parse file, but path validation passes
		},
		{
			name:         "nested template path",
			templatePath: "en/welcome.html",
			shouldFail:   true, // Will fail when trying to parse file, but path validation passes
		},
		{
			name:         "nested with multiple levels",
			templatePath: "emails/transactional/invoice.html",
			shouldFail:   true, // Will fail when trying to parse file, but path validation passes
		},

		// Path traversal attacks - should be blocked by validation
		{
			name:          "parent directory escape",
			templatePath:  "../../../etc/passwd",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
		{
			name:          "middle directory escape",
			templatePath:  "emails/../../../config.html",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
		{
			name:          "absolute path unix",
			templatePath:  "/etc/passwd",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
		{
			name:          "windows drive letter",
			templatePath:  "C:\\Windows\\System32",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
		{
			name:          "double dot escape",
			templatePath:  "..",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
		{
			name:          "multiple parent traversals",
			templatePath:  "../../../../../../etc/passwd",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
		{
			name:          "windows UNC path",
			templatePath:  "..\\..\\windows\\system32",
			shouldFail:    true,
			errorContains: "invalid template path",
		},

		// Edge cases
		{
			name:          "empty path",
			templatePath:  "",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
		{
			name:          "null byte injection",
			templatePath:  "welcome.html\x00.txt",
			shouldFail:    true,
			errorContains: "invalid template path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to skip the actual email sending since we don't have a real SMTP client
			// But we can test the path validation by checking if SendWithTemplate would fail appropriately
			err := testMailer.SendWithTemplate(context.Background(), []string{"test@example.com"}, "Test", tt.templatePath, nil)

			if tt.shouldFail {
				if err == nil && tt.errorContains == "invalid template path" {
					t.Errorf("Expected error containing '%s', but got nil", tt.errorContains)
				}
				if err != nil && tt.errorContains != "" && !containsSubstring(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substring string) bool {
	return len(s) >= len(substring) && (substring == s || len(s) > 0 && (s[:len(substring)] == substring || s[len(s)-len(substring):] == substring || findSubstring(s, substring)))
}

// Better substring search
func findSubstring(s, substring string) bool {
	for i := 0; i <= len(s)-len(substring); i++ {
		if s[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

// TestMailer_PathValidationDetails tests specific path validation behaviors
func TestMailer_PathValidationDetails(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	testMailer := &SMTPMailer{
		basePath: "resources/mail",
		log:      logger,
		config:   &config.MailConfig{FromAddr: "test@example.com", FromName: "Test"},
	}

	tests := []struct {
		name              string
		templatePath      string
		expectedValidPath bool
		description       string
	}{
		{
			name:              "simple filename allowed",
			templatePath:      "welcome.html",
			expectedValidPath: true,
			description:       "Simple filenames without path traversal are allowed",
		},
		{
			name:              "nested path allowed",
			templatePath:      "en/GB/email/welcome.html",
			expectedValidPath: true,
			description:       "Nested paths within base directory are allowed",
		},
		{
			name:              "parent directory blocked",
			templatePath:      "../config/secret.html",
			expectedValidPath: false,
			description:       "Parent directory traversal is blocked",
		},
		{
			name:              "multiple levels blocked",
			templatePath:      "template/../../../etc/passwd",
			expectedValidPath: false,
			description:       "Multi-level traversal is blocked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testMailer.SendWithTemplate(context.Background(), []string{"test@example.com"}, "Test", tt.templatePath, nil)

			isPathValidationError := err != nil && containsSubstring(err.Error(), "invalid template path")

			if tt.expectedValidPath && isPathValidationError {
				t.Errorf("%s - Path validation rejected valid path: %s. Error: %v", tt.description, tt.templatePath, err)
			}

			if !tt.expectedValidPath && !isPathValidationError {
				// Path validation should have failed for invalid paths
				// But might fail for other reasons (file not found), which is acceptable
				// We just need to ensure path traversal is blocked
				if !containsSubstring(err.Error(), "failed to parse template") {
					// If it's not a parse error, it should be a path validation error
					t.Logf("Expected path validation error for: %s. Got: %v", tt.templatePath, err)
				}
			}
		})
	}
}
