// storageSync is a worker receiving messages from localStorage published to NATS to resiliently sync everything to cloudStorage
package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/consumer"
)

type clientAuthInfoWriter struct {
	apiKey string
}

func (a *clientAuthInfoWriter) AuthenticateRequest(r runtime.ClientRequest, f strfmt.Registry) error {
	r.SetHeaderParam("Authorization", a.apiKey)
	return nil
}

func main() {
	// initialize logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().
		Timestamp().
		Str("service", "storageSync").
		Logger()

	// initialize local storage API client
	localTransportCfg := &client.TransportConfig{
		Host:     "localStorage",
		BasePath: "storage",
		Schemes:  []string{"https"},
	}
	localClient := client.NewHTTPClientWithConfig(strfmt.NewFormats(), localTransportCfg)

	// initialize cloud storage API client
	cloudTransportCfg := &client.TransportConfig{
		Host:     "cloudStorage",
		BasePath: "storage",
		Schemes:  []string{"https"},
	}
	cloudClient := client.NewHTTPClientWithConfig(strfmt.NewFormats(), cloudTransportCfg)

	// initizalize mock request authenticator
	auth := &clientAuthInfoWriter{"SECRETSECRETSECRETSECRETSECRETSE"}

	// initialize handlers
	handlers := consumer.NewHandlers(localClient.Storage, auth, cloudClient.Storage, auth, logger.With().Str("component", "sync/storage/consumer/handlers").Logger())

	// create nats/nats-streaming connection
	URLs := "tls://nats:secret@localNats:4242"
	ClusterID := "localNats"
	ClientID := "storageSync"
	ClientCert := "/certs/public.crt"
	ClientKey := "/certs/private.key"

	nc, err := nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
	if err != nil {
		log.Fatalln(err)
	}

	sc, err := stan.Connect(ClusterID, ClientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatalln(err)
	}

	// initalize consumer
	c := consumer.New(sc, handlers, time.Duration(time.Second), logger.With().Str("component", "sync/storage/consumer").Logger())
	c.StartSubscription(context.Background(), storageSync.FileNew)
	c.StartSubscription(context.Background(), storageSync.FileUpdate)
	c.StartSubscription(context.Background(), storageSync.FileDelete)

	// Run cleanup when SIGINT is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range signalChan {
			c.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
