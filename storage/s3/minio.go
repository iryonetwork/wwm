package s3

import (
	"github.com/iryonetwork/wwm/storage/s3/object"
	minio "github.com/minio/minio-go"
)

// iminio provides a wrapper for the GetObject method
type iminio struct {
	minio.Client
}

func (m iminio) GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (object.Object, error) {
	return m.Client.GetObject(bucketName, objectName, opts)
}

// BucketExists is a wrapper for original BucketExists method to fix incorrect handling of 404 returned as response
func (m iminio) BucketExists(bucketName string) (bool, error) {
	exists, err := m.Client.BucketExists(bucketName)
	if err != nil && minio.ToErrorResponse(err).Code == "NoSuchBucket" {
		return false, nil
	}
	return exists, err
}
