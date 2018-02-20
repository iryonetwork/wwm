// storageSync is a worker receiving messages from localStorage published to NATS to resiliently sync everything to cloudStorage
package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/metrics"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/consumer"
	"github.com/iryonetwork/wwm/utils"
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
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
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
	handlers := consumer.NewHandlers(localClient.Operations, auth, cloudClient.Operations, auth, logger.With().Str("component", "sync/storage/consumer/handlers").Logger())

	// create nats/nats-streaming connection
	URLs := "tls://nats:secret@localNats:4242"
	ClusterID := "localNats"
	ClientID := "storageSync"
	ClientCert := "/certs/public.crt"
	ClientKey := "/certs/private.key"
	var nc *nats.Conn
	var sc stan.Conn

	// retry connection to nats if unsuccesful
	err := utils.Retry(10, time.Duration(500*time.Millisecond), 2.0, logger.With().Str("connection", "nats").Logger(), func() error {
		var err error
		nc, err = nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
		return err
	})
	if err != nil {
		logger.Fatal().Msg("failed to connect to nats")
	}

	err = utils.Retry(10, time.Duration(500*time.Millisecond), 3.0, logger.With().Str("connection", "nats").Logger(), func() error {
		var err error
		sc, err = stan.Connect(ClusterID, ClientID, stan.NatsConn(nc))
		return err
	})
	if err != nil {
		logger.Fatal().Msg("failed to connect to nats-streaming")
	}

	// Create context with cancel func
	ctx, shutdown := context.WithCancel(context.Background())

	// Register metrics
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "consumer",
		Name:      "task_seconds",
		Help:      "Time taken to serve tasks",
	}, []string{"event", "ack"})
	prometheus.MustRegister(h)

	// initalize consumer
	cfg := consumer.Cfg{
		Connection: sc,
		AckWait:    time.Duration(time.Second),
		Handlers:   handlers,
	}
	c := consumer.New(context.Background(), cfg, logger.With().Str("component", "sync/storage/consumer").Logger(), h)
	c.StartSubscription(storageSync.FileNew)
	c.StartSubscription(storageSync.FileUpdate)
	c.StartSubscription(storageSync.FileDelete)

	go func() {
		err := metrics.ServePrometheusMetrics(ctx, ":9090", "")
		if err != nil {
			logger.Error().Err(err).Msg("prometheus metrics server failure")
		}
	}()

	// Run cleanup when sigint or sigterm is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		c.Close()
		shutdown()
		cleanupDone <- true
	}()
	<-cleanupDone
}
