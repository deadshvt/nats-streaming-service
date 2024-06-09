package middleware

import (
	"net/http"
	"time"

	"github.com/deadshvt/nats-streaming-service/pkg/logger"

	"github.com/rs/zerolog"
)

func Logging(lgr zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.LogWithParams(lgr, "Completed request", struct {
			Method string
			URI    string
			Took   time.Duration
		}{Method: r.Method, URI: r.RequestURI, Took: time.Since(start)})
	})
}
