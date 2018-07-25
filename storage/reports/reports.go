package reports

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/reports"
	"github.com/iryonetwork/wwm/storage/reports/db"
)

type (
	storage struct {
		ctx    context.Context
		logger zerolog.Logger
		db     db.DB
		gdb    *gorm.DB
	}
)

// New initializes a new instance of Storage
func New(ctx context.Context, gdb *gorm.DB, logger zerolog.Logger) (*storage, error) {
	s := &storage{
		ctx:    ctx,
		logger: logger.With().Str("component", "storage/reports").Logger(),
		db:     db.New(gdb),
		gdb:    gdb,
	}

	return s, nil
}

// Insert inserts new file or updated file data
func (s *storage) Insert(fileID, version, patientID string, timestamp strfmt.DateTime, data string) error {
	tx := s.db.Exec(
		"INSERT INTO \"files\" (file_id, version, patient_id, created_at, updated_at, data) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (file_id) DO UPDATE SET version = ?, updated_at = ?, data = ?",
		fileID,
		version,
		patientID,
		timestamp.String(),
		timestamp.String(),
		data,
		version,
		timestamp.String(),
		data,
	)

	if err := tx.GetError(); err != nil {
		s.logger.Error().Err(err).Msg("failed to save new file to database")
		return errors.Wrap(err, "failed to save new file to database")
	}

	return nil
}

// Remove deletes file data
func (s *storage) Remove(fileID string) error {
	tx := s.db.Delete(&reports.File{FileID: fileID})

	if err := tx.GetError(); err != nil {
		s.logger.Error().Err(err).Msg("failed to delete file from database")
		return errors.Wrap(err, "failed to delete file from database")
	}

	return nil
}

// Get looks up file by ID and version (optional)
func (s *storage) Get(fileID, version string) (*reports.File, error) {
	var file reports.File

	if err := s.db.Where("file_id = ?", fileID).First(&file).GetError(); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		s.logger.Error().Err(err).Msg("failed to get a file from database")
		return nil, errors.Wrap(err, "failed to get a file from database")
	}

	if version != "" && file.Version != version {
		return nil, nil
	}

	return &file, nil
}

// Exists check if file with given ID and version (optional) already exists in DB
func (s *storage) Exists(fileID, version string) (bool, error) {
	var file reports.File

	if err := s.db.Where("file_id = ?", fileID).Select("version").First(&file).GetError(); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		s.logger.Error().Err(err).Msg("failed to get a file from database")
		return false, errors.Wrap(err, "failed to get a file from database")
	}

	if version != "" && file.Version != version {
		return false, nil
	}

	return true, nil
}

// Find looks up files by patientID, data key values, createdStart, createdEnd
func (s *storage) Find(patientID string, dataKeyValues map[string]string, createdAtStart *strfmt.DateTime, createdAtEnd *strfmt.DateTime) (*[]reports.File, error) {
	// search for each separate token
	sqlWhere := []string{}
	sqlAttrs := []interface{}{}

	// add patient_id condition if present
	if patientID != "" {
		sqlWhere = append(sqlWhere, "patient_id = ?")
		sqlAttrs = append(sqlAttrs, patientID)
	}
	// add created_at start condition if present
	if createdAtStart != nil {
		sqlWhere = append(sqlWhere, "created_at >= ?")
		sqlAttrs = append(sqlAttrs, createdAtStart.String())
	}
	// add created_at end condition if present
	if createdAtEnd != nil {
		sqlWhere = append(sqlWhere, "created_at <= ?")
		sqlAttrs = append(sqlAttrs, createdAtEnd.String())
	}
	// add data key value conditions
	for key, value := range dataKeyValues {
		sqlWhere = append(sqlWhere, fmt.Sprintf("data->>'%s' = ?", key))
		sqlAttrs = append(sqlAttrs, value)
	}

	files := []reports.File{}
	err := s.db.Where(strings.Join(sqlWhere, " AND "), sqlAttrs...).Order("created_at asc").Find(&files).GetError()
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &files, nil
		}
		s.logger.Error().Err(err).Msg("failed to find matching files")
		return nil, errors.Wrap(err, "failed to find matching files")
	}

	return &files, nil
}

// Close closes the DB connection
func (s *storage) Close() error {
	return s.db.Close()
}
