package accountManager

import (
	"context"

	"github.com/iryonetwork/wwm/gen/models"
)

// Service describes actions supported by the accountManager service
type Service interface {
	Users(context.Context, string) ([]*models.User, error)
	User(context.Context, string) (*models.User, error)
	AddUser(context.Context, *models.User) (*models.User, error)
	UpdateUser(context.Context, *models.User) (*models.User, error)
	RemoveUser(context.Context, string) error
}

// Storage describes methods required from the storage used by the service
type Storage interface {
	GetUsers() ([]*models.User, error)
	GetUser(string) (*models.User, error)
	AddUser(*models.User) (*models.User, error)
	UpdateUser(*models.User) (*models.User, error)
	RemoveUser(string) error

	GetRoles() ([]*models.Role, error)
	GetRole(string) (*models.Role, error)
	AddRole(*models.Role) (*models.Role, error)
	UpdateRole(*models.Role) (*models.Role, error)
	RemoveRole(string) error

	GetRules() ([]*models.Rule, error)
	GetRule(string) (*models.Rule, error)
	AddRule(*models.Rule) (*models.Rule, error)
	UpdateRule(*models.Rule) (*models.Rule, error)
	RemoveRule(string) error
}

type accountManager struct {
	storage Storage
}
