package auth

import (
	"fmt"
	"os"
	"sync"

	"github.com/casbin/casbin"
	"github.com/go-openapi/swag"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
)

type Storage struct {
	db            *bolt.DB
	enforcer      *casbin.Enforcer
	encryptionKey []byte
	dbSync        *sync.RWMutex
	refreshRules  bool
	logger        zerolog.Logger
}

type InitData struct {
	Locations     []*models.Location
	Organizations []*models.Organization
	Clinics       []*models.Clinic
	Roles         []*models.Role
	Rules         []*models.Rule
	Users         []*models.User
	UserRoles     []*models.UserRole
}

var bucketUsers = []byte("users")
var bucketUsernames = []byte("usernames")
var bucketACLRules = []byte("rules")
var bucketRoles = []byte("roles")
var bucketLocations = []byte("locations")
var bucketLocationNames = []byte("locationNames")
var bucketOrganizations = []byte("organizations")
var bucketOrganizationNames = []byte("organizationNames")
var bucketClinics = []byte("clinics")
var bucketClinicNames = []byte("clinicNames")
var bucketUserRoles = []byte("userRoles")
var bucketUserIDUserRolesIndex = []byte("userIDUserRolesIndex")
var bucketRoleIDUserRolesIndex = []byte("roleIDUserRolesIndex")
var bucketDomainUserRolesIndex = []byte("domainUserRolesIndex")

var dbPermissions os.FileMode = 0666

// New returns a new instance of storage
func New(path string, key []byte, readOnly, refreshRules bool, logger zerolog.Logger) (*Storage, error) {
	logger = logger.With().Str("component", "storage/auth").Logger()
	logger.Debug().Msg("Initialize auth storage")
	if len(key) != 32 {
		return nil, fmt.Errorf("Encryption key must be 32 bytes long")
	}

	db, err := bolt.Open(key, path, dbPermissions, &bolt.Options{ReadOnly: readOnly})
	if err != nil {
		return nil, err
	}

	// initialize database
	if !readOnly {
		logger.Debug().Msg("Create db buckets")
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucketUsers)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketUsernames)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketRoles)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketUserRoles)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketUserIDUserRolesIndex)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketRoleIDUserRolesIndex)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketDomainUserRolesIndex)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketLocations)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketLocationNames)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketOrganizations)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketOrganizationNames)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketClinics)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketClinicNames)
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists(bucketACLRules)
			return err

		})
		if err != nil {
			return nil, err
		}
	}

	storage := &Storage{
		db:            db,
		encryptionKey: key,
		dbSync:        &sync.RWMutex{},
		refreshRules:  refreshRules,
		logger:        logger,
	}

	e, err := NewEnforcer(storage)
	if err != nil {
		return nil, err
	}

	storage.enforcer = e
	if readOnly {
		e.LoadPolicy()
	} else {
		storage.initializeRolesAndRules()
	}

	return storage, nil
}

func (s *Storage) initializeRolesAndRules() error {
	s.logger.Debug().Msg("Initialize roles and rules")
	_, err := s.GetRole(authCommon.EveryoneRole.ID)
	if err != nil {
		err := s.db.Update(func(tx *bolt.Tx) error {
			roleUUID, _ := uuid.FromString(authCommon.EveryoneRole.ID)
			data, _ := authCommon.EveryoneRole.MarshalBinary()

			return tx.Bucket(bucketRoles).Put(roleUUID.Bytes(), data)
		})
		if err != nil {
			return err
		}
	}

	_, err = s.GetRole(authCommon.AuthorRole.ID)
	if err != nil {
		err := s.db.Update(func(tx *bolt.Tx) error {
			roleUUID, _ := uuid.FromString(authCommon.AuthorRole.ID)
			data, _ := authCommon.AuthorRole.MarshalBinary()

			return tx.Bucket(bucketRoles).Put(roleUUID.Bytes(), data)
		})
		if err != nil {
			return err
		}
	}

	_, err = s.GetRole(authCommon.MemberRole.ID)
	if err != nil {
		err := s.db.Update(func(tx *bolt.Tx) error {
			roleUUID, _ := uuid.FromString(authCommon.MemberRole.ID)
			data, _ := authCommon.MemberRole.MarshalBinary()

			return tx.Bucket(bucketRoles).Put(roleUUID.Bytes(), data)
		})
		if err != nil {
			return err
		}
	}

	_, err = s.GetRole(authCommon.SuperadminRole.ID)
	if err != nil {
		err := s.db.Update(func(tx *bolt.Tx) error {
			roleUUID, _ := uuid.FromString(authCommon.SuperadminRole.ID)
			data, _ := authCommon.SuperadminRole.MarshalBinary()

			return tx.Bucket(bucketRoles).Put(roleUUID.Bytes(), data)
		})
		if err != nil {
			return err
		}
	}

	s.AddRule(&models.Rule{
		Subject:  &authCommon.EveryoneRole.ID,
		Action:   swag.Int64(Write),
		Resource: swag.String("/auth/login"),
	})

	s.AddRule(&models.Rule{
		Subject:  &authCommon.EveryoneRole.ID,
		Action:   swag.Int64(Write),
		Resource: swag.String("/api/auth/validate"),
	})

	s.AddRule(&models.Rule{
		Subject:  &authCommon.EveryoneRole.ID,
		Action:   swag.Int64(Read),
		Resource: swag.String("/api/auth/*"),
	})

	s.AddRule(&models.Rule{
		Subject:  &authCommon.EveryoneRole.ID,
		Action:   swag.Int64(Read | Write),
		Resource: swag.String("/api/auth/users/{self}*"),
	})

	s.AddRule(&models.Rule{
		Subject:  &authCommon.EveryoneRole.ID,
		Action:   swag.Int64(Read | Write),
		Resource: swag.String("/api/auth/users/me*"),
	})

	s.AddRule(&models.Rule{
		Subject:  &authCommon.SuperadminRole.ID,
		Action:   swag.Int64(Read | Write | Update | Delete),
		Resource: swag.String("*"),
	})

	s.AddRule(&models.Rule{
		Subject:  &authCommon.SuperadminRole.ID,
		Action:   swag.Int64(Read),
		Resource: swag.String("/frontend/dashboard*"),
	})

	s.AddRule(&models.Rule{
		Subject:  &authCommon.EveryoneRole.ID,
		Action:   swag.Int64(Read),
		Resource: swag.String("/frontend/dashboard/{self}"),
	})

	return nil
}

// Close closes the database
func (s *Storage) Close() error {
	return s.db.Close()
}

// LoadInitData inserts into database initial data, errors are generally ignored and only loggeda as info
func (s *Storage) LoadInitData(data InitData) {
	for _, location := range data.Locations {
		if location.ID == "" {
			_, err := s.AddLocation(location)
			if err != nil {
				s.logger.Info().Err(err).Msg("Location from init data could not be added")
			}
		} else {
			// check if already exists
			_, err := s.GetLocation(location.ID)
			if err == nil {
				s.logger.Info().Msg("Location from init data could not be added as location with that UUID already exists")
			} else {
				_, err := s.addLocation(location)
				if err != nil {
					s.logger.Info().Err(err).Msg("Location from init data could not be added")
				}
			}
		}
	}

	for _, organization := range data.Organizations {
		if organization.ID == "" {
			_, err := s.AddOrganization(organization)
			if err != nil {
				s.logger.Info().Err(err).Msg("Organization from init data could not be added")
			}
		} else {
			// check if already exists
			_, err := s.GetOrganization(organization.ID)
			if err == nil {
				s.logger.Info().Msg("Organization from init data could not be added as organization with that UUID already exists")
			} else {
				_, err := s.addOrganization(organization)
				if err != nil {
					s.logger.Info().Err(err).Msg("Organization from init data could not be added")
				}
			}
		}
	}

	for _, clinic := range data.Clinics {
		if clinic.ID == "" {
			_, err := s.AddClinic(clinic)
			if err != nil {
				s.logger.Info().Err(err).Msg("Clinic from init data could not be added")
			}
		} else {
			// check if already exists
			_, err := s.GetClinic(clinic.ID)
			if err == nil {
				s.logger.Info().Msg("Clinic from init data could not be added as clinic with that UUID already exists")
			} else {
				_, err := s.addClinic(clinic)
				if err != nil {
					s.logger.Info().Err(err).Msg("Clinic from init data could not be added")
				}
			}
		}
	}

	for _, role := range data.Roles {
		if role.ID == "" {
			_, err := s.AddRole(role)
			if err != nil {
				s.logger.Info().Err(err).Msg("Role from init data could not be added")
			}
		} else {
			// check if already exists
			_, err := s.GetRole(role.ID)
			if err == nil {
				s.logger.Info().Msg("Role from init data could not be added as role with that UUID already exists")
			} else {
				_, err := s.addRole(role)
				if err != nil {
					s.logger.Info().Err(err).Msg("Role from init data could not be added")
				}
			}
		}
	}

	for _, rule := range data.Rules {
		if rule.ID == "" {
			_, err := s.AddRule(rule)
			if err != nil {
				s.logger.Info().Err(err).Msg("Rule from init data could not be added")
			}
		} else {
			// check if already exists
			_, err := s.GetRule(rule.ID)
			if err == nil {
				s.logger.Info().Msg("Rule from init data could not be added as rule with that UUID already exists")
			} else {
				_, err := s.addRule(rule)
				if err != nil {
					s.logger.Info().Err(err).Msg("Rule from init data could not be added")
				}
			}
		}
	}

	for _, user := range data.Users {
		if user.ID == "" {
			_, err := s.AddUser(user)
			if err != nil {
				s.logger.Info().Err(err).Msg("User from init data could not be added")
			}
		} else {
			// check if already exists
			_, err := s.GetUser(user.ID)
			if err == nil {
				s.logger.Info().Msg("User from init data could not be added as user with that UUID already exists")
			} else {
				_, err := s.addUser(user)
				if err != nil {
					s.logger.Info().Err(err).Msg("User from init data could not be added")
				}
			}
		}
	}

	for _, userRole := range data.UserRoles {
		if userRole.ID == "" {
			_, err := s.AddUserRole(userRole)
			if err != nil {
				s.logger.Info().Err(err).Msg("User role from init data could not be added")
			}
		} else {
			// check if already exists
			_, err := s.GetUserRole(userRole.ID)
			if err == nil {
				s.logger.Info().Msg("User role from init data could not be added as user with that UUID already exists")
			} else {
				_, err := s.addUserRole(userRole)
				if err != nil {
					s.logger.Info().Err(err).Msg("User role from init data could not be added")
				}
			}
		}
	}
}
