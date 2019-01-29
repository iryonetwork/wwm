package authenticator

//go:generate ../../bin/mockgen.sh service/authenticator Service,AuthDataService,Enforcer $GOFILE

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/gobwas/glob"
	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/bcrypt"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/metrics"
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

// AuthDataService describes the functionality of AuthData service needed by authenicator service
type AuthDataService interface {
	UserByUsername(ctx context.Context, username string) (*models.User, error)
}

type Enforcer interface {
	Enforce(rvals ...interface{}) bool
	LoadPolicy() error
	HasPolicy(params ...interface{}) bool
	GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector
}

type service struct {
	domainType   string
	domainID     string
	authData     AuthDataService
	enforcer     Enforcer
	syncServices map[string]syncService
	logger       zerolog.Logger
}

// Login authenticates the user
func (a *service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.authData.UserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	permissions := a.validatePairs(user.ID, []*models.ValidationPair{{
		Actions:    swag.Int64(auth.Write),
		DomainType: swag.String(a.domainType),
		DomainID:   swag.String(a.domainID),
		Resource:   swag.String("/auth/login"),
	}})

	if !*permissions[0].Result {
		return "", utils.NewError(utils.ErrForbidden, "You do not have permission to log in")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		token, err := createTokenForUserID(&user.ID)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	return "", fmt.Errorf("User not found by username / password")
}

// Validate checks if the user has the capability to execute the specific
// actions on a resource
func (a *service) Validate(ctx context.Context, userID *string, queries []*models.ValidationPair) ([]*models.ValidationResult, error) {
	// validate queries
	for _, query := range queries {
		if query.Actions == nil || query.Resource == nil || query.DomainType == nil || query.DomainID == nil {
			return nil, utils.NewError(utils.ErrBadRequest, "Missing validation query parameters")
		}
	}

	if strings.HasPrefix(*userID, servicePrincipal) {
		keyID := (*userID)[len(servicePrincipal):]
		s := a.syncServices[keyID]
		results := make([]*models.ValidationResult, len(queries))
		for i, query := range queries {
			results[i] = &models.ValidationResult{
				Query:  query,
				Result: swag.Bool(s.glob.Match(*query.Resource)),
			}
		}
		return results, nil
	}

	return a.validatePairs(*userID, queries), nil
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
		a.logger.Error().Err(err).Msgf("failed to parse token: %v", tokenString)
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
			if ok && (request.URL.EscapedPath() == "/auth/validate" || s.glob.Match("/api"+request.URL.EscapedPath())) {
				return nil
			}
			return utils.NewError(utils.ErrForbidden, "You do not have permissions for this resource")
		}

		var action int64
		switch request.Method {
		case http.MethodPost:
			action = auth.Write
		case http.MethodPut:
			action = auth.Update
		case http.MethodDelete:
			action = auth.Delete
		default:
			action = auth.Read
		}

		result := a.validatePairs(*userID, []*models.ValidationPair{{
			DomainType: &a.domainType,
			DomainID:   &a.domainID,
			Actions:    &action,
			Resource:   swag.String("/api" + request.URL.EscapedPath()),
		}})

		if !*result[0].Result {
			return utils.NewError(utils.ErrForbidden, "You do not have permissions for this resource")
		}

		return nil
	})
}

// validatePairs validates ValidationPair array
func (a *service) validatePairs(subject string, actions []*models.ValidationPair) []*models.ValidationResult {
	results := make([]*models.ValidationResult, len(actions))

	for i, validation := range actions {
		var domain string
		if *validation.DomainType == authCommon.DomainTypeGlobal {
			domain = "*"
		} else {
			domain = fmt.Sprintf("%s.%s", *validation.DomainType, *validation.DomainID)
		}

		results[i] = &models.ValidationResult{
			Query:  validation,
			Result: swag.Bool(a.enforcer.Enforce(subject, domain, *validation.Resource, strconv.FormatInt(*validation.Actions, 10))),
		}
	}

	return results
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
func New(domainType, domainID string, authData AuthDataService, enforcer Enforcer, allowedServiceCertsAndPaths map[string][]string, logger zerolog.Logger) (Service, error) {
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
		domainType:   domainType,
		domainID:     domainID,
		authData:     authData,
		enforcer:     enforcer,
		syncServices: syncServices,
		logger:       logger,
	}, nil
}
