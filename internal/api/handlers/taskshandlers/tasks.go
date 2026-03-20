package taskshandlers

import (
	"todoapp/internal/services"
	"todoapp/pkg/logger"
)

type TasksHandler struct {
	taskService services.TasksServiceInterface
	logger logger.LoggerInterface
}

func NewTasksHandler(taskService services.TasksServiceInterface, logger logger.LoggerInterface) *TasksHandler {
	return &TasksHandler{
		taskService: taskService,
		logger: logger,
	}
}
