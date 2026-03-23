package taskshandlers

import (
	"todoapp/internal/services"
	"todoapp/internal/storage/cache"
	"todoapp/pkg/logger"
)

type TasksHandler struct {
	taskService  services.TasksServiceInterface
	logger       logger.LoggerInterface
	cacheService cache.CacheInterface
}

func NewTasksHandler(taskService services.TasksServiceInterface, logger logger.LoggerInterface, cacheService cache.CacheInterface) *TasksHandler {
	return &TasksHandler{
		taskService:  taskService,
		logger:       logger,
		cacheService: cacheService,
	}
}
