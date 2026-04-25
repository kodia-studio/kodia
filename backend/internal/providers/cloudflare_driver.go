package providers

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

// CloudflareProvider implements StorageProvider for Cloudflare R2
type CloudflareProvider struct {
	config     *StorageConfig
	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	validator  *FileValidator
	bucket     string
	basePath   string
	domain     string
}

// NewCloudflareProvider creates a new Cloudflare R2 storage provider
func NewCloudflareProvider(config *StorageConfig) (*CloudflareProvider, error) {
	if config.Cloudflare == nil {
		return nil, fmt.Errorf("Cloudflare config is required")
	}

	cf := config.Cloudflare

	if cf.AccountID == "" || cf.AccessKey == "" || cf.SecretKey == "" {
		return nil, fmt.Errorf("Cloudflare AccountID, AccessKey, and SecretKey are required")
	}

	// Create S3-compatible client for Cloudflare R2
	credProvider := credentials.NewStaticCredentialsProvider(
		cf.AccessKey,
		cf.SecretKey,
		"",
	)

	// Cloudflare R2 endpoint
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cf.AccountID)

	// Create AWS config
	awsCfg := aws.NewConfig()
	awsCfg.Credentials = credProvider
	awsCfg.Logger = logging.NewStandardLogger(nil)

	// Create R2 client
	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Credentials:  credProvider,
		Region:       "auto",
	})

	uploader := manager.NewUploader(client)
	downloader := manager.NewDownloader(client)

	validator := &FileValidator{
		MaxSize:      config.MaxFileSize,
		AllowedTypes: config.AllowedMimeTypes,
		AllowedExts:  config.AllowedExtensions,
	}

	basePath := config.Cloudflare.Path
	if basePath == "" {
		basePath = "uploads"
	}

	domain := config.Cloudflare.Domain
	if domain == "" {
		domain = fmt.Sprintf("%s.r2.cloudflarestorage.com", cf.AccountID)
	}

	return &CloudflareProvider{
		config:     config,
		client:     client,
		uploader:   uploader,
		downloader: downloader,
		validator:  validator,
		bucket:     config.Cloudflare.Bucket,
		basePath:   basePath,
		domain:     domain,
	}, nil
}

// Store uploads a file to Cloudflare R2
func (cp *CloudflareProvider) Store(ctx context.Context, file *FileUpload) (string, error) {
	if err := cp.validator.ValidateFile(file); err != nil {
		return "", err
	}

	// Generate key
	now := time.Now()
	filename := sanitizeFilename(file.Filename)
	key := filepath.Join(
		cp.basePath,
		fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day()),
		fmt.Sprintf("%d-%s", now.UnixNano(), filename),
	)

	// Upload to R2
	_, err := cp.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(cp.bucket),
		Key:         aws.String(key),
		Body:        file.Content,
		ContentType: aws.String(file.ContentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudflare R2: %w", err)
	}

	// Return key for later retrieval
	return key, nil
}

// StoreMultiple uploads multiple files to R2
func (cp *CloudflareProvider) StoreMultiple(ctx context.Context, files []*FileUpload) ([]string, error) {
	var keys []string

	for _, file := range files {
		key, err := cp.Store(ctx, file)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// Get retrieves file content from R2
func (cp *CloudflareProvider) Get(ctx context.Context, path string) ([]byte, error) {
	result, err := cp.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(cp.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object from R2: %w", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

// GetStream retrieves file as stream from R2
func (cp *CloudflareProvider) GetStream(ctx context.Context, path string) (io.Reader, error) {
	result, err := cp.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(cp.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object from R2: %w", err)
	}

	return result.Body, nil
}

// Delete removes a file from R2
func (cp *CloudflareProvider) Delete(ctx context.Context, path string) error {
	_, err := cp.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(cp.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object from R2: %w", err)
	}

	return nil
}

// DeleteMultiple removes multiple files from R2
func (cp *CloudflareProvider) DeleteMultiple(ctx context.Context, paths []string) error {
	for _, path := range paths {
		if err := cp.Delete(ctx, path); err != nil {
			return err
		}
	}
	return nil
}

// Exists checks if file exists in R2
func (cp *CloudflareProvider) Exists(ctx context.Context, path string) (bool, error) {
	_, err := cp.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(cp.bucket),
		Key:    aws.String(path),
	})

	if err == nil {
		return true, nil
	}

	// Check if 404 error
	if err.Error() == "NotFound" {
		return false, nil
	}

	return false, err
}

// GetURL returns the public URL for a file in R2
func (cp *CloudflareProvider) GetURL(ctx context.Context, path string) (string, error) {
	if cp.config.Cloudflare.PublicURL {
		return fmt.Sprintf("https://%s/%s", cp.domain, path), nil
	}

	// For private files, return a temporary URL
	return cp.GetTemporaryURL(ctx, path, 24*time.Hour)
}

// GetTemporaryURL returns a temporary signed URL
func (cp *CloudflareProvider) GetTemporaryURL(ctx context.Context, path string, duration time.Duration) (string, error) {
	// Cloudflare R2 doesn't have built-in pre-signed URLs like S3
	// For public buckets, return the public URL
	// For private buckets, you would need to use Cloudflare Workers or another solution
	return cp.GetURL(ctx, path)
}

// Copy copies a file within R2
func (cp *CloudflareProvider) Copy(ctx context.Context, sourcePath, destinationPath string) error {
	copySource := fmt.Sprintf("%s/%s", cp.bucket, sourcePath)

	_, err := cp.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(cp.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destinationPath),
	})

	if err != nil {
		return fmt.Errorf("failed to copy object in R2: %w", err)
	}

	return nil
}

// Close closes the Cloudflare provider
func (cp *CloudflareProvider) Close() error {
	return nil
}

// StoreImage stores an image with optional processing
func (cp *CloudflareProvider) StoreImage(ctx context.Context, file *FileUpload, options *ImageProcessingOptions) (string, error) {
	// For Cloudflare R2, consider using Cloudflare Image Optimization
	// For now, just store the image as-is
	return cp.Store(ctx, file)
}
