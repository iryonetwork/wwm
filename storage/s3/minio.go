package s3

import (
	"context"
	"io"

	minio "github.com/minio/minio-go"
)

// iminio provides a wrapper for the GetObject method
type iminio struct {
	minio.Client
}

func (m iminio) GetObjectWithContext(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
	return m.Client.GetObjectWithContext(ctx, bucketName, objectName, opts)
}

// BucketExists is a wrapper for original BucketExists method to fix incorrect handling of 404 returned as response
func (m iminio) BucketExists(bucketName string) (bool, error) {
	exists, err := m.Client.BucketExists(bucketName)
	if err != nil && minio.ToErrorResponse(err).Code == "NoSuchBucket" {
		return false, nil
	}
	return exists, err
}
