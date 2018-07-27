package filesDataExporter

//go:generate ../../bin/mockgen.sh reports/filesDataExporter Handlers $GOFILE

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/metrics"
)

type (
	// ExportResult defines file data export result
	ExportResult string

	// BatchFilesDataExporter defines public API of batch files data exporter
	BatchFilesDataExporter interface {
		// Export runs files data export for all the files since last successful run
		Export(ctx context.Context, lastSuccessfulRun time.Time) error
		// GetPrometheusMetricsCollection returns metrics to be registered for the component.
		GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector
	}
	// Sanitizer describes data sanitizer's public interface
	Sanitizer interface {
		// Sanitize sanitizes JSON string by encrypting values and/or removing certain keys
		Sanitize(context.Context, []byte) ([]byte, error)
	}

	// Handlers describes public API for reports data export handlers
	Handlers interface {
		// ExportFile export new files and file updates to reports storage.
		ExportFile(ctx context.Context, bucketID, fileID, version strfmt.UUID, timestamp strfmt.DateTime) (ExportResult, error)
		// ExportFileDelete export file deletion to reports storage.
		ExportFileDelete(ctx context.Context, bucketID, fileID strfmt.UUID, version strfmt.UUID, timestamp strfmt.DateTime) (ExportResult, error)
		// ListSourceBuckets lists all the buckets in source storage.
		ListSourceBuckets(ctx context.Context) ([]*models.BucketDescriptor, error)
		// ListSourceFiles lists all the files in the bucket of source storage including files marked as delete, ascending order by Created timestamp ensured.
		ListSourceFilesAsc(ctx context.Context, bucketID strfmt.UUID) ([]*models.FileDescriptor, error)
	}

	// Handler describes reports/filesDataExporter handler function
	Handler func(ctx context.Context, bucketID, fileID, version strfmt.UUID, created strfmt.DateTime) (ExportResult, error)

	// FileInfo defines file info struct
	FileInfo struct {
		BucketID string          `json:"bucketID,omitempty"`
		FileID   string          `json:"fileID,omitempty"`
		Version  string          `json:"version,omitempty"`
		Created  strfmt.DateTime `json:"created,omitempty"`
	}

	// FieldToSanitize defines specification for sanitizing individual field
	FieldToSanitize struct {
		Type                     string                 `json:"type"`
		EhrPath                  string                 `json:"ehrPath"`
		Transformation           string                 `json:"transformation,omitempty"`
		TransformationParameters map[string]interface{} `json:"transformationParameters,omitempty"`
		Items                    []FieldToSanitize      `json:"items,omitempty"`
	}
)

var ResultExported ExportResult = "exported"
var ResultError ExportResult = "error"
var ResultExportNotNeeded ExportResult = "exportNotNeeded"

func NewFileInfo() *FileInfo {
	return &FileInfo{}
}

func (f *FileInfo) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FileInfo) Unmarshal(m []byte) error {
	return json.Unmarshal(m, f)
}
