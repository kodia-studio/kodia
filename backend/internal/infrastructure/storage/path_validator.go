package storage

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidatePath ensures the given path is safe and doesn't allow directory traversal.
// Returns the cleaned path if valid, or an error if the path attempts to escape the base directory.
func ValidatePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Reject absolute paths
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths are not allowed")
	}

	// Reject paths starting with /
	if strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("paths must be relative")
	}

	// Reject Windows drive letters (C:, D:, etc.)
	if len(path) >= 2 && path[1] == ':' {
		return "", fmt.Errorf("absolute Windows paths are not allowed")
	}

	// Clean the path to resolve .. and . components
	cleanPath := filepath.Clean(path)

	// Reject if cleaning resulted in going up a directory
	if strings.HasPrefix(cleanPath, "..") || strings.Contains(cleanPath, "/..") {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}

	// Reject null bytes
	if strings.Contains(cleanPath, "\x00") {
		return "", fmt.Errorf("null bytes not allowed in path")
	}

	// Normalize separators (prevent mixed / and \ on Windows)
	normalizedPath := filepath.ToSlash(cleanPath)

	return normalizedPath, nil
}

// ValidatePathWithinBase validates a path and ensures it stays within the base directory.
// This provides defense-in-depth protection against path traversal.
func ValidatePathWithinBase(basePath, targetPath string) (string, error) {
	// First validate the path itself
	cleanPath, err := ValidatePath(targetPath)
	if err != nil {
		return "", err
	}

	// Construct the full path
	fullPath := filepath.Join(basePath, cleanPath)

	// Ensure the full path is within the base directory
	basePath = filepath.Clean(basePath)
	fullPath = filepath.Clean(fullPath)

	// Convert to absolute paths for comparison
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base path: %w", err)
	}

	absTarget, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve target path: %w", err)
	}

	// Ensure the target path starts with the base path
	// Add filepath.Separator to avoid matching partial directory names
	if !strings.HasPrefix(absTarget, absBase+string(filepath.Separator)) && absTarget != absBase {
		return "", fmt.Errorf("path escapes base directory: %s", targetPath)
	}

	return cleanPath, nil
}
