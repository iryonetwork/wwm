package consumer

//go:generate sh ../../../bin/mockgen.sh sync/storage/consumer Handlers $GOFILE

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/zerolog"

	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

type contextKey string

const subID contextKey = "ID"

type stanConsumer struct {
	ctx      context.Context
	conn     stan.Conn
	ackWait  time.Duration
	subs     []stan.Subscription
	subsLock sync.Mutex
	handlers Handlers
	logger   zerolog.Logger
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

func (c *stanConsumer) getMsgHandler(ctx context.Context, typ storageSync.EventType, h Handler) stan.MsgHandler {
	return func(msg *stan.Msg) {
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

		err = h(f)
		if err != nil {
			c.logger.Error().Err(err).
				Str("cmd", "MsgHandler").
				Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
				Msg("Failed handler invocation")

			return
		}

		// Acknowledge the message
		msg.Ack()
		c.logger.Debug().
			Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
			Str("cmd", "MsgHandler").
			Msgf("Acknowledged message: %s", msg)
	}
}

// New returns new consumer service with provided nats-streaming connection as underlying backend.
func New(ctx context.Context, sc stan.Conn, handlers Handlers, ackWait time.Duration, logger zerolog.Logger) storageSync.Consumer {
	c := &stanConsumer{
		ctx:      ctx,
		conn:     sc,
		handlers: handlers,
		ackWait:  ackWait,
		logger:   logger,
	}

	// Close if context is Done()
	go func() {
		<-ctx.Done()
		c.Close()
	}()

	return c
}
