package ports

import (
	"context"
	"io"
)

// StorageProvider defines the interface for file storage operations.
type StorageProvider interface {
	// Upload saves a file to the storage provider.
	Upload(ctx context.Context, path string, content io.Reader) (string, error)
	// Delete removes a file from the storage provider.
	Delete(ctx context.Context, path string) error
	// GetURL returns the public URL for a file.
	GetURL(ctx context.Context, path string) (string, error)
	// Exists checks if a file exists in the storage provider.
	Exists(ctx context.Context, path string) (bool, error)
}
