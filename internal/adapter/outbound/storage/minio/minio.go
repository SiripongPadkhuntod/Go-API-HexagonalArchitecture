package minio

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"hexagonalarchitecture/internal/core/port"
	"hexagonalarchitecture/internal/infrastructure/config"
)

var _ port.StoragePort = (*MinioStorage)(nil)

type MinioStorage struct {
	client *minio.Client
	config config.StorageConfig
}

func NewMinioStorage(cfg config.StorageConfig) (*MinioStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Make bucket if not exists
	ctx := context.Background()
	err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, cfg.BucketName)
		if errBucketExists == nil && exists {
			// Bucket already exists
		} else {
			return nil, err
		}
	}

	// Set bucket policy to public read
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, cfg.BucketName)

	err = client.SetBucketPolicy(ctx, cfg.BucketName, policy)
	if err != nil {
		return nil, err
	}

	return &MinioStorage{client: client, config: cfg}, nil
}

func (s *MinioStorage) UploadImage(ctx context.Context, bucketName, objectName string, data []byte, contentType string) (string, error) {
	if bucketName == "" {
		bucketName = s.config.BucketName
	}
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	protocol := "http"
	if s.config.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, s.config.Endpoint, bucketName, objectName), nil
}
