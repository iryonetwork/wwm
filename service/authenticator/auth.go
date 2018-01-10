package authenticator

//go:generate sh ../../bin/mockgen.sh service/authenticator Service,Storage $GOFILE

import (
	"context"
	"fmt"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/specs"
	"golang.org/x/crypto/bcrypt"
)

// Service describes the actions supported by the authenticator service
type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Validate(ctx context.Context, queries []*models.ValidationPair) ([]*models.ValidationResult, error)
	GetPublicKey(ctx context.Context, pubID string) (string, error)
	CreateTokenForUser(ctx context.Context, user *models.User) (string, error)
	GetUserFromToken(token string) (*models.User, error)
}

// Storage describes the functionality required for the service to function
type Storage interface {
	GetUserByUsername(string) (*models.User, error)
	FindACL(string, string, []specs.ACLRuleAction) ([]*specs.ACLRule, error)
	GetUserByID(id string) (*models.User, error)
}

type auth struct {
	storage Storage
}

// Login authenticates the user
func (a *auth) Login(_ context.Context, username, password string) (string, error) {
	user, err := a.storage.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		return createTokenForUser(user)
	}

	return "", fmt.Errorf("User not found by username / password")
}

// Validate checks if the user has the capability to execute the specific
// actions on a resource
func (a *auth) Validate(_ context.Context, queries []*models.ValidationPair) ([]*models.ValidationResult, error) {
	return nil, nil
}

// CreateTokenForUser creates a new token for user
func (a *auth) CreateTokenForUser(_ context.Context, user *models.User) (string, error) {
	return createTokenForUser(user)
}

// GetUserFromToken validates a token and returns the user if token is valid
func (a *auth) GetUserFromToken(token string) (*models.User, error) {
	userID, err := validateToken(token)
	if err != nil {
		return nil, err
	}

	return a.storage.GetUserByID(userID)
}

// GetPublicKey returns public key matching the pubID
func (a *auth) GetPublicKey(_ context.Context, pubID string) (string, error) {
	if pubID != keyID {
		return "", fmt.Errorf("Failed to find key with ID %s", keyID)
	}

	return publicKey, nil
}

// New returns a new instance of authenticator service
func New(storage Storage) Service {
	return &auth{storage: storage}
}
