/*
Object is a small package providing an interface for minio.Object to ease
testing. Interface is located in its own package due to issues with import
cycles in unit tests.
*/
package object

//go:generate sh ../../../bin/mockgen.sh storage/s3/object Object $GOFILE

import (
	"github.com/go-openapi/strfmt"

	minio "github.com/minio/minio-go"
)

// Object holds functions used on minio.Object
type Object interface {
	Stat() (minio.ObjectInfo, error)
	Read(b []byte) (n int, err error)
}

// NewObjectInfo holds data required to create a new file
type NewObjectInfo struct {
	Archetype   string
	Checksum    string
	Size        int64
	ContentType string
	Created     strfmt.DateTime
	Name        string
	Version     string
	Operation   string
}
