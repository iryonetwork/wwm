package authenticator

//go:generate sh ../../bin/mockgen.sh service/authenticator Service,Storage $GOFILE

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"golang.org/x/crypto/bcrypt"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/auth"
	"github.com/iryonetwork/wwm/utils"
)

// Service describes the actions supported by the authenticator service
type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Validate(ctx context.Context, userID *string, queries []*models.ValidationPair) ([]*models.ValidationResult, error)
	GetPublicKey(ctx context.Context, pubID string) (string, error)
	CreateTokenForUserID(ctx context.Context, userID *string) (string, error)
	GetUserIDFromToken(token string) (*string, error)
	Authorizer() runtime.Authorizer
}

// Storage describes the functionality required for the service to function
type Storage interface {
	GetUserByUsername(string) (*models.User, error)
	FindACL(subject string, actions []*models.ValidationPair) []*models.ValidationResult
	GetUser(id string) (*models.User, error)
}

type service struct {
	storage Storage
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
	return a.storage.FindACL(*userID, queries), nil
}

// CreateTokenForUser creates a new token for user
func (a *service) CreateTokenForUserID(_ context.Context, userID *string) (string, error) {
	return createTokenForUserID(userID)
}

// GetUserIDFromToken validates a token and returns the userID if token is valid
func (a *service) GetUserIDFromToken(token string) (*string, error) {
	userID, err := validateToken(token)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (a *service) Authorizer() runtime.Authorizer {
	return runtime.AuthorizerFunc(func(request *http.Request, principal interface{}) error {
		userID, ok := principal.(*string)
		if !ok {
			return fmt.Errorf("Principal type was '%T', expected '*string'", principal)
		}

		var action int64 = auth.Read
		if request.Method == http.MethodPost || request.Method == http.MethodPut {
			action = auth.Write
		} else if request.Method == http.MethodDelete {
			action = auth.Delete
		}

		result := a.storage.FindACL(*userID, []*models.ValidationPair{{
			Actions:  &action,
			Resource: swag.String(request.URL.EscapedPath()),
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

// New returns a new instance of authenticator service
func New(storage Storage) Service {
	return &service{storage: storage}
}
