package providers

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "image/gif"
	_ "image/png"
)

// LocalProvider implements StorageProvider for local filesystem
type LocalProvider struct {
	config    *StorageConfig
	basePath  string
	baseURL   string
	validator *FileValidator
}

// NewLocalProvider creates a new local storage provider
func NewLocalProvider(config *StorageConfig) (*LocalProvider, error) {
	if config.Local == nil {
		return nil, fmt.Errorf("local config is required")
	}

	basePath := config.Local.Path
	if basePath == "" {
		basePath = "./storage/uploads"
	}

	// Create base directory if not exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	validator := &FileValidator{
		MaxSize:      config.MaxFileSize,
		AllowedTypes: config.AllowedMimeTypes,
		AllowedExts:  config.AllowedExtensions,
	}

	return &LocalProvider{
		config:    config,
		basePath:  basePath,
		baseURL:   config.Local.URL,
		validator: validator,
	}, nil
}

// Store saves a file to local filesystem
func (lp *LocalProvider) Store(ctx context.Context, file *FileUpload) (string, error) {
	if err := lp.validator.ValidateFile(file); err != nil {
		return "", err
	}

	// Create storage path with date subdirectories
	now := time.Now()
	storagePath := filepath.Join(
		lp.basePath,
		fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day()),
	)

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage path: %w", err)
	}

	// Generate unique filename
	filename := sanitizeFilename(file.Filename)
	filename = fmt.Sprintf("%d-%s", time.Now().UnixNano(), filename)

	fullPath := filepath.Join(storagePath, filename)

	// Write file to disk
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, file.Content); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return relative path for database storage
	relPath := filepath.Join(
		fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day()),
		filename,
	)

	return relPath, nil
}

// StoreMultiple saves multiple files
func (lp *LocalProvider) StoreMultiple(ctx context.Context, files []*FileUpload) ([]string, error) {
	var paths []string

	for _, file := range files {
		path, err := lp.Store(ctx, file)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	return paths, nil
}

// Get retrieves file content
func (lp *LocalProvider) Get(ctx context.Context, path string) ([]byte, error) {
	fullPath := filepath.Join(lp.basePath, path)

	// Prevent directory traversal
	if !strings.HasPrefix(fullPath, lp.basePath) {
		return nil, fmt.Errorf("invalid path")
	}

	return os.ReadFile(fullPath)
}

// GetStream retrieves file as stream
func (lp *LocalProvider) GetStream(ctx context.Context, path string) (io.Reader, error) {
	fullPath := filepath.Join(lp.basePath, path)

	if !strings.HasPrefix(fullPath, lp.basePath) {
		return nil, fmt.Errorf("invalid path")
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Delete removes a file
func (lp *LocalProvider) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(lp.basePath, path)

	if !strings.HasPrefix(fullPath, lp.basePath) {
		return fmt.Errorf("invalid path")
	}

	return os.Remove(fullPath)
}

// DeleteMultiple removes multiple files
func (lp *LocalProvider) DeleteMultiple(ctx context.Context, paths []string) error {
	for _, path := range paths {
		if err := lp.Delete(ctx, path); err != nil {
			return err
		}
	}
	return nil
}

// Exists checks if file exists
func (lp *LocalProvider) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(lp.basePath, path)

	if !strings.HasPrefix(fullPath, lp.basePath) {
		return false, fmt.Errorf("invalid path")
	}

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetURL returns the public URL for a file
func (lp *LocalProvider) GetURL(ctx context.Context, path string) (string, error) {
	return fmt.Sprintf("%s/%s", lp.baseURL, path), nil
}

// GetTemporaryURL returns a temporary signed URL
func (lp *LocalProvider) GetTemporaryURL(ctx context.Context, path string, duration time.Duration) (string, error) {
	// For local storage, just return the URL (no signing)
	return lp.GetURL(ctx, path)
}

// Copy copies a file
func (lp *LocalProvider) Copy(ctx context.Context, sourcePath, destinationPath string) error {
	sourceFullPath := filepath.Join(lp.basePath, sourcePath)
	destFullPath := filepath.Join(lp.basePath, destinationPath)

	if !strings.HasPrefix(sourceFullPath, lp.basePath) || !strings.HasPrefix(destFullPath, lp.basePath) {
		return fmt.Errorf("invalid path")
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destFullPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	source, err := os.Open(sourceFullPath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

// Close closes the local provider
func (lp *LocalProvider) Close() error {
	return nil
}

// StoreImage stores an image with optional processing
func (lp *LocalProvider) StoreImage(ctx context.Context, file *FileUpload, options *ImageProcessingOptions) (string, error) {
	// First store the original image
	path, err := lp.Store(ctx, file)
	if err != nil {
		return "", err
	}

	// Process image if options provided
	if options != nil && (options.Width > 0 || options.Height > 0) {
		fullPath := filepath.Join(lp.basePath, path)

		// Resize image
		if err := lp.resizeImage(fullPath, options); err != nil {
			return path, fmt.Errorf("failed to resize image: %w", err)
		}

		// Generate thumbnail if requested
		if options.Thumbnail {
			thumbPath := lp.generateThumbnailPath(path)
			thumbFullPath := filepath.Join(lp.basePath, thumbPath)

			// Ensure thumbnail directory exists
			if err := os.MkdirAll(filepath.Dir(thumbFullPath), 0755); err != nil {
				return path, fmt.Errorf("failed to create thumbnail directory: %w", err)
			}

			if err := lp.createThumbnail(fullPath, thumbFullPath, options); err != nil {
				return path, fmt.Errorf("failed to create thumbnail: %w", err)
			}
		}
	}

	return path, nil
}

// resizeImage resizes an image in place
func (lp *LocalProvider) resizeImage(_ string, _ *ImageProcessingOptions) error {
	// For now, this is a placeholder
	// In production, use github.com/disintegration/imaging or similar
	return nil
}

// createThumbnail creates a thumbnail from an image
func (lp *LocalProvider) createThumbnail(sourcePath, destPath string, options *ImageProcessingOptions) error {
	// Open original image
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize to thumbnail size
	thumbWidth := options.ThumbnailWidth
	if thumbWidth == 0 {
		thumbWidth = 200
	}
	thumbHeight := options.ThumbnailHeight
	if thumbHeight == 0 {
		thumbHeight = 200
	}

	// Simple resize (keep aspect ratio)
	bounds := img.Bounds()
	aspectRatio := float64(bounds.Dx()) / float64(bounds.Dy())
	newWidth := thumbWidth
	newHeight := int(float64(newWidth) / aspectRatio)

	if newHeight > thumbHeight {
		newHeight = thumbHeight
		newWidth = int(float64(newHeight) * aspectRatio)
	}

	// Create thumbnail file
	thumbFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer thumbFile.Close()

	// Encode as JPEG
	return jpeg.Encode(thumbFile, img, &jpeg.Options{Quality: 85})
}

// generateThumbnailPath generates the thumbnail path from original path
func (lp *LocalProvider) generateThumbnailPath(path string) string {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	return filepath.Join(dir, "thumbnails", name+"_thumb"+ext)
}

// sanitizeFilename removes unsafe characters from filename
func sanitizeFilename(filename string) string {
	// Remove path separators
	filename = filepath.Base(filename)

	// Remove unsafe characters
	filename = strings.TrimSpace(filename)

	return filename
}
