package consumer

import (
	"bytes"

	"github.com/go-openapi/runtime"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client/storage"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

// Handlers describes public API for sync/storage event handlers
type Handlers interface {
	// SyncFile synchronizes new files and file updates to destination storage
	SyncFile(f *storageSync.FileInfo) error
	// SyncFileDelete synchronizes file deletion to destination storage.
	SyncFileDelete(f *storageSync.FileInfo) error
}

// Handler describes handler functions
type Handler func(f *storageSync.FileInfo) error

type handlers struct {
	source          *storage.Client
	sourceAuth      runtime.ClientAuthInfoWriter
	destination     *storage.Client
	destinationAuth runtime.ClientAuthInfoWriter
	logger          zerolog.Logger
}

// SyncFile synchronizes new files and file updates to destination storage
func (h *handlers) SyncFile(f *storageSync.FileInfo) error {
	// Get file from source storage
	var buf bytes.Buffer
	getParams := storage.NewFileGetVersionParams().WithBucket(f.BucketID).WithFileID(f.FileID).WithVersion(f.Version)
	resp, err := h.source.FileGetVersion(getParams, h.sourceAuth, &buf)

	if err != nil {
		if _, ok := err.(*storage.FileGetVersionNotFound); ok {
			h.logger.Error().Err(err).
				Str("bucket", f.BucketID).
				Str("fileID", f.FileID).
				Str("version", f.Version).
				Msg("File does not exist in source storage.")

			// File might have been already deleted; mark as succesful
			return nil
		}

		h.logger.Error().Err(err).
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("Error on trying to fetch file from source storage.")
		return err
	}

	// Check if sync is needed
	needsSync, err := h.needsSync(f, resp.XChecksum)
	if err != nil {
		return err
	}
	// Nothing to do
	if !needsSync {
		return nil
	}

	// Sync file
	syncParams := storage.NewSyncFileParams().WithBucket(f.BucketID).WithFileID(f.FileID).WithVersion(f.Version)

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
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("File already exists in destination storage")
	case created != nil:
		h.logger.Debug().
			Str("cmd", "SyncFile").
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("Succesfully synced file to destination storage")
	case err != nil:
		h.logger.Error().Err(err).
			Str("cmd", "SyncFile").
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("Failed to sync file to destination storage")
		switch err.(type) {
		case *storage.SyncFileConflict:
			// another attempt at sync should not be performed
			return nil
		default:
			return err
		}
	}

	return nil
}

// SyncFileDelete synchronizes file deletion to destination storage.
func (h *handlers) SyncFileDelete(f *storageSync.FileInfo) error {
	params := storage.NewSyncFileDeleteParams().WithBucket(f.BucketID).WithFileID(f.FileID).WithVersion(f.Version)

	_, err := h.destination.SyncFileDelete(params, h.destinationAuth)

	if err != nil {
		h.logger.Error().Err(err).
			Str("cmd", "SyncFileDelete").
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("Failed to sync file deletion")
		switch err.(type) {
		case *storage.SyncFileDeleteConflict:
			// another attempt at sync should not be performed
			return nil
		case *storage.SyncFileDeleteNotFound:
			// another attempt at sync should not be performed
			return nil
		default:
			return err
		}
	}

	h.logger.Debug().
		Str("cmd", "SyncFileDelete").
		Str("bucket", f.BucketID).
		Str("fileID", f.FileID).
		Str("version", f.Version).
		Msg("Succesfully synced file deletion to destination storage")

	return nil
}

// NewApiHandlers returns Handlers with cloudStorage and localStorage API used.
func NewHandlers(source *storage.Client, sourceAuth runtime.ClientAuthInfoWriter, destination *storage.Client, destinationAuth runtime.ClientAuthInfoWriter, logger zerolog.Logger) Handlers {
	return &handlers{
		source:          source,
		sourceAuth:      sourceAuth,
		destination:     destination,
		destinationAuth: destinationAuth,
		logger:          logger,
	}
}

func (h *handlers) needsSync(f *storageSync.FileInfo, sourceChecksum string) (bool, error) {
	// Verify in case file already exists in destination storage
	params := storage.NewSyncFileMetadataParams().WithBucket(f.BucketID).WithFileID(f.FileID).WithVersion(f.Version)
	resp, err := h.destination.SyncFileMetadata(params, h.destinationAuth)
	// File already exists
	if resp != nil {
		if resp.XChecksum != sourceChecksum {
			// There is conflict, log error
			h.logger.Error().
				Str("bucket", f.BucketID).
				Str("fileID", f.FileID).
				Str("version", f.Version).
				Msg("File already exists in destination storage and has different checksum.")
		}
		// Nothing to do
		return false, nil
	}
	// If file not found it needs sync, otherwise return error
	if _, ok := err.(*storage.SyncFileMetadataNotFound); !ok {
		return false, err
	}

	return true, nil
}
