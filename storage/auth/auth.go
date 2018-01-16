package auth

import (
	"github.com/casbin/casbin"
	"github.com/go-openapi/swag"
	"github.com/iryonetwork/wwm/gen/models"
	uuid "github.com/satori/go.uuid"

	bolt "github.com/coreos/bbolt"
)

type Storage struct {
	db       *bolt.DB
	enforcer *casbin.Enforcer
}

var bucketUsers = []byte("users")
var bucketUsernames = []byte("usernames")
var bucketACLRules = []byte("rules")
var bucketRoles = []byte("roles")

var everyoneRole = &models.Role{
	ID:    "338fae76-9859-4803-8441-c5c441319cfd",
	Name:  swag.String("everyone"),
	Users: []string{},
}

var adminRole = &models.Role{
	ID:    "3720198b-74ed-40de-a45e-8756f22e67d2",
	Name:  swag.String("admin"),
	Users: []string{},
}

// New returns a new instance of storage
func New(path string) (*Storage, error) {
	db, err := bolt.Open(path, 0x600, nil)
	if err != nil {
		return nil, err
	}

	// initialize database
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
		_, err = tx.CreateBucketIfNotExists(bucketACLRules)
		return err
	})
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db: db,
	}

	e, err := NewEnforcer(storage)
	if err != nil {
		return nil, err
	}

	storage.enforcer = e
	storage.initializeRolesAndRules()

	return storage, nil
}

func (s *Storage) initializeRolesAndRules() error {
	_, err := s.GetRole(everyoneRole.ID)
	if err != nil {
		err := s.db.Update(func(tx *bolt.Tx) error {
			roleUUID, _ := uuid.FromString(everyoneRole.ID)
			data, _ := everyoneRole.MarshalBinary()

			return tx.Bucket(bucketRoles).Put(roleUUID.Bytes(), data)
		})
		if err != nil {
			return err
		}
	}

	_, err = s.GetRole(adminRole.ID)
	if err != nil {
		err := s.db.Update(func(tx *bolt.Tx) error {
			roleUUID, _ := uuid.FromString(adminRole.ID)
			data, _ := adminRole.MarshalBinary()

			return tx.Bucket(bucketRoles).Put(roleUUID.Bytes(), data)
		})
		if err != nil {
			return err
		}
	}

	s.AddRule(&models.Rule{
		Subject:  &everyoneRole.ID,
		Action:   swag.Int64(Write),
		Resource: swag.String("/auth/login"),
	})

	s.AddRule(&models.Rule{
		Subject:  &everyoneRole.ID,
		Action:   swag.Int64(Write),
		Resource: swag.String("/auth/validate"),
	})

	s.AddRule(&models.Rule{
		Subject:  &everyoneRole.ID,
		Action:   swag.Int64(Read),
		Resource: swag.String("/auth/renew"),
	})

	s.AddRule(&models.Rule{
		Subject:  &adminRole.ID,
		Action:   swag.Int64(Read | Write | Delete),
		Resource: swag.String("*"),
	})

	return nil
}

// Close closes the database
func (s *Storage) Close() error {
	return s.db.Close()
}
