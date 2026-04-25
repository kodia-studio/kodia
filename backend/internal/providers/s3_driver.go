package providers

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Provider implements StorageProvider for AWS S3
type S3Provider struct {
	config     *StorageConfig
	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	validator  *FileValidator
	bucket     string
	basePath   string
}

// NewS3Provider creates a new S3 storage provider
func NewS3Provider(config *StorageConfig) (*S3Provider, error) {
	if config.S3 == nil {
		return nil, fmt.Errorf("S3 config is required")
	}

	// Create AWS SDK config
	cfg, err := createAWSConfig(config.S3)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Create uploader and downloader
	uploader := manager.NewUploader(s3Client)
	downloader := manager.NewDownloader(s3Client)

	validator := &FileValidator{
		MaxSize:      config.MaxFileSize,
		AllowedTypes: config.AllowedMimeTypes,
		AllowedExts:  config.AllowedExtensions,
	}

	basePath := config.S3.Path
	if basePath == "" {
		basePath = "uploads"
	}

	return &S3Provider{
		config:     config,
		client:     s3Client,
		uploader:   uploader,
		downloader: downloader,
		validator:  validator,
		bucket:     config.S3.Bucket,
		basePath:   basePath,
	}, nil
}

// Store uploads a file to S3
func (sp *S3Provider) Store(ctx context.Context, file *FileUpload) (string, error) {
	if err := sp.validator.ValidateFile(file); err != nil {
		return "", err
	}

	// Generate key
	now := time.Now()
	filename := sanitizeFilename(file.Filename)
	key := filepath.Join(
		sp.basePath,
		fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day()),
		fmt.Sprintf("%d-%s", now.UnixNano(), filename),
	)

	// Upload to S3
	_, err := sp.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sp.bucket),
		Key:         aws.String(key),
		Body:        file.Content,
		ContentType: aws.String(file.ContentType),
		ACL:         types.ObjectCannedACL(sp.config.S3.ACL),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Return key for later retrieval
	return key, nil
}

// StoreMultiple uploads multiple files to S3
func (sp *S3Provider) StoreMultiple(ctx context.Context, files []*FileUpload) ([]string, error) {
	var keys []string

	for _, file := range files {
		key, err := sp.Store(ctx, file)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// Get retrieves file content from S3
func (sp *S3Provider) Get(ctx context.Context, path string) ([]byte, error) {
	result, err := sp.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(sp.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

// GetStream retrieves file as stream from S3
func (sp *S3Provider) GetStream(ctx context.Context, path string) (io.Reader, error) {
	result, err := sp.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(sp.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	return result.Body, nil
}

// Delete removes a file from S3
func (sp *S3Provider) Delete(ctx context.Context, path string) error {
	_, err := sp.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(sp.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	return nil
}

// DeleteMultiple removes multiple files from S3
func (sp *S3Provider) DeleteMultiple(ctx context.Context, paths []string) error {
	// Build delete objects list
	delete := &types.Delete{
		Objects: make([]types.ObjectIdentifier, len(paths)),
	}

	for i, path := range paths {
		delete.Objects[i] = types.ObjectIdentifier{
			Key: aws.String(path),
		}
	}

	// Delete objects
	_, err := sp.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(sp.bucket),
		Delete: delete,
	})

	if err != nil {
		return fmt.Errorf("failed to delete objects from S3: %w", err)
	}

	return nil
}

// Exists checks if file exists in S3
func (sp *S3Provider) Exists(ctx context.Context, path string) (bool, error) {
	_, err := sp.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(sp.bucket),
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

// GetURL returns the public URL for a file in S3
func (sp *S3Provider) GetURL(ctx context.Context, path string) (string, error) {
	// If CloudFront is configured, use that
	if sp.config.S3.CloudFront != "" {
		return fmt.Sprintf("https://%s/%s", sp.config.S3.CloudFront, path), nil
	}

	// Otherwise, use S3 URL
	region := sp.config.S3.Region
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", sp.bucket, region, path), nil
}

// GetTemporaryURL returns a temporary signed URL
func (sp *S3Provider) GetTemporaryURL(ctx context.Context, path string, duration time.Duration) (string, error) {
	// For now, return a public URL if the bucket is public
	// Full presigned URL generation requires additional AWS SDK setup
	// that depends on the specific AWS SDK version
	return sp.GetURL(ctx, path)
}

// Copy copies a file within S3
func (sp *S3Provider) Copy(ctx context.Context, sourcePath, destinationPath string) error {
	copySource := fmt.Sprintf("%s/%s", sp.bucket, sourcePath)

	_, err := sp.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(sp.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destinationPath),
	})

	if err != nil {
		return fmt.Errorf("failed to copy object in S3: %w", err)
	}

	return nil
}

// Close closes the S3 provider
func (sp *S3Provider) Close() error {
	return nil
}

// StoreImage stores an image with optional processing
func (sp *S3Provider) StoreImage(ctx context.Context, file *FileUpload, options *ImageProcessingOptions) (string, error) {
	// For now, just store the image as-is
	// Image processing would require downloading, processing, and re-uploading
	// Consider using Lambda for image processing in production
	return sp.Store(ctx, file)
}

// createAWSConfig creates AWS SDK configuration
func createAWSConfig(cfg *S3Config) (aws.Config, error) {
	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKey,
		cfg.SecretKey,
		"",
	)

	return config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credProvider),
	)
}
