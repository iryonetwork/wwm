package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"github.com/agext/uuid"
	"github.com/go-openapi/strfmt"
	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/object"
	"github.com/rs/zerolog"
)

// Service describes storage service's public API.
type Service interface {
	// Checksum calculates the checksum of a given reader using sha256.Sum.
	Checksum(r io.Reader) (string, error)

	// FileList returns a list of latest versions of files. Older versions and
	// files markder as deleted are removed from the list.
	FileList(bucketID string) ([]*models.FileDescriptor, error)

	// FileGet returns the latest version of the file by returning the reader
	// and file details.
	FileGet(bucketID, fileID string) (io.ReadCloser, *models.FileDescriptor, error)

	// FileGetVersion returns a specific version of a file.
	FileGetVersion(bucketID, fileID, version string) (io.ReadCloser, *models.FileDescriptor, error)

	// FileListVersions returns a list of all modifications to a file.
	FileListVersions(bucker, fileID string) ([]*models.FileDescriptor, error)

	// FileNew creates a new file.
	FileNew(bucketID string, r io.Reader, contentType string, archetype string) (*models.FileDescriptor, error)

	// FileUpdate creates a new version of a file.
	FileUpdate(bucketID, fileID string, r io.Reader, contentType string, archetype string) (*models.FileDescriptor, error)

	// FileDelete marks file as deleted.
	FileDelete(bucketID, fileID string) error

	// FileSync syncs file with provided FileID and Version.
	FileSync(bucketID, fileID, version string, r io.Reader, contentType string, created strfmt.DateTime, archetype string) (*models.FileDescriptor, error)
}

var ErrAlreadyExists = errors.New("Item already exists")
var ErrAlreadyExistsConflict = errors.New("Item already exists and has differing checksum")

type service struct {
	s3          s3.Storage
	keyProvider s3.KeyProvider
	logger      zerolog.Logger
}

func (s *service) Checksum(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
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

	no := &object.NewObjectInfo{
		Archetype:   archetype,
		Size:        int64(buf.Len()),
		Checksum:    checksum,
		Created:     getTime(),
		ContentType: contentType,
		Version:     getUUID(),
		Name:        getUUID(),
		Operation:   string(s3.Write),
	}

	return s.s3.Write(bucketID, no, &buf)
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

	no := &object.NewObjectInfo{
		Archetype:   archetype,
		Checksum:    checksum,
		Size:        int64(buf.Len()),
		Created:     getTime(),
		ContentType: contentType,
		Version:     getUUID(),
		Name:        fileID,
		Operation:   string(s3.Write),
	}

	return s.s3.Write(bucketID, no, &buf)
}

func (s *service) FileDelete(bucketID, fileID string) error {
	// get the previous file
	_, fd, err := s.s3.Read(bucketID, fileID, "")
	if err != nil {
		return err
	}

	no := &object.NewObjectInfo{
		Archetype:   fd.Archetype,
		Checksum:    "",
		Size:        0,
		Created:     getTime(),
		ContentType: fd.ContentType,
		Version:     getUUID(),
		Name:        fileID,
		Operation:   string(s3.Delete),
	}

	_, err = s.s3.Write(bucketID, no, &bytes.Buffer{})
	return err
}

func (s *service) FileSync(bucketID, fileID, version string, r io.Reader, contentType string, created strfmt.DateTime, archetype string) (*models.FileDescriptor, error) {
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
		return fd, ErrAlreadyExists
	// Already exists and conflicts
	case err == nil && checksum != fd.Checksum:
		return nil, ErrAlreadyExistsConflict
	// Storage returned error and it is not "not found"
	case err != nil && err != s3.ErrNotFound:
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

// New returns a new instance of storage service
func New(s3 s3.Storage, keyProvider s3.KeyProvider, log zerolog.Logger) Service {
	return &service{s3: s3, keyProvider: keyProvider, logger: log}
}

var getUUID = func() string {
	return uuid.NewCrypto().String()
}

var getTime = func() strfmt.DateTime {
	return strfmt.DateTime(time.Now())
}
