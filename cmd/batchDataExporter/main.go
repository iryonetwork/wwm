package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/reports/filesDataExporter"
	"github.com/iryonetwork/wwm/reports/filesDataExporter/batch"
	"github.com/iryonetwork/wwm/service/serviceAuthenticator"
	"github.com/iryonetwork/wwm/storage/keyvalue"
	reportsStorage "github.com/iryonetwork/wwm/storage/reports"
	"github.com/iryonetwork/wwm/utils"
)

const (
	storageBucket string = "batchDataExporter"
	storageKey    string = "lastSuccessfulRun"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "batchDataExporter").
		Logger()

	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// get config
	cfg, err := getConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get config")
	}

	// initialize promethues metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// initialize last successful run with 0 value
	lastSuccessfulRun := time.Unix(0, 0)

	// initialize bolt key value storage to read last succesful run
	s, err := keyvalue.NewBolt(ctx, cfg.BoltDBFilepath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize key value storage")
	}
	// get metrics collection for key value storage and register in registry
	m := s.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
	}

	// read last succesful run
	storedTimestamp := s.Get(storageBucket, storageKey)
	if storedTimestamp != nil {
		timestamp, err := strfmt.ParseDateTime(string(storedTimestamp))
		if err == nil {
			lastSuccessfulRun = time.Time(timestamp)
		}
	}

	// initialize source storage API client
	source := runtimeClient.New(cfg.StorageHost, cfg.StoragePath, []string{"https"})
	source.Consumers = utils.ConsumersForSync()
	sourceClient := client.New(source, strfmt.Default)

	// connect to database
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require",
		cfg.DbUsername,
		cfg.DbPassword,
		cfg.PGHost,
		cfg.PGDatabase)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database connection")
	}
	db.LogMode(cfg.DbDetailedLog)

	// switch roles
	tx := db.Exec(fmt.Sprintf("SET ROLE '%s'", cfg.PGRole))
	if err := tx.Error; err != nil {
		logger.Fatal().Err(err).Msg("Failed to switch database roles")
	}

	// initialize storage
	storage, err := reportsStorage.New(ctx, db, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize reports storage")
	}

	// initialize data sanitizer
	key, err := base64.StdEncoding.DecodeString(cfg.DataEncryptionKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to decode data encryption key")
	}
	sanitizer, err := filesDataExporter.NewSanitizer(cfg.FieldsToSanitize.Slice, key, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize files data sanitizer")
	}

	// initialize request authenticator
	auth, err := serviceAuthenticator.New(cfg.CertPath, cfg.KeyPath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize storage API request authenticator")
	}

	// initialize handlers
	handlers := filesDataExporter.NewHandlers(sourceClient.Operations, auth, sanitizer, storage, logger)

	// initialize batch file data exporter
	e := batch.New(handlers, cfg.BucketsRateLimit, logger)

	// get prometheus metrics collection for batch sync and register in registry
	m = e.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
	}

	// initialize prometheus metrics pusher
	metricsPusher := push.New(cfg.PrometheusPushGatewayAddress, "batchDataExporter").Gatherer(metricsRegistry)

	// get current time to be saved as lastSuccesfulRun
	// do it before sync to account for anything that might have happened during sync duration
	startTime := strfmt.DateTime(time.Now())

	// Run export
	exitCh := make(chan error)
	go func() {
		exitCh <- e.Export(ctx, time.Time(lastSuccessfulRun))
	}()

	// Run cleanup when sigint or sigterm is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

Loop:
	for {
		select {
		case err := <-exitCh:
			if err != nil {
				logger.Error().Err(err).Msg("batch files data failed")
			} else {
				logger.Info().Msg("batch files data export successful")
				// save lastSuccesfulRun
				s.Update(storageBucket, storageKey, []byte(startTime.String()))
			}
			break Loop
		case <-signalChan:
			logger.Info().Msg("stopping batch files data export due to interrupt")
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
