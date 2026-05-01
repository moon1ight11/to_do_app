package taskshandlers

import (
	"todoapp/internal/services"
	"todoapp/internal/storage/cache"
	"todoapp/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type TasksHandler struct {
	taskService  services.TasksServiceInterface
	logger       logger.LoggerInterface
	cacheService cache.CacheInterface
	tracer       trace.Tracer
}

func NewTasksHandler(
	taskService services.TasksServiceInterface,
	logger logger.LoggerInterface,
	cacheService cache.CacheInterface,
	tracer trace.Tracer,
) *TasksHandler {
	return &TasksHandler{
		taskService:  taskService,
		logger:       logger,
		cacheService: cacheService,
		tracer:       tracer,
	}
}
