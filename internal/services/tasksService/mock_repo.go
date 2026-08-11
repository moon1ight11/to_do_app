package tasksservice

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"todoapp/internal/storage/repos/tasksrepos"
)

type MockTasksRepo struct {
	DBField               *sql.DB
	CreateTaskFunc        func(ctx context.Context, task tasksrepos.Task) error
	DeleteTaskFunc        func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error
	DeleteAllTasksFunc    func(ctx context.Context, ownerId uuid.UUID) error
	TasksByOwnerIdFunc    func(ctx context.Context, ownerId uuid.UUID) ([]tasksrepos.Task, error)
	SubtasksByTaskIdFunc  func(ctx context.Context, ownerId uuid.UUID, parentId uuid.UUID) ([]tasksrepos.Task, error)
	TaskByIdFunc          func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID) (tasksrepos.Task, error)
	UpdateTitleFunc       func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, title string, tx *sql.Tx) error
	UpdateDescriptionFunc func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, description string, tx *sql.Tx) error
	UpdateStartAtFunc     func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, start time.Time, tx *sql.Tx) error
	UpdateEndAtFunc       func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, end time.Time, tx *sql.Tx) error
	TaskCompletedFunc     func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error
	TaskUncompletedFunc   func(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error
}

func (m *MockTasksRepo) DB() *sql.DB { return m.DBField }

func (m *MockTasksRepo) CreateTask(ctx context.Context, task tasksrepos.Task) error {
	return m.CreateTaskFunc(ctx, task)
}

func (m *MockTasksRepo) DeleteTask(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	return m.DeleteTaskFunc(ctx, taskId, ownerId, tx)
}

func (m *MockTasksRepo) DeleteAllTasks(ctx context.Context, ownerId uuid.UUID) error {
	return m.DeleteAllTasksFunc(ctx, ownerId)
}

func (m *MockTasksRepo) TasksByOwnerId(ctx context.Context, ownerId uuid.UUID) ([]tasksrepos.Task, error) {
	return m.TasksByOwnerIdFunc(ctx, ownerId)
}

func (m *MockTasksRepo) SubtasksByTaskId(ctx context.Context, ownerId uuid.UUID, parentId uuid.UUID) ([]tasksrepos.Task, error) {
	return m.SubtasksByTaskIdFunc(ctx, ownerId, parentId)
}

func (m *MockTasksRepo) TaskById(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID) (tasksrepos.Task, error) {
	return m.TaskByIdFunc(ctx, taskId, ownerId)
}

func (m *MockTasksRepo) UpdateTitle(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, title string, tx *sql.Tx) error {
	return m.UpdateTitleFunc(ctx, taskId, ownerId, title, tx)
}

func (m *MockTasksRepo) UpdateDescription(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, description string, tx *sql.Tx) error {
	return m.UpdateDescriptionFunc(ctx, taskId, ownerId, description, tx)
}

func (m *MockTasksRepo) UpdateStartAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, start time.Time, tx *sql.Tx) error {
	return m.UpdateStartAtFunc(ctx, taskId, ownerId, start, tx)
}

func (m *MockTasksRepo) UpdateEndAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, end time.Time, tx *sql.Tx) error {
	return m.UpdateEndAtFunc(ctx, taskId, ownerId, end, tx)
}

func (m *MockTasksRepo) TaskCompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	return m.TaskCompletedFunc(ctx, taskId, ownerId, tx)
}

func (m *MockTasksRepo) TaskUncompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	return m.TaskUncompletedFunc(ctx, taskId, ownerId, tx)
}
