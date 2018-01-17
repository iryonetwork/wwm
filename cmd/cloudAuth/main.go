package main

//go:generate sh -c "mkdir -p ../../gen/auth/ && swagger generate server -A cloudAuth -t ../../gen/auth/ -f ../../docs/api/auth.yml --exclude-main --principal string"

import (
	"flag"
	"log"

	loads "github.com/go-openapi/loads"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/gen/auth/restapi"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations"
	"github.com/iryonetwork/wwm/service/accountManager"
	"github.com/iryonetwork/wwm/service/authenticator"
	"github.com/iryonetwork/wwm/storage/auth"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	createUsername := flag.String("username", "", "username to create")
	createPassword := flag.String("password", "", "password for new user")
	createEmail := flag.String("email", "", "email for new user")
	flag.Parse()

	// initialize storage
	// TODO: get key from vault
	key := []byte{0xe9, 0xf8, 0x2d, 0xf9, 0xc4, 0x14, 0xc1, 0x41, 0xdb, 0x87, 0x31, 0x1a, 0x95, 0x79, 0x5, 0xbf, 0x71, 0x12, 0x30, 0xd3, 0x2d, 0x8b, 0x59, 0x9d, 0x27, 0x13, 0xfa, 0x84, 0x55, 0x63, 0x64, 0x64}
	storage, err := auth.New("cloudAuth.db", key, false)
	if err != nil {
		log.Fatalln(err)
	}

	if *createUsername != "" {
		user := &models.User{
			Username: createUsername,
			Password: *createPassword,
			Email:    createEmail,
		}

		user, err := storage.AddUser(user)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = storage.AddUserToAdminRole(user.ID)
		if err != nil {
			log.Fatalln(err)
		}

		log.Fatalln("Created new user %s", *createUsername)
	}

	// initialize the service
	auth := authenticator.New(storage)
	account := accountManager.New(storage)

	api := operations.NewCloudAuthAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.TLSPort = 443
	server.TLSCertificate = "/certs/cloudAuth.pem"
	server.TLSCertificateKey = "/certs/cloudAuth-key.pem"
	server.EnabledListeners = []string{"https"}
	defer server.Shutdown()

	authHandlers := authenticator.NewHandlers(auth)

	api.TokenAuth = auth.GetUserIDFromToken
	api.APIAuthorizer = auth.Authorizer()
	api.AuthGetRenewHandler = authHandlers.GetRenew()
	api.AuthPostLoginHandler = authHandlers.PostLogin()
	api.AuthPostValidateHandler = authHandlers.PostValidate()

	api.UsersGetUsersHandler = getUsers(account)
	api.UsersGetUsersIDHandler = getUsersID(account)
	api.UsersPostUsersHandler = postUsers(account)
	api.UsersPutUsersIDHandler = putUsersID(account)
	api.UsersDeleteUsersIDHandler = deleteUsersID(account)

	api.RolesGetRolesHandler = getRoles(account)
	api.RolesGetRolesIDHandler = getRolesID(account)
	api.RolesPostRolesHandler = postRoles(account)
	api.RolesPutRolesIDHandler = putRolesID(account)
	api.RolesDeleteRolesIDHandler = deleteRolesID(account)

	api.RulesGetRulesHandler = getRules(account)
	api.RulesGetRulesIDHandler = getRulesID(account)
	api.RulesPostRulesHandler = postRules(account)
	api.RulesPutRulesIDHandler = putRulesID(account)
	api.RulesDeleteRulesIDHandler = deleteRulesID(account)

	server.SetHandler(api.Serve(nil))

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
