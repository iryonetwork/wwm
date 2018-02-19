package publisher

//go:generate ../../../bin/mockgen.sh sync/storage/publisher StanConnection $GOFILE

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

// Connection interface describes actions that have to be supported by underlying connection with nats-streaming. For testing purposes.
type StanConnection interface {
	Publish(subject string, data []byte) error
	Close() error
}

type Cfg struct {
	Connection      StanConnection
	Retries         int
	StartRetryWait  time.Duration
	RetryWaitFactor float32
}

type stanPublisher struct {
	ctx             context.Context
	conn            StanConnection
	retries         int
	startRetryWait  time.Duration
	retryWaitFactor float32
	wg              sync.WaitGroup
	logger          zerolog.Logger
	publishSeconds  prometheus.Histogram
	publishCalls    prometheus.Counter
}

type nullPublisher struct {
}

// Publish of nullPublisher does nothing.
func (p *nullPublisher) Publish(_ context.Context, _ storageSync.EventType, _ *storageSync.FileInfo) error {
	return nil
}

// PublishAsyncWithRetries of nullPublisher does nothing.
func (p *nullPublisher) PublishAsyncWithRetries(_ context.Context, _ storageSync.EventType, _ *storageSync.FileInfo) error {
	return nil
}

// Close of nullPublisher does nothing.
func (p *nullPublisher) Close() {
	return
}

// Publish pushes sync/storage event and returns synchronous response.
func (p *stanPublisher) Publish(_ context.Context, typ storageSync.EventType, f *storageSync.FileInfo) error {
	msg, err := f.Marshal()
	if err != nil {
		p.logger.Error().Err(err).
			Str("cmd", "Publish").
			Msg("Failed to marshal file info")

		return err
	}

	err = p.conn.Publish(string(typ), msg)
	if err != nil {
		p.logger.Error().Err(err).
			Str("cmd", "Publish").
			Msg("Failed to publish storage sync event")
		return err
	}
	return nil
}

// Publish starts goroutine that pushes sync/storage events and retries if publishing failed.
func (p *stanPublisher) PublishAsyncWithRetries(ctx context.Context, typ storageSync.EventType, f *storageSync.FileInfo) error {
	// Make sure we record duration metrics even if processing fails
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		p.publishSeconds.Observe(duration.Seconds())
	}()

	msg, err := f.Marshal()
	if err != nil {
		p.logger.Error().Err(err).
			Msg("Failed to marshal file descriptor")
		return err
	}

	p.wg.Add(1)
	go func() {
		var err error
		retryWait := p.startRetryWait
	RetryLoop:
		for i := 0; i < p.retries; i++ {
			select {
			case <-ctx.Done():
				p.logger.Debug().Err(ctx.Err()).
					Str("cmd", "PublishAsyncWithRetries").
					Str("type", string(typ)).
					Msg("Async publishing stopped due to context cancellation")
				break RetryLoop
			default:
				err = p.conn.Publish(string(typ), msg)
				p.publishCalls.Inc() // increase publish calls counter metrycs

				if err == nil {
					p.logger.Debug().
						Str("cmd", "PublishAsyncWithRetries").
						Str("type", string(typ)).
						Msgf("%s", msg)

					p.wg.Done()
					return
				}
				p.logger.Error().Err(err).
					Str("cmd", "PublishAsyncWithRetries").
					Msgf("Failed to publish storage sync event, retry in %s", retryWait)

				time.Sleep(retryWait)
				retryWait = time.Duration(float32(retryWait) * p.retryWaitFactor)
			}
		}
		if err != nil {
			// TODO: handle failure to publish, e.g. write messages to file that can be read later
			p.logger.Error().Err(err).
				Str("cmd", "PublishAsyncWithRetries").
				Msg("Failed to publish storage sync event, maximum number of retries reached")
		}
		p.wg.Done()
	}()

	return nil
}

// Close waits for all async publish routines to finish and closes underlying connection.
func (p *stanPublisher) Close() {
	p.wg.Wait()
	p.conn.Close()

	// unregister metrics
	prometheus.Unregister(p.publishSeconds)
	prometheus.Unregister(p.publishCalls)
}

// New returns new stanPublisher with provided nats-streaming connectiom as underlying backend.
func New(ctx context.Context, cfg Cfg, logger zerolog.Logger, durationMetrics prometheus.Histogram, countMetrics prometheus.Counter) storageSync.Publisher {
	p := &stanPublisher{
		ctx:             ctx,
		conn:            cfg.Connection,
		retries:         cfg.Retries,
		startRetryWait:  cfg.StartRetryWait,
		retryWaitFactor: cfg.RetryWaitFactor,
		logger:          logger,
		publishSeconds:  durationMetrics,
		publishCalls:    countMetrics,
	}
	// Close if context is Done()
	go func() {
		<-ctx.Done()
		p.Close()
	}()

	return p
}

// NewNullPublisher returns new nullPublisher for skipping publishing.
func NewNullPublisher(_ context.Context) storageSync.Publisher {
	return &nullPublisher{}
}
