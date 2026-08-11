package tasksservice

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"todoapp/internal/storage/repos/tasksrepos"
)

type tasksRepo interface {
	DB() *sql.DB
	CreateTask(ctx context.Context, task tasksrepos.Task) error
	DeleteTask(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error
	DeleteAllTasks(ctx context.Context, ownerId uuid.UUID) error
	TasksByOwnerId(ctx context.Context, ownerId uuid.UUID) ([]tasksrepos.Task, error)
	SubtasksByTaskId(ctx context.Context, ownerId uuid.UUID, parentId uuid.UUID) ([]tasksrepos.Task, error)
	TaskById(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID) (tasksrepos.Task, error)
	UpdateTitle(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, title string, tx *sql.Tx) error
	UpdateDescription(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, description string, tx *sql.Tx) error
	UpdateStartAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, start time.Time, tx *sql.Tx) error
	UpdateEndAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, end time.Time, tx *sql.Tx) error
	TaskCompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error
	TaskUncompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error
}
