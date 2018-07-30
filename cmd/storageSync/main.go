// storageSync is a worker receiving messages from localStorage published to NATS to resiliently sync everything to cloudStorage
package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	"github.com/iryonetwork/wwm/service/serviceAuthenticator"
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

	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// get config
	cfg, err := GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get config")
	}

	// initialize local storage API client
	local := runtimeClient.New(cfg.StorageHost, cfg.StoragePath, []string{"https"})
	local.Consumers = utils.ConsumersForSync()
	localClient := client.New(local, strfmt.Default)

	// initialize cloud storage API client
	cloud := runtimeClient.New(cfg.CloudStorageHost, cfg.CloudStoragePath, []string{"https"})
	cloud.Consumers = utils.ConsumersForSync()
	cloudClient := client.New(cloud, strfmt.Default)

	// initialize request authenticator
	auth, err := serviceAuthenticator.New(cfg.CertPath, cfg.KeyPath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize storage API request authenticator")
	}

	// initialize handlers
	handlers := storageSync.NewHandlers(localClient.Operations, auth, cloudClient.Operations, auth, logger)

	// create nats/nats-streaming connection
	URLs := fmt.Sprintf("tls://%s:%s@%s", cfg.NatsUsername, cfg.NatsSecret, cfg.NatsAddr)
	ClusterID := cfg.NatsClusterID
	ClientID := cfg.NatsClientID
	ClientCert := cfg.CertPath
	ClientKey := cfg.KeyPath
	var nc *nats.Conn
	var sc stan.Conn

	// retry connection to nats if unsuccesful
	err = utils.Retry(cfg.NatsConnRetries, cfg.NatsConnWait, cfg.NatsConnWaitFactor, logger.With().Str("connection", "nats").Logger(), func() error {
		var err error
		nc, err = nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
		return err
	})
	if err != nil {
		logger.Fatal().Msg("failed to connect to nats")
	}

	err = utils.Retry(cfg.NatsConnRetries, cfg.NatsConnWait, cfg.NatsConnWaitFactor, logger.With().Str("connection", "nats").Logger(), func() error {
		var err error
		sc, err = stan.Connect(ClusterID, ClientID, stan.NatsConn(nc))
		return err
	})
	if err != nil {
		logger.Fatal().Msg("failed to connect to nats-streaming")
	}

	// initalize consumer
	consumerCfg := consumer.Cfg{
		Connection:  sc,
		AckWait:     cfg.AckWait,
		MaxInflight: cfg.MaxInflight,
		Handlers:    handlers,
	}
	c := consumer.New(ctx, consumerCfg, logger)
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
	// create exit channel that is used to wait for all servers goroutines to exit orderly and carry the errors
	exitCh := make(chan error, 2)

	// start serving metrics
	go func() {
		exitCh <- metricsServer.ServePrometheusMetrics(ctx, fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.MetricsPort), cfg.MetricsNamespace, logger)
	}()
	// start serving status
	go func() {
		ss := statusServer.New(logger)
		exitCh <- ss.ListenAndServeHTTPs(ctx, fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.StatusPort), cfg.StatusNamespace, cfg.CertPath, cfg.KeyPath)
	}()

	// run cleanup when sigint or sigterm is received or error on starting server happened
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer cancelContext()

		for {
			select {
			case err := <-exitCh:
				logger.Info().Msg("exiting application because of exiting server goroutine")
				// pass error back to channel satisfy exit condition
				exitCh <- err
				return
			case <-signalChan:
				logger.Info().Msg("received interrupt")
				return
			}
		}
	}()

	<-ctx.Done()
	for i := 0; i < 2; i++ {
		err := <-exitCh
		if err != nil {
			logger.Debug().Err(err).Msg("gouroutine exit message")
		}
	}
}
