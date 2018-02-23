package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/iryonetwork/wwm/metrics"
	"github.com/iryonetwork/wwm/utils"
)

const requestSeconds metrics.ID = "requestSeconds"

// Metrics describes public methods of metrics middleware
type Metrics interface {
	Middleware(next http.Handler) http.Handler
	WithURLSanitize(sanitize utils.URLSanitize) Metrics
}

type apiMetrics struct {
	metricsCollection map[metrics.ID]prometheus.Collector
	urlSanitize       utils.URLSanitize
}

type codeRecordingResponseWriter struct {
	http.ResponseWriter
	code int
}

// Middleware wraps http.Handler and adds http request serving metrics
func (m *apiMetrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := strings.ToLower(r.Method)
		code := 500 // status code to use when completely fails

		path := r.URL.Path
		if m.urlSanitize != nil {
			path = m.urlSanitize(path)
		}

		// Make sure we record even if fails
		defer func() {
			duration := time.Since(start)
			m.metricsCollection[requestSeconds].(*prometheus.HistogramVec).
				With(prometheus.Labels{"path": path, "method": method, "code": fmt.Sprintf("%d", code)}).
				Observe(duration.Seconds())
		}()

		cw := &codeRecordingResponseWriter{w, 200}

		next.ServeHTTP(cw, r)

		code = cw.code
	})
}

// NewMetrics returns new metrics middleware instance
func NewMetrics(namespace string, subsystem string) Metrics {
	metricsCollection := make(map[metrics.ID]prometheus.Collector)
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_seconds",
		Help:      "Time taken to serve API request",
	}, []string{"path", "method", "code"})

	// Register metrics
	prometheus.MustRegister(h)
	metricsCollection[requestSeconds] = h

	return &apiMetrics{metricsCollection: metricsCollection}
}

// WithURLSanitize returns middleware metrics instance with specifice URLSanitize function
func (m *apiMetrics) WithURLSanitize(sanitize utils.URLSanitize) Metrics {
	m.urlSanitize = sanitize
	return m
}

// WriteHeader implementation to preserve status code for metrics middleware
func (w *codeRecordingResponseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
