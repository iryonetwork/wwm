package main

//go:generate sh -c "mkdir -p ../../gen/auth/ && swagger generate server -A cloudAuth -t ../../gen/auth/ -f ../../docs/api/auth.yml --exclude-main --principal string"

import (
	loads "github.com/go-openapi/loads"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/iryonetwork/wwm/gen/auth/restapi"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations"
	"github.com/iryonetwork/wwm/service/authSync"
	"github.com/iryonetwork/wwm/service/authenticator"
	"github.com/iryonetwork/wwm/storage/auth"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// initialize storage
	// TODO: get key from vault
	key := []byte{0xe9, 0xf8, 0x2d, 0xf9, 0xc4, 0x14, 0xc1, 0x41, 0xdb, 0x87, 0x31, 0x1a, 0x95, 0x79, 0x5, 0xbf, 0x71, 0x12, 0x30, 0xd3, 0x2d, 0x8b, 0x59, 0x9d, 0x27, 0x13, 0xfa, 0x84, 0x55, 0x63, 0x64, 0x64}
	// TODO: can't open read only if there is no initial file... should sync db from cloud on first run
	storage, err := auth.New("/wwm/localAuth.db", key, false)
	if err != nil {
		log.Fatalln(err)
	}

	// initialize the service
	auth, err := authenticator.New(storage, nil)
	if err != nil {
		log.Fatalln(err)
	}

	authSync, err := authSync.New(log.WithField("component", "authSync"), storage, "/certs/localAuthSync.pem", "/certs/localAuthSync-key.pem", "https://cloudAuth/auth/database")
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewCloudAuthAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.TLSPort = 443
	server.TLSCertificate = "/certs/localAuth.pem"
	server.TLSCertificateKey = "/certs/localAuth-key.pem"
	server.EnabledListeners = []string{"https"}
	defer server.Shutdown()

	authHandlers := authenticator.NewHandlers(auth)

	api.Logger = log.WithField("component", "server").Errorf
	api.TokenAuth = auth.GetPrincipalFromToken
	api.APIAuthorizer = auth.Authorizer()
	api.AuthGetRenewHandler = authHandlers.GetRenew()
	api.AuthPostLoginHandler = authHandlers.PostLogin()
	api.AuthPostValidateHandler = authHandlers.PostValidate()

	server.SetHandler(api.Serve(nil))

	gocron.Every(5).Minutes().Do(authSync.Sync)
	go gocron.Start()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}