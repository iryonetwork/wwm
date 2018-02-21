package main

//go:generate sh -c "mkdir -p ../../gen/storage && swagger generate server -A storage -t ../../gen/storage -f ../../docs/api/storage.yml --exclude-main --principal string"

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	loads "github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/restapi"
	"github.com/iryonetwork/wwm/gen/storage/restapi/operations"
	storage "github.com/iryonetwork/wwm/service/storage"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/mock"
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
	s3, err := s3.New(cfg, keys, logger.With().Str("component", "storage/s3").Logger())
	if err != nil {
		log.Fatalln(err)
	}

	// initialize the service
	service := storage.New(s3, keys, logger.With().Str("component", "service/storage").Logger())

	api := operations.NewStorageAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	server := restapi.NewServer(api)
	server.TLSHost = "0.0.0.0"
	server.TLSPort = 443
	server.EnabledListeners = []string{"https"}
	server.TLSCertificateKey = "/certs/cloudStorage-key.pem"
	server.TLSCertificate = "/certs/cloudStorage.pem"

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

type WildcardConsumer struct{}

func (w *WildcardConsumer) Consume(r io.Reader, in interface{}) error {
	b, err := ioutil.ReadAll(r)
	fmt.Println("WildcardConsumer::Consume", err, in, b)
	return nil
}
