package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// мидлвар для сбора метрик
func (m *Metrics) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// увеличиваем счетчик текущих запросов и гарантированно его уменьшаем при выходе
		m.RequestsInFlight.Inc()
		defer m.RequestsInFlight.Dec()

		// засекаем время начала обработки
		start := time.Now()

		// идем дальше в логику приложения
		c.Next()

		// высчитываем длительность запроса
		duration := time.Since(start).Seconds()

		// замечаем путь запроса и его метаданные
		endpoint := c.FullPath() // в Gin используем FullPath() для паттерна маршрута
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		// записываем длительность запроса и увеличиваем счетчик выполненных запросов
		m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
		m.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	}
}

// запись операции
func (m *Metrics) RecordOperation(operationType string) {
	m.OperationsTotal.WithLabelValues(operationType).Inc()
}

// запись ошибки
func (m *Metrics) RecordError(errorType, handler string) {
	m.ErrorsTotal.WithLabelValues(errorType, handler).Inc()
}

// регистрация эндпоинта прометея
func (m *Metrics) RegisterMetricsHandler(r *gin.Engine, path string) {
	r.GET(path, func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
}
