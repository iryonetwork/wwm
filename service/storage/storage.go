package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"time"

	"github.com/agext/uuid"
	"github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/object"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/utils"
)

// Service describes storage service's public API.
type Service interface {
	// Checksum calculates the checksum of a given reader using sha256.Sum.
	Checksum(r io.Reader) (string, error)

	// BucketList returns list of all the buckets.
	BucketList(ctx context.Context) ([]*models.BucketDescriptor, error)

	// FileList returns a list of latest versions of files. Older versions and
	// files marked as deleted are removed from the list.
	FileList(ctx context.Context, bucketID string) ([]*models.FileDescriptor, error)

	// FileGet returns the latest version of the file by returning the reader
	// and file details.
	FileGet(ctx context.Context, bucketID, fileID string) (io.ReadCloser, *models.FileDescriptor, error)

	// FileGetVersion returns a specific version of a file.
	FileGetVersion(ctx context.Context, bucketID, fileID, version string) (io.ReadCloser, *models.FileDescriptor, error)

	// FileListVersions returns a list of all modifications to a file.
	FileListVersions(ctx context.Context, bucketID, fileID string, createdAtSince, createdAtUntil *strfmt.DateTime) ([]*models.FileDescriptor, error)

	// FileNew creates a new file.
	FileNew(ctx context.Context, bucketID string, r io.Reader, contentType string, archetype string, labels []string) (*models.FileDescriptor, error)

	// FileUpdate creates a new version of a file.
	FileUpdate(ctx context.Context, bucketID, fileID string, r io.Reader, contentType string, archetype string, labels []string) (*models.FileDescriptor, error)

	// FileDelete marks file as deleted.
	FileDelete(ctx context.Context, bucketID, fileID string) error

	// SyncFileList returns a list of latest versions of files. Older versions are removed from the list.
	// Files marked as deleted are kept in the list.
	SyncFileList(ctx context.Context, bucketID string, createdAtSince, createdAtUntil *strfmt.DateTime) ([]*models.FileDescriptor, error)

	// SyncFile syncs file with provided fileID and version.
	SyncFile(ctx context.Context, bucketID, fileID, version string, r io.Reader, contentType string, created strfmt.DateTime, archetype string, labels []string) (*models.FileDescriptor, error)

	// SyncFileDelete sync file deletion.
	SyncFileDelete(ctx context.Context, bucketID, fileID, version string, created strfmt.DateTime) error
}

// Bucket or item was already deleted
var ErrDeleted = s3.ErrDeleted

// Item does not exists
var ErrNotFound = s3.ErrNotFound

// Item already exists
var ErrAlreadyExists = s3.ErrAlreadyExists

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

func (s *service) BucketList(ctx context.Context) ([]*models.BucketDescriptor, error) {
	// get the list and return
	return s.s3.ListBuckets(ctx)
}

func (s *service) FileList(ctx context.Context, bucketID string) ([]*models.FileDescriptor, error) {
	// init list to return
	list := []*models.FileDescriptor{}

	// check if bucket exists
	exists, err := s.s3.BucketExists(ctx, bucketID)
	if err != nil {
		s.logger.Info().Err(err).Str("bucket", bucketID).Msg("Failed to check if bucket exists")
		return nil, err
	}
	if !exists {
		return list, nil
	}

	// collect the list
	l, err := s.s3.List(ctx, bucketID, "")
	if err != nil {
		return nil, err
	}

	// extract only latest versions; latest version is already sorted
	// on top, add to return list; only include files with a write operation
	m := map[string]bool{}
	for _, f := range l {
		if _, ok := m[f.Name]; !ok {
			m[f.Name] = true
			if s3.Operation(f.Operation) == s3.Write {
				list = append(list, f)
			}
		}
	}

	return list, nil
}

func (s *service) FileGet(ctx context.Context, bucketID, fileID string) (io.ReadCloser, *models.FileDescriptor, error) {
	start := time.Now()
	rc, fd, err := s.s3.Read(ctx, bucketID, fileID, "")
	s.logger.Info().Str("method", "FileGet").Msgf("s3 read time %s", time.Since(start))

	return rc, fd, err
}

func (s *service) FileGetVersion(ctx context.Context, bucketID, fileID, version string) (io.ReadCloser, *models.FileDescriptor, error) {
	start := time.Now()
	rc, fd, err := s.s3.Read(ctx, bucketID, fileID, version)
	s.logger.Info().Str("method", "FileGetVersion").Msgf("s3 read time %s", time.Since(start))

	return rc, fd, err
}

func (s *service) FileListVersions(ctx context.Context, bucketID, fileID string, createdAtSince, createdAtUntil *strfmt.DateTime) ([]*models.FileDescriptor, error) {
	// init list to return
	list := []*models.FileDescriptor{}

	// check if bucket exists
	exists, err := s.s3.BucketExists(ctx, bucketID)
	if err != nil {
		s.logger.Info().Err(err).Str("bucket", bucketID).Msg("Failed to check if bucket exists")
		return nil, err
	}
	if !exists {
		return list, nil
	}

	l, err := s.s3.List(ctx, bucketID, fileID)
	if (createdAtSince == nil && createdAtUntil == nil) || err != nil {
		return l, err
	}

	// extract only versions fitting the created timestamp filtering specified
	for _, f := range l {
		if (createdAtSince == nil || (createdAtSince != nil && time.Time(*createdAtSince).Before(time.Time(f.Created)))) &&
			(createdAtUntil == nil || (createdAtUntil != nil && time.Time(*createdAtUntil).After(time.Time(f.Created)))) {
			list = append(list, f)
		}
	}

	return list, nil
}

func (s *service) FileNew(ctx context.Context, bucketID string, r io.Reader, contentType string, archetype string, labels []string) (*models.FileDescriptor, error) {
	err := s.EnsureBucket(ctx, bucketID)
	if err != nil {
		return nil, err
	}

	// calculate the checksum
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to calculate checksum")
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
		Labels:      labels,
	}

	start := time.Now()
	fd, err := s.s3.Write(ctx, bucketID, no, &buf)
	s.logger.Info().Str("method", "FileNew").Msgf("s3 write time %s", time.Since(start))

	if err == nil {
		s.publisher.PublishAsyncWithRetries(
			context.TODO(),
			storageSync.FileNew,
			&storageSync.FileInfo{BucketID: bucketID, FileID: fileID, Version: version, Created: fd.Created},
		)

		for _, label := range labels {
			err := s.updateFilesCollection(ctx, s3.Write, bucketID, label, fd)
			if err != nil {
				s.logger.Error().Err(err).Str("method", "FileNew").Msgf("failed to update files collection %s", label)
			}
		}
	}

	return fd, err
}

func (s *service) FileUpdate(ctx context.Context, bucketID, fileID string, r io.Reader, contentType string, archetype string, labels []string) (*models.FileDescriptor, error) {
	// get the previous file
	start := time.Now()
	_, old, err := s.s3.Read(ctx, bucketID, fileID, "")
	s.logger.Info().Str("method", "FileUpdate").Msgf("s3 read time %s", time.Since(start))

	if err != nil {
		return nil, err
	}

	// calculate the checksum
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Error().Err(err).Str("method", "FileUpdate").Msg("Failed to calculate checksum")
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
		Labels:      labels,
	}

	start = time.Now()
	fd, err := s.s3.Write(ctx, bucketID, no, &buf)
	s.logger.Info().Str("method", "FileUpdate").Msgf("s3 write time %s", time.Since(start))

	if err == nil {
		s.publisher.PublishAsyncWithRetries(
			context.TODO(),
			storageSync.FileUpdate,
			&storageSync.FileInfo{BucketID: bucketID, FileID: fileID, Version: version, Created: fd.Created},
		)

		for _, label := range labels {
			err := s.updateFilesCollection(ctx, s3.Write, bucketID, label, fd)
			if err != nil {
				s.logger.Error().Err(err).Str("method", "FileUpdate").Msgf("failed to update files collection %s", label)
			}
		}

		droppedLabels := utils.DiffSlice(old.Labels, labels)
		for _, label := range droppedLabels {
			err := s.updateFilesCollection(ctx, s3.Delete, bucketID, label, fd)
			if err != nil {
				s.logger.Error().Err(err).Str("method", "FileUpdate").Msgf("failed to update files collection %s", label)
			}
		}
	}
	return fd, err
}

func (s *service) FileDelete(ctx context.Context, bucketID, fileID string) error {
	// get the previous file
	start := time.Now()
	_, fd, err := s.s3.Read(ctx, bucketID, fileID, "")
	s.logger.Info().Str("method", "FileDelete").Msgf("s3 read time %s", time.Since(start))

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
		Labels:      fd.Labels,
	}

	start = time.Now()
	fd, err = s.s3.Write(ctx, bucketID, no, &bytes.Buffer{})
	s.logger.Info().Str("method", "FileDelete").Msgf("s3 write time %s", time.Since(start))

	if err == nil {
		s.publisher.PublishAsyncWithRetries(
			context.TODO(),
			storageSync.FileDelete,
			&storageSync.FileInfo{BucketID: bucketID, FileID: fileID, Version: version, Created: fd.Created},
		)
		for _, label := range fd.Labels {
			err := s.updateFilesCollection(ctx, s3.Delete, bucketID, label, fd)
			if err != nil {
				s.logger.Error().Err(err).Str("method", "FileDelete").Msgf("failed to update files collection %s", label)
			}
		}
	}
	return err
}

func (s *service) SyncFileList(ctx context.Context, bucketID string, createdAtSince, createdAtUntil *strfmt.DateTime) ([]*models.FileDescriptor, error) {
	// init list to return
	list := []*models.FileDescriptor{}

	// check if bucket exists
	exists, err := s.s3.BucketExists(ctx, bucketID)
	if err != nil {
		s.logger.Info().Err(err).Str("method", "SyncFileList").Str("bucket", bucketID).Msg("Failed to check if bucket exists")
		return nil, err
	}
	if !exists {
		return list, nil
	}

	// collect the list
	l, err := s.s3.List(ctx, bucketID, "")
	if err != nil {
		return nil, err
	}

	// extract only latest versions of files; latest version is already sorted
	// omit files that don't fit specified created timestamp filtering
	// on top, add to return list
	m := map[string]bool{}
	for _, f := range l {
		if _, ok := m[f.Name]; !ok {
			m[f.Name] = true
			if (createdAtSince == nil || (createdAtSince != nil && time.Time(*createdAtSince).Before(time.Time(f.Created)))) &&
				(createdAtUntil == nil || (createdAtUntil != nil && time.Time(*createdAtUntil).After(time.Time(f.Created)))) {
				list = append(list, f)
			}
		}
	}

	return list, nil
}

func (s *service) SyncFile(ctx context.Context, bucketID, fileID, version string, r io.Reader, contentType string, created strfmt.DateTime, archetype string, labels []string) (*models.FileDescriptor, error) {
	err := s.EnsureBucket(ctx, bucketID)
	if err != nil {
		return nil, err
	}

	// calculate the checksum
	var buf bytes.Buffer

	tee := io.TeeReader(r, &buf)
	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Error().Err(err).Str("method", "SyncFile").Msg("Failed to calculate checksum")
		return nil, err
	}

	// try to fetch
	start := time.Now()
	_, fd, err := s.s3.Read(ctx, bucketID, fileID, version)
	s.logger.Info().Str("method", "SyncFile").Msgf("s3 read time %s", time.Since(start))

	switch {
	// Already exists and does not conflict
	case err == nil && checksum == fd.Checksum:
		s.logger.Info().Str("method", "SyncFile").
			Msg("File already exists")
		return fd, ErrAlreadyExists
	// Already exists and conflicts
	case err == nil && checksum != fd.Checksum:
		s.logger.Info().Str("method", "SyncFile").
			Msg("File already exists and has conflicting checksum. Local file will be removed and replaced with sync file.")

		start = time.Now()
		err = s.s3.Delete(ctx, bucketID, fileID, version)
		s.logger.Info().Str("method", "SyncFile").Msgf("s3 delete time %s", time.Since(start))

		if err != nil {
			s.logger.Error().Err(err).Str("method", "SyncFile").
				Msg("Error while trying to delete file")
			return nil, err
		}
	// Storage returned error and it is not "not found"
	case err != nil && err != s3.ErrNotFound:
		s.logger.Error().Err(err).Str("method", "SyncFile").
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
		Labels:      labels,
	}

	start = time.Now()
	fd, err = s.s3.Write(ctx, bucketID, no, &buf)
	s.logger.Info().Str("method", "SyncFile").Msgf("s3 write time %s", time.Since(start))

	return fd, err
}

func (s *service) SyncFileDelete(ctx context.Context, bucketID, fileID, version string, created strfmt.DateTime) error {
	// get the previous file
	start := time.Now()
	_, fd, err := s.s3.Read(ctx, bucketID, fileID, "")
	s.logger.Info().Str("method", "SyncFileDelete").Msgf("s3 read time %s", time.Since(start))

	if err != nil {
		return err
	}

	// File was already deleted
	if fd.Operation == string(s3.Delete) {
		if fd.Version == version {
			s.logger.Debug().Str("method", "SyncFileDelete").
				Msg("File delete already synced")
			return nil
		}
		s.logger.Error().Str("method", "SyncFileDelete").
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
		Labels:      fd.Labels,
	}

	start = time.Now()
	_, err = s.s3.Write(ctx, bucketID, no, &bytes.Buffer{})
	s.logger.Info().Str("method", "SyncFileDelete").Msgf("s3 write time %s", time.Since(start))

	return err
}

func (s *service) EnsureBucket(ctx context.Context, bucketID string) error {
	// make sure bucket exists
	if err := s.s3.MakeBucket(ctx, bucketID); err != nil && err != s3.ErrAlreadyExists {
		s.logger.Error().Err(err).Str("bucket", bucketID).Msg("Failed to ensure bucket")
		return err
	}

	return nil
}

func (s *service) updateFilesCollection(ctx context.Context, operation s3.Operation, bucketID, label string, fd *models.FileDescriptor) error {
	var c *filesCollection

	start := time.Now()
	r, _, err := s.s3.Read(ctx, bucketID, label, "")
	s.logger.Info().Str("method", "updateFilesCollection").Msgf("s3 read time %s", time.Since(start))

	if err != nil {
		if err != s3.ErrNotFound {
			return err
		}
		c = &filesCollection{}
	} else {
		c, err = FilesCollection(r)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to parse file collection file")
			return err
		}
	}

	switch operation {
	case s3.Write:
		c.Update(fd)
	case s3.Delete:
		c.Remove(fd)
	}

	f, err := c.GetFile()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to generate file collection file")
		return err
	}

	// calculate the checksum
	var buf bytes.Buffer
	tee := io.TeeReader(f, &buf)
	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to calculate checksum")
		return err
	}

	fileID := label
	version := getUUID()
	no := &object.NewObjectInfo{
		Checksum:    checksum,
		Size:        int64(buf.Len()),
		Created:     getTime(),
		ContentType: "application/x-collection+json",
		Version:     version,
		Name:        fileID,
		Operation:   string(s3.Write),
		Labels:      []string{labelFilesCollection},
	}

	start = time.Now()
	fd, err = s.s3.Write(ctx, bucketID, no, &buf)
	s.logger.Info().Str("method", "updateFilesCollection").Msgf("s3 write time %s", time.Since(start))

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to write file collection file")
		return err
	}

	s.publisher.PublishAsyncWithRetries(
		context.TODO(),
		storageSync.FileUpdate,
		&storageSync.FileInfo{BucketID: bucketID, FileID: fileID, Version: fd.Version, Created: fd.Created},
	)

	return nil
}

// New returns a new instance of storage service
func New(s3 s3.Storage, keyProvider s3.KeyProvider, publisher storageSync.Publisher, logger zerolog.Logger) Service {
	logger.Error().Msg("test")
	logger = logger.With().Str("component", "service/storage").Logger()
	logger.Error().Msg("test")
	return &service{s3: s3, keyProvider: keyProvider, publisher: publisher, logger: logger}
}

var getUUID = func() string {
	return uuid.NewCrypto().String()
}

var getTime = func() strfmt.DateTime {
	return strfmt.DateTime(time.Now())
}
