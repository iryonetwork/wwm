package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/zerolog"

	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/publisher"
)

var (
	time1, _ = strfmt.ParseDateTime("2018-02-05T15:16:15.123Z")
	file     = &storageSync.FileInfo{"bucket", "file", "version"}
)

func main() {
	URLs := "tls://nats:secret@nats.iryo.local:4322"
	ClusterID := "nats_streaming_cluster"
	ClientID := "tmp_publisher"
	ClientCert := "/Users/mateuszkrasucki/go/src/github.com/iryonetwork/wwm/bin/tls/storageSyncPublisher.pem"
	ClientKey := "/Users/mateuszkrasucki/go/src/github.com/iryonetwork/wwm/bin/tls/storageSyncPublisher-key.pem"

	nc, _ := nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
	sc, _ := stan.Connect(ClusterID, ClientID, stan.NatsConn(nc))

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Str("component", "sync/storage/publisher").Str("natsClientID", ClientID).Timestamp().Logger()

	p := publisher.New(sc, 5, time.Duration(time.Second), 1.5, logger)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for {
			select {
			case <-time.Tick(2 * time.Second):
				p.PublishAsyncWithRetries(storageSync.FileNew, file)
			case <-signalChan:
				fmt.Printf("\nReceived an interrupt, closing connection...\n\n")
				p.Close()
				cleanupDone <- true
			}
		}
	}()
	<-cleanupDone
}
