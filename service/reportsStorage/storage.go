package reportsStorage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/agext/uuid"
	"github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/reportsStorage/models"
	storageModels "github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/object"
)

// Service describes storage service's public API.
type Service interface {
	// Checksum calculates the checksum of a given reader using sha256.Sum.
	Checksum(r io.Reader) (string, error)

	// ReportList returns a list of latest versions of files. Older versions and
	// files marked as deleted are removed from the list.
	ReportList(ctx context.Context, reportType string) ([]*models.ReportFileDescriptor, error)

	// ReportGet returns the latest version of the file by returning the reader
	// and file details.
	ReportGet(ctx context.Context, reportType, fileName string) (io.ReadCloser, *models.ReportFileDescriptor, error)

	// ReportGetVersion returns a specific version of a file.
	ReportGetVersion(ctx context.Context, reportType, fileName, version string) (io.ReadCloser, *models.ReportFileDescriptor, error)

	// ReportListVersions returns a list of all modifications to a file.
	ReportListVersions(ctx context.Context, reportType, fileName string, createdAtSince, createdAtUntil *strfmt.DateTime) ([]*models.ReportFileDescriptor, error)

	// ReportNew creates a new file.
	ReportNew(ctx context.Context, reportType string, r io.Reader, contentType string, dataSince *strfmt.DateTime, dataUntil strfmt.DateTime) (*models.ReportFileDescriptor, error)

	// ReportUpdate creates a new version of a file.
	ReportUpdate(ctx context.Context, reportType, fileName string, r io.Reader, contentType string) (*models.ReportFileDescriptor, error)

	// ReportDelete marks file as deleted.
	ReportDelete(ctx context.Context, reportType, fileName string) error
}

// Bucket or item was already deleted
var ErrDeleted = s3.ErrDeleted

// Item does not exists
var ErrNotFound = s3.ErrNotFound

// Item already exists
var ErrAlreadyExists = s3.ErrAlreadyExists

type service struct {
	s3     s3.Storage
	logger zerolog.Logger
}

func (s *service) Checksum(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

func (s *service) ReportList(ctx context.Context, reportType string) ([]*models.ReportFileDescriptor, error) {
	// init list to return
	list := []*models.ReportFileDescriptor{}

	// check if bucket exists
	exists, err := s.s3.BucketExists(ctx, reportType)
	if err != nil {
		s.logger.Info().Err(err).Str("bucket", reportType).Msg("Failed to check if bucket exists")
		return nil, err
	}
	if !exists {
		return list, nil
	}

	// collect the list
	l, err := s.s3.List(ctx, reportType, "")
	if err != nil {
		return nil, err
	}

	// extract only latest versions; latest version is already sorted
	// create ReportFileDescriptor objects
	m := map[string]bool{}
	for _, f := range l {
		if _, ok := m[f.Name]; !ok {
			m[f.Name] = true
			if s3.Operation(f.Operation) == s3.Write {
				rfd, err := getReportFileDescriptor(reportType, f)
				if err != nil {
					return nil, err
				}
				list = append(list, rfd)
			}
		}
	}

	return list, nil
}

func (s *service) ReportGet(ctx context.Context, reportType, fileName string) (io.ReadCloser, *models.ReportFileDescriptor, error) {
	start := time.Now()
	rc, fd, err := s.s3.Read(ctx, reportType, fileName, "")
	s.logger.Info().Str("method", "ReportGet").Msgf("s3 read time %s", time.Since(start))
	if err != nil {
		return nil, nil, err
	}

	rfd, err := getReportFileDescriptor(reportType, fd)
	return rc, rfd, err
}

func (s *service) ReportGetVersion(ctx context.Context, reportType, fileName, version string) (io.ReadCloser, *models.ReportFileDescriptor, error) {
	start := time.Now()
	rc, fd, err := s.s3.Read(ctx, reportType, fileName, version)
	s.logger.Info().Str("method", "ReportGetVersion").Msgf("s3 read time %s", time.Since(start))
	if err != nil {
		return nil, nil, err
	}

	rfd, err := getReportFileDescriptor(reportType, fd)
	return rc, rfd, err
}

func (s *service) ReportListVersions(ctx context.Context, reportType, fileName string, createdAtSince, createdAtUntil *strfmt.DateTime) ([]*models.ReportFileDescriptor, error) {
	// init list to return
	list := []*models.ReportFileDescriptor{}

	// check if bucket exists
	exists, err := s.s3.BucketExists(ctx, reportType)
	if err != nil {
		s.logger.Info().Err(err).Str("bucket", reportType).Msg("Failed to check if bucket exists")
		return nil, err
	}
	if !exists {
		return list, nil
	}

	l, err := s.s3.List(ctx, reportType, fileName)
	if err != nil {
		return nil, err
	}

	if createdAtSince == nil && createdAtUntil == nil {
		for _, f := range l {
			rfd, err := getReportFileDescriptor(reportType, f)
			if err != nil {
				return nil, err
			}
			list = append(list, rfd)
		}
	} else {
		// extract only versions fitting the created timestamp filtering specified
		for _, f := range l {
			if (createdAtSince == nil || (createdAtSince != nil && time.Time(*createdAtSince).Before(time.Time(f.Created)))) &&
				(createdAtUntil == nil || (createdAtUntil != nil && time.Time(*createdAtUntil).After(time.Time(f.Created)))) {
				rfd, err := getReportFileDescriptor(reportType, f)
				if err != nil {
					return nil, err
				}
				list = append(list, rfd)
			}
		}
	}
	return list, nil
}

func (s *service) ReportNew(ctx context.Context, reportType string, r io.Reader, contentType string, dataSince *strfmt.DateTime, dataUntil strfmt.DateTime) (*models.ReportFileDescriptor, error) {
	err := s.EnsureBucket(ctx, reportType)
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

	fileName := createFilename(dataSince, dataUntil)

	version := getUUID()
	no := &object.NewObjectInfo{
		Size:        int64(buf.Len()),
		Checksum:    checksum,
		Created:     getTime(),
		ContentType: contentType,
		Version:     version,
		Name:        fileName,
		Operation:   string(s3.Write),
	}

	start := time.Now()
	fd, err := s.s3.Write(ctx, reportType, no, &buf)
	s.logger.Info().Str("method", "ReportNew").Msgf("s3 write time %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	return getReportFileDescriptor(reportType, fd)
}

func (s *service) ReportUpdate(ctx context.Context, reportType, fileName string, r io.Reader, contentType string) (*models.ReportFileDescriptor, error) {
	// get the previous file
	start := time.Now()
	_, _, err := s.s3.Read(ctx, reportType, fileName, "")
	s.logger.Info().Str("method", "ReportUpdate").Msgf("s3 read time %s", time.Since(start))

	if err != nil {
		return nil, err
	}

	// calculate the checksum
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	checksum, err := s.Checksum(tee)
	if err != nil {
		s.logger.Error().Err(err).Str("method", "ReportUpdate").Msg("Failed to calculate checksum")
		return nil, err
	}

	version := getUUID()
	no := &object.NewObjectInfo{
		Checksum:    checksum,
		Size:        int64(buf.Len()),
		Created:     getTime(),
		ContentType: contentType,
		Version:     version,
		Name:        fileName,
		Operation:   string(s3.Write),
	}

	start = time.Now()
	fd, err := s.s3.Write(ctx, reportType, no, &buf)
	s.logger.Info().Str("method", "ReportUpdate").Msgf("s3 write time %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	return getReportFileDescriptor(reportType, fd)
}

func (s *service) ReportDelete(ctx context.Context, reportType, fileName string) error {
	// get the previous file
	start := time.Now()
	_, fd, err := s.s3.Read(ctx, reportType, fileName, "")
	s.logger.Info().Str("method", "ReportDelete").Msgf("s3 read time %s", time.Since(start))

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
		Name:        fileName,
		Operation:   string(s3.Delete),
		Labels:      fd.Labels,
	}

	start = time.Now()
	fd, err = s.s3.Write(ctx, reportType, no, &bytes.Buffer{})
	s.logger.Info().Str("method", "ReportDelete").Msgf("s3 write time %s", time.Since(start))

	return err
}

func (s *service) EnsureBucket(ctx context.Context, reportType string) error {
	// make sure bucket exists
	if err := s.s3.MakeBucket(ctx, reportType); err != nil && err != s3.ErrAlreadyExists {
		s.logger.Error().Err(err).Str("bucket", reportType).Msg("Failed to ensure bucket")
		return err
	}

	return nil
}

func getReportFileDescriptor(reportType string, fd *storageModels.FileDescriptor) (*models.ReportFileDescriptor, error) {
	var dataSince strfmt.DateTime
	var dataUntil strfmt.DateTime

	coverageDatesString, err := base64.StdEncoding.DecodeString(fd.Name)
	if err != nil {
		return nil, err
	}
	coverageDates := strings.Split(string(coverageDatesString), "/")
	if len(coverageDates) > 1 {
		dataSince, err = strfmt.ParseDateTime(coverageDates[0])
		if err != nil {
			return nil, err
		}
		dataUntil, err = strfmt.ParseDateTime(coverageDates[1])
		if err != nil {
			return nil, err
		}
	} else {
		dataUntil, err = strfmt.ParseDateTime(coverageDates[0])
		if err != nil {
			return nil, err
		}
	}

	return &models.ReportFileDescriptor{
		Checksum:    fd.Checksum,
		ContentType: fd.ContentType,
		Created:     fd.Created,
		ReportType:  reportType,
		DataSince:   dataSince,
		DataUntil:   dataUntil,
		Name:        fd.Name,
		Operation:   fd.Operation,
		Path:        fd.Path,
		Size:        fd.Size,
		Version:     fd.Version,
	}, nil
}

func createFilename(dataSince *strfmt.DateTime, dataUntil strfmt.DateTime) string {
	if dataSince == nil {
		return base64.StdEncoding.EncodeToString([]byte(dataUntil.String()))
	}

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s/%s", dataSince.String(), dataUntil.String())))
}

// New returns a new instance of storage service
func New(s3 s3.Storage, logger zerolog.Logger) Service {
	logger = logger.With().Str("component", "service/reportsStorage").Logger()
	return &service{s3: s3, logger: logger}
}

var getUUID = func() string {
	return uuid.NewCrypto().String()
}

var getTime = func() strfmt.DateTime {
	return strfmt.DateTime(time.Now())
}
