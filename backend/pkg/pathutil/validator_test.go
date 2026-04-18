package pathutil

import (
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantError bool
		wantPath  string
	}{
		// Valid paths
		{
			name:      "simple filename",
			path:      "welcome.html",
			wantError: false,
			wantPath:  "welcome.html",
		},
		{
			name:      "nested directory",
			path:      "emails/welcome.html",
			wantError: false,
			wantPath:  "emails/welcome.html",
		},
		{
			name:      "deeply nested",
			path:      "en/templates/email/welcome.html",
			wantError: false,
			wantPath:  "en/templates/email/welcome.html",
		},
		{
			name:      "filename with dots",
			path:      "invoice.template.html",
			wantError: false,
			wantPath:  "invoice.template.html",
		},

		// Path traversal attacks
		{
			name:      "simple parent directory traversal",
			path:      "../etc/passwd",
			wantError: true,
		},
		{
			name:      "multiple parent directory traversal",
			path:      "../../../../../../etc/passwd",
			wantError: true,
		},
		{
			name:      "parent directory in middle",
			path:      "emails/../../../config.html",
			wantError: true,
		},
		{
			name:      "windows UNC path",
			path:      "..\\..\\windows\\system32",
			wantError: true,
		},

		// Invalid paths
		{
			name:      "empty path",
			path:      "",
			wantError: true,
		},
		{
			name:      "absolute path unix",
			path:      "/etc/passwd",
			wantError: true,
		},
		{
			name:      "absolute path windows",
			path:      "C:\\Windows\\System32",
			wantError: true,
		},
		{
			name:      "path starting with slash",
			path:      "/emails/welcome.html",
			wantError: true,
		},
		{
			name:      "null byte injection",
			path:      "welcome.html\x00.txt",
			wantError: true,
		},
		{
			name:      "double dot",
			path:      "..",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePath(tt.path)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePath() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && got != tt.wantPath {
				t.Errorf("ValidatePath() got = %v, want %v", got, tt.wantPath)
			}
		})
	}
}

func TestValidatePathWithinBase(t *testing.T) {
	baseDir := "/var/mail/templates"

	tests := []struct {
		name       string
		basePath   string
		targetPath string
		wantError  bool
		wantPath   string
	}{
		// Valid paths within base
		{
			name:       "simple template in base",
			basePath:   baseDir,
			targetPath: "welcome.html",
			wantError:  false,
			wantPath:   "welcome.html",
		},
		{
			name:       "nested template within base",
			basePath:   baseDir,
			targetPath: "en/welcome.html",
			wantError:  false,
			wantPath:   "en/welcome.html",
		},
		{
			name:       "deeply nested within base",
			basePath:   baseDir,
			targetPath: "emails/transactional/invoice.html",
			wantError:  false,
			wantPath:   "emails/transactional/invoice.html",
		},

		// Path traversal attacks
		{
			name:       "escape to parent",
			basePath:   baseDir,
			targetPath: "../../../etc/passwd",
			wantError:  true,
		},
		{
			name:       "escape with nested traversal",
			basePath:   baseDir,
			targetPath: "emails/../../../config.html",
			wantError:  true,
		},
		{
			name:       "windows style escape",
			basePath:   baseDir,
			targetPath: "..\\..\\windows\\system32",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePathWithinBase(tt.basePath, tt.targetPath)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePathWithinBase() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && got != tt.wantPath {
				t.Errorf("ValidatePathWithinBase() got = %v, want %v", got, tt.wantPath)
			}
		})
	}
}

func BenchmarkValidatePath(b *testing.B) {
	paths := []string{
		"welcome.html",
		"emails/welcome.html",
		"en/templates/email/welcome.html",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			ValidatePath(path)
		}
	}
}

func BenchmarkValidatePathWithinBase(b *testing.B) {
	baseDir := "/var/mail/templates"
	paths := []string{
		"welcome.html",
		"emails/welcome.html",
		"en/templates/email/welcome.html",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			ValidatePathWithinBase(baseDir, path)
		}
	}
}
