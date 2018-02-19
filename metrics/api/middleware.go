package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/iryonetwork/wwm/utils"
)

type Metrics interface {
	Middleware(next http.Handler) http.Handler
	WithURLSanitize(sanitize utils.URLSanitize) Metrics
}

type metrics struct {
	requestSeconds *prometheus.HistogramVec
	urlSanitize    utils.URLSanitize
}

type codeRecordingResponseWriter struct {
	http.ResponseWriter
	code int
}

func (m *metrics) Middleware(next http.Handler) http.Handler {
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
			m.requestSeconds.
				With(prometheus.Labels{"path": path, "method": method, "code": fmt.Sprintf("%d", code)}).
				Observe(duration.Seconds())
		}()

		cw := &codeRecordingResponseWriter{w, 200}

		next.ServeHTTP(cw, r)

		code = cw.code
	})
}

func NewMetrics(namespace string, subsystem string) Metrics {
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_seconds",
		Help:      "Time taken to serve API request",
	}, []string{"path", "method", "code"})

	// Register metrics
	prometheus.MustRegister(h)

	return &metrics{requestSeconds: h}
}

func (m *metrics) WithURLSanitize(sanitize utils.URLSanitize) Metrics {
	m.urlSanitize = sanitize
	return m
}

func (w *codeRecordingResponseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
