package main

//go:generate sh -c "mkdir -p ../../gen/storage && swagger generate server -A storage -t ../../gen/storage -f ../../docs/api/storage.yml --exclude-main --principal string"

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	loads "github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
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
	"github.com/iryonetwork/wwm/storage/s3/mock"
	"github.com/iryonetwork/wwm/sync/storage/publisher"
	"github.com/iryonetwork/wwm/utils"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "cloudStorage").
		Logger()

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load swagger spec")
		return
	}

	// initialize keyProvider
	ctrl := gomock.NewController(nil)
	keys := mock.NewMockKeyProvider(ctrl)
	keys.EXPECT().Get(gomock.Any()).AnyTimes().Return("SECRETSECRETSECRETSECRETSECRETSE", nil)

	// initialize storage
	cfg := &s3.Config{
		Endpoint:     "cloudMinio:9000",
		AccessKey:    "cloud",
		AccessSecret: "cloudminio",
		Secure:       true,
		Region:       "us-east-1",
	}
	s3, err := s3.New(cfg, keys, logger)
	if err != nil {
		log.Fatalln(err)
	}

	// initialize the service
	service := storage.New(s3, keys, publisher.NewNullPublisher(context.Background()), logger)

	// initialize authorizer
	auth := authorizer.New("https://cloudAuth/auth/validate", logger.With().Str("component", "service/authorizer").Logger())

	api := operations.NewStorageAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	server := restapi.NewServer(api)
	server.TLSHost = "0.0.0.0"
	server.TLSPort = 443
	server.EnabledListeners = []string{"https"}
	server.TLSCertificateKey = "/certs/cloudStorage-key.pem"
	server.TLSCertificate = "/certs/cloudStorage.pem"

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
	api.SyncFileMetadataHandler = storageHandlers.SyncFileMetadata()
	api.SyncFileHandler = storageHandlers.SyncFile()
	api.SyncFileDeleteHandler = storageHandlers.SyncFileDelete()
	api.SyncBucketListHandler = storageHandlers.SyncBucketList()
	api.SyncFileListHandler = storageHandlers.SyncFileList()
	api.SyncFileListVersionsHandler = storageHandlers.SyncFileListVersions()

	api.RegisterConsumer("*/*", &WildcardConsumer{})

	// initialize metrics middleware
	m := APIMetrics.NewMetrics("api", "").
		WithURLSanitize(utils.WhitelistURLSanitize([]string{"storage", "versions", "sync"}))

	// set handler with middlewares
	apiHandler := logMW.APILogMiddleware(api.Serve(nil), logger)
	apiHandler = m.Middleware(apiHandler)

	server.SetHandler(apiHandler)

	// Start servers
	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	// waitGroup for all main go routines
	var wg sync.WaitGroup
	// create error channel
	errCh := make(chan error)

	// start serving metrics
	go func() {
		wg.Add(1)
		defer wg.Done()
		errCh <- metricsServer.ServePrometheusMetrics(ctx, ":9090", "storage", logger)
	}()
	// start serving status
	go func() {
		wg.Add(1)
		defer wg.Done()
		ss := statusServer.New(logger)
		errCh <- ss.ListenAndServeHTTPs(ctx, "cloudStorage:4433", "", "/certs/cloudStorage.pem", "/certs/cloudStorage-key.pem")
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

type WildcardConsumer struct{}

func (w *WildcardConsumer) Consume(r io.Reader, in interface{}) error {
	b, err := ioutil.ReadAll(r)
	fmt.Println("WildcardConsumer::Consume", err, in, b)
	return nil
}
