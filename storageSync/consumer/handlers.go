package consumer

import (
	"fmt"
	"math/rand"

	"github.com/go-openapi/runtime"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3"
)

// Handlers describes public API for storageSync event handlers
type Handlers interface {
	FileNew(fd *models.FileDescriptor) error
	FileUpdate(fd *models.FileDescriptor) error
	FileDelete(fd *models.FileDescriptor) error
}

// Handler describes handler functions
type Handler func(fd *models.FileDescriptor) error

type handlers struct {
	storage          s3.Storage
	storageAPIClient *client.Storage
	auth             runtime.ClientAuthInfoWriter
	logger           zerolog.Logger
}

// FileNew uploads new file to cloud storage.
func (h *handlers) FileNew(fd *models.FileDescriptor) error {
	h.logger.Info().Msg("FileNew API handler is not yet implemented! random success")
	if rand.Float32() > 0.5 {
		return nil
	}
	return fmt.Errorf("error")
}

// FileUpdate uploads updated file to cloud storage.
func (h *handlers) FileUpdate(fd *models.FileDescriptor) error {
	h.logger.Info().Msg("FileUpdate API handler is not yet implemented! random success")
	if rand.Float32() > 0.5 {
		return nil
	}
	return fmt.Errorf("error")
}

// FileDelete deletes file to cloud storage.
func (h *handlers) FileDelete(fd *models.FileDescriptor) error {
	h.logger.Info().Msg("FileDelete API handler is not yet implemented! random success")
	if rand.Float32() > 0.5 {
		return nil
	}
	return fmt.Errorf("error")
}

// NewApiHandlers returns Handlers with cloudStorage and localStorage API used.
func NewHandlers(storage s3.Storage, storageAPIClient *client.Storage, auth runtime.ClientAuthInfoWriter, logger zerolog.Logger) Handlers {
	return &handlers{
		storage:          storage,
		storageAPIClient: storageAPIClient,
		auth:             auth,
		logger:           logger,
	}
}
