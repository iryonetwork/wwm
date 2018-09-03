package statusReporter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/status"
)

// type for environment group name
type Environment string

// Environment values
const (
	Local    Environment = "local"
	Cloud    Environment = "cloud"
	External Environment = "external"
)

// Respone defines status reporter response
type Response struct {
	Status   status.Value     `json:"status"`
	Local    *status.Response `json:"local,omitempty"`
	Cloud    *status.Response `json:"cloud,omitempty"`
	External *status.Response `json:"external,omitempty"`
}

type StatusReporter struct {
	cloud    map[string]status.Component
	local    map[string]status.Component
	external map[string]status.Component
	logger   zerolog.Logger
}

// Status returns status response
func (r *StatusReporter) Status() *Response {
	st := status.OK

	localResp := r.EnvironmentStatus(Local)()
	if localResp != nil && localResp.Status.Int() > st.Int() {
		st = localResp.Status
	}

	// cloud and external error trigger warning
	cloudResp := r.EnvironmentStatus(Cloud)()
	if cloudResp != nil && cloudResp.Status.Int() > st.Int() {
		st = status.Warning
	}
	externalResp := r.EnvironmentStatus(External)()
	if externalResp != nil && externalResp.Status.Int() > st.Int() {
		st = status.Warning
	}

	return &Response{
		Status:   st,
		Local:    localResp,
		Cloud:    cloudResp,
		External: externalResp,
	}
}

// EnvironmentStatus returns Status function for environment
func (r *StatusReporter) EnvironmentStatus(env Environment) func() *status.Response {
	var components map[string]status.Component
	switch env {
	case Cloud:
		components = r.cloud
	case Local:
		components = r.local
	case External:
		components = r.external
	}

	return func() *status.Response {
		if len(components) == 0 {
			return nil
		}

		st := status.OK
		componentsResp := make(map[string]*status.Response)

		for id, c := range components {
			resp := c.Status()
			if resp != nil {
				if resp.Status.Int() > st.Int() {
					st = resp.Status
				}
				componentsResp[id] = resp
			}
		}

		return &status.Response{
			Status:     st,
			Components: componentsResp,
		}
	}
}

// AddURLComponent adds status component to the reporter
func (r *StatusReporter) AddComponent(env Environment, id string, c status.Component) {
	switch env {
	case Local:
		r.local[id] = c
	case Cloud:
		r.cloud[id] = c
	case External:
		r.external[id] = c
	}
}

// Handler returns http.Handler to facilitate HTTP server of status reporter with JSON responses
func (r *StatusReporter) Handler(prefix string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/%s", prefix), handlerFunc(func() interface{} { return r.Status() }))
	mux.HandleFunc(fmt.Sprintf("/%s/", prefix), handlerFunc(func() interface{} { return r.Status() }))
	mux.HandleFunc(fmt.Sprintf("/%s/%s", prefix, Local), handlerFunc(func() interface{} { return r.EnvironmentStatus(Local) }))
	mux.HandleFunc(fmt.Sprintf("/%s/%s/", prefix, Local), handlerFunc(func() interface{} { return r.EnvironmentStatus(Local) }))
	mux.HandleFunc(fmt.Sprintf("/%s/%s", prefix, Cloud), handlerFunc(func() interface{} { return r.EnvironmentStatus(Cloud) }))
	mux.HandleFunc(fmt.Sprintf("/%s/%s/", prefix, Cloud), handlerFunc(func() interface{} { return r.EnvironmentStatus(Cloud) }))
	mux.HandleFunc(fmt.Sprintf("/%s/%s", prefix, External), handlerFunc(func() interface{} { return r.EnvironmentStatus(External) }))
	mux.HandleFunc(fmt.Sprintf("/%s/%s/", prefix, External), handlerFunc(func() interface{} { return r.EnvironmentStatus(External) }))

	return mux
}

// handlerFunc is a generic handler to faciliate StatusReporter HTTP serving with JSON responses
func handlerFunc(f func() interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		status := f()
		if status == nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorChecker.LogError(json.NewEncoder(rw).Encode(f()))
	}
}

// NewStatusReporter returns instance of status reporter service
func New(logger zerolog.Logger) *StatusReporter {
	logger = logger.With().Str("component", "service/statusReporter").Logger()

	return &StatusReporter{
		local:    make(map[string]status.Component),
		cloud:    make(map[string]status.Component),
		external: make(map[string]status.Component),
		logger:   logger,
	}
}
