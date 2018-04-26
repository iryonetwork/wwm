package main

//go:generate sh -c "mkdir -p ../../gen/waitlist/ && swagger generate server -A waitlist -t ../../gen/waitlist/ -f ../../docs/api/waitlist.yml --exclude-main --principal string"

import (
	"encoding/base64"
	"fmt"
	"os"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/restapi"
	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations"
	"github.com/iryonetwork/wwm/service/authorizer"
	"github.com/iryonetwork/wwm/storage/waitlist"
	"github.com/iryonetwork/wwm/utils"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "waitlist").
		Logger()

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

	auth := authorizer.New(cfg.DomainType, cfg.DomainID, fmt.Sprintf("https://%s/%s/validate", cfg.AuthHost, cfg.AuthPath), logger)

	api := operations.NewWaitlistAPI(swaggerSpec)
	api.ServeError = utils.ServeError
	api.TokenAuth = auth.GetPrincipalFromToken
	api.APIAuthorizer = auth.Authorizer()

	server := restapi.NewServer(api)
	server.TLSHost = cfg.ServerHost
	server.TLSPort = cfg.ServerPort
	server.TLSCertificate = flags.Filename(cfg.CertPath)
	server.TLSCertificateKey = flags.Filename(cfg.KeyPath)
	server.EnabledListeners = []string{"https"}
	defer server.Shutdown()

	h := &handlers{s: storage}

	api.WaitlistDeleteListIDHandler = h.WaitlistDeleteListID()
	api.WaitlistGetHandler = h.WaitlistGet()
	api.WaitlistPostHandler = h.WaitlistPost()
	api.WaitlistPutListIDHandler = h.WaitlistPutListID()

	api.ItemDeleteListIDItemIDHandler = h.ItemDeleteListIDItemID()
	api.ItemGetListIDHandler = h.ItemGetListID()
	api.ItemPostListIDHandler = h.ItemPostListID()
	api.ItemPutListIDItemIDHandler = h.ItemPutListIDItemID()

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(api.Serve(nil))

	server.SetHandler(handler)

	if err := server.Serve(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
