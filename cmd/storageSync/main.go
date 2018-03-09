// storageSync is a worker receiving messages from localStorage published to NATS to resiliently sync everything to cloudStorage
package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	statusServer "github.com/iryonetwork/wwm/status/server"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/consumer"
	"github.com/iryonetwork/wwm/utils"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "storageSync").
		Logger()

	// initialize local storage API client
	local := runtimeClient.New("localStorage", "storage", []string{"https"})
	local.Consumers = utils.ConsumersForSync()
	localClient := client.New(local, strfmt.Default)

	// initialize cloud storage API client
	cloud := runtimeClient.New("cloudStorage", "storage", []string{"https"})
	cloud.Consumers = utils.ConsumersForSync()
	cloudClient := client.New(cloud, strfmt.Default)

	// initialize request authenticator
	auth, err := storageSync.NewRequestAuthenticator("/certs/public.crt", "/certs/private.key", logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize storage API request authenticator")
	}

	// initialize handlers
	handlers := storageSync.NewHandlers(localClient.Operations, auth, cloudClient.Operations, auth, logger)

	// create nats/nats-streaming connection
	URLs := "tls://nats:secret@localNats:4242"
	ClusterID := "localNats"
	ClientID := "storageSync"
	ClientCert := "/certs/public.crt"
	ClientKey := "/certs/private.key"
	var nc *nats.Conn
	var sc stan.Conn

	// retry connection to nats if unsuccesful
	err = utils.Retry(10, time.Duration(500*time.Millisecond), 2.0, logger.With().Str("connection", "nats").Logger(), func() error {
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
	ctx, cancelContext := context.WithCancel(context.Background())

	// initalize consumer
	cfg := consumer.Cfg{
		Connection: sc,
		AckWait:    time.Duration(time.Second),
		Handlers:   handlers,
	}
	c := consumer.New(ctx, cfg, logger)
	// Register metrics
	m := c.GetPrometheusMetricsCollection()
	for _, metric := range m {
		prometheus.MustRegister(metric)
		defer prometheus.Unregister(metric)
	}

	// Start subscriptions
	c.StartSubscription(storageSync.FileNew)
	c.StartSubscription(storageSync.FileUpdate)
	c.StartSubscription(storageSync.FileDelete)

	// Start servers
	errCh := make(chan error)
	// waitGroup for all main go routines
	var wg sync.WaitGroup

	// start serving metrics
	go func() {
		wg.Add(1)
		defer wg.Done()

		errCh <- metricsServer.ServePrometheusMetrics(ctx, ":9090", "", logger)
	}()

	// start serving status
	go func() {
		wg.Add(1)
		defer wg.Done()

		ss := statusServer.New(logger)
		errCh <- ss.ListenAndServeHTTPs(ctx, "storageSync:4433", "", "/certs/public.crt", "/certs/private.key")
	}()

	// run cleanup when sigint or sigterm is received or error on starting server happened
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		wg.Add(1)
		defer cancelContext()
		defer wg.Done()

		select {
		case err := <-errCh:
			logger.Error().Err(err).Msg("failed to start server")
		case <-signalChan:
			logger.Error().Msg("received interrupt")
		}
	}()

	wg.Wait()
}
