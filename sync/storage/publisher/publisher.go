package publisher

//go:generate sh ../../../bin/mockgen.sh sync/storage/publisher StanConnection $GOFILE

import (
	"sync"
	"time"

	"github.com/rs/zerolog"

	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

// Connection interface describes actions that have to be supported by underlying connection with nats-streaming. For testing purposes.
type StanConnection interface {
	Publish(subject string, data []byte) error
	Close() error
}

type stanPublisher struct {
	conn            StanConnection
	retries         int
	startRetryWait  time.Duration
	retryWaitFactor float32
	wg              sync.WaitGroup
	logger          zerolog.Logger
}

type nullPublisher struct {
}

// Publish of nullPublisher does nothing.
func (p *nullPublisher) Publish(typ storageSync.EventType, f *storageSync.FileInfo) error {
	return nil
}

// PublishAsyncWithRetries of nullPublisher does nothing.
func (p *nullPublisher) PublishAsyncWithRetries(typ storageSync.EventType, f *storageSync.FileInfo) error {
	return nil
}

// Close of nullPublisher does nothing.
func (p *nullPublisher) Close() {
	return
}

// Publish pushes sync/storage event and returns synchronous response.
func (p *stanPublisher) Publish(typ storageSync.EventType, f *storageSync.FileInfo) error {
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
func (p *stanPublisher) PublishAsyncWithRetries(typ storageSync.EventType, f *storageSync.FileInfo) error {
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
		for i := 0; i < p.retries; i++ {
			err = p.conn.Publish(string(typ), msg)
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
}

// New returns new stanPublisher with provided nats-streaming connectiom as underlying backend.
func New(sc StanConnection, retries int, startRetryWait time.Duration, retryWaitFactor float32, logger zerolog.Logger) storageSync.Publisher {
	return &stanPublisher{
		conn:            sc,
		logger:          logger,
		retries:         retries,
		startRetryWait:  startRetryWait,
		retryWaitFactor: retryWaitFactor,
	}
}

// NewNullPublisher returns new nullPublisher for skipping publishing.
func NewNullPublisher() storageSync.Publisher {
	return &nullPublisher{}
}
