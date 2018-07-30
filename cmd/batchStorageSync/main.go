package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/service/serviceAuthenticator"
	"github.com/iryonetwork/wwm/storage/keyvalue"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/batch"
	"github.com/iryonetwork/wwm/utils"
)

const (
	storageBucket string = "batchStorageSync"
	storageKey    string = "lastSuccessfulRun"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "batchStorageSync").
		Logger()

	// Create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// get config
	cfg, err := GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get config")
	}

	// initialize promethues metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// initialize last successful run with 0 value
	lastSuccessfulRun := time.Unix(0, 0)

	// initialize bolt key value storage to read last succesful run
	storage, err := keyvalue.NewBolt(ctx, cfg.BoltDBFilepath, logger)
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

	// initialize batchStorageSync
	batchCfg := batch.Cfg{
		BucketsRateLimit:        cfg.BucketsRateLimit,
		FilesPerBucketRateLimit: cfg.FilesPerBucketRateLimit,
	}
	s := batch.New(handlers, batchCfg, logger)

	// get prometheus metrics collection for batch sync and register in registry
	m = s.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
	}

	// initialize prometheus metrics pusher
	metricsPusher := push.New(cfg.PrometheusPushGatewayAddress, "batchStorageSync").Gatherer(metricsRegistry)

	// get current time to be saved as lastSuccesfulRun
	// do it before sync to account for anything that might have happened during sync duration
	startTime := strfmt.DateTime(time.Now())

	// Run sync
	exitCh := make(chan error)
	go func() {
		exitCh <- s.Sync(ctx, time.Time(lastSuccessfulRun))
	}()

	// Run cleanup when sigint or sigterm is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

Loop:
	for {
		select {
		case err := <-exitCh:
			if err != nil {
				logger.Error().Err(err).Msg("batch sync failed")
			} else {
				logger.Info().Msg("batch sync successfull")
				// save lastSuccesfulRun
				storage.Update(storageBucket, storageKey, []byte(startTime.String()))
			}
			break Loop
		case <-signalChan:
			logger.Info().Msg("stopping batch sync due to interrupt")
			cancelContext()
			break Loop
		}
	}

	// push metrics to the push gateway
	err = metricsPusher.Add()
	if err != nil {
		logger.Error().Err(err).Msg("failed to push metrics to push gateway")
	}
}
