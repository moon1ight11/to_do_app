package taskshandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services/tasksservice"
)

type TasksHandler struct {
	taskService *tasksservice.TasksService
	jwtService  jwt.TokenService
}

func NewTasksHandler(taskService *tasksservice.TasksService, jwtService jwt.TokenService) *TasksHandler {
	return &TasksHandler{
		taskService: taskService,
		jwtService:  jwtService}
}
