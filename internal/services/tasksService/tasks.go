package tasksservice

import (
	"todoapp/internal/storage/repos/tasksrepos"
)

type TasksService struct {
	tasksRepo *tasksrepos.Repo
}

func NewTasksService(tasksRepo *tasksrepos.Repo) *TasksService {
	return &TasksService{tasksRepo: tasksRepo}
}