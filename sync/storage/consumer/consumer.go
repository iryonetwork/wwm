package consumer

//go:generate sh ../../../bin/mockgen.sh sync/storage/consumer Handlers $GOFILE

import (
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/zerolog"

	storageSync "github.com/iryonetwork/wwm/sync/storage"
)

type stanConsumer struct {
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
	var mh stan.MsgHandler
	switch typ {
	case storageSync.FileNew:
		mh = c.getMsgHandler(ID, typ, c.handlers.FileNew)
	case storageSync.FileUpdate:
		mh = c.getMsgHandler(ID, typ, c.handlers.FileUpdate)
	case storageSync.FileDelete:
		mh = c.getMsgHandler(ID, typ, c.handlers.FileDelete)
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
	for _, sub := range c.subs {
		sub.Close()
	}
	c.conn.Close()
}

func (c *stanConsumer) getMsgHandler(ID int, typ storageSync.EventType, h Handler) stan.MsgHandler {
	return func(msg *stan.Msg) {
		c.logger.Info().
			Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
			Msgf("Received message: %s", msg)

		f := storageSync.NewFileInfo()
		err := f.Unmarshal(msg.Data)
		if err != nil {
			c.logger.Error().Err(err).
				Str("cmd", "getMsgHandler").
				Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
				Msg("Failed to unmarshal message")

			return
		}

		err = h(f)
		if err != nil {
			c.logger.Error().Err(err).
				Str("cmd", "getMsgHandler").
				Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
				Msg("Failed handler invocation")

			return
		}

		// Acknowledge the message
		msg.Ack()

		c.logger.Info().
			Str("subscription", fmt.Sprintf("%s:%d", typ, ID)).
			Str("cmd", "getMsgHandler").
			Msgf("Acknowledged message: %s", msg)
	}
}

// New returns new consumer service with provided nats-streaming connection as underlying backend.
func New(sc stan.Conn, handlers Handlers, ackWait time.Duration, logger zerolog.Logger) storageSync.Consumer {
	return &stanConsumer{
		conn:     sc,
		handlers: handlers,
		ackWait:  ackWait,
		logger:   logger,
	}
}
