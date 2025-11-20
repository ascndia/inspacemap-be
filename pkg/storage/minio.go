package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MinIOProvider struct {
	client *s3.Client
}

// NewMinIOProvider inisialisasi client S3
func NewMinIOProvider(endpoint, accessKey, secretKey, region string) *MinIOProvider {
	// Konfigurasi Custom untuk MinIO/Localstack
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint, // e.g. http://localhost:9000
					SigningRegion: region,
				}, nil
			},
		)),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// ForcePathStyle = true wajib untuk MinIO
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &MinIOProvider{client: client}
}

// GetPresignedPutURL generates URL upload langsung
func (m *MinIOProvider) GetPresignedPutURL(ctx context.Context, bucket, key, contentType string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(m.client)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})

	if err != nil {
		return "", fmt.Errorf("failed to presign request: %w", err)
	}

	return req.URL, nil
}

// DeleteObject hapus file fisik
func (m *MinIOProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := m.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
