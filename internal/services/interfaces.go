package services

import (
	"context"
	"time"
	"todoapp/internal/api/models"
	"github.com/google/uuid"
)

type UsersServiceInterface interface {
	AddUser(ctx context.Context, user models.UserAuth) (uuid.UUID, error)
	CheckAndGetUser(ctx context.Context, user models.UserAuth) (models.UserRequest, error)
	CheckNameAndEmail(ctx context.Context, userName string, userEmail string) (bool, error)
	GetUser(ctx context.Context, userId uuid.UUID) (models.UserRequest, error)
	UpdateUser(ctx context.Context, name *string, pass *string, email *string, userId uuid.UUID) error
	DeleteUser(ctx context.Context, userId uuid.UUID) error
}

type TasksServiceInterface interface {
	CreateTask(ctx context.Context, task models.Task) error
	GetAllTasks(ctx context.Context, userId uuid.UUID) ([]models.Task, error)
	GetOneTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (models.Task, error)
	ChangeTask(ctx context.Context,
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

type SettingsServiceInterface interface {
	GetSettings(ctx context.Context, userId uuid.UUID) (models.Setting, error)
	UpdateSettings(ctx context.Context, userId uuid.UUID, duration *float64, tz *string) error
}
