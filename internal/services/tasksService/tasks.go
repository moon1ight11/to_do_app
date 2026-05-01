package tasksservice

import (
	"todoapp/internal/storage/repos/tasksrepos"

	"go.opentelemetry.io/otel/trace"
)

type TasksService struct {
	tasksRepo *tasksrepos.Repo
	tracer    trace.Tracer
}

func NewTasksService(tasksRepo *tasksrepos.Repo, tracer trace.Tracer) *TasksService {
	return &TasksService{
		tasksRepo: tasksRepo,
		tracer:    tracer,
	}
}
