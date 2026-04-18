package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	kodia_config "github.com/kodia-studio/kodia/pkg/config"
)

// S3StorageProvider implements ports.StorageProvider using AWS S3 or MinIO.
type S3StorageProvider struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

// NewS3StorageProvider creates a new S3StorageProvider.
func NewS3StorageProvider(cfg *kodia_config.Config) (*S3StorageProvider, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.Storage.Endpoint != "" {
			return aws.Endpoint{
				URL:           cfg.Storage.Endpoint,
				SigningRegion: cfg.Storage.Region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Storage.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Storage.AccessID,
			cfg.Storage.SecretKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Storage.Endpoint != "" {
			o.UsePathStyle = true // Required for MinIO
		}
	})

	publicURL := cfg.Storage.PublicURL
	if publicURL == "" {
		if cfg.Storage.Endpoint != "" {
			publicURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(cfg.Storage.Endpoint, "/"), cfg.Storage.Bucket)
		} else {
			publicURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.Storage.Bucket, cfg.Storage.Region)
		}
	}

	return &S3StorageProvider{
		client:    client,
		bucket:    cfg.Storage.Bucket,
		publicURL: publicURL,
	}, nil
}

func (p *S3StorageProvider) Upload(ctx context.Context, path string, content io.Reader) (string, error) {
	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
		Body:   content,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}
	return path, nil
}

func (p *S3StorageProvider) Delete(ctx context.Context, path string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}
	return nil
}

func (p *S3StorageProvider) GetURL(ctx context.Context, path string) (string, error) {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(p.publicURL, "/"), path), nil
}

func (p *S3StorageProvider) Exists(ctx context.Context, path string) (bool, error) {
	_, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		// In AWS SDK v2, we check the error type or message for 404
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
