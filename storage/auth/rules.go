package auth

import (
	"strconv"

	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

// GetRules returns all rules
func (s *Storage) GetRules() ([]*models.Rule, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	rules := []*models.Rule{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		rules, err = s.getRulesWithTx(tx)
		return err
	})

	return rules, err
}

// getRulesWithTx gets rules from the database within passed bolt transaction
func (s *Storage) getRulesWithTx(tx *bolt.Tx) ([]*models.Rule, error) {
	rules := []*models.Rule{}

	b := tx.Bucket(bucketACLRules)

	err := b.ForEach(func(_, data []byte) error {
		rule := &models.Rule{}
		err := rule.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		rules = append(rules, rule)
		return nil
	})

	return rules, err
}

// GetRule returns rule by the id
func (s *Storage) GetRule(id string) (*models.Rule, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	rule := &models.Rule{}
	// look up the rule
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		rule, err = s.getRuleWithTx(tx, id)
		return err
	})

	return rule, err
}

// getRuleWithTx gets rule from the database within passed bolt transaction
func (s *Storage) getRuleWithTx(tx *bolt.Tx, id string) (*models.Rule, error) {
	ruleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	rule := &models.Rule{}

	// read rule by id
	data := tx.Bucket(bucketACLRules).Get(ruleUUID.Bytes())
	if data == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find rule by id = '%s'", id)
	}

	// decode the rule
	err = rule.UnmarshalBinary(data)

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

// AddRule generates new UUID,  adds rule to the database and updates related entities
func (s *Storage) AddRule(rule *models.Rule) (*models.Rule, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// generate ID
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	rule.ID = id.String()

	return s.addRule(rule)
}

func (s *Storage) addRule(rule *models.Rule) (*models.Rule, error) {
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

	var addedRule *models.Rule
	err = s.db.Update(func(tx *bolt.Tx) error {
		var err error
		// insert rule
		addedRule, err = s.insertRuleWithTx(tx, rule)
		return err
	})

	// refresh policy
	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return rule, err
}

// UpdateRule updates the rule and related entities in the database
func (s *Storage) UpdateRule(rule *models.Rule) (*models.Rule, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var updatedRule *models.Rule
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error
		// check if rule exists
		_, err = s.getRuleWithTx(tx, rule.ID)
		if err != nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find rule by id = '%s'", rule.ID)
		}

		err = s.checkSubject(*rule.Subject)
		if err != nil {
			return err
		}

		updatedRule, err = s.insertRuleWithTx(tx, rule)
		return err
	})

	// refresh policy
	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return rule, err
}

// insertRuleWithTx updates rule in the database within passed bolt transaction (does not update related entities)
func (s *Storage) insertRuleWithTx(tx *bolt.Tx, rule *models.Rule) (*models.Rule, error) {
	// get ID as UUID
	ruleUUID, err := uuid.FromString(rule.ID)
	if err != nil {
		return nil, err
	}

	data, err := rule.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// update rule
	err = tx.Bucket(bucketACLRules).Put(ruleUUID.Bytes(), data)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// RemoveRule removes rule by id from the database and updates related entities
func (s *Storage) RemoveRule(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := s.getRuleWithTx(tx, id)
		if err != nil {
			return err
		}

		return s.removeRuleWithTx(tx, id)
	})

	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return err
}

// removeRuleWithTx removes rule from the database within passed bolt transaction (does not update related entities)
func (s *Storage) removeRuleWithTx(tx *bolt.Tx, id string) error {
	ruleUUID, _ := uuid.FromString(id)

	return tx.Bucket(bucketACLRules).Delete(ruleUUID.Bytes())
}
