package authenticator

//go:generate ../../bin/mockgen.sh service/authenticator Service,Storage $GOFILE

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/gobwas/glob"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/bcrypt"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/auth"
	"github.com/iryonetwork/wwm/utils"
)

// Service describes the actions supported by the authenticator service
type Service interface {
	// Login returns token that will be used for next requests or error if username/password is wrong
	Login(ctx context.Context, username, password string) (string, error)

	// Validate checks if user has permissions for specified paths and operations
	Validate(ctx context.Context, userID *string, queries []*models.ValidationPair) ([]*models.ValidationResult, error)

	// GetPublicKey returns public key
	GetPublicKey(ctx context.Context, pubID string) (string, error)

	// CreateTokenForUserID return token that is used for authentication
	CreateTokenForUserID(ctx context.Context, userID *string) (string, error)

	// GetPrincipalFromToken returns user ID if token is valid
	GetPrincipalFromToken(token string) (*string, error)

	// Authorizer returns function that checks if user is authorized to make a request
	Authorizer() runtime.Authorizer
}

// Storage describes the functionality required for the service to function
type Storage interface {
	GetUserByUsername(string) (*models.User, error)
	FindACL(subject string, actions []*models.ValidationPair) []*models.ValidationResult
	GetUser(id string) (*models.User, error)
}

type service struct {
	storage      Storage
	syncServices map[string]syncService
	logger       zerolog.Logger
}

// Login authenticates the user
func (a *service) Login(_ context.Context, username, password string) (string, error) {
	user, err := a.storage.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	permissions := a.storage.FindACL(user.ID, []*models.ValidationPair{{
		Actions:  swag.Int64(auth.Write),
		Resource: swag.String("/auth/login"),
	}})

	if !permissions[0].Result {
		return "", utils.NewError(utils.ErrForbidden, "You do not have permission to log in")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		return createTokenForUserID(&user.ID)
	}

	return "", fmt.Errorf("User not found by username / password")
}

// Validate checks if the user has the capability to execute the specific
// actions on a resource
func (a *service) Validate(_ context.Context, userID *string, queries []*models.ValidationPair) ([]*models.ValidationResult, error) {
	if strings.HasPrefix(*userID, servicePrincipal) {
		keyID := (*userID)[len(servicePrincipal):]
		s, _ := a.syncServices[keyID]

		results := make([]*models.ValidationResult, len(queries))
		for i, query := range queries {
			results[i] = &models.ValidationResult{
				Query:  query,
				Result: s.glob.Match(strings.TrimPrefix(*query.Resource, "/api")),
			}
		}
		return results, nil
	}

	return a.storage.FindACL(*userID, queries), nil
}

// CreateTokenForUser creates a new token for user
func (a *service) CreateTokenForUserID(_ context.Context, userID *string) (string, error) {
	return createTokenForUserID(userID)
}

const servicePrincipal = "__service__"

// GetPrincipalFromToken validates a token and returns the userID for user tokens
// or returns "__service__<KeyID>" for tokens used in cloud sync
func (a *service) GetPrincipalFromToken(tokenString string) (*string, error) {
	principal := ""

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*Claims)

		if s, ok := a.syncServices[claims.KeyID]; ok {
			principal = servicePrincipal + claims.KeyID
			return s.publicKey, nil
		}

		if claims.KeyID == keyID {
			private, err := getPrivateKey()
			if err != nil {
				return nil, err
			}
			principal = claims.Subject
			return private.Public(), nil
		}

		return nil, fmt.Errorf("Signing key not found")
	})

	if err != nil {
		return swag.String(""), err
	}

	_, ok := token.Claims.(*Claims)
	if !token.Valid || !ok {
		return swag.String(""), fmt.Errorf("Token is invalid")
	}

	return &principal, nil
}

func (a *service) Authorizer() runtime.Authorizer {
	return runtime.AuthorizerFunc(func(request *http.Request, principal interface{}) error {
		userID, ok := principal.(*string)
		if !ok {
			return fmt.Errorf("Principal type was '%T', expected '*string'", principal)
		}

		// allow access for service operations without checking ACL
		if strings.HasPrefix(*userID, servicePrincipal) {
			keyID := (*userID)[len(servicePrincipal):]
			s, ok := a.syncServices[keyID]
			if ok && (request.URL.EscapedPath() == "/auth/validate" || s.glob.Match(request.URL.EscapedPath())) {
				return nil
			}
			return utils.NewError(utils.ErrForbidden, "You do not have permissions for this resource")
		}

		var action int64 = auth.Read
		if request.Method == http.MethodPost || request.Method == http.MethodPut {
			action = auth.Write
		} else if request.Method == http.MethodDelete {
			action = auth.Delete
		}

		result := a.storage.FindACL(*userID, []*models.ValidationPair{{
			Actions:  &action,
			Resource: swag.String("/api" + request.URL.EscapedPath()),
		}})

		if !result[0].Result {
			return utils.NewError(utils.ErrForbidden, "You do not have permissions for this resource")
		}

		return nil
	})
}

// GetPublicKey returns public key matching the pubID
func (a *service) GetPublicKey(_ context.Context, pubID string) (string, error) {
	if pubID != keyID {
		return "", fmt.Errorf("Failed to find key with ID %s", keyID)
	}

	return publicKey, nil
}

type syncService struct {
	publicKey crypto.PublicKey
	glob      glob.Glob
}

// New returns a new instance of authenticator service
func New(storage Storage, allowedServiceCertsAndPaths map[string][]string, logger zerolog.Logger) (Service, error) {
	logger = logger.With().Str("component", "service/authenticator").Logger()
	logger.Debug().Msg("Initialize authenticator service")

	syncServices := map[string]syncService{}

	for cert, paths := range allowedServiceCertsAndPaths {
		content, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(content)
		if block == nil {
			return nil, fmt.Errorf("Invalid PEM file")
		}

		c, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		thumb, err := acme.JWKThumbprint(c.PublicKey)
		if err != nil {
			return nil, err
		}

		g, err := glob.Compile("{" + strings.Join(paths, ",") + "}")
		if err != nil {
			return nil, err
		}

		s := syncService{
			publicKey: c.PublicKey,
			glob:      g,
		}

		syncServices[thumb] = s
	}

	return &service{
		storage:      storage,
		syncServices: syncServices,
		logger:       logger,
	}, nil
}
