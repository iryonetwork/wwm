// storageSync is a worker receiving messages from localStorage published to NATS to resiliently sync everything to cloudStorage
package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/storageSync"
	"github.com/iryonetwork/wwm/storageSync/consumer"
)

func main() {
	URLs := "tls://nats:secret@nats.iryo.local:4322"
	ClusterID := "nats_streaming_cluster"
	ClientID := "storageSyncConsumer"
	ClientCert := "/Users/mateuszkrasucki/go/src/github.com/iryonetwork/wwm/bin/tls/storageSyncConsumer.pem"
	ClientKey := "/Users/mateuszkrasucki/go/src/github.com/iryonetwork/wwm/bin/tls/storageSyncConsumer-key.pem"

	nc, _ := nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
	sc, _ := stan.Connect(ClusterID, ClientID, stan.NatsConn(nc))

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Str("component", "storageSync/consumer").Str("natsClientID", ClientID).Timestamp().Logger()
	handlers := consumer.NewHandlers(nil, nil, nil, logger)

	c, err := consumer.New(sc, handlers, time.Duration(time.Second), logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize consumer")
	}
	c.StartSubscription(storageSync.FileNew)
	c.StartSubscription(storageSync.FileUpdate)
	c.StartSubscription(storageSync.FileDelete)

	// Run cleanup when SIGINT is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			c.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
