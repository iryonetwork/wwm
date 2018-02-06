package storageSync

//go:generate sh ../bin/mockgen.sh storageSync Publisher,Consumer $GOFILE

import (
	"github.com/iryonetwork/wwm/gen/storage/models"
)

// EventType defines event type
type EventType string

// EventType constants
const (
	FileNew    = "file.new"
	FileUpdate = "file.update"
	FileDelete = "file.delete"
)

// Publisher describes storageSync publisher public API
type Publisher interface {
	// Publish pushes storageSync event and returns synchronous response.
	Publish(typ EventType, fd *models.FileDescriptor) error
	// Publish starts goroutine that pushes storageSync event and retries if publishing failed.
	PublishAsyncWithRetries(typ EventType, fd *models.FileDescriptor) error
	// Close waits for all async publish routines to finish and closes underlying connection.
	Close()
}

// Consumer describes methods used by storageSync/consumer service.
type Consumer interface {
	// StartConsumer starts consumer following service configration.
	StartSubscription(typ EventType) error
	// Returns number of subscriptions within consumer instance.
	GetNumberOfSubsriptions() int
	// Close closes all service consumers.
	Close()
}
