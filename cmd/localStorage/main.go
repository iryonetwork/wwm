package main

//go:generate sh -c "mkdir -p ../../gen/storage && swagger generate server -A storage -t ../../gen/storage -f ../../docs/api/storage.yml --exclude-main --principal string"

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	loads "github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/restapi"
	"github.com/iryonetwork/wwm/gen/storage/restapi/operations"
	"github.com/iryonetwork/wwm/metrics"
	APIMetrics "github.com/iryonetwork/wwm/metrics/api"
	storage "github.com/iryonetwork/wwm/service/storage"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/mock"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/publisher"
	"github.com/iryonetwork/wwm/utils"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "localStorage").
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
		Endpoint:     "localMinio:9000",
		AccessKey:    "local",
		AccessSecret: "localminio",
		Secure:       true,
		Region:       "us-east-1",
	}
	s3, err := s3.New(cfg, keys, logger.With().Str("component", "storage/s3").Logger())
	if err != nil {
		log.Fatalln(err)
	}

	// initialize storageSync publisher
	// create nats/nats-streaming connection
	URLs := "tls://nats:secret@localNats:4242"
	ClusterID := "localNats"
	ClientID := "localStorage"
	ClientCert := "/certs/localStoragePublisher.pem"
	ClientKey := "/certs/localStoragePublisher-key.pem"
	var nc *nats.Conn
	var sc publisher.StanConnection
	var p storageSync.Publisher

	// Connect to NATS
	// retry connectng to nats if unsuccesful
	err = utils.Retry(5, time.Duration(500*time.Millisecond), 3.0, logger.With().Str("connect", "nats").Logger(), func() error {
		var err error
		nc, err = nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
		return err
	})

	// Connect to NATS-Streaming if NATS connection succesful
	if err == nil {
		// retry connecting to nats-straming if unsuccesful
		err = utils.Retry(5, time.Duration(500*time.Millisecond), 3.0, logger.With().Str("connect", "nats-streaming").Logger(), func() error {
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
		// Register metrics
		h := prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "publisher",
			Name:      "publish_seconds",
			Help:      "Time taken to publish task",
		})
		prometheus.MustRegister(h)

		c := prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "publisher",
			Name:      "publish_calls",
			Help:      "Number of publish calls to nats-streaming",
		})
		prometheus.MustRegister(c)

		cfg := publisher.Cfg{
			Connection:      sc,
			Retries:         5,
			StartRetryWait:  time.Duration(10 * time.Second),
			RetryWaitFactor: 2.0,
		}
		l := logger.With().Str("component", "sync/storage/publisher").Logger()
		p = publisher.New(context.Background(), cfg, l, h, c)
	}
	defer p.Close()

	// initialize the servicex
	service := storage.New(s3, keys, p, logger.With().Str("component", "service/storage").Logger())

	api := operations.NewStorageAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	server := restapi.NewServer(api)
	server.TLSHost = "0.0.0.0"
	server.TLSPort = 443
	server.EnabledListeners = []string{"https"}
	server.TLSCertificateKey = "/certs/localStorage-key.pem"
	server.TLSCertificate = "/certs/localStorage.pem"

	storageHandlers := storage.NewHandlers(service, logger.With().Str("component", "service/storage/handlers").Logger())

	api.TokenAuth = storageHandlers.GetUserIDFromToken
	api.APIAuthorizer = storageHandlers.Authorizer()
	api.StorageFileListHandler = storageHandlers.FileList()
	api.StorageFileGetHandler = storageHandlers.FileGet()
	api.StorageFileGetVersionHandler = storageHandlers.FileGetVersion()
	api.StorageFileListVersionsHandler = storageHandlers.FileListVersions()
	api.StorageFileNewHandler = storageHandlers.FileNew()
	api.StorageFileUpdateHandler = storageHandlers.FileUpdate()
	api.StorageFileDeleteHandler = storageHandlers.FileDelete()

	api.RegisterConsumer("*/*", &WildcardConsumer{})

	// initialize metrics middleware
	m := APIMetrics.NewMetrics("api", "").
		WithURLSanitize(utils.WhitelistURLSanitize([]string{"storage", "versions", "sync"}))

	// set handler with middlewares
	apiHandler := apiLogMiddleware(api.Serve(nil), logger.With().Str("component", "logMW").Logger())
	apiHandler = m.Middleware(apiHandler)

	server.SetHandler(apiHandler)

	// Start servers
	errCh := make(chan error)
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(errCh)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		errCh <- metrics.ServePrometheusMetrics(context.Background(), ":9090", "storage")
	}()

	go func() {
		wg.Add(1)
		defer server.Shutdown()
		defer wg.Done()
		errCh <- server.Serve()
	}()

	for err := range errCh {
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}
}

func apiLogMiddleware(next http.Handler, logger zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("New request")
		next.ServeHTTP(w, r)
	})
}

type WildcardConsumer struct{}

func (w *WildcardConsumer) Consume(r io.Reader, in interface{}) error {
	b, err := ioutil.ReadAll(r)
	fmt.Println("WildcardConsumer::Consume", err, in, b)
	return nil
}
