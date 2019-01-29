package authDataManager

//go:generate ../../bin/mockgen.sh service/authDataManager Storage $GOFILE

import (
	"context"
	"io"

	"github.com/rs/zerolog"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
)

// Service describes actions supported by the authDataManager service
type Service interface {
	// Users returns all users
	Users(ctx context.Context) ([]*models.User, error)

	// UserByUsername returns user by username
	UserByUsername(ctx context.Context, username string) (*models.User, error)

	// Users returns user by its ID
	User(ctx context.Context, id string) (*models.User, error)

	// UserRoleIDs fetches IDs of roles that the user has been assigned (with optional domain filtering).
	UserRoleIDs(ctx context.Context, id string, domainType, domainID *string) ([]string, error)

	// UserOrganizationIDs fetches IDs of organizations at which the user has been assigned a role (with optional role ID filtering).
	UserOrganizationIDs(ctx context.Context, id string, roleID *string) ([]string, error)

	// UserClinicIDs fetches IDs of clinics at which the user has been assigned a role (with optional role ID filtering).
	UserClinicIDs(ctx context.Context, id string, roleID *string) ([]string, error)

	// UserLocationIDs fetches IDs of locations at which the user has been assigned a role (with optional role ID filtering); both locations of clinics and locations at which user has been assigned a role manually are returned.
	UserLocationIDs(ctx context.Context, id string, roleID *string) ([]string, error)

	// AddUser adds user
	AddUser(ctx context.Context, user *models.User) (*models.User, error)

	// UpdateUser updates user
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)

	// RemoveUser removes user by its ID
	RemoveUser(ctx context.Context, id string) error

	// Roles returns all roles
	Roles(ctx context.Context) ([]*models.Role, error)

	// Role returns role by its ID
	Role(ctx context.Context, id string) (*models.Role, error)

	// RoleUsers fetches list of IDs of users that have been assigned the role (with optional domain filtering).
	RoleUserIDs(ctx context.Context, id string, domainType *string, domainID *string) ([]string, error)

	// AddRole creates new role
	AddRole(ctx context.Context, role *models.Role) (*models.Role, error)

	// UpdateRole updates role
	UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error)

	// RemoveRole removes role by its ID
	RemoveRole(ctx context.Context, id string) error

	// Rules returns all ACL rules
	Rules(ctx context.Context) ([]*models.Rule, error)

	// Rule returns rule by its ID
	Rule(ctx context.Context, id string) (*models.Rule, error)

	// AddRule creates new rule
	AddRule(ctx context.Context, rule *models.Rule) (*models.Rule, error)

	// UpdateRule updates rule
	UpdateRule(ctx context.Context, rule *models.Rule) (*models.Rule, error)

	// RemoveRule removes rule by its ID
	RemoveRule(ctx context.Context, id string) error

	// Organizations returns all organizations
	Organizations(ctx context.Context) ([]*models.Organization, error)

	// Organizations returns organization by its ID
	Organization(ctx context.Context, id string) (*models.Organization, error)

	// OrganizationLocationIDs returns IDs of organization's locations by organization's ID
	OrganizationLocationIDs(ctx context.Context, id string) ([]string, error)

	// AddOrganization creates new organization
	AddOrganization(ctx context.Context, organization *models.Organization) (*models.Organization, error)

	// UpdateOrganization updates organization
	UpdateOrganization(ctx context.Context, organization *models.Organization) (*models.Organization, error)

	// RemoveOrganization removes organization by its ID
	RemoveOrganization(ctx context.Context, id string) error

	// Clinics returns all clinics
	Clinics(ctx context.Context) ([]*models.Clinic, error)

	// Clinics returns clinic by its ID
	Clinic(ctx context.Context, id string) (*models.Clinic, error)

	// AddClinic creates new clinic
	AddClinic(ctx context.Context, clinic *models.Clinic) (*models.Clinic, error)

	// UpdateClinic updates clinic
	UpdateClinic(ctx context.Context, clinic *models.Clinic) (*models.Clinic, error)

	// RemoveClinic removes clinic by its ID
	RemoveClinic(ctx context.Context, id string) error

	// Locations returns all locations
	Locations(ctx context.Context) ([]*models.Location, error)

	// Locations returns location by its ID
	Location(ctx context.Context, id string) (*models.Location, error)

	// LocationOrganizationIDs returns IDs of location's organizations by location's ID
	LocationOrganizationIDs(ctx context.Context, id string) ([]string, error)

	// LocationUserIDs fetches list of IDs of users that have been assigned a role at the location (with optional role ID filtering); both users of clinics associated with the locations and users that have been assigned a role at the location manually are returned.
	LocationUserIDs(ctx context.Context, id string, roleID *string) ([]string, error)

	// AddLocation creates new location
	AddLocation(ctx context.Context, location *models.Location) (*models.Location, error)

	// UpdateLocation updates location
	UpdateLocation(ctx context.Context, location *models.Location) (*models.Location, error)

	// RemoveLocation removes location by its ID
	RemoveLocation(ctx context.Context, id string) error

	// FindUserRoles returns user roles based on filtering query parameters.
	FindUserRoles(ctx context.Context, userID *string, roleID *string, domainType *string, domainID *string) ([]*models.UserRole, error)

	// UserRole returns user role by its ID
	UserRole(ctx context.Context, id string) (*models.UserRole, error)

	// AddRole creates a new user role
	AddUserRole(ctx context.Context, userRole *models.UserRole) (*models.UserRole, error)

	// RemoveUserRole removes user role by its ID
	RemoveUserRole(ctx context.Context, id string) error

	// DomainUserIDs fetches list of IDs of users that have been assigned a role at the domain (with optional role ID filtering).
	DomainUserIDs(ctx context.Context, domainType, domainID, roleID *string) ([]string, error)

	// DBChecksum fetches checksum of underlying database
	DBChecksum() ([]byte, error)

	// WriteDBTo writes the whole underlying database to a writer
	WriteDBTo(writer io.Writer) (int64, error)
}

// Storage describes methods required from the storage used by the service
type Storage interface {
	GetUsers() ([]*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	GetUser(id string) (*models.User, error)
	AddUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	RemoveUser(id string) error

	GetRoles() ([]*models.Role, error)
	GetRole(id string) (*models.Role, error)
	AddRole(role *models.Role) (*models.Role, error)
	UpdateRole(role *models.Role) (*models.Role, error)
	RemoveRole(id string) error

	GetRules() ([]*models.Rule, error)
	GetRule(id string) (*models.Rule, error)
	AddRule(rule *models.Rule) (*models.Rule, error)
	UpdateRule(rule *models.Rule) (*models.Rule, error)
	RemoveRule(id string) error

	GetOrganizations() ([]*models.Organization, error)
	GetOrganization(id string) (*models.Organization, error)
	GetOrganizationClinics(id string) ([]*models.Clinic, error)
	GetOrganizationLocationIDs(id string) ([]string, error)
	AddOrganization(organization *models.Organization) (*models.Organization, error)
	UpdateOrganization(organization *models.Organization) (*models.Organization, error)
	RemoveOrganization(id string) error

	GetClinics() ([]*models.Clinic, error)
	GetClinic(id string) (*models.Clinic, error)
	GetClinicOrganization(id string) (*models.Organization, error)
	GetClinicLocation(id string) (*models.Location, error)
	AddClinic(clinic *models.Clinic) (*models.Clinic, error)
	UpdateClinic(clinic *models.Clinic) (*models.Clinic, error)
	RemoveClinic(string) error

	GetLocations() ([]*models.Location, error)
	GetLocation(id string) (*models.Location, error)
	GetLocationClinics(id string) ([]*models.Clinic, error)
	GetLocationOrganizationIDs(id string) ([]string, error)
	AddLocation(location *models.Location) (*models.Location, error)
	UpdateLocation(location *models.Location) (*models.Location, error)
	RemoveLocation(id string) error

	GetUserRoles() ([]*models.UserRole, error)
	GetUserRole(id string) (*models.UserRole, error)
	GetUserRoleByContent(userID string, roleID string, domainType string, domainID string) (*models.UserRole, error)
	FindUserRoles(userID *string, roleID *string, domainType *string, domainID *string) ([]*models.UserRole, error)
	AddUserRole(userRole *models.UserRole) (*models.UserRole, error)
	RemoveUserRole(id string) error

	GetChecksum() ([]byte, error)
	WriteTo(writer io.Writer) (int64, error)
}

type authDataManager struct {
	storage Storage
	logger  zerolog.Logger
}

// New returns a new instance of auth data manager service
func New(storage Storage, logger zerolog.Logger) Service {
	logger.Debug().Msg("Initialize auth data manager service")

	return &authDataManager{
		storage: storage,
		logger:  logger,
	}
}

// Users returns all users
func (a *authDataManager) Users(_ context.Context) ([]*models.User, error) {
	return a.storage.GetUsers()
}

// UserByUsername returns user by username
func (a *authDataManager) UserByUsername(_ context.Context, username string) (*models.User, error) {
	return a.storage.GetUserByUsername(username)
}

// User returns user by ID
func (a *authDataManager) User(_ context.Context, userID string) (*models.User, error) {
	return a.storage.GetUser(userID)
}

// UserRoleIDs fetches IDs of roles that the user has been assigned (with optional domain filtering).
func (a *authDataManager) UserRoleIDs(_ context.Context, id string, domainType, domainID *string) ([]string, error) {
	// create roles map to avoid duplicates
	roleIDsMap := make(map[string]bool)

	// find user roles
	userRoles, err := a.storage.FindUserRoles(&id, nil, domainType, domainID)
	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := roleIDsMap[*userRole.RoleID]
		if !ok {
			roleIDsMap[*userRole.RoleID] = true
		}
	}

	// create response
	roleIDs := []string{}
	for roleID := range roleIDsMap {
		roleIDs = append(roleIDs, roleID)
	}

	return roleIDs, nil
}

// UserOrganizationIDs fetches IDs of organizations at which the user has been assigned a role (with optional role ID filtering).
func (a *authDataManager) UserOrganizationIDs(_ context.Context, id string, roleID *string) ([]string, error) {
	// create organizations map to avoid duplicates
	organizationIDsMap := make(map[string]bool)

	// find user roles
	userRoles, err := a.storage.FindUserRoles(&id, roleID, &authCommon.DomainTypeOrganization, nil)

	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := organizationIDsMap[*userRole.DomainID]
		if !ok {
			organizationIDsMap[*userRole.DomainID] = true
		}
	}

	// create response
	organizationIDs := []string{}
	for organizationID := range organizationIDsMap {
		organizationIDs = append(organizationIDs, organizationID)
	}

	return organizationIDs, nil
}

// UserClinicIDs fetches IDs of clinics at which the user has been assigned a role (with optional role ID filtering).
func (a *authDataManager) UserClinicIDs(_ context.Context, id string, roleID *string) ([]string, error) {
	// create clinics map to avoid duplicates
	clinicIDsMap := make(map[string]bool)

	// find roles
	userRoles, err := a.storage.FindUserRoles(&id, roleID, &authCommon.DomainTypeClinic, nil)
	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := clinicIDsMap[*userRole.DomainID]
		if !ok {
			clinicIDsMap[*userRole.DomainID] = true
		}
	}

	// create response
	clinicIDs := []string{}
	for clinicID := range clinicIDsMap {
		clinicIDs = append(clinicIDs, clinicID)
	}

	return clinicIDs, nil
}

// UserLocationIDs fetches IDs of locations at which the user has been assigned a role (with optional role ID filtering); both locations of clinics and locations at which user has been assigned a role manually are returned.
func (a *authDataManager) UserLocationIDs(_ context.Context, id string, roleID *string) ([]string, error) {
	// create locations map to avoid duplicates
	locationIDsMap := make(map[string]bool)

	// find user roles for location domain type
	userRoles, err := a.storage.FindUserRoles(&id, roleID, &authCommon.DomainTypeLocation, nil)
	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := locationIDsMap[*userRole.DomainID]
		if !ok {
			locationIDsMap[*userRole.DomainID] = true
		}
	}

	// create clinics map to avoid fetching clinic twice
	clinicFetchedMap := make(map[string]bool)

	// find user roles for clinic domain type
	userRoles, err = a.storage.FindUserRoles(&id, roleID, &authCommon.DomainTypeClinic, nil)
	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := clinicFetchedMap[*userRole.DomainID]
		if !ok {
			// fetch clinic
			clinic, err := a.storage.GetClinic(*userRole.DomainID)
			if err != nil {
				return nil, err
			}
			// put clinic in the map
			clinicFetchedMap[clinic.ID] = true

			// check if location already in the map
			_, ok := locationIDsMap[*clinic.Location]
			if !ok {
				locationIDsMap[*clinic.Location] = true
			}
		}
	}

	// create response
	locationIDs := []string{}
	for locationID := range locationIDsMap {
		locationIDs = append(locationIDs, locationID)
	}

	return locationIDs, nil
}

// AddUser creates new user
func (a *authDataManager) AddUser(_ context.Context, user *models.User) (*models.User, error) {
	return a.storage.AddUser(user)
}

// UpdateUser updates user
func (a *authDataManager) UpdateUser(_ context.Context, user *models.User) (*models.User, error) {
	return a.storage.UpdateUser(user)
}

// RemoveUser removes user
func (a *authDataManager) RemoveUser(_ context.Context, userID string) error {
	return a.storage.RemoveUser(userID)
}

// Roles returns all roles
func (a *authDataManager) Roles(_ context.Context) ([]*models.Role, error) {
	return a.storage.GetRoles()
}

// Role returns role by ID
func (a *authDataManager) Role(_ context.Context, roleID string) (*models.Role, error) {
	return a.storage.GetRole(roleID)
}

// RoleUserIDs fetches list of IDs of users that have been assigned the role (with optional domain filtering).
func (a *authDataManager) RoleUserIDs(_ context.Context, id string, domainType *string, domainID *string) ([]string, error) {
	// create users map to avoid duplicates
	userIDsMap := make(map[string]bool)

	// find user roles
	userRoles, err := a.storage.FindUserRoles(nil, &id, domainType, domainID)
	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := userIDsMap[*userRole.UserID]
		if !ok {
			userIDsMap[*userRole.UserID] = true
		}
	}

	// create response
	userIDs := []string{}
	for userID := range userIDsMap {
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// AddRole creates new role
func (a *authDataManager) AddRole(_ context.Context, role *models.Role) (*models.Role, error) {
	return a.storage.AddRole(role)
}

// UpdateRole updates role
func (a *authDataManager) UpdateRole(_ context.Context, role *models.Role) (*models.Role, error) {
	return a.storage.UpdateRole(role)
}

// RemoveRole removes role
func (a *authDataManager) RemoveRole(_ context.Context, roleID string) error {
	return a.storage.RemoveRole(roleID)
}

// Rules returns all rule
func (a *authDataManager) Rules(_ context.Context) ([]*models.Rule, error) {
	return a.storage.GetRules()
}

// Rule returns rule by ID
func (a *authDataManager) Rule(_ context.Context, ruleID string) (*models.Rule, error) {
	return a.storage.GetRule(ruleID)
}

// AddRule creates new rule
func (a *authDataManager) AddRule(_ context.Context, rule *models.Rule) (*models.Rule, error) {
	return a.storage.AddRule(rule)
}

// UpdateRule updates rule
func (a *authDataManager) UpdateRule(_ context.Context, rule *models.Rule) (*models.Rule, error) {
	return a.storage.UpdateRule(rule)
}

// RemoveRule removes rule
func (a *authDataManager) RemoveRule(_ context.Context, ruleID string) error {
	return a.storage.RemoveRule(ruleID)
}

// Organizations returns all organizations
func (a *authDataManager) Organizations(_ context.Context) ([]*models.Organization, error) {
	return a.storage.GetOrganizations()
}

// Organization returns organization by ID
func (a *authDataManager) Organization(_ context.Context, organizationID string) (*models.Organization, error) {
	return a.storage.GetOrganization(organizationID)
}

// OrganizationLocationIDs returns organization's locations by organization's ID
func (a *authDataManager) OrganizationLocationIDs(_ context.Context, organizationID string) ([]string, error) {
	return a.storage.GetOrganizationLocationIDs(organizationID)
}

// AddOrganization creates new organization
func (a *authDataManager) AddOrganization(_ context.Context, organization *models.Organization) (*models.Organization, error) {
	return a.storage.AddOrganization(organization)
}

// UpdateOrganization updates organization
func (a *authDataManager) UpdateOrganization(_ context.Context, organization *models.Organization) (*models.Organization, error) {
	return a.storage.UpdateOrganization(organization)
}

// RemoveOrganization removes organization
func (a *authDataManager) RemoveOrganization(_ context.Context, organizationID string) error {
	return a.storage.RemoveOrganization(organizationID)
}

// Clinics returns all clinics
func (a *authDataManager) Clinics(_ context.Context) ([]*models.Clinic, error) {
	return a.storage.GetClinics()
}

// Clinic returns clinic by ID
func (a *authDataManager) Clinic(_ context.Context, clinicID string) (*models.Clinic, error) {
	return a.storage.GetClinic(clinicID)
}

// AddClinic creates new clinic
func (a *authDataManager) AddClinic(_ context.Context, clinic *models.Clinic) (*models.Clinic, error) {
	return a.storage.AddClinic(clinic)
}

// UpdateClinic updates clinic
func (a *authDataManager) UpdateClinic(_ context.Context, clinic *models.Clinic) (*models.Clinic, error) {
	return a.storage.UpdateClinic(clinic)
}

// RemoveClinic removes clinic
func (a *authDataManager) RemoveClinic(_ context.Context, clinicID string) error {
	return a.storage.RemoveClinic(clinicID)
}

// Locations returns all locations
func (a *authDataManager) Locations(_ context.Context) ([]*models.Location, error) {
	return a.storage.GetLocations()
}

// Location returns location by ID
func (a *authDataManager) Location(_ context.Context, locationID string) (*models.Location, error) {
	return a.storage.GetLocation(locationID)
}

// LocationOrganizationIDs returns IDs of location's organizations by location's ID
func (a *authDataManager) LocationOrganizationIDs(_ context.Context, locationID string) ([]string, error) {
	return a.storage.GetLocationOrganizationIDs(locationID)
}

// LocationUserIDs fetches list of IDs of users that have been assigned a role at the location (with optional role ID filtering); both users of clinics associated with the locations and users that have been assigned a role at the location manually are returned.
func (a *authDataManager) LocationUserIDs(_ context.Context, id string, roleID *string) ([]string, error) {
	// create users map to avoid duplicates
	userIDsMap := make(map[string]bool)

	// find user roles
	userRoles, err := a.storage.FindUserRoles(nil, roleID, &authCommon.DomainTypeLocation, &id)
	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := userIDsMap[*userRole.UserID]
		if !ok {
			userIDsMap[*userRole.UserID] = true
		}
	}

	// find all the clinics at the location
	clinics, err := a.storage.GetLocationClinics(id)
	if err != nil {
		return nil, err
	}
	for _, clinic := range clinics {
		// find user roles for each clinic
		userRoles, err = a.storage.FindUserRoles(nil, roleID, &authCommon.DomainTypeClinic, &clinic.ID)
		if err != nil {
			return nil, err
		}

		for _, userRole := range userRoles {
			// check if already in the map
			_, ok := userIDsMap[*userRole.UserID]
			if !ok {
				userIDsMap[*userRole.UserID] = true
			}
		}
	}

	// create response
	userIDs := []string{}
	for userID := range userIDsMap {
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// AddLocation creates new location
func (a *authDataManager) AddLocation(_ context.Context, location *models.Location) (*models.Location, error) {
	return a.storage.AddLocation(location)
}

// UpdateLocation updates location
func (a *authDataManager) UpdateLocation(_ context.Context, location *models.Location) (*models.Location, error) {
	return a.storage.UpdateLocation(location)
}

// RemoveLocation removes location
func (a *authDataManager) RemoveLocation(_ context.Context, locationID string) error {
	return a.storage.RemoveLocation(locationID)
}

// FindUserRoles returns user roles based on filtering query parameters.
func (a *authDataManager) FindUserRoles(_ context.Context, userID *string, roleID *string, domainType *string, domainID *string) ([]*models.UserRole, error) {
	return a.storage.FindUserRoles(userID, roleID, domainType, domainID)
}

// UserRole returns user role by its ID
func (a *authDataManager) UserRole(_ context.Context, id string) (*models.UserRole, error) {
	return a.storage.GetUserRole(id)
}

// AddRole creates a new user role
func (a *authDataManager) AddUserRole(_ context.Context, userRole *models.UserRole) (*models.UserRole, error) {
	return a.storage.AddUserRole(userRole)
}

// RemoveUserRole removes user role by its ID
func (a *authDataManager) RemoveUserRole(_ context.Context, id string) error {
	return a.storage.RemoveUserRole(id)
}

// DomainUserIDs fetches list of IDs of users that have been assigned a role at the domain (with optional role ID filtering).
func (a *authDataManager) DomainUserIDs(_ context.Context, domainType, domainID, roleID *string) ([]string, error) {
	// create users map to avoid duplicates
	userIDsMap := make(map[string]bool)

	// find user roles
	userRoles, err := a.storage.FindUserRoles(nil, roleID, domainType, domainID)
	if err != nil {
		return nil, err
	}

	for _, userRole := range userRoles {
		// check if already in the map
		_, ok := userIDsMap[*userRole.UserID]
		if !ok {
			userIDsMap[*userRole.UserID] = true
		}
	}

	// create response
	userIDs := []string{}
	for userID := range userIDsMap {
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// DBChecksum fetches checksum of underlying database
func (a *authDataManager) DBChecksum() ([]byte, error) {
	return a.storage.GetChecksum()
}

// WriteDBTo writes the whole underlying database to a writer
func (a *authDataManager) WriteDBTo(writer io.Writer) (int64, error) {
	return a.storage.WriteTo(writer)
}
