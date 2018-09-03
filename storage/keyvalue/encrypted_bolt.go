package keyvalue

import (
	"context"
	"fmt"
	"time"

	"github.com/iryonetwork/wwm/log/errorChecker"
	bolt "github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/metrics"
)

type encryptedBoltKeyValue struct {
	ctx               context.Context
	db                *bolt.DB
	logger            zerolog.Logger
	metricsCollection map[metrics.ID]prometheus.Collector
}

// Add item to in-memory key-value storage
func (s *encryptedBoltKeyValue) Add(bucket string, key string, value []byte) error {
	// Make sure we record duration metrics even if processing fails
	start := time.Now()
	success := false
	defer func() {
		duration := time.Since(start)
		s.metricsCollection[operationSeconds].(*prometheus.HistogramVec).
			With(prometheus.Labels{"operation": operationAdd, "success": fmt.Sprintf("%t", success)}).
			Observe(duration.Seconds())
	}()

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			s.logger.Error().Err(err).Str("bucket", bucket).Msg("failed to ensure bucket")
			return errors.Wrapf(err, "failed to ensure bucket %s", bucket)
		}

		if s.Get(bucket, key) != nil {
			s.logger.Error().Str("bucket", bucket).Str("key", key).Msg("key already exists in bucket")
			return errors.Errorf("key %s already exists in bucket %s", key, bucket)
		}

		return b.Put([]byte(key), value)
	})

	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucket).Str("key", key).Msg("failed to add key")
		return errors.Wrapf(err, "failed to add key %s to bucket %s", key, bucket)
	}

	success = true
	return nil
}

// Update item in in-memory key-value storage
func (s *encryptedBoltKeyValue) Update(bucket string, key string, value []byte) error {
	// Make sure we record duration metrics even if processing fails
	start := time.Now()
	success := false
	defer func() {
		duration := time.Since(start)
		s.metricsCollection[operationSeconds].(*prometheus.HistogramVec).
			With(prometheus.Labels{"operation": operationAdd, "success": fmt.Sprintf("%t", success)}).
			Observe(duration.Seconds())
	}()

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			s.logger.Error().Err(err).Str("bucket", bucket).Msg("failed to ensure bucket")
			return errors.Wrapf(err, "failed to ensure bucket %s", bucket)
		}

		return b.Put([]byte(key), value)
	})

	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucket).Str("key", key).Msg("failed to update key")
		return errors.Wrapf(err, "failed to update key %s from bucket %s", key, bucket)
	}

	success = true
	return nil
}

// Get item from in-memory key-value storage
func (s *encryptedBoltKeyValue) Get(bucket string, key string) []byte {
	// Make sure we record duration metrics even if processing fails
	start := time.Now()
	success := false
	defer func() {
		duration := time.Since(start)
		s.metricsCollection[operationSeconds].(*prometheus.HistogramVec).
			With(prometheus.Labels{"operation": operationGet, "success": fmt.Sprintf("%t", success)}).
			Observe(duration.Seconds())
	}()

	var val []byte

	errorChecker.LogError(s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		val = b.Get([]byte(key))
		return nil
	}))

	success = true
	return val
}

// Delete item from in-memory key-value storage
func (s *encryptedBoltKeyValue) Delete(bucket string, key string) error {
	// Make sure we record duration metrics even if processing fails
	start := time.Now()
	success := false
	defer func() {
		duration := time.Since(start)
		s.metricsCollection[operationSeconds].(*prometheus.HistogramVec).
			With(prometheus.Labels{"operation": operationGet, "success": fmt.Sprintf("%t", success)}).
			Observe(duration.Seconds())
	}()

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		return b.Delete([]byte(key))
	})

	if err != nil {
		s.logger.Error().Err(err).Str("bucket", bucket).Str("key", key).Msg("failed to delete key")
		return errors.Wrapf(err, "failed to delete key %s from bucket %s", key, bucket)
	}

	success = true
	return nil
}

// Close releases DB.
func (s *encryptedBoltKeyValue) Close() error {
	err := s.db.Close()
	if err != nil {
		s.logger.Error().Err(err).Msg("error while closing db")
		return errors.Wrap(err, "error while closing db")
	}

	return nil
}

//getPrometheusMetricsCollection returns all prometheus metrics collectors that should be registered to expose metrics of component
func (s *encryptedBoltKeyValue) GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector {
	return s.metricsCollection
}

// NewBolt creates & returns new bolt based key-value storage
func NewEncryptedBolt(ctx context.Context, key []byte, filepath string, logger zerolog.Logger) (Storage, error) {
	logger = logger.With().Str("component", "storage/keyvalue").Logger()

	b, err := bolt.Open(key, filepath, 0644, nil)
	if err != nil {
		logger.Error().Err(err).Str("filepath", filepath).Msg("failed to initialize encrypted bolt key value storage")
		return nil, errors.Wrapf(err, "failed to initialie encrypted bolt key value storage with %s", filepath)
	}

	metricsCollection := make(map[metrics.ID]prometheus.Collector)
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "encyrpted_bolt_key_value",
		Name:      "operation_seconds",
		Help:      "Time taken to execute key-value storage operation",
	}, []string{"operation", "success"})
	metricsCollection[operationSeconds] = h

	s := &encryptedBoltKeyValue{ctx: ctx, db: b, logger: logger, metricsCollection: metricsCollection}

	// close on context done
	go func() {
		<-ctx.Done()
		s.Close()
	}()

	return s, nil
}
