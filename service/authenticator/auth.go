package authenticator

//go:generate ../../bin/mockgen.sh service/authenticator Service,AuthDataService,Enforcer $GOFILE

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type Claims struct {
	KeyID string `json:"kid"`
	jwt.StandardClaims
}

const servicePrincipal = "__service__"

var tokenExpiersIn = time.Duration(15) * time.Minute

type service struct {
	domainType         string
	domainID           string
	authData           AuthDataService
	enforcer           Enforcer
	syncServices       map[string]syncService
	jwtKeyID           string
	jwtPrivateKey      *rsa.PrivateKey
	jwtPublicKeyString string
	logger             zerolog.Logger
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
		token, err := a.CreateTokenForUserID(ctx, &user.ID)
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

		if claims.KeyID == a.jwtKeyID {
			principal = claims.Subject
			return a.jwtPrivateKey.Public(), nil
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
	if pubID != a.jwtKeyID {
		return "", fmt.Errorf("Failed to find key with ID %s", a.jwtKeyID)
	}

	return a.jwtPublicKeyString, nil
}

// CreateTokenForUserID creates a new token from user ID
func (a *service) CreateTokenForUserID(_ context.Context, id *string) (string, error) {
	// compose the claims
	claims := &Claims{
		KeyID: a.jwtKeyID,
		StandardClaims: jwt.StandardClaims{
			Subject:   *id,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenExpiersIn).Unix(),
		},
	}

	// create the token
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(a.jwtPrivateKey)
}

type syncService struct {
	publicKey crypto.PublicKey
	glob      glob.Glob
}

// New returns a new instance of authenticator service
func New(domainType, domainID string, authData AuthDataService, enforcer Enforcer, jwtPrivateKeyPath string, allowedServiceCertsAndPaths map[string][]string, logger zerolog.Logger) (Service, error) {
	logger = logger.With().Str("component", "service/authenticator").Logger()
	logger.Debug().Msg("Initialize authenticator service")

	// read jwt signing keys
	jwtPrivateKey, jwtPublicKeyString, jwtPublicKeyThumb, err := parseJwtKeys(jwtPrivateKeyPath)
	if err != nil {
		return nil, err
	}

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
		domainType:         domainType,
		domainID:           domainID,
		authData:           authData,
		enforcer:           enforcer,
		syncServices:       syncServices,
		jwtKeyID:           jwtPublicKeyThumb,
		jwtPrivateKey:      jwtPrivateKey,
		jwtPublicKeyString: jwtPublicKeyString,
		logger:             logger,
	}, nil
}

func parseJwtKeys(jwtPrivateKeyPath string) (*rsa.PrivateKey, string, string, error) {
	// read jwt signing key
	privateKeyRaw, err := ioutil.ReadFile(jwtPrivateKeyPath)
	if err != nil {
		return nil, "", "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyRaw)
	if err != nil {
		return nil, "", "", err
	}

	publicKey, ok := (privateKey.Public()).(*rsa.PublicKey)
	if !ok {
		return nil, "", "", fmt.Errorf("failed to get public key out of private key")
	}
	publicKeyThumb, err := acme.JWKThumbprint(publicKey)
	if err != nil {
		return nil, "", "", err
	}

	var pemPublicKeyBlock = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}
	publicKeyBuf := new(bytes.Buffer)
	err = pem.Encode(publicKeyBuf, pemPublicKeyBlock)
	if err != nil {
		return nil, "", "", err
	}

	return privateKey, publicKeyBuf.String(), publicKeyThumb, nil
}
