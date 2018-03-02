package log

import (
	"net/http"

	"github.com/rs/zerolog"
)

func APILogMiddleware(next http.Handler, logger zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("New request")
		next.ServeHTTP(w, r)
	})
}
