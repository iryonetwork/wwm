package auth

import (
	"strconv"

	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/utils"
)

// GetRules returns all rules
func (s *Storage) GetRules() ([]*models.Rule, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	rules := []*models.Rule{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketACLRules)

		return b.ForEach(func(_, data []byte) error {
			data, err := s.decrypt(data)
			if err != nil {
				return err
			}

			rule := &models.Rule{}
			err = rule.UnmarshalBinary(data)
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
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

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

		data, err = s.decrypt(data)
		if err != nil {
			return err
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
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	err := s.checkSubject(*rule.Subject)
	if err != nil {
		return nil, err
	}

	eft := "allow"
	if rule.Deny {
		eft = "deny"
	}
	if s.enforcer.HasPolicy(*rule.Subject, *rule.Resource, strconv.FormatInt(*rule.Action, 10), eft) {
		return nil, utils.NewError(utils.ErrBadRequest, "Rule with that parameters already exist")
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

		data, err = s.encrypt(data)
		if err != nil {
			return err
		}

		// insert rule
		return tx.Bucket(bucketACLRules).Put(id.Bytes(), data)
	})

	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return rule, err
}

// UpdateRule updates the rule
func (s *Storage) UpdateRule(rule *models.Rule) (*models.Rule, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

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

		data, err = s.encrypt(data)
		if err != nil {
			return err
		}

		// update rule
		return bRules.Put(ruleUUID.Bytes(), data)
	})

	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return rule, err
}

// RemoveRule removes rule by id
func (s *Storage) RemoveRule(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	_, err := s.GetRule(id)
	if err != nil {
		return err
	}

	ruleUUID, _ := uuid.FromString(id)

	err = s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketACLRules).Delete(ruleUUID.Bytes())
	})

	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return err
}
