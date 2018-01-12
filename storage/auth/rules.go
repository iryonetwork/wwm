package auth

import (
	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/utils"
)

// GetRules returns all rules
func (s *Storage) GetRules() ([]*models.Rule, error) {
	rules := []*models.Rule{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketACLRules)

		return b.ForEach(func(k, v []byte) error {
			rule := &models.Rule{}
			err := rule.UnmarshalBinary(v)
			if err != nil {
				return err
			}

			rules = append(rules, rule)
			return nil
		})
	})

	return rules, err
}

// GetRule returns rule by the id
func (s *Storage) GetRule(id string) (*models.Rule, error) {
	ruleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	rule := &models.Rule{}

	// look up the rule
	err = s.db.View(func(tx *bolt.Tx) error {
		// read rule by id
		data := tx.Bucket(bucketACLRules).Get(ruleUUID.Bytes())
		if data == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find rule by id = '%s'", id)
		}

		// decode the rule
		return rule.UnmarshalBinary(data)
	})

	return rule, err
}

// checkSubject checks if user or role exists in database
func (s *Storage) checkSubject(subject string) error {
	_, err := s.GetUser(subject)
	if err != nil {
		if e, ok := err.(utils.Error); !ok || e.Code() != utils.ErrNotFound {
			return err
		}

		_, err := s.GetRole(subject)
		if err != nil {
			if e, ok := err.(utils.Error); !ok || e.Code() != utils.ErrNotFound {
				return err
			}

			return utils.NewError(utils.ErrBadRequest, "Failed to find user or role '%s'", subject)
		}
	}
	return nil
}

// AddRule adds rule to the database
func (s *Storage) AddRule(rule *models.Rule) (*models.Rule, error) {
	err := s.checkSubject(*rule.Subject)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		// generatu uuid
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		rule.ID = id.String()

		data, err := rule.MarshalBinary()
		if err != nil {
			return err
		}

		// insert rule
		return tx.Bucket(bucketACLRules).Put(id.Bytes(), data)
	})

	return rule, err
}

// UpdateRule updates the rule
func (s *Storage) UpdateRule(rule *models.Rule) (*models.Rule, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// get buckets
		bRules := tx.Bucket(bucketACLRules)

		// check if rule exists
		ruleUUID, err := uuid.FromString(rule.ID)
		if err != nil {
			return utils.NewError(utils.ErrBadRequest, err.Error())
		}

		if bRules.Get(ruleUUID.Bytes()) == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find rule by id = '%s'", rule.ID)
		}

		err = s.checkSubject(*rule.Subject)
		if err != nil {
			return err
		}

		data, err := rule.MarshalBinary()
		if err != nil {
			return err
		}

		// update rule
		return bRules.Put(ruleUUID.Bytes(), data)
	})

	return rule, err
}

// RemoveRule removes rule by id
func (s *Storage) RemoveRule(id string) error {
	_, err := s.GetRule(id)
	if err != nil {
		return err
	}

	ruleUUID, _ := uuid.FromString(id)

	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketACLRules).Delete(ruleUUID.Bytes())
	})
}
