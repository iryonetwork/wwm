package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/log"
	"github.com/iryonetwork/wwm/metrics/api"
)

// ServePrometheusMetrics starts prometheus metrics server
func ServePrometheusMetrics(ctx context.Context, addr string, namespace string, logger zerolog.Logger) error {
	logger = logger.With().Str("component", "metrics/server").Logger()

	// initialize metrics middleware
	m := api.NewMetrics("metrics", "")

	path := "/metrics"
	if namespace != "" {
		path = fmt.Sprintf("/%s/metrics", namespace)
	}

	mux := http.NewServeMux()
	mux.Handle(path, prometheus.Handler())
	s := &http.Server{
		Addr:    addr,
		Handler: m.Middleware(log.APILogMiddleware(mux, logger)),
	}

	go func() {
		<-ctx.Done()
		s.Shutdown(ctx)
	}()

	logger.Info().Msgf("Starting metrics server at %s%s", addr, path)

	return s.ListenAndServe()
}
