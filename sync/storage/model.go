package storage

//go:generate sh ../../bin/mockgen.sh sync/storage Publisher,Consumer $GOFILE

import (
	"encoding/json"
)

// EventType defines event type
type EventType string

// EventType constants
const (
	FileNew    EventType = "file.new"
	FileUpdate EventType = "file.update"
	FileDelete EventType = "file.delete"
)

// Publisher describes sync/storage publisher public methods.
type Publisher interface {
	// Publish pushes sync/storage event and returns synchronous response.
	Publish(typ EventType, f *FileInfo) error
	// Publish starts goroutine that pushes storageSync event and retries if publishing failed.
	PublishAsyncWithRetries(typ EventType, f *FileInfo) error
	// Close waits for all async publish routines to finish and closes underlying connection.
	Close()
}

// Consumer describes public methods of consumer used by storageSync service.
type Consumer interface {
	// StartConsumer starts consumer following service configration.
	StartSubscription(typ EventType) error
	// Returns number of subscriptions within consumer instance.
	GetNumberOfSubsriptions() int
	// Close closes all service consumers.
	Close()
}

type FileInfo struct {
	BucketID string `json:"bucketID,omitempty"`
	FileID   string `json:"fileID,omitempty"`
	Version  string `json:"version,omitempty"`
}

func NewFileInfo() *FileInfo {
	return &FileInfo{}
}

func (f *FileInfo) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FileInfo) Unmarshal(m []byte) error {
	return json.Unmarshal(m, f)
}
