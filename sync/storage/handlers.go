package storage

//go:generate sh ../../bin/mockgen.sh sync/storage Handlers $GOFILE

import (
	"bytes"
	"context"
	"sort"

	"github.com/go-openapi/runtime"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client/operations"
	"github.com/iryonetwork/wwm/gen/storage/models"
)

// Handlers describes public API for sync/storage event handlers
type Handlers interface {
	// SyncFile synchronizes new files and file updates to destination storage
	SyncFile(ctx context.Context, bucketID, fileID, version string, timestamp strfmt.DateTime) (SyncResult, error)
	// SyncFileDelete synchronizes file deletion to destination operations.
	SyncFileDelete(ctx context.Context, bucketID, fileID, version string, timestamp strfmt.DateTime) (SyncResult, error)
	// ListSourceBuckets lists all the buckets in source storage.
	ListSourceBuckets(ctx context.Context) ([]*models.BucketDescriptor, error)
	// ListSourceFiles lists all the files in the bucket of source storage including files marked as delete, ascending order by Created timestamp ensured.
	ListSourceFilesAsc(ctx context.Context, bucketID string) ([]*models.FileDescriptor, error)
	// ListSourceFileVersions lists all the file versions in the source storage ascending order by Created timestamp ensured.
	ListSourceFileVersionsAsc(ctx context.Context, bucketID, fileID string) ([]*models.FileDescriptor, error)
	// ListDestinationFileVersions lists all the file versions in the destination storage ascending order by Created timestamp ensured.
	ListDestinationFileVersionsAsc(ctx context.Context, bucketID, fileID string) ([]*models.FileDescriptor, error)
}

// Handler describes sync/storage sync handler function
type Handler func(ctx context.Context, bucketID, fileID, version string, created strfmt.DateTime) (SyncResult, error)

type handlers struct {
	source          *operations.Client
	sourceAuth      runtime.ClientAuthInfoWriter
	destination     *operations.Client
	destinationAuth runtime.ClientAuthInfoWriter
	logger          zerolog.Logger
}

// SyncFile synchronizes new files and file updates to destination storage
func (h *handlers) SyncFile(ctx context.Context, bucketID, fileID, version string, timestamp strfmt.DateTime) (SyncResult, error) {
	// Get file from source storage
	var buf bytes.Buffer

	getParams := operations.NewFileGetVersionParams().
		WithBucket(bucketID).
		WithFileID(fileID).
		WithVersion(version).
		WithContext(ctx)
	resp, err := h.source.FileGetVersion(getParams, h.sourceAuth, &buf)

	if err != nil {
		if _, ok := err.(*operations.FileGetVersionNotFound); ok {
			h.logger.Error().Err(err).
				Str("bucket", bucketID).
				Str("fileID", fileID).
				Str("version", version).
				Msg("File does not exist in source operations.")

			// File might have been already deleted; mark as succesful
			return ResultSyncNotNeeded, nil
		}

		h.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Error on trying to fetch file from source operations.")
		return ResultError, err
	}

	// Check if sync is needed
	needsSync, err := h.needsSync(ctx, bucketID, fileID, version, resp.XChecksum)
	if err != nil {
		return ResultError, err
	}
	// Nothing to do
	if !needsSync {
		// all is good but nothing was synced
		return ResultSyncNotNeeded, nil
	}

	// Sync file
	syncParams := operations.NewSyncFileParams().
		WithBucket(bucketID).
		WithFileID(fileID).
		WithVersion(version).
		WithContext(ctx).
		WithCreated(resp.XCreated)
	if resp.XArchetype != "" {
		syncParams.SetArchetype(&resp.XArchetype)
	}
	syncParams.SetContentType(resp.ContentType)
	syncParams.SetFile(runtime.NamedReader("FileReader", &buf))
	ok, created, err := h.destination.SyncFile(syncParams, h.destinationAuth)

	switch {
	case ok != nil:
		h.logger.Debug().
			Str("cmd", "SyncFile").
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("File already exists in destination storage")
	case created != nil:
		h.logger.Debug().
			Str("cmd", "SyncFile").
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Succesfully synced file to destination storage")
	case err != nil:
		h.logger.Error().Err(err).
			Str("cmd", "SyncFile").
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Failed to sync file to destination storage")
		switch err.(type) {
		case *operations.SyncFileConflict:
			// another attempt at sync should not be performed
			return ResultConflict, err
		default:
			return ResultError, err
		}
	}

	return ResultSynced, nil
}

// SyncFileDelete synchronizes file deletion to destination operations.
func (h *handlers) SyncFileDelete(ctx context.Context, bucketID, fileID, version string, timestamp strfmt.DateTime) (SyncResult, error) {
	params := operations.NewSyncFileDeleteParams().
		WithBucket(bucketID).
		WithFileID(fileID).
		WithVersion(version).
		WithCreated(timestamp).
		WithContext(ctx)
	_, err := h.destination.SyncFileDelete(params, h.destinationAuth)

	if err != nil {
		h.logger.Error().Err(err).
			Str("cmd", "SyncFileDelete").
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Failed to sync file deletion")
		switch err.(type) {
		case *operations.SyncFileDeleteConflict:
			// another attempt at sync should not be performed
			return ResultConflict, err
		case *operations.SyncFileDeleteNotFound:
			// another attempt at sync should not be performed
			return ResultSyncNotNeeded, err
		default:
			return ResultError, err
		}
	}

	h.logger.Debug().
		Str("cmd", "SyncFileDelete").
		Str("bucket", bucketID).
		Str("fileID", fileID).
		Str("version", version).
		Msg("Succesfully synced file deletion to destination storage")

	return ResultSynced, nil
}

// ListSourceBuckets lists all the buckets in source storage.
func (h *handlers) ListSourceBuckets(ctx context.Context) ([]*models.BucketDescriptor, error) {
	return h.listBuckets(ctx, h.source, h.sourceAuth)
}

// ListSourceFiles lists all the files in the bucket of source storage including files marked as delete.
func (h *handlers) ListSourceFilesAsc(ctx context.Context, bucketID string) ([]*models.FileDescriptor, error) {
	return h.listFilesAsc(ctx, h.source, h.sourceAuth, bucketID)
}

// ListSourceFileVersions lists all the file versions in the source storage.
func (h *handlers) ListSourceFileVersionsAsc(ctx context.Context, bucketID, fileID string) ([]*models.FileDescriptor, error) {
	return h.listFileVersionsAsc(ctx, h.source, h.sourceAuth, bucketID, fileID)
}

// ListDestinationFileVersions lists all the file versions in the destination storage.
func (h *handlers) ListDestinationFileVersionsAsc(ctx context.Context, bucketID, fileID string) ([]*models.FileDescriptor, error) {
	return h.listFileVersionsAsc(ctx, h.destination, h.destinationAuth, bucketID, fileID)
}

// NewApiHandlers returns Handlers with cloudStorage and localStorage API used.
func NewHandlers(source *operations.Client, sourceAuth runtime.ClientAuthInfoWriter, destination *operations.Client, destinationAuth runtime.ClientAuthInfoWriter, logger zerolog.Logger) Handlers {
	return &handlers{
		source:          source,
		sourceAuth:      sourceAuth,
		destination:     destination,
		destinationAuth: destinationAuth,
		logger:          logger,
	}
}

func (h *handlers) needsSync(ctx context.Context, bucketID, fileID, version, sourceChecksum string) (bool, error) {
	// Verify in case file already exists in destination storage
	params := operations.NewSyncFileMetadataParams().
		WithBucket(bucketID).
		WithFileID(fileID).
		WithVersion(version).
		WithContext(ctx)
	resp, err := h.destination.SyncFileMetadata(params, h.destinationAuth)

	// File already exists
	if resp != nil {
		if resp.XChecksum != sourceChecksum {
			// There is conflict, log error
			h.logger.Error().
				Str("bucket", bucketID).
				Str("fileID", fileID).
				Str("version", version).
				Msg("File already exists in destination storage and has different checksum.")
		}
		// Nothing to do
		return false, nil
	}
	// If file not found it needs sync, otherwise return error
	if _, ok := err.(*operations.SyncFileMetadataNotFound); !ok {
		return false, err
	}

	return true, nil
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

	// ensure ascending order by created time
	files := resp.Payload
	sort.Sort(ascByCreated(files))

	return files, nil
}

func (h *handlers) listFileVersionsAsc(ctx context.Context, c *operations.Client, auth runtime.ClientAuthInfoWriter, bucketID, fileID string) ([]*models.FileDescriptor, error) {
	params := operations.NewFileListVersionsParams().
		WithBucket(bucketID).
		WithFileID(fileID).
		WithContext(ctx)
	resp, err := c.FileListVersions(params, auth)

	if err != nil {
		// If not found return empty, otherwise return error
		if _, ok := err.(*operations.FileListVersionsNotFound); !ok {
			return nil, err
		}

		return []*models.FileDescriptor{}, nil
	}

	// ensure ascending order by created time
	files := resp.Payload
	sort.Sort(ascByCreated(files))

	return files, nil
}
