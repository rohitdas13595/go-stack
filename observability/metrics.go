package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gostack_http_requests_total",
		Help: "HTTP requests",
	}, []string{"method", "path", "code"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gostack_http_request_duration_seconds",
		Help:    "HTTP latency",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

// MetricsHandler returns /metrics handler.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// HTTPMiddleware records Prometheus metrics.
func HTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			path := r.URL.Path
			httpRequests.WithLabelValues(r.Method, path, strconv.Itoa(rw.status)).Inc()
			httpDuration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
