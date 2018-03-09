package main

//go:generate sh -c "mkdir -p ../../gen/auth/ && swagger generate server -A cloudAuth -t ../../gen/auth/ -f ../../docs/api/auth.yml --exclude-main --principal string"

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	loads "github.com/go-openapi/loads"
	"github.com/jasonlvhit/gocron"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/auth/restapi"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations"
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

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load swagger spec")
	}

	// initialize storage
	// TODO: get key from vault
	key := []byte{0xe9, 0xf8, 0x2d, 0xf9, 0xc4, 0x14, 0xc1, 0x41, 0xdb, 0x87, 0x31, 0x1a, 0x95, 0x79, 0x5, 0xbf, 0x71, 0x12, 0x30, 0xd3, 0x2d, 0x8b, 0x59, 0x9d, 0x27, 0x13, 0xfa, 0x84, 0x55, 0x63, 0x64, 0x64}

	dbPath := "/data/localAuth.db"
	// if there is no database file download it from cloud
	_, err = os.Stat(dbPath)
	if _, err := os.Stat(dbPath); err != nil {
		storage, err := auth.New(dbPath, key, false, false, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to initialize auth storage")
		}

		authSync, err := authSync.New(storage, "/certs/localAuthSync.pem", "/certs/localAuthSync-key.pem", "https://cloudAuth/auth/database", logger.With().Str("component", "service/authSync-initialCloudDownload").Logger())
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

	// initialize the service
	auth, err := authenticator.New(storage, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize authenticator service")
	}

	authSync, err := authSync.New(storage, "/certs/localAuthSync.pem", "/certs/localAuthSync-key.pem", "https://cloudAuth/auth/database", logger.With().Str("component", "service/authSync").Logger())
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize authSync service")
	}
	api := operations.NewCloudAuthAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	server := restapi.NewServer(api)
	server.TLSPort = 443
	server.TLSCertificate = "/certs/localAuth.pem"
	server.TLSCertificateKey = "/certs/localAuth-key.pem"
	server.EnabledListeners = []string{"https"}
	defer server.Shutdown()

	authHandlers := authenticator.NewHandlers(auth)

	serverLogger := logger.WithLevel(zerolog.InfoLevel).Str("component", "server")
	api.Logger = serverLogger.Msgf
	api.TokenAuth = auth.GetPrincipalFromToken
	api.APIAuthorizer = auth.Authorizer()
	api.AuthGetRenewHandler = authHandlers.GetRenew()
	api.AuthPostLoginHandler = authHandlers.PostLogin()
	api.AuthPostValidateHandler = authHandlers.PostValidate()

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(api.Serve(nil))

	server.SetHandler(handler)

	gocron.Every(5).Minutes().Do(authSync.Sync)
	go gocron.Start()

	// Start servers
	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	// waitGroup for all main go routines
	var wg sync.WaitGroup
	// create error channel
	errCh := make(chan error)

	// start serving status
	go func() {
		wg.Add(1)
		defer wg.Done()

		ss := statusServer.New(logger)
		errCh <- ss.ListenAndServeHTTPs(ctx, "localAuth:4433", "", "/certs/localAuth.pem", "/certs/localAuth-key.pem")
	}()
	// start serving API
	go func() {
		wg.Add(1)
		defer wg.Done()
		defer server.Shutdown()

		localErrCh := make(chan error)
		go func() {
			localErrCh <- server.Serve()
		}()

		select {
		case err := <-localErrCh:
			localErrCh <- err
		case <-ctx.Done():
			//do nothing except deferred cleanup
		}
	}()

	// run cleanup when sigint or sigterm is received or error on starting server happened
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		wg.Add(1)
		defer cancelContext()
		defer wg.Done()

		select {
		case err := <-errCh:
			logger.Error().Err(err).Msg("failed to start server")
		case <-signalChan:
			logger.Error().Msg("received interrupt")
		}
	}()

	wg.Wait()
}
