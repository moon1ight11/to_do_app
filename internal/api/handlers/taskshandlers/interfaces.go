package taskshandlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type taskService interface {
	CreateTask(ctx context.Context, task models.Task) error
	GetAllTasks(ctx context.Context, userId uuid.UUID) ([]models.Task, error)
	GetOneTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (models.Task, error)
	ChangeTask(
		ctx context.Context,
		taskId uuid.UUID,
		userId uuid.UUID,
		title *string,
		description *string,
		startAt *time.Time,
		endAt *time.Time,
		completedAt *bool,
	) error
	DeleteTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error
}
