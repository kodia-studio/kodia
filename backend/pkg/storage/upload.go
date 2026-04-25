package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/kodia-studio/kodia/internal/providers"
)

// UploadManager handles file uploads with validation and processing
type UploadManager struct {
	provider  providers.StorageProvider
	validator *providers.FileValidator
	config    *providers.StorageConfig
}

// NewUploadManager creates a new upload manager
func NewUploadManager(provider providers.StorageProvider, config *providers.StorageConfig) *UploadManager {
	validator := &providers.FileValidator{
		MaxSize:      config.MaxFileSize,
		AllowedTypes: config.AllowedMimeTypes,
		AllowedExts:  config.AllowedExtensions,
	}

	return &UploadManager{
		provider:  provider,
		validator: validator,
		config:    config,
	}
}

// Upload handles a single file upload
func (um *UploadManager) Upload(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	// Convert to FileUpload
	file, err := providers.FromMultipartFile(fh)
	if err != nil {
		return "", err
	}

	// Validate
	if err := um.validator.ValidateFile(file); err != nil {
		return "", err
	}

	// Store
	path, err := um.provider.Store(ctx, file)
	if err != nil {
		return "", err
	}

	return path, nil
}

// UploadImage uploads and processes an image
func (um *UploadManager) UploadImage(ctx context.Context, fh *multipart.FileHeader, options *providers.ImageProcessingOptions) (string, error) {
	// Convert to FileUpload
	file, err := providers.FromMultipartFile(fh)
	if err != nil {
		return "", err
	}

	// Validate it's an image
	if !isImageType(file.ContentType) {
		return "", fmt.Errorf("file must be an image (JPEG, PNG, GIF, WebP)")
	}

	// Validate
	if err := um.validator.ValidateFile(file); err != nil {
		return "", err
	}

	// Check if provider supports image processing
	if imgProvider, ok := um.provider.(providers.StorageProviderWithImages); ok {
		return imgProvider.StoreImage(ctx, file, options)
	}

	// Fall back to regular storage
	return um.provider.Store(ctx, file)
}

// UploadMultiple uploads multiple files
func (um *UploadManager) UploadMultiple(ctx context.Context, fhs []*multipart.FileHeader) ([]string, error) {
	var files []*providers.FileUpload

	for _, fh := range fhs {
		file, err := providers.FromMultipartFile(fh)
		if err != nil {
			return nil, err
		}

		if err := um.validator.ValidateFile(file); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return um.provider.StoreMultiple(ctx, files)
}

// Delete deletes a file
func (um *UploadManager) Delete(ctx context.Context, path string) error {
	return um.provider.Delete(ctx, path)
}

// DeleteMultiple deletes multiple files
func (um *UploadManager) DeleteMultiple(ctx context.Context, paths []string) error {
	return um.provider.DeleteMultiple(ctx, paths)
}

// GetURL gets the public URL for a file
func (um *UploadManager) GetURL(ctx context.Context, path string) (string, error) {
	return um.provider.GetURL(ctx, path)
}

// GetContent retrieves file content
func (um *UploadManager) GetContent(ctx context.Context, path string) ([]byte, error) {
	return um.provider.Get(ctx, path)
}

// Exists checks if file exists
func (um *UploadManager) Exists(ctx context.Context, path string) (bool, error) {
	return um.provider.Exists(ctx, path)
}

// Move moves a file from source to destination
func (um *UploadManager) Move(ctx context.Context, sourcePath, destinationPath string) error {
	// Copy file
	if err := um.provider.Copy(ctx, sourcePath, destinationPath); err != nil {
		return err
	}

	// Delete original
	return um.provider.Delete(ctx, sourcePath)
}

// Helper functions

func isImageType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

// UploadResult represents the result of a file upload
type UploadResult struct {
	Path      string `json:"path"`
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mime_type"`
	UploadedAt int64 `json:"uploaded_at"`
}

// BuildUploadResult builds an UploadResult from a FileUpload and path
func (um *UploadManager) BuildUploadResult(ctx context.Context, file *providers.FileUpload, path string) (*UploadResult, error) {
	url, err := um.GetURL(ctx, path)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Path:       path,
		URL:        url,
		Filename:   file.Filename,
		Size:       file.Size,
		MimeType:   file.ContentType,
		UploadedAt: ctx.Value("timestamp").(int64),
	}, nil
}

// FileInfo represents file information
type FileInfo struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	URL      string `json:"url"`
	Exists   bool   `json:"exists"`
}

// GetFileInfo retrieves information about a file
func (um *UploadManager) GetFileInfo(ctx context.Context, path string) (*FileInfo, error) {
	exists, err := um.provider.Exists(ctx, path)
	if err != nil {
		return nil, err
	}

	url, err := um.GetURL(ctx, path)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Path:     path,
		Filename: filepath.Base(path),
		URL:      url,
		Exists:   exists,
	}, nil
}
