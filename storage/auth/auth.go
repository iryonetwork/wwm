package auth

import (
	"github.com/iryonetwork/wwm/specs"

	bolt "github.com/coreos/bbolt"
)

type Storage struct {
	db *bolt.DB
}

var bucketUsers = []byte("users")
var bucketUsernames = []byte("usernames")
var bucketACLRules = []byte("rules")
var bucketRoles = []byte("groups")

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
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

// Close closes the database
func (s *Storage) Close() error {
	return s.db.Close()
}

// FindACL loads all the matching rules
func (s *Storage) FindACL(userID, resource string, actions []specs.ACLRuleAction) ([]*specs.ACLRule, error) {
	return nil, nil
}
