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
	client           *s3.Client
	externalEndpoint string
	accessKey        string
	secretKey        string
}

// NewMinIOProvider inisialisasi client S3
func NewMinIOProvider(endpoint, accessKey, secretKey, region string) *MinIOProvider {
	return NewMinIOProviderWithExternal(endpoint, accessKey, secretKey, region, endpoint)
}

// NewMinIOProviderWithExternal inisialisasi client S3 dengan endpoint terpisah untuk URL eksternal
func NewMinIOProviderWithExternal(internalEndpoint, accessKey, secretKey, region, externalEndpoint string) *MinIOProvider {
	// Konfigurasi untuk internal communication (backend to MinIO)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           internalEndpoint, // Internal endpoint for backend operations
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

	return &MinIOProvider{
		client:           client,
		externalEndpoint: externalEndpoint,
		accessKey:        accessKey,
		secretKey:        secretKey,
	}
}

// GetPresignedPutURL generates URL upload langsung (untuk frontend)
func (m *MinIOProvider) GetPresignedPutURL(ctx context.Context, bucket, key, contentType string, expiry time.Duration) (string, error) {
	// Buat client terpisah dengan external endpoint untuk URL yang bisa diakses frontend
	externalCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // Region tidak penting untuk presigned URL
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(m.accessKey, m.secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           m.externalEndpoint, // External endpoint for frontend access
					SigningRegion: "us-east-1",
				}, nil
			},
		)),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create external config: %w", err)
	}

	externalClient := s3.NewFromConfig(externalCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(externalClient)

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
