package main

//go:generate sh -c "mkdir -p ../../gen/waitlist/ && swagger generate server -A waitlist -t ../../gen/waitlist/ -f ../../docs/api/waitlist.yml --exclude-main --principal string"

import (
	"os"

	loads "github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/restapi"
	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations"
	"github.com/iryonetwork/wwm/storage/waitlist"
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

	// TODO: don't hardcode this key
	key := []byte{0x8c, 0x7b, 0x71, 0x7f, 0xd9, 0x13, 0xaf, 0xef, 0x5d, 0xcb, 0x18, 0x84, 0xc9, 0x9c, 0xc, 0x44, 0x61, 0x8b, 0xa6, 0xa9, 0x78, 0x69, 0x31, 0x0, 0x21, 0x55, 0x51, 0x22, 0xc2, 0xf4, 0xa0, 0xe3}

	// initialize the service
	storage, err := waitlist.New("/data/waitlist.db", key, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize waitlist storage")
	}

	api := operations.NewWaitlistAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.TLSPort = 443
	server.TLSCertificate = "/certs/public.crt"
	server.TLSCertificateKey = "/certs/private.key"
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

	api.TokenAuth = func(token string) (*string, error) {
		return swag.String("test"), nil
	}

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(api.Serve(nil))

	server.SetHandler(handler)

	if err := server.Serve(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
