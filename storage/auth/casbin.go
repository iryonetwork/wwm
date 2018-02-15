package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/casbin/casbin"
	casbinmodel "github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/rs/zerolog"
)

// Persmissions
const (
	Read   = 1
	Write  = 1 << 1
	Delete = 1 << 2
)

// NewAdapter returns new Adapter
func NewAdapter(storage *Storage) *Adapter {
	logger := storage.logger.With().Str("component", "storage/auth/casbin").Logger()
	logger.Debug().Msg("Initialize casbin bolt adapter")

	return &Adapter{
		s:      storage,
		logger: logger,
	}
}

// Adapter is casbin adapter to bbolt
type Adapter struct {
	s      *Storage
	logger zerolog.Logger
}

// LoadPolicy loads policy from database
func (a *Adapter) LoadPolicy(model casbinmodel.Model) error {
	a.logger.Debug().Msg("Load policy from database")
	rules, err := a.s.GetRules()
	if err != nil {
		return err
	}
	roles, err := a.s.GetRoles()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		eft := "allow"
		if rule.Deny {
			eft = "deny"
		}
		persist.LoadPolicyLine(fmt.Sprintf("p, %s, %s, %d, %s", *rule.Subject, *rule.Resource, *rule.Action, eft), model)
	}

	for _, role := range roles {
		for _, user := range role.Users {
			persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s", user, role.ID), model)
		}
	}

	return nil
}

// SavePolicy saves policy to database
func (a *Adapter) SavePolicy(model casbinmodel.Model) error {
	return errors.New("not implemented")
}

// AddPolicy adds a policy rule to the storage
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}

func binaryMatch(val1, val2 int64) bool {
	return val1&val2 == val1
}

// BinaryMatchFunc does and bitwise match
func BinaryMatchFunc(args ...interface{}) (interface{}, error) {
	val1, err := strconv.ParseInt(args[0].(string), 10, 64)
	if err != nil {
		return nil, err
	}
	val2, err := strconv.ParseInt(args[1].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return (bool)(binaryMatch(val1, val2)), nil
}

// SelfMatch replaces {self} in key2 with subject and then compares it with key1
// this is useful when you want to apply policy for current subject
// e.g. allow user to edit its profile, but not profiles of other users
func SelfMatch(key1, key2, subject string) bool {
	key2 = strings.Replace(key2, "{self}", subject, -1)
	return key1 == key2
}

// SelfMatchFunc is the wrapper for SelfMatch.
func SelfMatchFunc(args ...interface{}) (interface{}, error) {
	key1 := args[0].(string)
	key2 := args[1].(string)
	subject := args[2].(string)

	return (bool)(SelfMatch(key1, key2, subject)), nil
}

// NewEnforcer returns new casbin enforcer
func NewEnforcer(storage *Storage) (*casbin.Enforcer, error) {
	m := casbin.NewModel(`[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && (keyMatch(r.obj, p.obj) || selfMatch(r.obj, p.obj, r.sub)) && binaryMatch(r.act, p.act)`)

	a := NewAdapter(storage)
	e := casbin.NewEnforcer(m, a, false)
	e.AddFunction("binaryMatch", BinaryMatchFunc)
	e.AddFunction("selfMatch", SelfMatchFunc)

	err := e.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return e, nil
}

// FindACL loads all the matching rules
func (s *Storage) FindACL(subject string, actions []*models.ValidationPair) []*models.ValidationResult {
	results := make([]*models.ValidationResult, len(actions), len(actions))

	for i, validation := range actions {
		results[i] = &models.ValidationResult{
			Query:  validation,
			Result: s.enforcer.Enforce(subject, *validation.Resource, strconv.FormatInt(*validation.Actions, 10)),
		}
	}

	return results
}
