package main

//go:generate sh -c "mkdir -p ../../gen/storage/ && swagger generate client -A storage -t ../../gen/storage/ -f ../../docs/api/storage.yml --principal string"

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-openapi/runtime"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	"github.com/iryonetwork/wwm/gen/storage/client/operations"
	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/reports/generator"
	"github.com/iryonetwork/wwm/reports/generator/xlsxWriter"
	"github.com/iryonetwork/wwm/service/serviceAuthenticator"
	"github.com/iryonetwork/wwm/storage/keyvalue"
	reportsStorage "github.com/iryonetwork/wwm/storage/reports"
)

const (
	fileUUIDStorageBucket string = "batchReportGenerator"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "batchReportGenerator").
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

	// initialize reports storage API client
	rc := runtimeClient.New(cfg.StorageHost, cfg.StoragePath, []string{"https"})
	filesStorageClient := client.New(rc, strfmt.Default)

	// initialize request authenticator
	auth, err := serviceAuthenticator.New(cfg.CertPath, cfg.KeyPath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize reports storage API request authenticator")
	}

	// initalize generator
	g, err := generator.New(storage, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize generator")
	}

	// initialize prometheus metrics pusher
	metricsPusher := push.New(cfg.PrometheusPushGatewayAddress, "batchReportGenerator").Gatherer(metricsRegistry)

	// initialize bolt key value storage to read filenames
	s, err := keyvalue.NewBolt(ctx, cfg.BoltDBFilepath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize key value storage")
	}
	// get metrics collection for key value storage and register in registry
	m := s.GetPrometheusMetricsCollection()
	for _, metric := range m {
		metricsRegistry.MustRegister(metric)
	}

	for _, spec := range cfg.ReportSpecs.Slice {
		// initialize xlsx writer
		writer, err := xlsxWriter.New(spec, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to initialize xlsx file writer")
		}
		defer os.Remove(writer.Filename())

		// run generator
		generated, err := g.Generate(ctx, writer, spec, nil, nil)
		if err != nil {
			logger.Fatal().Err(err).Msgf("failed to generate report file %s", spec.Type)
		}

		if generated {
			// flush writer
			writer.Flush()
			if writer.Error() != nil {
				logger.Fatal().Err(err).Msgf("failed to flush xlsx report writer")
			}

			// read file and turn content into buffer
			b, err := ioutil.ReadFile(writer.Filename())
			if err != nil {
				logger.Fatal().Err(err).Msgf("failed to read temp file")
			}
			buf := bytes.NewBuffer(b)

			// get fileUUID for this type of report
			var fileUUID string
			fileUUIDBytes := s.Get(fileUUIDStorageBucket, spec.Type)
			if fileUUIDBytes != nil {
				fileUUID = string(fileUUIDBytes)
				// if file UUID found, update file
				fileUpdateParams := operations.NewFileUpdateParams().
					WithBucket(strfmt.UUID(cfg.ReportsBucketUUID)).
					WithFileID(fileUUID).
					WithContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet").
					WithLabels([]string{spec.Type}).
					WithFile(runtime.NamedReader("reader", buf))

				_, err := filesStorageClient.Operations.FileUpdate(fileUpdateParams, auth)
				if err != nil {
					logger.Fatal().Err(err).Msgf("failed to upload report file %s", spec.Type)
				}
			} else {
				// if file UUID not found, create new file
				// upload file to storage
				fileNewParams := operations.NewFileNewParams().
					WithBucket(strfmt.UUID(cfg.ReportsBucketUUID)).
					WithContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet").
					WithLabels([]string{spec.Type}).
					WithFile(runtime.NamedReader("reader", buf))

				ok, err := filesStorageClient.Operations.FileNew(fileNewParams, auth)
				if err != nil {
					logger.Fatal().Err(err).Msgf("failed to upload report file %s", spec.Type)
				}

				fileUUID = ok.Payload.Name

				// Update file UUID for this report type in key value storage
				errorChecker.LogError(s.Update(fileUUIDStorageBucket, spec.Type, []byte(fileUUID)))
			}

			logger.Info().Msgf("report %s was uploaded as file %s", spec.Type, fileUUID)
		}
	}

	// push metrics to the push gateway
	err = metricsPusher.Add()
	if err != nil {
		logger.Error().Err(err).Msg("failed to push metrics to push gateway")
	}
}
