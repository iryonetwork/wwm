package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/metrics"
	"github.com/iryonetwork/wwm/reports/filesDataExporter"
	"github.com/iryonetwork/wwm/utils"
)

type (
	exportError struct {
		id  string // identifier of resource that failed to sync
		err error
	}

	batchDataExporter struct {
		handlers          filesDataExporter.Handlers
		bucketsRateLimit  int
		logger            zerolog.Logger
		metricsCollection map[metrics.ID]prometheus.Collector
	}
)

const exportSeconds metrics.ID = "exportSeconds"

const labelFilesCollection = "filesCollection"

var bucketsToSkip = map[string]bool{"encounters": true, "patients": true}

// Export runs files data export for all the files since `dataSince` timestamp
func (s *batchDataExporter) Export(ctx context.Context, dataSince time.Time) error {
	bucketsRateLimit := make(chan bool, s.bucketsRateLimit)

	buckets, err := s.handlers.ListSourceBuckets(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list source buckets")
		return errors.Wrap(err, "failed to list source buckets")
	}

	numberOfBuckets := len(buckets)

	ch := make(chan *exportError)
	for _, b := range buckets {
		if _, ok := bucketsToSkip[b.Name]; !ok {
			go s.exportBucket(ctx, dataSince, b.Name, ch, bucketsRateLimit)
		} else {
			numberOfBuckets--
		}
	}

	var errCount int
	for i := 0; i < numberOfBuckets; i++ {
		syncErr := <-ch
		if syncErr != nil {
			s.logger.Error().Err(syncErr.err).Str("bucket", syncErr.id).Msg("failed to sync")
			errCount++
		}
	}

	if errCount > 0 {
		s.logger.Error().Msgf("%d failure(s) out of %d bucket(s) to export", errCount, len(buckets))
		return errors.Errorf("%d failure(s) out of %d bucket(s) to export", errCount, len(buckets))
	}

	return nil
}

// GetPrometheusMetricsCollection returns all prometheus metrics collectors to be registered
func (s *batchDataExporter) GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector {
	return s.metricsCollection
}

func New(handlers filesDataExporter.Handlers, bucketsRateLimit int, logger zerolog.Logger) filesDataExporter.BatchFilesDataExporter {
	logger = logger.With().Str("component", "reports/filesDataExporter/batch").Logger()

	metricsCollection := make(map[metrics.ID]prometheus.Collector)
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "batch",
		Name:      "file_export_seconds",
		Help:      "Time taken to export file",
	}, []string{"operation", "success", "result"})
	metricsCollection[exportSeconds] = h

	return &batchDataExporter{
		handlers:          handlers,
		bucketsRateLimit:  bucketsRateLimit,
		logger:            logger,
		metricsCollection: metricsCollection,
	}
}

func (s *batchDataExporter) exportBucket(ctx context.Context, dataSince time.Time, bucketID string, errCh chan *exportError, rateLimit chan bool) {
	rateLimit <- true

	files, err := s.handlers.ListSourceFilesAsc(ctx, bucketID, strfmt.DateTime(dataSince))
	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucketID).Msg("failed to list source files")
		errCh <- &exportError{bucketID, errors.Wrap(err, fmt.Sprintf("failed to list source files in bucket %s", bucketID))}
		return
	}

	var exportCount int
	var errCount int

	for _, f := range files {
		select {
		case <-ctx.Done():
			s.logger.Error().Str("bucket", bucketID).Msg("aborting bucket export due to context cancellation")
			errCh <- &exportError{bucketID, errors.Wrap(ctx.Err(), fmt.Sprintf("aborting bucket export due to context cancellation"))}
			return
		default:
			if !utils.SliceContains(f.Labels, labelFilesCollection) && time.Time(f.Created).After(dataSince) {
				exportCount++
				err := s.exportFile(ctx, bucketID, f)
				if err != nil {
					errCount++
				}
			}
		}
	}

	if errCount > 0 {
		s.logger.Error().Str("bucket", bucketID).Msgf("%d failure(s) out of %d file(s) to export", errCount, exportCount)
		errCh <- &exportError{bucketID, errors.Errorf("%d failure(s) out of %d file(s) to export in bucket %s", errCount, exportCount, bucketID)}
		<-rateLimit
		return
	}

	<-rateLimit
	errCh <- nil
}

func (s *batchDataExporter) exportFile(ctx context.Context, bucketID string, f *models.FileDescriptor) error {
	// Make sure we record duration metrics even if processing fails, set default values for labels
	start := time.Now()
	success := false
	result := filesDataExporter.ResultExportNotNeeded
	defer func() {
		duration := time.Since(start)
		s.metricsCollection[exportSeconds].(*prometheus.HistogramVec).
			With(prometheus.Labels{"operation": string(f.Operation), "success": fmt.Sprintf("%t", success), "result": string(result)}).
			Observe(duration.Seconds())
	}()

	var err error
	switch f.Operation {
	case models.FileDescriptorOperationW:
		result, err = s.handlers.ExportFile(ctx, bucketID, f.Name, f.Version, f.Created)
	case models.FileDescriptorOperationD:
		result, err = s.handlers.ExportFileDelete(ctx, bucketID, f.Name, f.Version, f.Created)
	}

	if err != nil {
		s.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("file", f.Name).
			Str("version", f.Version).
			Str("operation", string(f.Operation)).
			Msg("failed to export")
	} else {
		success = true
		s.logger.Info().
			Str("bucket", bucketID).
			Str("file", f.Name).
			Str("version", f.Version).
			Str("operation", string(f.Operation)).
			Msg("successfully exported")
	}

	return err
}
