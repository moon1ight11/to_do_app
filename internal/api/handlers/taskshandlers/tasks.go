package taskshandlers

import (
	"todoapp/internal/services"
	"todoapp/pkg/logger"
)

type TasksHandler struct {
	taskService services.TasksServiceInterface
	logger *logger.Logger
}

func NewTasksHandler(taskService services.TasksServiceInterface, logger *logger.Logger) *TasksHandler {
	return &TasksHandler{
		taskService: taskService,
		logger: logger,
	}
}
