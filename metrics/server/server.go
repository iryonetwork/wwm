package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/iryonetwork/wwm/metrics/api"
	"github.com/prometheus/client_golang/prometheus"
)

// ServePrometheusMetrics starts prometheus metrics server
func ServePrometheusMetrics(ctx context.Context, addr string, namespace string) error {
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
		Handler: m.Middleware(mux),
	}

	go func() {
		<-ctx.Done()
		s.Shutdown(ctx)
	}()

	log.Printf("Starting metrics server at %s%s", addr, path)

	return s.ListenAndServe()
}
