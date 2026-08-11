package tasksservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace/noop"

	"todoapp/internal/api/models"
	"todoapp/internal/storage/repos/tasksrepos"
)

func ptr[T any](v T) *T { return &v }

func TestCreateTask_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockRepo := &MockTasksRepo{
		CreateTaskFunc: func(ctx context.Context, task tasksrepos.Task) error {
			return nil
		},
	}

	svc := NewTasksService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	task := models.Task{Title: ptr("test"), OwnerId: userId}
	err := svc.CreateTask(ctx, task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateTask_EndBeforeStart(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	start := time.Now()
	end := start.Add(-time.Hour)

	mockRepo := &MockTasksRepo{}
	svc := NewTasksService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	task := models.Task{Title: ptr("test"), OwnerId: userId, StartAt: &start, EndAt: &end}
	err := svc.CreateTask(ctx, task)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAllTasks_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	taskId := uuid.New()

	mockRepo := &MockTasksRepo{
		TasksByOwnerIdFunc: func(ctx context.Context, ownerId uuid.UUID) ([]tasksrepos.Task, error) {
			return []tasksrepos.Task{{Id: &taskId, Title: ptr("test"), OwnerId: userId}}, nil
		},
		SubtasksByTaskIdFunc: func(ctx context.Context, ownerId uuid.UUID, parentId uuid.UUID) ([]tasksrepos.Task, error) {
			return []tasksrepos.Task{}, nil
		},
	}

	svc := NewTasksService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	tasks, err := svc.GetAllTasks(ctx, userId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestGetAllTasks_Error(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockRepo := &MockTasksRepo{
		TasksByOwnerIdFunc: func(ctx context.Context, ownerId uuid.UUID) ([]tasksrepos.Task, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewTasksService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	_, err := svc.GetAllTasks(ctx, userId)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetOneTask_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	taskId := uuid.New()

	mockRepo := &MockTasksRepo{
		TaskByIdFunc: func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID) (tasksrepos.Task, error) {
			return tasksrepos.Task{Id: &taskId, Title: ptr("test"), OwnerId: userId}, nil
		},
		SubtasksByTaskIdFunc: func(ctx context.Context, ownerId uuid.UUID, parentId uuid.UUID) ([]tasksrepos.Task, error) {
			return []tasksrepos.Task{}, nil
		},
	}

	svc := NewTasksService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	task, err := svc.GetOneTask(ctx, taskId, userId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *task.Id != taskId {
		t.Errorf("expected taskId %v, got %v", taskId, *task.Id)
	}
}

func TestGetOneTask_NotFound(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	taskId := uuid.New()

	mockRepo := &MockTasksRepo{
		TaskByIdFunc: func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID) (tasksrepos.Task, error) {
			return tasksrepos.Task{}, errors.New("not found")
		},
	}

	svc := NewTasksService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	_, err := svc.GetOneTask(ctx, taskId, userId)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
