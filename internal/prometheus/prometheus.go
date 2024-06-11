package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	httpRequestTotal    *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		httpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"path"}),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "Duration of HTTP requests in seconds",
			},
			[]string{"path"}),
	}
}

func (m *Metrics) Register() {
	prometheus.MustRegister(m.httpRequestTotal)
	prometheus.MustRegister(m.httpRequestDuration)
}

func (m *Metrics) Handler(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(m.httpRequestDuration.WithLabelValues(path))
		defer timer.ObserveDuration()
		m.httpRequestTotal.WithLabelValues(path).Inc()
		next.ServeHTTP(w, r)
	}
}
