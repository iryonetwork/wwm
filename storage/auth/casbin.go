package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/casbin/casbin"
	casbinmodel "github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/gobwas/glob"
	"github.com/rs/zerolog"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
)

// Persmissions
const (
	Read   = 1
	Update = 1 << 1
	Write  = 1 << 2
	Delete = 1 << 3
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

	userRoles, err := a.s.GetUserRoles()
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

	for _, userRole := range userRoles {
		switch *userRole.DomainType {
		case authCommon.DomainTypeOrganization:
			// if it's wildcard role for organization domain type, iterate through all organization and load role for all of them
			if *userRole.DomainID == authCommon.DomainIDWildcard {
				organizations, err := a.s.GetOrganizations()
				if err != nil {
					return err
				}
				for _, organization := range organizations {
					persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeOrganization, organization.ID)), model)
				}
			} else {
				persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeOrganization, *userRole.DomainID)), model)
			}
		case authCommon.DomainTypeClinic:
			// if it's wildcard role for clinic domain type, iterate through all clinics and load role for all of them and for corresponding location
			if *userRole.DomainID == authCommon.DomainIDWildcard {
				clinics, err := a.s.GetClinics()
				if err != nil {
					return err
				}
				for _, clinic := range clinics {
					persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeClinic, clinic.ID)), model)
					// for clinic all user roles apply also for clinic's location
					persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeLocation, *clinic.Location)), model)
				}
			} else {
				clinic, err := a.s.GetClinic(*userRole.DomainID)
				if err != nil {
					return err
				}
				persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeClinic, clinic.ID)), model)
				// for clinic all user roles apply also for clinic's location
				persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeLocation, *clinic.Location)), model)
			}
		case authCommon.DomainTypeLocation:
			// if it's wildcard role for location domain type, iterate through all locations and load role for all of them
			if *userRole.DomainID == authCommon.DomainIDWildcard {
				locations, err := a.s.GetLocations()
				if err != nil {
					return err
				}
				for _, location := range locations {
					persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeLocation, location.ID)), model)
				}
			} else {
				persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeOrganization, *userRole.DomainID)), model)
			}
		case authCommon.DomainTypeUser:
			// if it's wildcard role for user domain type, iterate through all users and load role for all of them
			if *userRole.DomainID == authCommon.DomainIDWildcard {
				users, err := a.s.GetUsers()
				if err != nil {
					return err
				}
				for _, user := range users {
					persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeUser, user.ID)), model)
				}
			} else {
				persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", authCommon.DomainTypeUser, *userRole.DomainID)), model)
			}
		case authCommon.DomainTypeGlobal:
			persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, "*"), model)
		default:
			persist.LoadPolicyLine(fmt.Sprintf("g, %s, %s, %s", *userRole.UserID, *userRole.RoleID, fmt.Sprintf("%s.%s", *userRole.DomainType, *userRole.DomainID)), model)
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

type wildcardMatch struct {
	globs sync.Map
}

func (w *wildcardMatch) Match(args ...interface{}) (interface{}, error) {
	request := args[0].(string)
	policy := args[1].(string)

	var rule glob.Glob
	g, ok := w.globs.Load(policy)
	if !ok {
		newGlob, err := glob.Compile(policy)
		if err != nil {
			return nil, err
		}

		w.globs.Store(policy, newGlob)
		rule = newGlob
	} else {
		rule = g.(glob.Glob)
	}

	return rule.Match(request), nil
}

// NewEnforcer returns new casbin enforcer
func NewEnforcer(storage *Storage) (*casbin.Enforcer, error) {
	m := casbin.NewModel(`[request_definition]
r = sub, dom, obj, act
[dom actual location]

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = (g(r.sub, p.sub, r.dom) ||  g(r.sub, p.sub, "*")) && (wildcardMatch(r.obj, p.obj) || selfMatch(r.obj, p.obj, r.sub)) && binaryMatch(r.act, p.act)`)

	a := NewAdapter(storage)
	e := casbin.NewEnforcer(m, a, false)
	e.AddFunction("binaryMatch", BinaryMatchFunc)
	e.AddFunction("selfMatch", SelfMatchFunc)

	w := &wildcardMatch{}
	e.AddFunction("wildcardMatch", w.Match)

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
		var domain string
		if *validation.DomainType == authCommon.DomainTypeGlobal {
			domain = "*"
		} else {
			domain = fmt.Sprintf("%s.%s", *validation.DomainType, *validation.DomainID)
		}
		results[i] = &models.ValidationResult{
			Query:  validation,
			Result: s.enforcer.Enforce(subject, domain, *validation.Resource, strconv.FormatInt(*validation.Actions, 10)),
		}
	}

	return results
}
