package taskshandlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type MockTaskService struct {
	CreateTaskFunc  func(ctx context.Context, task models.Task) error
	GetAllTasksFunc func(ctx context.Context, userId uuid.UUID) ([]models.Task, error)
	GetOneTaskFunc  func(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (models.Task, error)
	ChangeTaskFunc  func(
		ctx context.Context,
		taskId uuid.UUID,
		userId uuid.UUID,
		title *string,
		description *string,
		startAt *time.Time,
		endAt *time.Time,
		completedAt *bool,
	) error
	DeleteTaskFunc func(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error
}

func (m *MockTaskService) CreateTask(ctx context.Context, task models.Task) error {
	return m.CreateTaskFunc(ctx, task)
}

func (m *MockTaskService) GetAllTasks(ctx context.Context, userId uuid.UUID) ([]models.Task, error) {
	return m.GetAllTasksFunc(ctx, userId)
}

func (m *MockTaskService) GetOneTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (models.Task, error) {
	return m.GetOneTaskFunc(ctx, taskId, userId)
}

func (m *MockTaskService) ChangeTask(
	ctx context.Context,
	taskId uuid.UUID,
	userId uuid.UUID,
	title *string,
	description *string,
	startAt *time.Time,
	endAt *time.Time,
	completedAt *bool,
) error {
	return m.ChangeTaskFunc(ctx, taskId, userId, title, description, startAt, endAt, completedAt)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
	return m.DeleteTaskFunc(ctx, taskId, userId)
}
