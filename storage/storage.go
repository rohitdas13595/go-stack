package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Disk is local filesystem storage.
type Disk struct {
	Root string
}

func (d *Disk) Put(ctx context.Context, path string, r io.Reader, _ int64) (string, error) {
	full := filepath.Join(d.Root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return full, nil
}

func (d *Disk) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	full := filepath.Join(d.Root, filepath.FromSlash(path))
	return os.Open(full)
}

func (d *Disk) Delete(ctx context.Context, path string) error {
	full := filepath.Join(d.Root, filepath.FromSlash(path))
	return os.Remove(full)
}

// S3Store stores objects in S3-compatible bucket.
type S3Store struct {
	Client *s3.Client
	Bucket string
	PublicBaseURL string
}

// NewS3FromEnv builds client from AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_REGION.
func NewS3FromEnv(ctx context.Context, bucket, region, endpoint string) (*S3Store, error) {
	var opts []func(*config.LoadOptions) error
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}
	if endpoint != "" {
		opts = append(opts, config.WithBaseEndpoint(endpoint))
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})
	return &S3Store{Client: cli, Bucket: bucket}, nil
}

// NewS3Static creates S3 client with static credentials (e.g. MinIO).
func NewS3Static(ctx context.Context, bucket, region, endpoint, key, secret string) (*S3Store, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, "")),
	)
	if err != nil {
		return nil, err
}
	cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	return &S3Store{Client: cli, Bucket: bucket}, nil
}

func (s *S3Store) Put(ctx context.Context, key string, body io.Reader, size int64) (string, error) {
	if s.Client == nil {
		return "", fmt.Errorf("storage: nil client")
	}
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return "", err
	}
	if s.PublicBaseURL != "" {
		return fmt.Sprintf("%s/%s", stringsTrimSlash(s.PublicBaseURL), key), nil
	}
	return "s3://" + s.Bucket + "/" + key, nil
}

func stringsTrimSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

func (s *S3Store) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (s *S3Store) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

// SignedURL is a placeholder; production should use PresignGetObject.
func (s *S3Store) SignedURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	_ = ttl
	return "", fmt.Errorf("storage: SignedURL not implemented")
}
