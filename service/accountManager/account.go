package accountManager

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/auth/models"
)

// Service describes actions supported by the accountManager service
type Service interface {
	Users(context.Context, string) ([]*models.User, error)
	User(context.Context, string) (*models.User, error)
	AddUser(context.Context, *models.User) (*models.User, error)
	UpdateUser(context.Context, *models.User) (*models.User, error)
	RemoveUser(context.Context, string) error

	Roles(context.Context, string) ([]*models.Role, error)
	Role(context.Context, string) (*models.Role, error)
	AddRole(context.Context, *models.Role) (*models.Role, error)
	UpdateRole(context.Context, *models.Role) (*models.Role, error)
	RemoveRole(context.Context, string) error

	Rules(context.Context, string) ([]*models.Rule, error)
	Rule(context.Context, string) (*models.Rule, error)
	AddRule(context.Context, *models.Rule) (*models.Rule, error)
	UpdateRule(context.Context, *models.Rule) (*models.Rule, error)
	RemoveRule(context.Context, string) error
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
	logger  zerolog.Logger
}

// New returns a new instance of authenticator service
func New(storage Storage, logger zerolog.Logger) Service {
	logger.Debug().Msg("Initialize account manager service")

	return &accountManager{
		storage: storage,
		logger:  logger,
	}
}

// Users returns all users
func (a *accountManager) Users(_ context.Context, _ string) ([]*models.User, error) {
	return a.storage.GetUsers()
}

// User returns user by ID
func (a *accountManager) User(_ context.Context, userID string) (*models.User, error) {
	return a.storage.GetUser(userID)
}

// AddUser creates new user
func (a *accountManager) AddUser(_ context.Context, user *models.User) (*models.User, error) {
	return a.storage.AddUser(user)
}

// UpdateUser updates user
func (a *accountManager) UpdateUser(_ context.Context, user *models.User) (*models.User, error) {
	return a.storage.UpdateUser(user)
}

// RemoveUser removes user
func (a *accountManager) RemoveUser(_ context.Context, userID string) error {
	return a.storage.RemoveUser(userID)
}

// Roles returns all roles
func (a *accountManager) Roles(_ context.Context, _ string) ([]*models.Role, error) {
	return a.storage.GetRoles()
}

// Role returns role by ID
func (a *accountManager) Role(_ context.Context, roleID string) (*models.Role, error) {
	return a.storage.GetRole(roleID)
}

// AddRole creates new role
func (a *accountManager) AddRole(_ context.Context, role *models.Role) (*models.Role, error) {
	return a.storage.AddRole(role)
}

// UpdateRole updates role
func (a *accountManager) UpdateRole(_ context.Context, role *models.Role) (*models.Role, error) {
	return a.storage.UpdateRole(role)
}

// RemoveRole removes role
func (a *accountManager) RemoveRole(_ context.Context, roleID string) error {
	return a.storage.RemoveRole(roleID)
}

// Rules returns all rule
func (a *accountManager) Rules(_ context.Context, _ string) ([]*models.Rule, error) {
	return a.storage.GetRules()
}

// Rule returns rule by ID
func (a *accountManager) Rule(_ context.Context, ruleID string) (*models.Rule, error) {
	return a.storage.GetRule(ruleID)
}

// AddRule creates new rule
func (a *accountManager) AddRule(_ context.Context, rule *models.Rule) (*models.Rule, error) {
	return a.storage.AddRule(rule)
}

// UpdateRule updates rule
func (a *accountManager) UpdateRule(_ context.Context, rule *models.Rule) (*models.Rule, error) {
	return a.storage.UpdateRule(rule)
}

// RemoveRule removes rule
func (a *accountManager) RemoveRule(_ context.Context, ruleID string) error {
	return a.storage.RemoveRule(ruleID)
}
