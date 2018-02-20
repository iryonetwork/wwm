package main

//go:generate sh -c "mkdir -p ../../gen/storage && swagger generate server -A storage -t ../../gen/storage -f ../../docs/api/storage.yml --exclude-main --principal string"

import (
	// "log"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	loads "github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/restapi"
	"github.com/iryonetwork/wwm/gen/storage/restapi/operations"
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

	// retry connection to nats if unsuccesful
	err = retry(5, time.Duration(500*time.Millisecond), 3.0, logger.With().Str("connect", "nats").Logger(), func() error {
		var err error
		nc, err = nats.Connect(URLs, nats.ClientCert(ClientCert, ClientKey))
		return err
	})
	if err == nil {
		err = retry(5, time.Duration(500*time.Millisecond), 3.0, logger.With().Str("connect", "nats-streaming").Logger(), func() error {
			var err error
			sc, err = stan.Connect(ClusterID, ClientID, stan.NatsConn(nc))
			return err
		})
	}
	if err == nil {
		p = publisher.New(context.Background(), sc, 5, time.Duration(10*time.Second), 2.0, logger.With().Str("component", "sync/storage/publisher").Logger())
	} else {
		// start localStorage with null publisher
		p = publisher.NewNullPublisher(context.Background())
		logger.Error().Msg("storage service will be started with null storage sync publisher due to failed nats-streaming connection attempts")
	}

	// initialize the service
	service := storage.New(s3, keys, p, logger.With().Str("component", "service/storage").Logger())

	api := operations.NewStorageAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	server := restapi.NewServer(api)
	server.TLSHost = "0.0.0.0"
	server.TLSPort = 443
	server.EnabledListeners = []string{"https"}
	server.TLSCertificateKey = "/certs/localStorage-key.pem"
	server.TLSCertificate = "/certs/localStorage.pem"

	defer server.Shutdown()

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

	server.SetHandler(apiLogMiddleware(api.Serve(nil), logger.With().Str("component", "logMW").Logger()))

	if err := server.Serve(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func apiLogMiddleware(next http.Handler, logger zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("New request")
		next.ServeHTTP(w, r)
	})
}

// retry helper method to sanely retry to connect
func retry(attempts int, sleep time.Duration, factor float32, logger zerolog.Logger, toRetry func() error) (err error) {
	for i := 0; ; i++ {
		err = toRetry()
		if err == nil {
			return nil
		}

		if i >= (attempts - 1) {
			break
		}

		logger.Error().Err(err).Msgf("retry number %d in %s", i+1, sleep)
		time.Sleep(sleep)
		sleep = time.Duration(float32(sleep) * factor) // increase time to sleep by factor
	}
	logger.Error().Msgf("failed to complete in %d retries", attempts)

	return err
}

type WildcardConsumer struct{}

func (w *WildcardConsumer) Consume(r io.Reader, in interface{}) error {
	b, err := ioutil.ReadAll(r)
	fmt.Println("WildcardConsumer::Consume", err, in, b)
	return nil
}
