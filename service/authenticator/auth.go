package authenticator

//go:generate ../../mockgen.sh $GOFILE

import (
	"context"

	"github.com/iryonetwork/wwm/specs"
)

// Service describes the actions supported by the authenticator service
type Service interface {
	Login(ctx context.Context, username, password string) (bool, error)
	Validate(ctx context.Context, resource, action string) (bool, error)
}

// Storage describes the functionality required for the service to function
type Storage interface {
	GetUser(string) (specs.User, error)
	FindACL(string, string, []specs.ACLAction) ([]specs.ACL, error)
}

type auth struct {
	storage Storage
}

// Login authenticates the user
func (a auth) Login(ctx context.Context, username, password string) (bool, error) {
	return true, nil
}

// Validate checks if the user has the capability to execute the specific
// actions on a resource
func (a auth) Validate(ctx context.Context, resource, action string) (bool, error) {
	return true, nil
}
