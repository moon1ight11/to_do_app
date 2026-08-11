package tasksservice

import "go.opentelemetry.io/otel/trace"

type TasksService struct {
	tasksRepo tasksRepo
	tracer    trace.Tracer
}

func NewTasksService(tasksRepo tasksRepo, tracer trace.Tracer) *TasksService {
	return &TasksService{
		tasksRepo: tasksRepo,
		tracer:    tracer,
	}
}
