package main

// goswagger code for spec discovery.yml is generated in ../localDiscovery

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/discovery/restapi"
	"github.com/iryonetwork/wwm/gen/discovery/restapi/operations"
	APIMetrics "github.com/iryonetwork/wwm/metrics/api"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	"github.com/iryonetwork/wwm/service/authorizer"
	discoveryService "github.com/iryonetwork/wwm/service/discovery"
	statusServer "github.com/iryonetwork/wwm/status/server"
	discoveryStorage "github.com/iryonetwork/wwm/storage/discovery"
	"github.com/iryonetwork/wwm/utils"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "cloudDiscovery").
		Logger()

	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// get config
	cfg, err := getConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get config")
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load swagger spec")
		return
	}

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
	db.LogMode(true)

	// switch roles
	tx := db.Exec(fmt.Sprintf("SET ROLE '%s'", cfg.PGRole))
	if err := tx.Error; err != nil {
		logger.Fatal().Err(err).Msg("Failed to switch database roles")
	}

	// initialize storage
	storage, err := discoveryStorage.New(ctx, db, "", logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize discovery storage")
	}

	// initialize the service
	service := discoveryService.New(ctx, storage, nil, logger)

	discoveryHandlers := discoveryService.NewHandlers(service, logger)

	auth := authorizer.New(cfg.DomainType, cfg.DomainID, fmt.Sprintf("https://%s/%s/validate", cfg.AuthHost, cfg.AuthPath), logger)

	api := operations.NewDiscoveryAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	api.TokenAuth = auth.GetPrincipalFromToken
	api.APIAuthorizer = auth.Authorizer()
	api.Logger = logger.Info().Str("component", "api").Msgf

	api.QueryHandler = discoveryHandlers.Query()
	api.CreateHandler = discoveryHandlers.Create()
	api.UpdateHandler = discoveryHandlers.Update()
	api.DeleteHandler = discoveryHandlers.Delete()
	api.FetchHandler = discoveryHandlers.Fetch()
	api.LinkHandler = discoveryHandlers.Link()
	api.UnlinkHandler = discoveryHandlers.Unlink()
	api.CodesGetHandler = discoveryHandlers.CodesGet()
	api.CodeGetHandler = discoveryHandlers.CodeGet()

	server := restapi.NewServer(api)
	server.TLSHost = cfg.ServerHost
	server.TLSPort = cfg.ServerPort
	server.TLSCertificate = flags.Filename(cfg.CertPath)
	server.TLSCertificateKey = flags.Filename(cfg.KeyPath)
	server.EnabledListeners = []string{"https"}

	// initialize metrics middleware
	m := APIMetrics.NewMetrics("api", "").
		WithURLSanitize(utils.WhitelistURLSanitize([]string{"storage", "versions", "sync"}))

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(api.Serve(nil))
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
	for i := 0; i < 2; i++ {
		err := <-exitCh
		if err != nil {
			logger.Debug().Err(err).Msg("goroutine exit message")
		}
	}
}
