package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"github.com/agext/uuid"
	"github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/object"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

// Service describes storage service's public API.
type Service interface {
	// Checksum calculates the checksum of a given reader using sha256.Sum.
	Checksum(r io.Reader) (string, error)

	// BucketList returns list of all the buckets.
	BucketList() ([]*models.BucketDescriptor, error)

	// FileList returns a list of latest versions of files. Older versions and
	// files marked as deleted are removed from the list.
	FileList(bucketID string) ([]*models.FileDescriptor, error)

	// FileGet returns the latest version of the file by returning the reader
	// and file details.
	FileGet(bucketID, fileID string) (io.ReadCloser, *models.FileDescriptor, error)

	// FileGetVersion returns a specific version of a file.
	FileGetVersion(bucketID, fileID, version string) (io.ReadCloser, *models.FileDescriptor, error)

	// FileListVersions returns a list of all modifications to a file.
	FileListVersions(bucketID, fileID string) ([]*models.FileDescriptor, error)

	// FileNew creates a new file.
	FileNew(bucketID string, r io.Reader, contentType string, archetype string) (*models.FileDescriptor, error)

	// FileUpdate creates a new version of a file.
	FileUpdate(bucketID, fileID string, r io.Reader, contentType string, archetype string) (*models.FileDescriptor, error)

	// FileDelete marks file as deleted.
	FileDelete(bucketID, fileID string) error

	// SyncFileList returns a list of latest versions of files. Older versions are removed from the list.
	// Files marked as deleted are kept in the list.
	SyncFileList(bucketID string) ([]*models.FileDescriptor, error)

	// SyncFile syncs file with provided fileID and version.
	SyncFile(bucketID, fileID, version string, r io.Reader, contentType string, created strfmt.DateTime, archetype string) (*models.FileDescriptor, error)

	// SyncFileDelete sync file deletion.
	SyncFileDelete(bucketID, fileID, version string, created strfmt.DateTime) error
}

// Bucket or item was already deleted
var ErrDeleted = s3.ErrDeleted

// Item does not exists
var ErrNotFound = s3.ErrNotFound

// Item already exists
var ErrAlreadyExists = s3.ErrAlreadyExists

// Item already exists and conflicts
var ErrAlreadyExistsConflict = errors.New("Item already exists and its checksum is different")

type service struct {
	s3          s3.Storage
	keyProvider s3.KeyProvider
	publisher   storageSync.Publisher
	logger      zerolog.Logger
}

func (s *service) Checksum(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

func (s *service) BucketList() ([]*models.BucketDescriptor, error) {
	// get the list and return
	return s.s3.ListBuckets()
}

func (s *service) FileList(bucketID string) ([]*models.FileDescriptor, error) {
	// make sure bucket exists
	if err := s.s3.MakeBucket(bucketID); err != nil && err != s3.ErrAlreadyExists {
		s.logger.Info().Err(err).Str("bucket", bucketID).Msg("Failed to ensure bucket")
		return nil, err
	}

	// collect the list
	list, err := s.s3.List(bucketID, "")
	if err != nil {
		return nil, err
	}

	// extract only latest versions in a map; latest version is already sorted
	// on top
	m := map[string]*models.FileDescriptor{}
	for _, f := range list {
		if _, ok := m[f.Name]; !ok {
			m[f.Name] = f
		}
	}

	// extract a list out of a map; only include files with a write operation
	list = []*models.FileDescriptor{}
	for _, f := range m {
		if s3.Operation(f.Operation) == s3.Write {
			list = append(list, f)
		}
	}

	return list, nil
}

func (s *service) FileGet(bucketID, fileID string) (io.ReadCloser, *models.FileDescriptor, error) {
	return s.s3.Read(bucketID, fileID, "")
}

func (s *service) FileGetVersion(bucketID, fileID, version string) (io.ReadCloser, *models.FileDescriptor, error) {
	return s.s3.Read(bucketID, fileID, version)
}

func (s *service) FileListVersions(bucketID, fileID string) ([]*models.FileDescriptor, error) {
	return s.s3.List(bucketID, fileID)
}

func (s *service) FileNew(bucketID string, r io.Reader, contentType string, archetype string) (*models.FileDescriptor, error) {
	// make sure bucket exists
	if err := s.s3.MakeBucket(bucketID); err != nil && err != s3.ErrAlreadyExists {
		s.logger.Info().Err(err).Str("bucket", bucketID).Msg("Failed to ensure bucket")
		return nil, err
	}

	// calculate the checksum
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Info().Err(err).Msg("Failed to calculate checksum")
		return nil, err
	}

	fileID := getUUID()
	version := getUUID()
	no := &object.NewObjectInfo{
		Archetype:   archetype,
		Size:        int64(buf.Len()),
		Checksum:    checksum,
		Created:     getTime(),
		ContentType: contentType,
		Version:     version,
		Name:        fileID,
		Operation:   string(s3.Write),
	}

	fd, err := s.s3.Write(bucketID, no, &buf)
	if err == nil {
		s.publisher.PublishAsyncWithRetries(
			context.Background(),
			storageSync.FileNew,
			&storageSync.FileInfo{BucketID: bucketID, FileID: fileID, Version: version},
		)
	}
	return fd, err
}

func (s *service) FileUpdate(bucketID, fileID string, r io.Reader, contentType string, archetype string) (*models.FileDescriptor, error) {
	// get the previous file
	_, _, err := s.s3.Read(bucketID, fileID, "")
	if err != nil {
		return nil, err
	}

	// calculate the checksum
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Info().Err(err).Msg("Failed to calculate checksum")
		return nil, err
	}

	version := getUUID()
	no := &object.NewObjectInfo{
		Archetype:   archetype,
		Checksum:    checksum,
		Size:        int64(buf.Len()),
		Created:     getTime(),
		ContentType: contentType,
		Version:     version,
		Name:        fileID,
		Operation:   string(s3.Write),
	}

	fd, err := s.s3.Write(bucketID, no, &buf)
	if err == nil {
		s.publisher.PublishAsyncWithRetries(
			context.Background(),
			storageSync.FileUpdate,
			&storageSync.FileInfo{BucketID: bucketID, FileID: fileID, Version: version},
		)
	}
	return fd, err
}

func (s *service) FileDelete(bucketID, fileID string) error {
	// get the previous file
	_, fd, err := s.s3.Read(bucketID, fileID, "")
	if err != nil {
		return err
	}

	version := getUUID()
	no := &object.NewObjectInfo{
		Archetype:   fd.Archetype,
		Checksum:    "",
		Size:        0,
		Created:     getTime(),
		ContentType: fd.ContentType,
		Version:     version,
		Name:        fileID,
		Operation:   string(s3.Delete),
	}

	_, err = s.s3.Write(bucketID, no, &bytes.Buffer{})
	if err == nil {
		s.publisher.PublishAsyncWithRetries(
			context.Background(),
			storageSync.FileDelete,
			&storageSync.FileInfo{BucketID: bucketID, FileID: fileID, Version: version},
		)
	}
	return err
}

func (s *service) SyncFileList(bucketID string) ([]*models.FileDescriptor, error) {
	// make sure bucket exists
	if err := s.s3.MakeBucket(bucketID); err != nil && err != s3.ErrAlreadyExists {
		s.logger.Info().Err(err).Str("bucket", bucketID).Msg("Failed to ensure bucket")
		return nil, err
	}

	// collect the list
	list, err := s.s3.List(bucketID, "")
	if err != nil {
		return nil, err
	}

	// extract only latest versions in a map; latest version is already sorted
	// on top
	m := map[string]*models.FileDescriptor{}
	for _, f := range list {
		if _, ok := m[f.Name]; !ok {
			m[f.Name] = f
		}
	}

	// extract a list out of a map;
	list = []*models.FileDescriptor{}
	for _, f := range m {
		list = append(list, f)
	}

	return list, nil
}

func (s *service) SyncFile(bucketID, fileID, version string, r io.Reader, contentType string, created strfmt.DateTime, archetype string) (*models.FileDescriptor, error) {
	// calculate the checksum
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Info().Err(err).Msg("Failed to calculate checksum")
		return nil, err
	}

	// try to fetch
	_, fd, err := s.s3.Read(bucketID, fileID, version)

	switch {
	// Already exists and does not conflict
	case err == nil && checksum == fd.Checksum:
		s.logger.Info().
			Msg("File already exists")
		return fd, ErrAlreadyExists
	// Already exists and conflicts
	case err == nil && checksum != fd.Checksum:
		s.logger.Error().
			Msg("File already exists and has conflicting checksum")
		return nil, ErrAlreadyExistsConflict
	// Storage returned error and it is not "not found"
	case err != nil && err != s3.ErrNotFound:
		s.logger.Error().Err(err).
			Msg("Error while trying to read file")
		return nil, err
	}

	no := &object.NewObjectInfo{
		Archetype:   archetype,
		Checksum:    checksum,
		Size:        int64(buf.Len()),
		Created:     created,
		ContentType: contentType,
		Version:     version,
		Name:        fileID,
		Operation:   string(s3.Write),
	}

	return s.s3.Write(bucketID, no, &buf)
}

func (s *service) SyncFileDelete(bucketID, fileID, version string, created strfmt.DateTime) error {
	// get the previous file
	_, fd, err := s.s3.Read(bucketID, fileID, "")
	if err != nil {
		return err
	}

	// File was already deleted
	if fd.Operation == string(s3.Delete) {
		if fd.Version == version {
			s.logger.Debug().
				Msg("File delete already synced")
			return nil
		}
		s.logger.Error().
			Msg("File already deleted and delete has conflicting version")
		return ErrDeleted
	}

	// Write delete object
	no := &object.NewObjectInfo{
		Archetype:   fd.Archetype,
		Checksum:    "",
		Size:        0,
		Created:     created,
		ContentType: fd.ContentType,
		Version:     version,
		Name:        fileID,
		Operation:   string(s3.Delete),
	}

	_, err = s.s3.Write(bucketID, no, &bytes.Buffer{})
	return err
}

// New returns a new instance of storage service
func New(s3 s3.Storage, keyProvider s3.KeyProvider, publisher storageSync.Publisher, log zerolog.Logger) Service {
	return &service{s3: s3, keyProvider: keyProvider, publisher: publisher, logger: log}
}

var getUUID = func() string {
	return uuid.NewCrypto().String()
}

var getTime = func() strfmt.DateTime {
	return strfmt.DateTime(time.Now())
}
