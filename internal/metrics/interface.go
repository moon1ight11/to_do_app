package metrics

import (
	"github.com/gin-gonic/gin"
)

type MetricsInterface interface {
	Middleware() gin.HandlerFunc
	RecordOperation(operationType string)
	RecordError(errorType, handler string)
	RegisterMetricsHandler(r *gin.Engine, path string)
}
