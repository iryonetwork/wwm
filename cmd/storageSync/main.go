// storageSync is a worker receiving messages from localStorage published to NATS to resiliently sync everything to cloudStorage
package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/mock"
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

	// initialize keyProvider
	ctrl := gomock.NewController(nil)
	keys := mock.NewMockKeyProvider(ctrl)
	keys.EXPECT().Get(gomock.Any()).AnyTimes().Return("SECRETSECRETSECRETSECRETSECRETSE", nil)

	// initialize storage
	s3Cfg := &s3.Config{
		Endpoint:     "localMinio:9000",
		AccessKey:    "local",
		AccessSecret: "localminio",
		Secure:       true,
		Region:       "us-east-1",
	}
	s3, err := s3.New(s3Cfg, keys, logger.With().Str("component", "storage/s3").Logger())
	if err != nil {
		log.Fatalln(err)
	}

	// initialize remote storage client
	transportCfg := &client.TransportConfig{
		Host:     "iryo.local",
		BasePath: "storage",
		Schemes:  []string{"https"},
	}
	client := client.NewHTTPClientWithConfig(strfmt.NewFormats(), transportCfg)

	// initizalize mock request authenticator
	auth := &clientAuthInfoWriter{"SECRETSECRETSECRETSECRETSECRETSE"}

	// initialize handlers
	handlers := consumer.NewHandlers(s3, client.Storage, auth, logger.With().Str("component", "sync/storage/consumer/handlers").Logger())

	// create nats/nats-streaming connection
	URLs := "tls://nats:secret@nats.iryo.local:4322"
	ClusterID := "localnats"
	ClientID := "storageSyncConsumer"
	ClientCert := "/Users/mateuszkrasucki/go/src/github.com/iryonetwork/wwm/bin/tls/storageSyncConsumer.pem"
	ClientKey := "/Users/mateuszkrasucki/go/src/github.com/iryonetwork/wwm/bin/tls/storageSyncConsumer-key.pem"
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
