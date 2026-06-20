package middleware

import (
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"time"
)

var metricServerReq = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "http_server",
	Subsystem: "requests",
	Name:      "total",
	Help:      "The total number of http requests.",
	Labels:    []string{"path", "method", "status_code"},
})

var metricServerDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
	Namespace: "http_server",
	Subsystem: "requests",
	Name:      "duration_ms",
	Help:      "The duration of http requests.",
	Labels:    []string{"path", "method"},
	Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
})

func PrometheusMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next(rw, r)
			dur := time.Since(start).Milliseconds()
			metricServerDur.Observe(float64(dur), r.URL.Path, r.Method)
			metricServerReq.Inc(r.URL.Path, r.Method, http.StatusText(rw.statusCode))
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
