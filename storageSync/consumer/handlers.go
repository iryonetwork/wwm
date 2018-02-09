package consumer

import (
	"github.com/go-openapi/runtime"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client/storage"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storageSync"
)

// Handlers describes public API for storageSync event handlers
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
	return nil
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
	return nil
}
