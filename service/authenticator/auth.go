package authenticator

//go:generate ../../mockgen.sh $GOFILE

import (
	"context"
	"fmt"

	"github.com/iryonetwork/wwm/specs"
)

// Service describes the actions supported by the authenticator service
type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Validate(ctx context.Context, queries []*specs.ValidationPair) ([]*specs.ValidationResult, error)
}

// Storage describes the functionality required for the service to function
type Storage interface {
	GetUserByUsername(string) (*specs.User, error)
	FindACL(string, string, []specs.ACLRuleAction) ([]*specs.ACLRule, error)
}

type auth struct {
	storage Storage
}

// Login authenticates the user
func (a *auth) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.storage.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// @TODO !!fix!! ...
	if user.Password == password {
		return createTokenForUser(user)
	}

	return "", fmt.Errorf("User not found by username / password")
}

// Validate checks if the user has the capability to execute the specific
// actions on a resource
func (a *auth) Validate(ctx context.Context, queries []*specs.ValidationPair) ([]*specs.ValidationResult, error) {
	return nil, nil
}

// New returns a new instance of authenticator service
func New(storage Storage) Service {
	return &auth{storage: storage}
}
