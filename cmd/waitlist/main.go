package main

//go:generate sh -c "mkdir -p ../../gen/waitlist/ && swagger generate server -A waitlist -t ../../gen/waitlist/ -f ../../docs/api/waitlist.yml --exclude-main --principal string"

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/restapi"
	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations"
	logMW "github.com/iryonetwork/wwm/log"
	APIMetrics "github.com/iryonetwork/wwm/metrics/api"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	"github.com/iryonetwork/wwm/service/authorizer"
	statusServer "github.com/iryonetwork/wwm/status/server"
	"github.com/iryonetwork/wwm/storage/waitlist"
	"github.com/iryonetwork/wwm/utils"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "waitlist").
		Logger()

	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load swagger spec")
		return
	}

	// get config
	cfg, err := getConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get config")
	}

	// initialize the service
	key, err := base64.StdEncoding.DecodeString(cfg.StorageEncryptionKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to decode storage encryption key")
	}
	storage, err := waitlist.New(cfg.BoltDBFilepath, key, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize waitlist storage")
	}
	_, err = storage.EnsureDefaultList(cfg.DefaultListID, cfg.DefaultListName)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ensure default list")
	}

	// run migrations
	err = storage.MigrateVitalSigns()
	if err != nil {
		logger.Error().Err(err).Msg("Vital sings migration failed")
	}

	auth := authorizer.New(cfg.DomainType, cfg.DomainID, fmt.Sprintf("https://%s/%s/validate", cfg.AuthHost, cfg.AuthPath), logger)

	api := operations.NewWaitlistAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	api.TokenAuth = auth.GetPrincipalFromToken
	api.APIAuthorizer = auth.Authorizer()

	server := restapi.NewServer(api)
	server.Host = cfg.ServerHost
	server.Port = cfg.ServerPortHTTP
	server.TLSHost = cfg.ServerHost
	server.TLSPort = cfg.ServerPortHTTPS
	server.TLSCertificate = flags.Filename(cfg.CertPath)
	server.TLSCertificateKey = flags.Filename(cfg.KeyPath)
	server.EnabledListeners = []string{"http", "https"}
	defer server.Shutdown()

	h := &handlers{s: storage}

	api.WaitlistDeleteListIDHandler = h.WaitlistDeleteListID()
	api.WaitlistGetHandler = h.WaitlistGet()
	api.WaitlistPostHandler = h.WaitlistPost()
	api.WaitlistPutListIDHandler = h.WaitlistPutListID()

	api.ItemDeleteListIDItemIDHandler = h.ItemDeleteListIDItemID()
	api.ItemGetListIDHandler = h.ItemGetListID()
	api.ItemGetListIDHistoryHandler = h.ItemGetListIDHistory()
	api.ItemPostListIDHandler = h.ItemPostListID()
	api.ItemPutListIDItemIDHandler = h.ItemPutListIDItemID()
	api.ItemPutListIDItemIDTopHandler = h.ItemPutListIDItemIDTop()
	api.ItemPutListIDItemIDReopenHandler = h.ItemPutListIDItemIDReopen()

	// initialize metrics middleware
	apiMetrics := APIMetrics.NewMetrics("api", "").WithURLSanitize(utils.WhitelistURLSanitize([]string{}))

	// set API handler with middlewares
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(api.Serve(nil))
	handler = logMW.APILogMiddleware(handler, logger)
	handler = apiMetrics.Middleware(handler)
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
