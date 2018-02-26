/*
Package s3 provides an abstraction layer for file management on S3 compatible
storage. Minio's go library is used to provide basic S3 compatibility.

Features

S3 supports following features:

    - listing files inside a bucket
    - creating new files
    - reading files
    - encrypting all files using an external key provider

Encryption

To support encryption s3 requires an external key provider that can provide the
storage correct key for the current bucket / user ID.

Storing metadata

Metadata is stored inside the file name. The end file name on S3 storage will
look like this

	FILENAME.VERSION.OPERATION.TIMESTAMP.CHECKSUM.ARCHETYPE
	-- 40 --.- 1-40-.--- 1 ---.-- 13 ---.-- 44 --.--- * ---

Filenames on S3 are limited to around 1024 bytes meaning that the last archetype
value can be up to 886 characters long.

@TODO How to add new values to file name?
*/
package s3

//go:generate sh ../../bin/mockgen.sh storage/s3 Storage,KeyProvider,Minio $GOFILE

import (
	"fmt"
	"io"
	"regexp"
	"sort"

	"github.com/go-openapi/strfmt"
	"github.com/minio/minio-go/pkg/encrypt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3/object"
	"github.com/minio/minio-go"
)

const metaArchetype = "x-archetype"
const metaCreated = "x-created"
const metaChecksum = "x-checksum"

// Storage provides an interface for s3 public functions
type Storage interface {
	BucketExists(bucketID string) (bool, error)
	MakeBucket(bucketID string) error
	ListBuckets() ([]*models.BucketDescriptor, error)
	List(bucketID, prefix string) ([]*models.FileDescriptor, error)
	Read(bucketID, fileID, version string) (io.ReadCloser, *models.FileDescriptor, error)
	Write(bucketID string, newFile *object.NewObjectInfo, r io.Reader) (*models.FileDescriptor, error)
}

// KeyProvider lists methods required for reading encryption keys
type KeyProvider interface {
	Get(string) (string, error)
	Read(b []byte) (n int, err error)
}

// Minio interface describes functions used in minio-go package for mocking
// purposes.
type Minio interface {
	MakeBucket(bucketName, location string) error
	BucketExists(bucketName string) (bool, error)
	ListBuckets() ([]minio.BucketInfo, error)
	ListObjectsV2(bucketName, prefix string, recursive bool, doneCh <-chan struct{}) <-chan minio.ObjectInfo
	GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (object.Object, error)
	GetEncryptedObject(bucketName, objectName string, encryptMaterials encrypt.Materials) (io.ReadCloser, error)
	PutObject(bucketName, objectName string, reader io.Reader, objectSize int64,
		opts minio.PutObjectOptions) (n int64, err error)
	PutEncryptedObject(bucketName, objectName string, reader io.Reader, encryptMaterials encrypt.Materials) (n int64, err error)
}

var nameVersionRE = regexp.MustCompile("^(.*)\\.(\\d+)$")

// Config holds all details required to connect to an S3 storage
type Config struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	Secure       bool
	Region       string
}

type s3storage struct {
	cfg    *Config
	client Minio
	keys   KeyProvider
	logger zerolog.Logger
}

// Operation represents a single character operation
type Operation string

// Write represents write operation
const Write Operation = Operation(models.FileDescriptorOperationW)

// Delete represents read operation
const Delete Operation = Operation(models.FileDescriptorOperationD)

// ErrAlreadyExists indicates bucket or file already exists
var ErrAlreadyExists = errors.New("Item already exists")

// ErrNotFound indicates file or bucket were not found
var ErrNotFound = errors.New("File not found")

// ErrDeleted indicates file or bucket were already deleted
var ErrDeleted = errors.New("File was deleted")

// New creates a new instance of s3 storage
func New(cfg *Config, keys KeyProvider, logger zerolog.Logger) (Storage, error) {
	c, err := minio.NewWithRegion(cfg.Endpoint, cfg.AccessKey, cfg.AccessSecret, cfg.Secure, cfg.Region)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize minio with region")
	}

	obj := &s3storage{
		cfg:    cfg,
		client: iminio{*c},
		keys:   keys,
		logger: logger,
	}

	return obj, nil
}

// Check if bycket already exits
func (s *s3storage) BucketExists(bucketID string) (bool, error) {
	s.logger.Debug().Str("cmd", "s3::BucketExists").Msgf("('%s')", bucketID)

	exists, err := s.client.BucketExists(bucketID)
	if err != nil {
		return false, errors.Wrap(err, "Failed to check if bucket exists")
	}

	return exists, nil
}

// MakeBucket creates a bucket, return ErrAlreadyExists if bucket already exists
func (s *s3storage) MakeBucket(bucketID string) error {
	s.logger.Debug().Str("cmd", "s3::MakeBucket").Msgf("('%s')", bucketID)

	exists, err := s.client.BucketExists(bucketID)
	if err != nil {
		return errors.Wrap(err, "Failed to check if bucket exists")
	}
	if exists {
		return ErrAlreadyExists
	}

	if !exists {
		if err := s.client.MakeBucket(bucketID, s.cfg.Region); err != nil {
			return errors.Wrap(err, "Failed to create a new bucket")
		}
	}

	return nil

}

// ListBuckets returns a list of buckets
func (s *s3storage) ListBuckets() ([]*models.BucketDescriptor, error) {
	s.logger.Debug().Str("cmd", "s3::ListBuckets")

	b, err := s.client.ListBuckets()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to list buckets")
	}

	buckets := []*models.BucketDescriptor{}
	for _, info := range b {
		bd, err := bucketInfoToBucketDescriptor(info)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to convert bucketInfo to bucketDescriptor")
		}
		buckets = append(buckets, bd)
	}

	return buckets, nil
}

// List returns a list of files stored inside a bucket
func (s *s3storage) List(bucketID, prefix string) ([]*models.FileDescriptor, error) {
	s.logger.Debug().Str("cmd", "s3::List").Msgf("('%s', '%s')", bucketID, prefix)

	// Check if bucket exists first
	exists, err := s.client.BucketExists(bucketID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to check if bucket exists")
	}
	if !exists {
		// Nothing to list
		return []*models.FileDescriptor{}, nil
	}

	ch := make(chan struct{})
	defer close(ch)
	infos := s.client.ListObjectsV2(bucketID, prefix, false, ch)

	files := []*models.FileDescriptor{}
	for info := range infos {
		if info.Err != nil {
			return nil, errors.Wrap(info.Err, "Failed to read object from a list")
		}

		fd, err := objectInfoToFileDescriptor(info, bucketID)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to convert object to fileDescriptor")
		}

		files = append(files, fd)
	}

	sort.Sort(byCreated(files))
	return files, nil
}

// Read fetches contents from the storage
func (s *s3storage) Read(bucketID, fileID, version string) (io.ReadCloser, *models.FileDescriptor, error) {
	s.logger.Debug().Str("cmd", "s3::Read").Msgf("('%s', '%s', '%s')", bucketID, fileID, version)

	// find the file
	prefix := fmt.Sprintf("%s.", fileID)
	if version != "" {
		prefix += fmt.Sprintf("%s.", version)
	}
	list, err := s.List(bucketID, prefix)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to list files")
	}
	if len(list) == 0 {
		return nil, nil, ErrNotFound
	}
	md, err := metadataFromFileDescriptor(list[0])
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to parse metadata from fileDescriptor")
	}

	// read the key
	em, err := getCBCKey(bucketID, s.keys)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to set CBC key")
	}

	// fetch the file
	reader, err := s.client.GetEncryptedObject(bucketID, md.String(), em)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to fetch enc. object")
	}

	return reader, list[0], nil
}

// Write creates a new file in the storage
func (s *s3storage) Write(bucketID string, newFile *object.NewObjectInfo, r io.Reader) (*models.FileDescriptor, error) {
	s.logger.Debug().Str("cmd", "s3::Write").Msgf("('%s', '%+v', reader)", bucketID, newFile)

	// validate operation
	op := Operation(newFile.Operation)
	if op != Write && op != Delete {
		return nil, fmt.Errorf("Received an invalid operation '%s'", op)
	}

	// get the key
	em, err := getCBCKey(bucketID, s.keys)
	if err != nil {
		s.logger.Info().Err(err).Msg("Failed to set the CBC key")
	}

	// // compose the put options
	// opts := minio.PutObjectOptions{
	// 	ContentType:      newFile.ContentType,
	// 	EncryptMaterials: em,
	// }

	// collect meta data
	meta, err := metadataFromNewFile(newFile)
	if err != nil {
		s.logger.Info().Err(err).Msg("Failed to collect metadata from new file")
	}

	// upload the file
	_, err = s.client.PutEncryptedObject(bucketID, meta.String(), r, em)
	if err != nil {
		s.logger.Info().Err(err).Msg("Failed to call PutObject")
	}

	// generate the file descriptor
	fd := &models.FileDescriptor{
		Name:        newFile.Name,
		Version:     newFile.Version,
		Archetype:   newFile.Archetype,
		ContentType: newFile.ContentType,
		Checksum:    newFile.Checksum,
		Created:     newFile.Created,
		Path:        fmt.Sprintf("%s/%s/%s", bucketID, meta.filename, meta.version),
		Size:        newFile.Size,
		Operation:   string(op),
	}

	return fd, nil
}

func objectInfoToFileDescriptor(info minio.ObjectInfo, bucketID string) (*models.FileDescriptor, error) {
	meta, err := metadataFromKey(info.Key)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to extract metadata from key")
	}

	// copy basic data
	fd := &models.FileDescriptor{
		Size:        info.Size,
		ContentType: info.ContentType,
		Path:        fmt.Sprintf("%s/%s/%s", bucketID, meta.filename, meta.version),
		Name:        meta.filename,
		Version:     meta.version,
		Checksum:    meta.checksum,
		Created:     strfmt.DateTime(meta.created),
		Archetype:   meta.archetype,
		Operation:   string(meta.operation),
	}

	return fd, nil
}

func bucketInfoToBucketDescriptor(info minio.BucketInfo) (*models.BucketDescriptor, error) {
	// copy
	bd := &models.BucketDescriptor{
		Name:    info.Name,
		Created: strfmt.DateTime(info.CreationDate),
	}

	return bd, nil
}

func getCBCKey(bucketID string, keys KeyProvider) (encrypt.Materials, error) {
	// read the key
	secret, err := keys.Get(bucketID)
	if err != nil {
		return nil, err
	}

	// create the materials
	return encrypt.NewCBCSecureMaterials(encrypt.NewSymmetricKey([]byte(secret)))
}
