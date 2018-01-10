package auth

import (
	"fmt"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/specs"
	"golang.org/x/crypto/bcrypt"

	bolt "github.com/coreos/bbolt"
)

type Storage struct {
	db *bolt.DB
}

var bucketUsers = []byte("users")
var bucketUsernames = []byte("usernames")
var bucketACLRules = []byte("rules")

// New returns a new instance of storage
func New() (*Storage, error) {
	db, err := bolt.Open("/tmp/auth.db", 0x600, nil)
	if err != nil {
		return nil, err
	}

	// add initial user
	err = db.Update(func(tx *bolt.Tx) error {
		// get buckets
		bUsers, err := tx.CreateBucketIfNotExists(bucketUsers)
		if err != nil {
			return err
		}
		bUsernames, err := tx.CreateBucketIfNotExists(bucketUsernames)
		if err != nil {
			return err
		}

		// get sample user
		id, username, data, err := getSampleUser()
		if err != nil {
			return err
		}

		// insert user
		err = bUsers.Put(id, data)
		if err != nil {
			return err
		}

		// insert username
		return bUsernames.Put(username, id)
	})
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

// FindACL loads all the matching rules
func (s *Storage) FindACL(userID, resource string, actions []specs.ACLRuleAction) ([]*specs.ACLRule, error) {
	return nil, nil
}

// GetUserByUsername returns user by the username
func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	// look up the user
	err := s.db.View(func(tx *bolt.Tx) error {
		// find the user in the usernames bucket
		id := tx.Bucket(bucketUsernames).Get([]byte(username))
		if id == nil {
			return fmt.Errorf("Failed to find username %s", username)
		}

		// read user by id
		data := tx.Bucket(bucketUsers).Get(id)
		if data == nil {
			return fmt.Errorf("Failed to find user by username %s (id = %s)", username, id)
		}

		// decode the user
		err := user.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}

// GetUserByID returns user by the ID
func (s *Storage) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}

	// look up the user
	err := s.db.View(func(tx *bolt.Tx) error {
		// read user by id
		data := tx.Bucket(bucketUsers).Get([]byte(id))
		if data == nil {
			return fmt.Errorf("Failed to find user by id = %s", id)
		}

		// decode the user
		err := user.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}

func getSampleUser() ([]byte, []byte, []byte, error) {
	password, err := bcrypt.GenerateFromPassword([]byte("password"), 0)
	if err != nil {
		return nil, nil, nil, err
	}

	user := &models.User{
		ID:       "SOME-ID",
		Username: "username",
		Password: string(password),
		Email:    "info@iryo.io",
	}

	// encode the object
	data, err := user.MarshalBinary()

	return []byte(user.ID), []byte(user.Username), data, err
}
