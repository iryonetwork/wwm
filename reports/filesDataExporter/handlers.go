package filesDataExporter

import (
	"bytes"
	"context"
	"sort"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client/operations"
	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/reports"
)

type handlers struct {
	source         *operations.Client
	sourceAuth     runtime.ClientAuthInfoWriter
	dataSanitizer  Sanitizer
	reportsStorage reports.Storage
	logger         zerolog.Logger
}

// ExportFile export new files and file updates to reports storage.
func (h *handlers) ExportFile(ctx context.Context, bucketID, fileID, version string, timestamp strfmt.DateTime) (ExportResult, error) {
	// Get file from source storage
	var buf bytes.Buffer

	getParams := operations.NewFileGetVersionParams().
		WithBucket(bucketID).
		WithFileID(fileID).
		WithVersion(version).
		WithContext(ctx)
	_, err := h.source.FileGetVersion(getParams, h.sourceAuth, &buf)

	if err != nil {
		if _, ok := err.(*operations.FileGetVersionNotFound); ok {
			h.logger.Error().Err(err).
				Str("bucket", bucketID).
				Str("fileID", fileID).
				Str("version", version).
				Msg("File does not exist in source storage.")

			// File might have been already deleted; mark as succesful
			return ResultExportNotNeeded, nil
		}

		h.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Error on trying to fetch file from source storage.")
		return ResultError, err
	}

	// Check if export is needed
	needsExport, err := h.needsExport(ctx, fileID, version)
	if err != nil {
		return ResultError, err
	}
	// Nothing to do
	if !needsExport {
		// all is good but nothing was exported
		return ResultExportNotNeeded, nil
	}

	data, err := h.dataSanitizer.Sanitize(ctx, buf.Bytes())
	if err != nil {
		h.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Failed to sanitize file's data")
		return ResultError, err
	}

	err = h.reportsStorage.Insert(fileID, version, bucketID, timestamp, string(data))
	if err != nil {
		h.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Failed to save file's data in reports storage")
		return ResultError, err
	}

	return ResultExported, nil
}

// ExportFileDelete export file deletion to reports storage.
func (h *handlers) ExportFileDelete(ctx context.Context, bucketID, fileID, _ string, _ strfmt.DateTime) (ExportResult, error) {
	err := h.reportsStorage.Remove(fileID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Msg("Failed to delete file data from reports storage")
		return ResultError, err
	}

	return ResultExported, nil
}

// ListSourceBuckets lists all the buckets in source storage.
func (h *handlers) ListSourceBuckets(ctx context.Context) ([]*models.BucketDescriptor, error) {
	return h.listBuckets(ctx, h.source, h.sourceAuth)
}

// ListSourceFiles lists all the files in the bucket of source storage including files marked as deleted.
func (h *handlers) ListSourceFilesAsc(ctx context.Context, bucketID string) ([]*models.FileDescriptor, error) {
	return h.listFilesAsc(ctx, h.source, h.sourceAuth, bucketID)
}

// NewApiHandlers returns Handlers with cloudStorage and localStorage API used.
func NewHandlers(source *operations.Client, sourceAuth runtime.ClientAuthInfoWriter, dataSanitizer Sanitizer, reportsStorage reports.Storage, logger zerolog.Logger) Handlers {
	logger = logger.With().Str("component", "reports/filesDataExporter/handlers").Logger()

	return &handlers{
		source:         source,
		sourceAuth:     sourceAuth,
		dataSanitizer:  dataSanitizer,
		reportsStorage: reportsStorage,
		logger:         logger,
	}
}

func (h *handlers) needsExport(ctx context.Context, fileID, version string) (bool, error) {
	// Verify in case file already exists in destination storage
	exists, err := h.reportsStorage.Exists(fileID, version)
	if err != nil {
		return false, err
	}

	return !exists, nil
}

func (h *handlers) listBuckets(ctx context.Context, c *operations.Client, auth runtime.ClientAuthInfoWriter) ([]*models.BucketDescriptor, error) {
	params := operations.NewSyncBucketListParams().WithContext(ctx)
	resp, err := c.SyncBucketList(params, auth)

	if err != nil {
		// If not found return empty, otherwise return error
		if _, ok := err.(*operations.SyncBucketListNotFound); !ok {
			return nil, err
		}

		return []*models.BucketDescriptor{}, nil
	}

	return resp.Payload, nil
}

func (h *handlers) listFilesAsc(ctx context.Context, c *operations.Client, auth runtime.ClientAuthInfoWriter, bucketID string) ([]*models.FileDescriptor, error) {
	params := operations.NewSyncFileListParams().WithBucket(bucketID).WithContext(ctx)
	resp, err := c.SyncFileList(params, auth)

	if err != nil {
		// If not found return empty, otherwise return error
		if _, ok := err.(*operations.SyncFileListNotFound); !ok {
			return nil, err
		}

		return []*models.FileDescriptor{}, nil
	}

	// ensure descending order by created time
	files := resp.Payload
	sort.Sort(descByCreated(files))

	return files, nil
}
