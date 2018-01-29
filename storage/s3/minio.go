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
