package main

//go:generate sh -c "mkdir -p ../../gen/storage && swagger generate server -A storage -t ../../gen/storage -f ../../docs/api/storage.yml --exclude-main --principal string"

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/restapi"
	"github.com/iryonetwork/wwm/gen/storage/restapi/operations"
	logMW "github.com/iryonetwork/wwm/log"
	APIMetrics "github.com/iryonetwork/wwm/metrics/api"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	"github.com/iryonetwork/wwm/service/authorizer"
	storage "github.com/iryonetwork/wwm/service/storage"
	statusServer "github.com/iryonetwork/wwm/status/server"
	"github.com/iryonetwork/wwm/storage/s3"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/publisher"
	"github.com/iryonetwork/wwm/utils"
	"github.com/iryonetwork/wwm/utils/keyProvider"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "localStorage").
		Logger()

	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// get config
	cfg, err := GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get config")
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load swagger spec")
		return
	}

	// initialize keyProvider
	key, err := base64.StdEncoding.DecodeString(cfg.StorageEncryptionKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to decode storage encryption key")
	}
	keys := keyProvider.New(string(key))

	// initialize storage
	s3cfg := &s3.Config{
		Endpoint:     cfg.S3Endpoint,
		AccessKey:    cfg.S3AccessKey,
		AccessSecret: cfg.S3Secret,
		Secure:       true,
		Region:       cfg.S3Region,
	}
	s3, err := s3.New(s3cfg, keys, logger)
	if err != nil {
		log.Fatalln(err)
	}

	// initialize storageSync publisher
	// create nats/nats-streaming connection
	URLs := fmt.Sprintf("tls://%s:%s@%s", cfg.NatsUsername, cfg.NatsSecret, cfg.NatsAddr)
	ClusterID := cfg.NatsClusterID
	ClientID := cfg.NatsClientID
	ClientCert := cfg.CertPath
	ClientKey := cfg.KeyPath
	var nc *nats.Conn
	var sc publisher.StanConnection
	var p storageSync.Publisher

	// Connect to NATS
	// retry connectng to nats if unsuccesful
	err = utils.Retry(cfg.NatsConnRetries, cfg.NatsConnWait, cfg.NatsConnWaitFactor, logger.With().Str("connect", "nats").Logger(), func() error {
		var err error
		nc, err = nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
		return err
	})

	// Connect to NATS-Streaming if NATS connection succesful
	if err == nil {
		// retry connecting to nats-straming if unsuccesful
		err = utils.Retry(cfg.NatsConnRetries, cfg.NatsConnWait, cfg.NatsConnWaitFactor, logger.With().Str("connect", "nats-streaming").Logger(), func() error {
			var err error
			sc, err = stan.Connect(ClusterID, ClientID, stan.NatsConn(nc))
			return err
		})
	}

	// Initialize publisher
	if err != nil {
		// if connection to nats-streaming was unsuccesful use null publisher
		p = publisher.NewNullPublisher(context.Background())
		logger.Error().Msg("storage service will be started with null storage sync publisher due to failed nats-streaming connection attempts")
	} else {
		// if connection to nats-streaming was succesful use nats-streaming publisher
		cfg := publisher.Cfg{
			Connection:      sc,
			Retries:         5,
			StartRetryWait:  time.Duration(10 * time.Second),
			RetryWaitFactor: 2.0,
		}
		p = publisher.New(context.Background(), cfg, logger)
		// Register metrics
		m := p.GetPrometheusMetricsCollection()
		for _, metric := range m {
			prometheus.MustRegister(metric)
			defer prometheus.Unregister(metric)
		}
	}
	defer p.Close()

	// initialize the servicex
	service := storage.New(s3, keys, p, logger)

	// initialize authorizer
	auth := authorizer.New(cfg.DomainType, cfg.DomainID, fmt.Sprintf("https://%s/%s/validate", cfg.AuthHost, cfg.AuthPath), logger.With().Str("component", "service/authorizer").Logger())

	api := operations.NewStorageAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	server := restapi.NewServer(api)
	server.TLSHost = cfg.ServerHost
	server.TLSPort = cfg.ServerPort
	server.EnabledListeners = []string{"https"}
	server.TLSCertificateKey = flags.Filename(cfg.KeyPath)
	server.TLSCertificate = flags.Filename(cfg.CertPath)

	storageHandlers := storage.NewHandlers(service, logger)

	serverLogger := logger.WithLevel(zerolog.InfoLevel).Str("component", "server")
	api.Logger = serverLogger.Msgf
	api.TokenAuth = auth.GetPrincipalFromToken
	api.APIAuthorizer = auth.Authorizer()
	api.FileListHandler = storageHandlers.FileList()
	api.FileGetHandler = storageHandlers.FileGet()
	api.FileGetVersionHandler = storageHandlers.FileGetVersion()
	api.FileListVersionsHandler = storageHandlers.FileListVersions()
	api.FileNewHandler = storageHandlers.FileNew()
	api.FileUpdateHandler = storageHandlers.FileUpdate()
	api.FileDeleteHandler = storageHandlers.FileDelete()
	api.SyncBucketListHandler = storageHandlers.SyncBucketList()
	api.SyncFileListHandler = storageHandlers.SyncFileList()

	api.RegisterConsumer("*/*", &WildcardConsumer{})

	// initialize metrics middleware
	m := APIMetrics.NewMetrics("api", "").
		WithURLSanitize(utils.WhitelistURLSanitize([]string{"storage", "versions", "sync"}))

	// set API handler with middlewares
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(api.Serve(nil))
	handler = logMW.APILogMiddleware(handler, logger)
	handler = m.Middleware(handler)

	server.SetHandler(handler)

	// Start servers
	// create exit channel that is used to wait for all servers goroutines to exit orederly and carry the errors
	exitCh := make(chan error, 3)

	// start serving metrics
	go func() {
		exitCh <- metricsServer.ServePrometheusMetrics(ctx, fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.MetricsPort), cfg.MetricsNamespace, logger)
	}()
	// start serving status
	go func() {
		ss := statusServer.New(logger)
		exitCh <- ss.ListenAndServeHTTPs(ctx, fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.StatusPort), cfg.StatusNamespace, cfg.CertPath, cfg.KeyPath)
	}()
	// start serving API
	go func() {
		defer server.Shutdown()

		errCh := make(chan error)
		go func() {
			errCh <- server.Serve()
		}()

		for {
			select {
			case err := <-errCh:
				exitCh <- err
				return
			case <-ctx.Done():
				exitCh <- fmt.Errorf("API server exiting because of cancelled context")
				// do nothing, shutdown is deferred
				return
			}
		}
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
	for i := 0; i < 3; i++ {
		err := <-exitCh
		if err != nil {
			logger.Debug().Err(err).Msg("gouroutine exit message")
		}
	}
}

type WildcardConsumer struct{}

func (w *WildcardConsumer) Consume(r io.Reader, in interface{}) error {
	b, err := ioutil.ReadAll(r)
	fmt.Println("WildcardConsumer::Consume", err, in, b)
	return nil
}
