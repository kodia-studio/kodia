package storage

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
			path:      "document.pdf",
			wantError: false,
			wantPath:  "document.pdf",
		},
		{
			name:      "nested directory",
			path:      "uploads/user-123/profile.jpg",
			wantError: false,
			wantPath:  "uploads/user-123/profile.jpg",
		},
		{
			name:      "path with dots in filename",
			path:      "archive.backup.tar.gz",
			wantError: false,
			wantPath:  "archive.backup.tar.gz",
		},
		{
			name:      "path with numbers and underscores",
			path:      "user_data_2024_01_15.csv",
			wantError: false,
			wantPath:  "user_data_2024_01_15.csv",
		},
		{
			name:      "nested path with mixed separators normalized",
			path:      "uploads/docs/file.txt",
			wantError: false,
			wantPath:  "uploads/docs/file.txt",
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
			path:      "uploads/../../../etc/passwd",
			wantError: true,
		},
		{
			name:      "dot-dot-slash attack",
			path:      "files/..\\..\\windows\\system32",
			wantError: true,
		},
		{
			name:      "windows UNC path",
			path:      "..\\..\\system32",
			wantError: true,
		},
		{
			name:      "encoded path traversal attempt",
			path:      "uploads/..%2f..%2fetc",
			wantError: true, // Reject paths that look like traversal even if encoded
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
			wantError: true, // Reject Windows drive letters
		},
		{
			name:      "path starting with slash",
			path:      "/uploads/file.txt",
			wantError: true,
		},
		{
			name:      "null byte injection",
			path:      "file.txt\x00.exe",
			wantError: true,
		},
		{
			name:      "single dot",
			path:      ".",
			wantError: false,
			wantPath:  ".",
		},
		{
			name:      "double dot",
			path:      "..",
			wantError: true,
		},
		{
			name:      "mixed separators in attack",
			path:      "..\\..\\..\\etc/passwd",
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
	baseDir := "/var/uploads"

	tests := []struct {
		name      string
		basePath  string
		targetPath string
		wantError bool
		wantPath  string
	}{
		// Valid paths within base
		{
			name:       "simple file in base",
			basePath:   baseDir,
			targetPath: "document.pdf",
			wantError:  false,
			wantPath:   "document.pdf",
		},
		{
			name:       "nested directory within base",
			basePath:   baseDir,
			targetPath: "user-123/profile.jpg",
			wantError:  false,
			wantPath:   "user-123/profile.jpg",
		},
		{
			name:       "deeply nested within base",
			basePath:   baseDir,
			targetPath: "2024/01/15/document.pdf",
			wantError:  false,
			wantPath:   "2024/01/15/document.pdf",
		},

		// Path traversal attacks
		{
			name:       "escape to parent",
			basePath:   baseDir,
			targetPath: "../etc/passwd",
			wantError:  true,
		},
		{
			name:       "escape with nested traversal",
			basePath:   baseDir,
			targetPath: "user/../../../etc/passwd",
			wantError:  true,
		},
		{
			name:       "multiple traversals",
			basePath:   baseDir,
			targetPath: "../../../../../../etc/passwd",
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

// Benchmark tests for performance
func BenchmarkValidatePath(b *testing.B) {
	validPaths := []string{
		"document.pdf",
		"user-123/profile.jpg",
		"uploads/docs/2024/01/file.pdf",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range validPaths {
			ValidatePath(path)
		}
	}
}

func BenchmarkValidatePathWithinBase(b *testing.B) {
	baseDir := "/var/uploads"
	validPaths := []string{
		"document.pdf",
		"user-123/profile.jpg",
		"uploads/docs/2024/01/file.pdf",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range validPaths {
			ValidatePathWithinBase(baseDir, path)
		}
	}
}
