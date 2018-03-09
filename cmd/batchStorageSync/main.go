package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/storage/keyvalue"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/batch"
	"github.com/iryonetwork/wwm/utils"
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

	// initialize bolt key value storage to read last succesful run
	storage, err := keyvalue.NewBolt(ctx, boltFilepath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize key value storage")
	}
	// get metrics collection for key value storage and register in registry
	m := storage.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
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

	// initialize batchStorageSync
	s := batch.New(handlers, logger)
	// get prometheus metrics collection for batch sync and register in registry
	m = s.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
	}

	// initialize prometheus metrics pusher
	metricsPusher := push.New("http://localPrometheusPushGateway:9091", "batchStorageSync").Gatherer(metricsRegistry)

	// get current time to be saved as lastSuccesfulRun
	// do it before sync to account for anything that might have happened during sync duration
	startTime := strfmt.DateTime(time.Now())

	// waitGroup for all main go routines
	var wg sync.WaitGroup

	// Run sync
	errCh := make(chan error)
	go func() {
		wg.Add(1)
		defer wg.Done()

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

	wg.Wait()

	// push metrics to the push gateway
	err = metricsPusher.Add()
	if err != nil {
		logger.Error().Err(err).Msg("failed to push metrics to push gateway")
	}
}
