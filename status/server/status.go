package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/metrics/api"
	"github.com/iryonetwork/wwm/status"
)

type statusServer struct {
	components map[string]status.Component
	server     *http.Server
	logger     zerolog.Logger
}

// AddService adds component to be checked on status call
func (s *statusServer) AddComponent(name string, component status.Component) {
	s.components[name] = component
}

// Status returns status response
func (s *statusServer) Status() *status.Response {
	currentStatus := status.OK
	componentsResp := make(map[string]*status.Response)
	for id, c := range s.components {
		cResp := c.Status()

		// higher the StatusValue the worse the status, combined response will return worst status of all components
		if cResp != nil && cResp.Status.Int() > currentStatus.Int() {
			currentStatus = cResp.Status
		}
		componentsResp[id] = cResp
	}

	return &status.Response{
		Status:     currentStatus,
		Components: componentsResp,
	}
}

// ServeHTTP serves status API response
func (s *statusServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	errorChecker.LogError(json.NewEncoder(rw).Encode(s.Status()))
}

// ListenAndServeHTTP starts serving status endpoint
func (s *statusServer) ListenAndServeHTTPs(ctx context.Context, addr string, namespace string, certFile, keyFile string) error {
	if s.server != nil {
		s.server.Close()
	}

	// initialize metrics middleware
	m := api.NewMetrics("status", "")

	path := "/status"
	if namespace != "" {
		path = fmt.Sprintf("/%s/status", namespace)
	}

	mux := http.NewServeMux()
	mux.Handle(path, s)
	s.server = &http.Server{
		Addr:    addr,
		Handler: m.Middleware(mux),
	}

	go func() {
		<-ctx.Done()
		s.Close()
	}()

	s.logger.Info().Msgf("Starting status server at %s%s", addr, path)

	return s.server.ListenAndServeTLS(certFile, keyFile)
}

func (s *statusServer) Close() error {
	return s.server.Close()
}

func New(logger zerolog.Logger) *statusServer {
	logger = logger.With().Str("component", "status/server").Logger()

	components := make(map[string]status.Component)

	return &statusServer{logger: logger, components: components}
}
