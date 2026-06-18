package port

import "context"

type StoragePort interface {
	UploadImage(ctx context.Context, bucketName, objectName string, data []byte, contentType string) (string, error)
}
