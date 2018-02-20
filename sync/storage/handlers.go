package storage

//go:generate sh ../../bin/mockgen.sh sync/storage Handlers $GOFILE

import (
	"bytes"

	"github.com/go-openapi/runtime"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client/operations"
	"github.com/iryonetwork/wwm/gen/storage/models"
)

// Handlers describes public API for sync/storage event handlers
type Handlers interface {
	// SyncFile synchronizes new files and file updates to destination storage
	SyncFile(bucketID, fileID, version string) error
	// SyncFileDelete synchronizes file deletion to destination operations.
	SyncFileDelete(bucketID, fileID, version string) error
	// ListSourceBuckets lists all the buckets in source storage.
	ListSourceBuckets() ([]*models.BucketDescriptor, error)
	// ListSourceFiles lists all the files in the bucket of source storage including files marked as delete.
	ListSourceFiles(bucketID string) ([]*models.FileDescriptor, error)
	// ListSourceFileVersions lists all the file versions in the source storage.
	ListSourceFileVersions(bucketID, fileID string) ([]*models.FileDescriptor, error)
	// ListDestinationFileVersions lists all the file versions in the destination storage.
	ListDestinationFileVersions(bucketID, fileID string) ([]*models.FileDescriptor, error)
}

// Handler describes sync/storage handler function
type Handler func(bucketID, fileID, version string) error

type handlers struct {
	source          *operations.Client
	sourceAuth      runtime.ClientAuthInfoWriter
	destination     *operations.Client
	destinationAuth runtime.ClientAuthInfoWriter
	logger          zerolog.Logger
}

// SyncFile synchronizes new files and file updates to destination storage
func (h *handlers) SyncFile(bucketID, fileID, version string) error {
	// Get file from source storage
	var buf bytes.Buffer
	getParams := operations.NewFileGetVersionParams().WithBucket(bucketID).WithFileID(fileID).WithVersion(version)
	resp, err := h.source.FileGetVersion(getParams, h.sourceAuth, &buf)

	if err != nil {
		if _, ok := err.(*operations.FileGetVersionNotFound); ok {
			h.logger.Error().Err(err).
				Str("bucket", bucketID).
				Str("fileID", fileID).
				Str("version", version).
				Msg("File does not exist in source operations.")

			// File might have been already deleted; mark as succesful
			return nil
		}

		h.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("fileID", fileID).
			Str("version", version).
			Msg("Error on trying to fetch file from source operations.")
		return err
	}

	// Check if sync is needed
	needsSync, err := h.needsSync(bucketID, fileID, version, resp.XChecksum)
	if err != nil {
		return err
	}
	// Nothing to do
	if !needsSync {
		return nil
	}

	// Sync file
	syncParams := operations.NewSyncFileParams().WithBucket(bucketID).WithFileID(fileID).WithVersion(version)

	if resp.XArchetype != "" {
		syncParams.SetArchetype(&resp.XArchetype)
	}
	syncParams.SetContentType(resp.ContentType)
	syncParams.SetCreated(resp.XCreated)
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
			return nil
		default:
			return err
		}
	}

	return nil
}

// SyncFileDelete synchronizes file deletion to destination operations.
func (h *handlers) SyncFileDelete(bucketID, fileID, version string) error {
	params := operations.NewSyncFileDeleteParams().WithBucket(bucketID).WithFileID(fileID).WithVersion(version)

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
			return nil
		case *operations.SyncFileDeleteNotFound:
			// another attempt at sync should not be performed
			return nil
		default:
			return err
		}
	}

	h.logger.Debug().
		Str("cmd", "SyncFileDelete").
		Str("bucket", bucketID).
		Str("fileID", fileID).
		Str("version", version).
		Msg("Succesfully synced file deletion to destination storage")

	return nil
}

// ListSourceBuckets lists all the buckets in source storage.
func (h *handlers) ListSourceBuckets() ([]*models.BucketDescriptor, error) {
	return h.listBuckets(h.source, h.sourceAuth)
}

// ListSourceFiles lists all the files in the bucket of source storage including files marked as delete.
func (h *handlers) ListSourceFiles(bucketID string) ([]*models.FileDescriptor, error) {
	return h.listFiles(h.source, h.sourceAuth, bucketID)
}

// ListSourceFileVersions lists all the file versions in the source storage.
func (h *handlers) ListSourceFileVersions(bucketID, fileID string) ([]*models.FileDescriptor, error) {
	return h.listFileVersions(h.source, h.sourceAuth, bucketID, fileID)
}

// ListDestinationFileVersions lists all the file versions in the destination storage.
func (h *handlers) ListDestinationFileVersions(bucketID, fileID string) ([]*models.FileDescriptor, error) {
	return h.listFileVersions(h.destination, h.destinationAuth, bucketID, fileID)
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

func (h *handlers) needsSync(bucketID, fileID, version, sourceChecksum string) (bool, error) {
	// Verify in case file already exists in destination storage
	params := operations.NewSyncFileMetadataParams().WithBucket(bucketID).WithFileID(fileID).WithVersion(version)
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

func (h *handlers) listBuckets(c *operations.Client, auth runtime.ClientAuthInfoWriter) ([]*models.BucketDescriptor, error) {
	resp, err := c.SyncBucketList(nil, auth)

	if err != nil {
		// If not found return empty, otherwise return error
		if _, ok := err.(*operations.SyncBucketListNotFound); !ok {
			return nil, err
		}

		return []*models.BucketDescriptor{}, nil
	}

	return resp.Payload, nil
}

func (h *handlers) listFiles(c *operations.Client, auth runtime.ClientAuthInfoWriter, bucketID string) ([]*models.FileDescriptor, error) {
	params := operations.NewSyncFileListParams().WithBucket(bucketID)
	resp, err := c.SyncFileList(params, auth)

	if err != nil {
		// If not found return empty, otherwise return error
		if _, ok := err.(*operations.SyncFileListNotFound); !ok {
			return nil, err
		}

		return []*models.FileDescriptor{}, nil
	}

	return resp.Payload, nil
}

func (h *handlers) listFileVersions(c *operations.Client, auth runtime.ClientAuthInfoWriter, bucketID, fileID string) ([]*models.FileDescriptor, error) {
	params := operations.NewFileListVersionsParams().WithBucket(bucketID).WithFileID(fileID)
	resp, err := c.FileListVersions(params, auth)

	if err != nil {
		// If not found return empty, otherwise return error
		if _, ok := err.(*operations.FileListVersionsNotFound); !ok {
			return nil, err
		}

		return []*models.FileDescriptor{}, nil
	}

	return resp.Payload, nil
}
