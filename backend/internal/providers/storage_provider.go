package providers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

// StorageProvider defines the interface for file storage operations
type StorageProvider interface {
	// Store saves a file and returns the path/URL
	Store(ctx context.Context, file *FileUpload) (string, error)

	// StoreMultiple saves multiple files
	StoreMultiple(ctx context.Context, files []*FileUpload) ([]string, error)

	// Get retrieves file content
	Get(ctx context.Context, path string) ([]byte, error)

	// GetStream retrieves file content as stream
	GetStream(ctx context.Context, path string) (io.Reader, error)

	// Delete removes a file
	Delete(ctx context.Context, path string) error

	// DeleteMultiple removes multiple files
	DeleteMultiple(ctx context.Context, paths []string) error

	// Exists checks if file exists
	Exists(ctx context.Context, path string) (bool, error)

	// GetURL returns the public URL for a file
	GetURL(ctx context.Context, path string) (string, error)

	// GetTemporaryURL returns a temporary signed URL (for private files)
	GetTemporaryURL(ctx context.Context, path string, duration time.Duration) (string, error)

	// Copy copies a file to another path
	Copy(ctx context.Context, sourcePath, destinationPath string) error

	// Close closes the storage provider
	Close() error
}

// FileUpload represents a file to be uploaded
type FileUpload struct {
	// File content
	Content io.Reader

	// Original filename
	Filename string

	// MIME type
	ContentType string

	// File size in bytes
	Size int64

	// Destination path (including filename)
	Path string

	// Optional metadata
	Metadata map[string]string

	// Access level: public or private
	AccessLevel AccessLevel

	// Optional: custom headers
	Headers map[string]string
}

// AccessLevel defines file access level
type AccessLevel string

const (
	AccessPublic  AccessLevel = "public"
	AccessPrivate AccessLevel = "private"
)

// StorageConfig contains configuration for the storage provider
type StorageConfig struct {
	Driver           string          `mapstructure:"driver" default:"local"`
	Local            *LocalConfig    `mapstructure:"local"`
	S3               *S3Config       `mapstructure:"s3"`
	Cloudflare       *CloudflareConfig `mapstructure:"cloudflare"`
	MaxFileSize      int64           `mapstructure:"max_file_size" default:"52428800"` // 50MB
	AllowedMimeTypes []string        `mapstructure:"allowed_mime_types"`
	AllowedExtensions []string       `mapstructure:"allowed_extensions"`
	EnableImageProcessing bool        `mapstructure:"enable_image_processing" default:"true"`
}

// LocalConfig contains local storage configuration
type LocalConfig struct {
	Path      string `mapstructure:"path" default:"./storage/uploads"`
	URL       string `mapstructure:"url" default:"http://localhost:8080/storage"`
	Symlink   bool   `mapstructure:"symlink" default:"false"`
}

// S3Config contains AWS S3 configuration
type S3Config struct {
	Region      string `mapstructure:"region" default:"us-east-1"`
	Bucket      string `mapstructure:"bucket"`
	AccessKey   string `mapstructure:"access_key"`
	SecretKey   string `mapstructure:"secret_key"`
	Endpoint    string `mapstructure:"endpoint"` // Optional for S3-compatible services
	Path        string `mapstructure:"path" default:"uploads"`
	ACL         string `mapstructure:"acl" default:"public-read"`
	CloudFront  string `mapstructure:"cloudfront"` // Optional CloudFront distribution
}

// CloudflareConfig contains Cloudflare R2 configuration
type CloudflareConfig struct {
	AccountID   string `mapstructure:"account_id"`
	AccessKey   string `mapstructure:"access_key"`
	SecretKey   string `mapstructure:"secret_key"`
	Bucket      string `mapstructure:"bucket"`
	Domain      string `mapstructure:"domain"`
	Path        string `mapstructure:"path" default:"uploads"`
	PublicURL   bool   `mapstructure:"public_url" default:"true"`
}

// FileValidator validates files before storage
type FileValidator struct {
	MaxSize      int64
	AllowedTypes []string
	AllowedExts  []string
}

// ValidateFile validates file metadata
func (fv *FileValidator) ValidateFile(file *FileUpload) error {
	if file == nil {
		return fmt.Errorf("file is required")
	}

	if file.Filename == "" {
		return fmt.Errorf("filename is required")
	}

	if file.ContentType == "" {
		return fmt.Errorf("content type is required")
	}

	if fv.MaxSize > 0 && file.Size > fv.MaxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", fv.MaxSize)
	}

	if len(fv.AllowedTypes) > 0 && !isAllowedType(file.ContentType, fv.AllowedTypes) {
		return fmt.Errorf("file type %s is not allowed", file.ContentType)
	}

	if len(fv.AllowedExts) > 0 && !hasAllowedExtension(file.Filename, fv.AllowedExts) {
		return fmt.Errorf("file extension for %s is not allowed", file.Filename)
	}

	return nil
}

// NewStorageProvider creates a new storage provider based on config
func NewStorageProvider(config *StorageConfig) (StorageProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("storage config is required")
	}

	switch config.Driver {
	case "local":
		if config.Local == nil {
			return nil, fmt.Errorf("local config is required for local driver")
		}
		return NewLocalProvider(config)

	case "s3":
		if config.S3 == nil {
			return nil, fmt.Errorf("S3 config is required for s3 driver")
		}
		return NewS3Provider(config)

	case "cloudflare", "r2":
		if config.Cloudflare == nil {
			return nil, fmt.Errorf("Cloudflare config is required for cloudflare driver")
		}
		return NewCloudflareProvider(config)

	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", config.Driver)
	}
}

// FromMultipartFile converts multipart.FileHeader to FileUpload
func FromMultipartFile(fh *multipart.FileHeader) (*FileUpload, error) {
	if fh == nil {
		return nil, fmt.Errorf("file header is required")
	}

	file, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &FileUpload{
		Content:     file,
		Filename:    fh.Filename,
		ContentType: fh.Header.Get("Content-Type"),
		Size:        fh.Size,
		Metadata:    make(map[string]string),
		AccessLevel: AccessPublic,
		Headers:     make(map[string]string),
	}, nil
}

// Helper functions

func isAllowedType(contentType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if contentType == allowed || (contentType[:len(allowed)] == allowed && contentType[len(allowed)] == '/') {
			return true
		}
	}
	return false
}

func hasAllowedExtension(filename string, allowedExts []string) bool {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			ext := filename[i:]
			for _, allowed := range allowedExts {
				if ext == allowed || ext == "."+allowed {
					return true
				}
			}
			return false
		}
	}
	return false
}

// ImageProcessingOptions defines options for image processing
type ImageProcessingOptions struct {
	// Resize dimensions
	Width  int
	Height int
	// Quality (1-100)
	Quality int
	// Generate thumbnail
	Thumbnail bool
	// Thumbnail size
	ThumbnailWidth  int
	ThumbnailHeight int
}

// StorageProviderWithImages extends StorageProvider with image processing
type StorageProviderWithImages interface {
	StorageProvider
	// StoreImage stores an image with optional processing
	StoreImage(ctx context.Context, file *FileUpload, options *ImageProcessingOptions) (string, error)
}
