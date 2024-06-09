package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"path"},
	)
)

func Init() {
	prometheus.MustRegister(httpRequestTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func Handler(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(path))
		defer timer.ObserveDuration()
		httpRequestTotal.WithLabelValues(path).Inc()
		next.ServeHTTP(w, r)
	}
}
