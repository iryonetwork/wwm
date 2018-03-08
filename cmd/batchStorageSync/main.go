package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/storage/keyvalue"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/batch"
)

var (
	boltFilepath  string = "/data/batchStorageSync.db"
	storageBucket string = "batchStorageSync"
	storageKey    string = "lastSuccessfulRun"
)

func main() {
	// Create context with cancel func
	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "batchStorageSync").
		Logger()

	// initialize promethues metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// initialize last successful run with 0 value
	lastSuccessfulRun := time.Unix(0, 0)

	// get metrics collection for key value storage and register in registry
	m := keyvalue.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
	}

	// initialize bolt key value storage to read last succesful run
	storage, err := keyvalue.NewBolt(ctx, boltFilepath, logger.With().Str("component", "storage/keyvalue").Logger(), m)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initiazlie key value storage")
	}
	// read last succesful run
	storedTimestamp := storage.Get(storageBucket, storageKey)
	if storedTimestamp != nil {
		timestamp, err := strfmt.ParseDateTime(string(storedTimestamp))
		if err == nil {
			lastSuccessfulRun = time.Time(timestamp)
		}
	}

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

	// initialize request authenticator
	auth, err := storageSync.NewRequestAuthenticator("/certs/public.crt", "/certs/private.key", logger.With().Str("component", "sync/storage/auth").Logger())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initiazlie storage API request authenticator")
	}

	// initialize handlers
	handlers := storageSync.NewHandlers(localClient.Operations, auth, cloudClient.Operations, auth, logger.With().Str("component", "sync/storage/handlers").Logger())

	// get metrics collection for key value storage and register in registry

	m = batch.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
	}
	// initialize batchStorageSync
	s := batch.New(handlers, logger, m)

	// initialize prometheus metrics pusher
	metricsPusher := push.New("http://localPrometheusPushGateway:9091", "batchStorageSync").Gatherer(metricsRegistry)

	// get current time to be saved as lastSuccesfulRun
	// do it before sync to account for anything that might have happened during sync duration
	startTime := strfmt.DateTime(time.Now())

	// Run sync
	errCh := make(chan error)
	go func() {
		err := s.Sync(ctx, time.Time(lastSuccessfulRun))
		errCh <- err
	}()

	// Run cleanup when sigint or sigterm is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil {
			logger.Error().Err(err).Msg("batch sync failed")
		} else {
			logger.Info().Msg("batch sync successfull")
			// save lastSuccesfulRun
			storage.Update(storageBucket, storageKey, []byte(startTime.String()))
		}
	case <-signalChan:
		logger.Info().Msg("stopping batch sync due to interrupt")
	}

	// push metrics to the push gateway
	err = metricsPusher.Add()
	if err != nil {
		logger.Error().Err(err).Msg("failed to push metrics to push gateway")
	}
}
