package consumer

import (
	"github.com/go-openapi/runtime"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client/storage"
	"github.com/iryonetwork/wwm/storage/s3"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

// Handlers describes public API for sync/storage event handlers
type Handlers interface {
	FileNew(f *storageSync.FileInfo) error
	FileUpdate(f *storageSync.FileInfo) error
	FileDelete(f *storageSync.FileInfo) error
}

// Handler describes handler functions
type Handler func(f *storageSync.FileInfo) error

type handlers struct {
	localStorage        s3.Storage
	remoteStorageClient *storage.Client
	auth                runtime.ClientAuthInfoWriter
	logger              zerolog.Logger
}

// FileNew uploads new file to cloud storage.
func (h *handlers) FileNew(f *storageSync.FileInfo) error {
	return h.fileSync(f)
}

// FileUpdate uploads updated file to cloud storage.
func (h *handlers) FileUpdate(f *storageSync.FileInfo) error {
	return h.fileSync(f)
}

// FileDelete deletes file to cloud storage.
func (h *handlers) FileDelete(f *storageSync.FileInfo) error {
	params := storage.NewFileDeleteParams().WithBucket(f.BucketID).WithFileID(f.FileID)

	_, err := h.remoteStorageClient.FileDelete(params, h.auth)

	return err
}

// NewApiHandlers returns Handlers with cloudStorage and localStorage API used.
func NewHandlers(localStorage s3.Storage, remoteStorageClient *storage.Client, auth runtime.ClientAuthInfoWriter, logger zerolog.Logger) Handlers {
	return &handlers{
		localStorage:        localStorage,
		remoteStorageClient: remoteStorageClient,
		auth:                auth,
		logger:              logger,
	}
}

func (h *handlers) fileSync(f *storageSync.FileInfo) error {
	params := storage.NewFileSyncParams().WithBucket(f.BucketID).WithFileID(f.FileID).WithVersion(f.Version)

	r, fd, err := h.localStorage.Read(f.BucketID, f.FileID, f.Version)
	if err != nil {
		if err == s3.ErrNotFound {
			h.logger.Info().
				Str("bucket", f.BucketID).
				Str("fileID", f.FileID).
				Str("version", f.Version).
				Msg("File does not exist in local storage.")

			// File might have been already deleted
			return nil
		}
		return err
	}

	if fd.Archetype != "" {
		params.SetArchetype(&fd.Archetype)
	}
	params.SetContentType(fd.ContentType)
	params.SetCreated(fd.Created)
	params.SetFile(runtime.NamedReader("FileReader", r))

	ok, created, err := h.remoteStorageClient.FileSync(params, h.auth)

	switch {
	case ok != nil:
		h.logger.Info().
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("File already exists in remote storage.")
	case created != nil:
		h.logger.Debug().
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("Succesfully synchronized file to remote storage")
	case err != nil:
		h.logger.Error().Err(err).
			Str("bucket", f.BucketID).
			Str("fileID", f.FileID).
			Str("version", f.Version).
			Msg("Failed to synchornize file to remote storage")
	}

	return err
}
