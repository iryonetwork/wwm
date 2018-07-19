package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/metrics"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

type Cfg struct {
	BucketsRateLimit        int
	FilesPerBucketRateLimit int
}

type syncError struct {
	id  string // identifier of resource that failed to sync
	err error
}

type batchStorageSync struct {
	handlers                storageSync.Handlers
	bucketsRateLimit        int
	filesPerBucketRateLimit int
	logger                  zerolog.Logger
	metricsCollection       map[metrics.ID]prometheus.Collector
}

const syncSeconds metrics.ID = "syncSeconds"

func (s *batchStorageSync) Sync(ctx context.Context, lastSuccessfulRun time.Time) error {
	bucketRateLimit := make(chan bool, s.bucketsRateLimit)

	buckets, err := s.handlers.ListSourceBuckets(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list source buckets")
		return errors.Wrap(err, "failed to list source buckets")
	}

	ch := make(chan *syncError)
	for _, b := range buckets {
		go s.syncBucket(ctx, lastSuccessfulRun, b.Name, ch, bucketRateLimit)
	}

	var errCount int
	for i := 0; i < len(buckets); i++ {
		syncErr := <-ch
		if syncErr != nil {
			s.logger.Error().Err(syncErr.err).Str("bucket", syncErr.id).Msg("failed to sync")
			errCount++
		}
	}

	if errCount > 0 {
		s.logger.Error().Msgf("%d failure(s) out of %d bucket(s) to sync", errCount, len(buckets))
		return errors.Errorf("%d failure(s) out of %d bucket(s) to sync", errCount, len(buckets))
	}

	return nil
}

// GetPrometheusMetricsCollection returns all prometheus metrics collectors to be registered
func (s *batchStorageSync) GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector {
	return s.metricsCollection
}

func New(handlers storageSync.Handlers, cfg Cfg, logger zerolog.Logger) storageSync.BatchSync {
	logger = logger.With().Str("component", "sync/storage/batch").Logger()

	metricsCollection := make(map[metrics.ID]prometheus.Collector)
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "batch",
		Name:      "file_version_sync_seconds",
		Help:      "Time taken to sync file",
	}, []string{"operation", "success", "result"})
	metricsCollection[syncSeconds] = h

	return &batchStorageSync{
		handlers:                handlers,
		bucketsRateLimit:        cfg.BucketsRateLimit,
		filesPerBucketRateLimit: cfg.FilesPerBucketRateLimit,
		logger:                  logger,
		metricsCollection:       metricsCollection,
	}
}

func (s *batchStorageSync) syncBucket(ctx context.Context, lastSuccessfulRun time.Time, bucketID string, errCh chan *syncError, rateLimit chan bool) {
	rateLimit <- true

	fileRateLimit := make(chan bool, s.filesPerBucketRateLimit)

	files, err := s.handlers.ListSourceFilesAsc(ctx, bucketID)
	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucketID).Msg("failed to list source files")
		errCh <- &syncError{bucketID, errors.Wrap(err, fmt.Sprintf("failed to list source files in bucket %s", bucketID))}
		return
	}

	var syncCount int
	var errCount int
	ch := make(chan *syncError)

	for _, f := range files {
		if time.Time(f.Created).After(lastSuccessfulRun) {
			syncCount++
			go s.syncFile(ctx, lastSuccessfulRun, bucketID, f.Name, ch, fileRateLimit)
		}
	}

	for i := 0; i < syncCount; i++ {
		syncErr := <-ch
		if syncErr != nil {
			s.logger.Error().Err(syncErr.err).Str("bucket", bucketID).Str("file", syncErr.id).Msg("failed to sync")
			errCount++
		}
	}

	if errCount > 0 {
		s.logger.Error().Str("bucket", bucketID).Msgf("%d failure(s) out of %d file(s) to sync", errCount, syncCount)
		errCh <- &syncError{bucketID, errors.Errorf("%d failure(s) out of %d file(s) to sync in bucket %s", errCount, syncCount, bucketID)}
		<-rateLimit
		return
	}
	<-rateLimit

	errCh <- nil
}

func (s *batchStorageSync) syncFile(ctx context.Context, lastSuccessfulRun time.Time, bucketID, fileID string, errCh chan *syncError, rateLimit chan bool) {
	rateLimit <- true

	versions, err := s.handlers.ListSourceFileVersionsAsc(ctx, bucketID, fileID)
	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucketID).Str("file", fileID).Msg("failed to list source versions")
		errCh <- &syncError{fileID, errors.Wrap(err, fmt.Sprintf("failed to list source versions of file %s in bucket %s", fileID, bucketID))}
		return
	}

	var syncCount int
	var errCount int

	for _, f := range versions {
		select {
		case <-ctx.Done():
			s.logger.Error().Str("bucket", bucketID).Str("file", fileID).Msg("aborting file sync due to context cancellation")
			errCh <- &syncError{fileID, errors.Wrap(ctx.Err(), fmt.Sprintf("aborting file sync due to context cancellation"))}
			return
		default:
			if time.Time(f.Created).After(lastSuccessfulRun) {
				syncCount++
				err := s.syncFileVersion(ctx, bucketID, fileID, f)
				if err != nil {
					errCount++
				}
			}
		}
	}

	if errCount > 0 {
		s.logger.Error().Str("bucket", bucketID).Str("file", fileID).Msgf("%d failure(s) out of %d version(s) to sync", errCount, syncCount)
		errCh <- &syncError{fileID, errors.Errorf("%d failure(s) out of %d version(s) to sync for file %s in bucket %s", errCount, syncCount, fileID, bucketID)}
		<-rateLimit
		return
	}

	<-rateLimit
	errCh <- nil
}

func (s *batchStorageSync) syncFileVersion(ctx context.Context, bucketID, fileID string, f *models.FileDescriptor) error {
	// Make sure we record duration metrics even if processing fails, set default values for labels
	start := time.Now()
	success := false
	result := storageSync.ResultSyncNotNeeded
	defer func() {
		duration := time.Since(start)
		s.metricsCollection[syncSeconds].(*prometheus.HistogramVec).
			With(prometheus.Labels{"operation": string(f.Operation), "success": fmt.Sprintf("%t", success), "result": string(result)}).
			Observe(duration.Seconds())
	}()

	var err error
	switch f.Operation {
	case models.FileDescriptorOperationW:
		result, err = s.handlers.SyncFile(ctx, bucketID, fileID, f.Version, f.Created)
	case models.FileDescriptorOperationD:
		result, err = s.handlers.SyncFileDelete(ctx, bucketID, fileID, f.Version, f.Created)
	}

	if err != nil {
		s.logger.Error().Err(err).
			Str("bucket", bucketID).
			Str("file", fileID).
			Str("version", f.Version).
			Str("operation", string(f.Operation)).
			Msg("failed to sync")
	} else {
		success = true
		s.logger.Info().
			Str("bucket", bucketID).
			Str("file", fileID).
			Str("version", f.Version).
			Str("operation", string(f.Operation)).
			Msg("successfully synced")
	}

	return err
}
