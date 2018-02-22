package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

type batchStorageSync struct {
	handlers storageSync.Handlers
	logger   zerolog.Logger
}

func (s *batchStorageSync) Sync(ctx context.Context, lastSuccessfulRun time.Time) error {
	buckets, err := s.handlers.ListSourceBuckets(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list source buckets")
		return errors.Wrap(err, "failed to list source buckets")
	}

	errChannels := make(map[string]chan error)
	for _, b := range buckets {
		ch := make(chan error)
		errChannels[b.Name] = ch

		go s.syncBucket(ctx, lastSuccessfulRun, b.Name, ch)
	}

	var errCount int
	for bucketID, ch := range errChannels {
		err := <-ch
		if err != nil {
			s.logger.Error().Err(err).Str("bucket", bucketID).Msg("failed to sync")
			errCount++
		}
	}

	if errCount > 0 {
		s.logger.Error().Msgf("%d failure(s) out of %d bucket(s) to sync", errCount, len(errChannels))
		return errors.Errorf("%d failure(s) out of %d bucket(s) to sync", errCount, len(errChannels))
	}

	return nil
}

func New(handlers storageSync.Handlers, logger zerolog.Logger) storageSync.BatchSync {
	return &batchStorageSync{
		handlers: handlers,
		logger:   logger,
	}
}

func (s *batchStorageSync) syncBucket(ctx context.Context, lastSuccessfulRun time.Time, bucketID string, errCh chan error) {
	files, err := s.handlers.ListSourceFiles(ctx, bucketID)
	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucketID).Msg("failed to list source files")
		errCh <- errors.Wrap(err, fmt.Sprintf("failed to list source files in bucket %s", bucketID))
		return
	}

	errChannels := make(map[string]chan error)
	for _, f := range files {
		if time.Time(f.Created).After(lastSuccessfulRun) {
			ch := make(chan error)
			errChannels[f.Name] = ch

			go s.syncFile(ctx, lastSuccessfulRun, bucketID, f.Name, ch)
		}
	}

	var errCount int
	for fileID, ch := range errChannels {
		err := <-ch
		if err != nil {
			s.logger.Error().Err(err).Str("bucket", bucketID).Str("file", fileID).Msg("failed to sync")
			errCount++
		}
	}

	if errCount > 0 {
		s.logger.Error().Str("bucket", bucketID).Msgf("%d failure(s) out of %d file(s) to sync", errCount, len(errChannels))
		errCh <- errors.Errorf("%d failure(s) out of %d file(s) to sync in bucket %s", errCount, len(errChannels), bucketID)
		return
	}
	errCh <- nil
}

func (s *batchStorageSync) syncFile(ctx context.Context, lastSuccessfulRun time.Time, bucketID, fileID string, errCh chan error) {
	versions, err := s.handlers.ListSourceFileVersions(ctx, bucketID, fileID)
	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucketID).Str("file", fileID).Msg("failed to list source versions")
		errCh <- errors.Wrap(err, fmt.Sprintf("failed to list source versions of file %s in bucket %s", fileID, bucketID))
		return
	}

	var syncCount int
	var errCount int

	for _, f := range versions {
		select {
		case <-ctx.Done():
			s.logger.Error().Str("bucket", bucketID).Str("file", fileID).Msg("aborting file sync due to context cancellation")
			errCh <- errors.Wrap(ctx.Err(), fmt.Sprintf("aborting file sync due to context cancellation"))
			return
		default:
			if time.Time(f.Created).After(lastSuccessfulRun) {
				syncCount++
				switch f.Operation {
				case models.FileDescriptorOperationW:
					err = s.handlers.SyncFile(ctx, bucketID, fileID, f.Version, f.Created)
					if err != nil {
						errCount++
						s.logger.Error().Err(err).Str("bucket", bucketID).Str("file", fileID).Str("version", f.Version).Msg("failed to sync")
					} else {
						s.logger.Info().Str("bucket", bucketID).Str("file", fileID).Str("version", f.Version).Msg("successfully synced")
					}
				case models.FileDescriptorOperationD:
					err = s.handlers.SyncFileDelete(ctx, bucketID, fileID, f.Version, f.Created)
					if err != nil {
						errCount++
						s.logger.Error().Err(err).Str("bucket", bucketID).Str("file", fileID).Str("version", f.Version).Msg("failed to sync deletion")
					} else {
						s.logger.Info().Str("bucket", bucketID).Str("file", fileID).Str("version", f.Version).Msg("successfully synced deletion")
					}
				}
			}
		}
	}

	if errCount > 0 {
		s.logger.Error().Str("bucket", bucketID).Str("file", fileID).Msgf("%d failure(s) out of %d version(s) to sync", errCount, syncCount)
		errCh <- errors.Errorf("%d failure(s) out of %d version(s) to sync for file %s in bucket %s", errCount, syncCount, fileID, bucketID)
		return
	}
	errCh <- nil
}
