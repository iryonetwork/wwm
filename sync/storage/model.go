package storage

//go:generate ../../bin/mockgen.sh sync/storage Publisher,Consumer,Handlers $GOFILE

import (
	"context"
	"encoding/json"
	"time"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/iryonetwork/wwm/metrics"
)

// EventType defines event type
type EventType string

// SyncResult defines synchronization result
type SyncResult string

// EventType constants
const (
	FileNew    EventType = "file.new"
	FileUpdate EventType = "file.update"
	FileDelete EventType = "file.delete"
)

// Publisher describes sync/storage publisher public methods.
type Publisher interface {
	// Publish pushes sync/storage event and returns synchronous response.
	Publish(ctx context.Context, typ EventType, f *FileInfo) error
	// Publish starts goroutine that pushes storageSync event and retries if publishing failed.
	PublishAsyncWithRetries(ctx context.Context, typ EventType, f *FileInfo) error
	// Close waits for all async publish routines to finish and closes underlying connection.
	Close()
	// GetPrometheusMetricsCollection returns metrics to be registered for the component.
	GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector
}

// Consumer describes public methods of consumer used by storageSync service.
type Consumer interface {
	// StartConsumer starts consumer following service configration.
	StartSubscription(typ EventType) error
	// Returns number of subscriptions within consumer instance.
	GetNumberOfSubsriptions() int
	// Close closes all service consumers.
	Close()
	// GetPrometheusMetricsCollection returns metrics to be registered for the component.
	GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector
}

type BatchSync interface {
	Sync(ctx context.Context, lastSuccessfulRun time.Time) error
	// GetPrometheusMetricsCollection returns metrics to be registered for the component.
	GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector
}

type FileInfo struct {
	BucketID string          `json:"bucketID,omitempty"`
	FileID   string          `json:"fileID,omitempty"`
	Version  string          `json:"version,omitempty"`
	Created  strfmt.DateTime `json:"created,omitempty"`
}

var ResultSynced SyncResult = "synced"
var ResultConflict SyncResult = "conflict"
var ResultError SyncResult = "error"
var ResultSyncNotNeeded SyncResult = "syncNotNeeded"

func NewFileInfo() *FileInfo {
	return &FileInfo{}
}

func (f *FileInfo) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FileInfo) Unmarshal(m []byte) error {
	return json.Unmarshal(m, f)
}
