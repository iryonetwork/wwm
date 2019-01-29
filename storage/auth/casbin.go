package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/casbin/casbin"
	casbinmodel "github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/gobwas/glob"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/metrics"
)

// Persmissions
const (
	Read   = 1
	Write  = 1 << 1
	Delete = 1 << 2
	Update = 1 << 3
)

const loadPolicySeconds metrics.ID = "load_policy_seconds"
const enforceSeconds metrics.ID = "enforce_seconds"

// NewAdapter returns new Adapter
func NewAdapter(storage *Storage, logger zerolog.Logger) *Adapter {
	logger = logger.With().Str("component", "storage/auth/casbinAdapter").Logger()
	logger.Debug().Msg("Initialize casbin bolt adapter")

	metricsCollection := make(map[metrics.ID]prometheus.Collector)
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "casbin_adapter",
		Name:      "load_policy_seconds",
		Help:      "Time taken to load casbin policy",
	})
	metricsCollection[loadPolicySeconds] = h

	return &Adapter{
		s:                 storage,
		logger:            logger,
		metricsCollection: metricsCollection,
	}
}

// Adapter is casbin adapter to bbolt
type Adapter struct {
	s                 *Storage
	logger            zerolog.Logger
	metricsCollection map[metrics.ID]prometheus.Collector
}

// LoadPolicy loads policy from database
func (a *Adapter) LoadPolicy(model casbinmodel.Model) error {
	// Make sure we record duration metrics even if processing fails
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		a.metricsCollection[loadPolicySeconds].(prometheus.Histogram).Observe(duration.Seconds())
	}()

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

// GetPrometheusMetricsCollection returns all prometheus metrics collectors to be registered
func (a *Adapter) GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector {
	return a.metricsCollection
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

// SelfReplace replaces {self} in key2 with subject and then returns it for comparison.
// this is useful when you want to apply policy for current subject
// e.g. allow user to edit its profile, but not profiles of other users
func SelfReplace(key, subject string) string {
	key = strings.Replace(key, "{self}", subject, -1)
	return key
}

// SelfReplaceFunc is the wrapper for SelfReplace.
func SelfReplaceFunc(args ...interface{}) (interface{}, error) {
	key := args[0].(string)
	subject := args[1].(string)

	return (string)(SelfReplace(key, subject)), nil
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

type enforcer struct {
	*casbin.Enforcer
	adapter           *Adapter
	metricsCollection map[metrics.ID]prometheus.Collector
}

// NewEnforcer returns new casbin enforcer
func NewEnforcer(storage *Storage, logger zerolog.Logger) (Enforcer, error) {
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
m = (g(r.sub, p.sub, r.dom) ||  g(r.sub, p.sub, "*")) && (wildcardMatch(r.obj, p.obj) || wildcardMatch(r.obj, selfReplace(p.obj, r.sub))) && binaryMatch(r.act, p.act)`)

	a := NewAdapter(storage, logger)
	e := casbin.NewEnforcer(m, a, false)
	e.AddFunction("binaryMatch", BinaryMatchFunc)
	e.AddFunction("selfReplace", SelfReplaceFunc)

	w := &wildcardMatch{}
	e.AddFunction("wildcardMatch", w.Match)

	err := e.LoadPolicy()
	if err != nil {
		return nil, err
	}

	metricsCollection := make(map[metrics.ID]prometheus.Collector)
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "casbin_enforcer",
		Name:      "enforce_seconds",
		Help:      "Time taken to enforce policy",
	})
	metricsCollection[enforceSeconds] = h

	return &enforcer{e, a, metricsCollection}, nil
}

// Enforce is a wrapper for casbin Enforce method made to measure execution time
func (e *enforcer) Enforce(rvals ...interface{}) bool {
	// Make sure we record duration metrics even if processing fails
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.metricsCollection[enforceSeconds].(prometheus.Histogram).Observe(duration.Seconds())
	}()

	return (e.Enforcer).Enforce(rvals...)
}

// GetPrometheusMetricsCollection returns all prometheus metrics collectors to be registered
func (e *enforcer) GetPrometheusMetricsCollection() map[metrics.ID]prometheus.Collector {
	collection := e.metricsCollection
	for k, v := range e.adapter.GetPrometheusMetricsCollection() {
		collection[k] = v
	}

	return collection
}
