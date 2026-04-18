package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kodia-studio/kodia/pkg/config"
)

// LocalStorageProvider implements ports.StorageProvider using the local filesystem.
type LocalStorageProvider struct {
	baseDir  string
	publicURL string
}

// NewLocalStorageProvider creates a new LocalStorageProvider.
func NewLocalStorageProvider(cfg *config.Config) *LocalStorageProvider {
	// Ensure directory exists
	_ = os.MkdirAll(cfg.Storage.LocalDir, 0755)

	return &LocalStorageProvider{
		baseDir:   cfg.Storage.LocalDir,
		publicURL: cfg.App.BaseURL + "/uploads", // Default convention
	}
}

func (p *LocalStorageProvider) Upload(ctx context.Context, path string, content io.Reader) (string, error) {
	fullPath := filepath.Join(p.baseDir, path)
	
	// Ensure subdirectories exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, content); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return path, nil
}

func (p *LocalStorageProvider) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(p.baseDir, path)
	return os.Remove(fullPath)
}

func (p *LocalStorageProvider) GetURL(ctx context.Context, path string) (string, error) {
	return fmt.Sprintf("%s/%s", p.publicURL, path), nil
}

func (p *LocalStorageProvider) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(p.baseDir, path)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
