package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/go-nats-streaming"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/metrics"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

type contextKey string

const subID contextKey = "ID"
const taskSeconds metrics.ID = "taskSeconds"

type Cfg struct {
	Connection stan.Conn
	AckWait    time.Duration
	Handlers   storageSync.Handlers
}

type stanConsumer struct {
	ctx               context.Context
	conn              stan.Conn
	ackWait           time.Duration
	handlers          storageSync.Handlers
	subs              []stan.Subscription
	subsLock          sync.Mutex
	logger            zerolog.Logger
	metricsCollection map[metrics.ID]prometheus.Collector
}

// Start starts new nats-streaming queue subscription.
func (c *stanConsumer) StartSubscription(typ storageSync.EventType) error {
	c.subsLock.Lock()

	// ID is a sequential number of subscription within consumer.
	ID := len(c.subs) + 1
	ctx := context.WithValue(c.ctx, subID, ID)
	var mh stan.MsgHandler
	switch typ {
	case storageSync.FileNew:
		mh = c.getMsgHandler(ctx, typ, c.handlers.SyncFile)
	case storageSync.FileUpdate:
		mh = c.getMsgHandler(ctx, typ, c.handlers.SyncFile)
	case storageSync.FileDelete:
		mh = c.getMsgHandler(ctx, typ, c.handlers.SyncFileDelete)
	default:
		c.subsLock.Unlock()
		return fmt.Errorf("Invalid event type")
	}

	// Subscribe to subject:EventType, queueGroup:EventType, durableName:EventType
	sub, err := c.conn.QueueSubscribe(
		string(typ),
		string(typ),
		mh,
		stan.SetManualAckMode(),
		stan.AckWait(c.ackWait),
		stan.DurableName(string(typ)),
	)

	if err != nil {
		c.logger.Error().Err(err).
			Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
			Str("cmd", "StartSubscription").
			Msg("Failed to start nats-streaming subscription")
	} else {
		c.subs = append(c.subs, sub)
	}
	c.subsLock.Unlock()

	return err
}

// Returns number of subscriptions within consumer instance.
func (c *stanConsumer) GetNumberOfSubsriptions() int {
	return len(c.subs)
}

// Close closes nats-streaming connection
func (c *stanConsumer) Close() {
	c.subsLock.Lock()
	for _, sub := range c.subs {
		sub.Close()
	}
	c.subs = []stan.Subscription{}
	c.subsLock.Unlock()
	c.conn.Close()
}

func (c *stanConsumer) getMsgHandler(ctx context.Context, typ storageSync.EventType, h storageSync.Handler) stan.MsgHandler {
	return func(msg *stan.Msg) {
		// Make sure we record duration metrics even if processing fails
		start := time.Now()
		ack := false
		defer func() {
			duration := time.Since(start)
			c.metricsCollection[taskSeconds].(*prometheus.HistogramVec).
				With(prometheus.Labels{"event": string(typ), "ack": fmt.Sprintf("%t", ack)}).
				Observe(duration.Seconds())
		}()

		ID := ctx.Value(subID).(int)
		c.logger.Debug().
			Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
			Msgf("Received message: %s", msg)

		f := storageSync.NewFileInfo()
		err := f.Unmarshal(msg.Data)
		if err != nil {
			c.logger.Error().Err(err).
				Str("cmd", "MsgHandler").
				Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
				Msg("Failed to unmarshal message")

			return
		}

		err = h(ctx, f.BucketID, f.FileID, f.Version, f.Created)
		if err != nil {
			c.logger.Error().Err(err).
				Str("cmd", "MsgHandler").
				Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
				Msg("Failed handler invocation")

			return
		}

		// Change ack variable value for metrics
		ack = true
		// Acknowledge the message
		msg.Ack()
		c.logger.Debug().
			Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
			Str("cmd", "MsgHandler").
			Msgf("Acknowledged message: %s", msg)
	}
}

// GetPrometheusMetricsCollection returns all prometheus metrics collectors needed to initalize instance of consumer (for registration)
func GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector {
	metricsCollection := make(map[metrics.ID]prometheus.Collector)
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "consumer",
		Name:      "task_seconds",
		Help:      "Time taken to serve tasks",
	}, []string{"event", "ack"})
	metricsCollection[taskSeconds] = h

	return metricsCollection
}

// New returns new consumer service with provided nats-streaming connection as underlying backend.
func New(ctx context.Context, cfg Cfg, logger zerolog.Logger, metricsCollection map[metrics.ID]prometheus.Collector) storageSync.Consumer {
	c := &stanConsumer{
		ctx:               ctx,
		conn:              cfg.Connection,
		handlers:          cfg.Handlers,
		ackWait:           cfg.AckWait,
		logger:            logger,
		metricsCollection: metricsCollection,
	}

	// Close if context is Done()
	go func() {
		<-ctx.Done()
		c.Close()
	}()

	return c
}
