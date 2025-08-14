package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	minioCreds "github.com/minio/minio-go/v7/pkg/credentials"

	cfg "github.com/vyve/vyve-backend/internal/config"
)

// Storage defines the storage interface
type Storage interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
	GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
	ListObjects(ctx context.Context, prefix string) ([]string, error)
}

// S3Storage implements Storage using AWS S3
type S3Storage struct {
	client *s3.Client
	bucket string
	region string
}

// NewS3Storage creates a new S3 storage
func NewS3Storage(cfg cfg.AWSConfig) (Storage, error) {
	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Storage{
		client: client,
		bucket: cfg.S3Bucket,
		region: cfg.Region,
	}, nil
}

// Upload uploads data to S3
func (s *S3Storage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return s.GetURL(key), nil
}

// Download downloads data from S3
func (s *S3Storage) Download(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

// Delete deletes an object from S3
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// GetURL returns the public URL for an object
func (s *S3Storage) GetURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
}

// GeneratePresignedURL generates a presigned URL for an object
func (s *S3Storage) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := presigner.PresignGetObject(ctx, input, func(po *s3.PresignOptions) { po.Expires = expiration })
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return result.URL, nil
}

// ListObjects lists objects with a prefix
func (s *S3Storage) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	var keys []string
	paginator := s3.NewListObjectsV2Paginator(s.client, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, object := range output.Contents {
			keys = append(keys, *object.Key)
		}
	}

	return keys, nil
}

// MinIOStorage implements Storage using MinIO
type MinIOStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
	secure   bool
}

// NewMinIOStorage creates a new MinIO storage
func NewMinIOStorage(cfg cfg.StorageConfig) (Storage, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  minioCreds.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Check if bucket exists, create if not
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{
			Region: cfg.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinIOStorage{
		client:   client,
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		secure:   cfg.UseSSL,
	}, nil
}

// Upload uploads data to MinIO
func (m *MinIOStorage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	
	_, err := m.client.PutObject(ctx, m.bucket, key, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	return m.GetURL(key), nil
}

// Download downloads data from MinIO
func (m *MinIOStorage) Download(ctx context.Context, key string) ([]byte, error) {
	object, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}
	defer object.Close()

	return io.ReadAll(object)
}

// Delete deletes an object from MinIO
func (m *MinIOStorage) Delete(ctx context.Context, key string) error {
	err := m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from MinIO: %w", err)
	}
	return nil
}

// GetURL returns the public URL for an object
func (m *MinIOStorage) GetURL(key string) string {
	// For MinIO, construct the URL manually
	protocol := "http"
	if m.secure {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, m.endpoint, m.bucket, key)
}

// GeneratePresignedURL generates a presigned URL for an object
func (m *MinIOStorage) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucket, key, expiration, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// ListObjects lists objects with a prefix
func (m *MinIOStorage) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	
	objectCh := m.client.ListObjects(ctx, m.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}
		keys = append(keys, object.Key)
	}

	return keys, nil
}