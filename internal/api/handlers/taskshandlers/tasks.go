package taskshandlers

import (
	"todoapp/internal/metrics"
	"todoapp/internal/storage/cache"
	"todoapp/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type TasksHandler struct {
	taskService  taskService
	logger       logger.LoggerInterface
	cacheService cache.CacheInterface
	tracer       trace.Tracer
	metrics      *metrics.Metrics
}

func NewTasksHandler(
	taskService taskService,
	logger logger.LoggerInterface,
	cacheService cache.CacheInterface,
	tracer trace.Tracer,
	metrics *metrics.Metrics,
) *TasksHandler {
	return &TasksHandler{
		taskService:  taskService,
		logger:       logger,
		cacheService: cacheService,
		tracer:       tracer,
		metrics:      metrics,
	}
}
