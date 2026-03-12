package taskshandlers

import (
	"todoapp/internal/services/tasksservice"
)

type TasksHandler struct {
	taskService *tasksservice.TasksService
}

func NewTasksHandler(taskService *tasksservice.TasksService) *TasksHandler {
	return &TasksHandler{
		taskService: taskService,
	}
}
