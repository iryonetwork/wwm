package authenticator

import "context"

// Service describes the actions supported by the authenticator service
type Service interface {
	Login(ctx context.Context, username, password string) (bool, error)
	Validate(ctx context.Context, resource, action string) (bool, error)
}
