package accountManager

import (
	"context"

	"github.com/iryonetwork/wwm/specs"
)

// Service describes actions supported by the accountManager service
type Service interface {
	Users(context.Context, *specs.GetUsersRequest) ([]*specs.User, error)
	User(context.Context, string) (*specs.User, error)
	AddUser(context.Context, *specs.User) (*specs.User, error)
	UpdateUser(context.Context, *specs.User) (*specs.User, error)
	RemoveUser(context.Context, string) error

	Rules(context.Context, *specs.GetACLRulesRequest) ([]*specs.ACLRule, error)
	Rule(context.Context, string) (*specs.ACLRule, error)
	AddRule(context.Context, *specs.ACLRule) (*specs.ACLRule, error)
	UpdateRule(context.Context, *specs.ACLRule) (*specs.ACLRule, error)
	RemoveRule(context.Context, string) error
}

// Storage describes methods required from the storage used by the service
type Storage interface {
	GetUsers(*specs.GetUsersRequest) ([]*specs.User, error)
	GetUser(string) (*specs.User, error)
	GetUserBy(string) (*specs.User, error)
	AddUser(*specs.User) (*specs.User, error)
	UpdateUser(*specs.User) (*specs.User, error)
	RemoveUser(string) error

	GetRules(*specs.GetACLRulesRequest) ([]*specs.ACLRule, error)
	GetRule(string) (*specs.ACLRule, error)
	AddRule(*specs.ACLRule) (*specs.ACLRule, error)
	UpdateRule(*specs.ACLRule) (*specs.ACLRule, error)
	RemoveRule(string) error
}

type accountManager struct {
	storage Storage
}
