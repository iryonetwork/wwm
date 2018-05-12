package main

//go:generate sh -c "mkdir -p ../../gen/auth/ && swagger generate server -A cloudAuth -t ../../gen/auth/ -f ../../docs/api/auth.yml --exclude-main --principal string"

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	loads "github.com/go-openapi/loads"
	"github.com/jasonlvhit/gocron"
	flags "github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/auth/restapi"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations"
	logMW "github.com/iryonetwork/wwm/log"
	APIMetrics "github.com/iryonetwork/wwm/metrics/api"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	"github.com/iryonetwork/wwm/service/authDataManager"
	"github.com/iryonetwork/wwm/service/authSync"
	"github.com/iryonetwork/wwm/service/authenticator"
	statusServer "github.com/iryonetwork/wwm/status/server"
	"github.com/iryonetwork/wwm/storage/auth"
	"github.com/iryonetwork/wwm/utils"
)

func main() {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "localAuth").
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
	}

	// initialize storage
	key, err := base64.StdEncoding.DecodeString(cfg.StorageEncryptionKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to decode storage encryption key")
	}

	dbPath := cfg.BoltDBFilepath
	// if there is no database file download it from cloud
	_, err = os.Stat(dbPath)
	if _, err := os.Stat(dbPath); err != nil {
		storage, err := auth.New(dbPath, key, false, false, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to initialize auth storage")
		}

		authSync, err := authSync.New(storage, cfg.AuthSyncCertPath, cfg.AuthSyncKeyPath, fmt.Sprintf("https://%s/%s/database", cfg.CloudAuthHost, cfg.CloudAuthPath), logger.With().Str("component", "service/authSync-initialCloudDownload").Logger())
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to initialize authSync service")
		}

		for try := 1; try <= 5; try++ {
			err := authSync.Sync()
			if err == nil {
				break
			}
			if try == 5 {
				os.Remove(dbPath)
				logger.Fatal().Err(err).Msg("Failed to sync database from cloud after 5 tries")
			}
			time.Sleep(time.Duration(try*3) * time.Second)
		}
		storage.Close()
	}

	storage, err := auth.New(dbPath, key, true, true, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize auth storage")
	}

	// register metrics collected by storage
	m := storage.GetPrometheusMetricsCollection()
	for _, metric := range m {
		prometheus.MustRegister(metric)
		defer prometheus.Unregister(metric)
	}

	// initialize the services
	auth, err := authenticator.New(cfg.DomainType, cfg.DomainID, storage, cfg.ServiceCertsAndPaths.Map, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize authenticator service")
	}
	authSync, err := authSync.New(storage, cfg.AuthSyncCertPath, cfg.AuthSyncKeyPath, fmt.Sprintf("https://%s/%s/database", cfg.CloudAuthHost, cfg.CloudAuthPath), logger.With().Str("component", "service/authSync").Logger())
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize authSync service")
	}
	authData := authDataManager.New(storage, logger.With().Str("component", "service/authDataManager").Logger())

	// setup API
	api := operations.NewCloudAuthAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	server := restapi.NewServer(api)
	server.Host = cfg.ServerHost
	server.Port = cfg.ServerPortHTTP
	server.TLSHost = cfg.ServerHost
	server.TLSPort = cfg.ServerPortHTTPS
	server.TLSCertificate = flags.Filename(cfg.CertPath)
	server.TLSCertificateKey = flags.Filename(cfg.KeyPath)
	server.EnabledListeners = []string{"http", "https"}
	defer server.Shutdown()

	authHandlers := authenticator.NewHandlers(auth)
	authDataHandlers := authDataManager.NewHandlers(authData)

	serverLogger := logger.WithLevel(zerolog.InfoLevel).Str("component", "server")
	api.Logger = serverLogger.Msgf
	api.TokenAuth = auth.GetPrincipalFromToken
	api.APIAuthorizer = auth.Authorizer()
	api.GetRenewHandler = authHandlers.GetRenew()
	api.PostLoginHandler = authHandlers.PostLogin()
	api.PostValidateHandler = authHandlers.PostValidate()

	api.GetUsersHandler = authDataHandlers.GetUsers()
	api.GetUsersIDHandler = authDataHandlers.GetUsersID()
	api.GetUsersIDRolesHandler = authDataHandlers.GetUsersIDRoles()
	api.GetUsersIDOrganizationsHandler = authDataHandlers.GetUsersIDOrganizations()
	api.GetUsersIDClinicsHandler = authDataHandlers.GetUsersIDClinics()
	api.GetUsersIDLocationsHandler = authDataHandlers.GetUsersIDLocations()
	api.GetUsersMeHandler = authDataHandlers.GetUsersMe()
	api.GetUsersMeRolesHandler = authDataHandlers.GetUsersMeRoles()
	api.GetUsersMeOrganizationsHandler = authDataHandlers.GetUsersMeOrganizations()
	api.GetUsersMeClinicsHandler = authDataHandlers.GetUsersMeClinics()
	api.GetUsersMeLocationsHandler = authDataHandlers.GetUsersMeLocations()

	api.GetRolesHandler = authDataHandlers.GetRoles()
	api.GetRolesIDHandler = authDataHandlers.GetRolesID()
	api.GetRolesIDUsersHandler = authDataHandlers.GetRolesIDUsers()

	api.GetRulesHandler = authDataHandlers.GetRules()
	api.GetRulesIDHandler = authDataHandlers.GetRulesID()

	api.GetClinicsHandler = authDataHandlers.GetClinics()
	api.GetClinicsIDHandler = authDataHandlers.GetClinicsID()
	api.GetClinicsIDUsersHandler = authDataHandlers.GetClinicsIDUsers()

	api.GetOrganizationsHandler = authDataHandlers.GetOrganizations()
	api.GetOrganizationsIDHandler = authDataHandlers.GetOrganizationsID()
	api.GetOrganizationsIDLocationsHandler = authDataHandlers.GetOrganizationsIDLocations()
	api.GetOrganizationsIDUsersHandler = authDataHandlers.GetOrganizationsIDUsers()

	api.GetLocationsHandler = authDataHandlers.GetLocations()
	api.GetLocationsIDHandler = authDataHandlers.GetLocationsID()
	api.GetLocationsIDOrganizationsHandler = authDataHandlers.GetLocationsIDOrganizations()
	api.GetLocationsIDUsersHandler = authDataHandlers.GetLocationsIDUsers()

	api.GetUserRolesHandler = authDataHandlers.GetUserRoles()
	api.GetUserRolesIDHandler = authDataHandlers.GetUserRolesID()

	// initialize metrics middleware
	apiMetrics := APIMetrics.NewMetrics("api", "").
		WithURLSanitize(utils.WhitelistURLSanitize([]string{
			"login",
			"validate",
			"renew",
			"users",
			"roles",
			"clinics",
			"locations",
			"organizations",
			"userRoles",
			"rules",
			"database",
		}))

	// set handler with middlewares
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(api.Serve(nil))
	handler = logMW.APILogMiddleware(handler, logger)
	handler = apiMetrics.Middleware(handler)
	server.SetHandler(handler)

	gocron.Every(5).Minutes().Do(authSync.Sync)
	go gocron.Start()

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
			logger.Debug().Err(err).Msg("goroutine exit message")
		}
	}
}
