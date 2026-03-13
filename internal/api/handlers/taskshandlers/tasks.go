package taskshandlers

import (
	"todoapp/internal/services/tasksservice"
	"todoapp/pkg/logger"
)

type TasksHandler struct {
	taskService *tasksservice.TasksService
	logger *logger.Logger
}

func NewTasksHandler(taskService *tasksservice.TasksService, logger *logger.Logger) *TasksHandler {
	return &TasksHandler{
		taskService: taskService,
		logger: logger,
	}
}
