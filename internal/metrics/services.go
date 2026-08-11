package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (m *Metrics) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.RequestsInFlight.Inc()
		defer m.RequestsInFlight.Dec()

		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		endpoint := c.FullPath()
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
		m.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	}
}

func (m *Metrics) RecordOperation(operationType string) {
	m.OperationsTotal.WithLabelValues(operationType).Inc()
}

func (m *Metrics) RecordError(errorType, handler string) {
	m.ErrorsTotal.WithLabelValues(errorType, handler).Inc()
}

func (m *Metrics) RegisterMetricsHandler(r *gin.Engine, path string) {
	r.GET(path, func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
}
