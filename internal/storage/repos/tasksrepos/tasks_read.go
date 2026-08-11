package tasksrepos

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (db *Repo) TasksByOwnerId(ctx context.Context, ownerId uuid.UUID) ([]Task, error) {
	ctx, span := db.tracer.Start(ctx, "repo.TasksByOwnerId")
	defer span.End()

	query := `
		SELECT id, title, description, owner_id, start_at, end_at, completed_at IS NOT NULL as completed
		FROM todo_app.tasks
		WHERE owner_id = $1 AND parent_task_id IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.DB().QueryContext(ctx, query, ownerId)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("tasksrepos.TasksByOwnerId: query: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.OwnerId,
			&task.StartAt,
			&task.EndAt,
			&task.CompletedAt,
		); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("tasksrepos.TasksByOwnerId: scan: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("tasksrepos.TasksByOwnerId: rows: %w", err)
	}

	return tasks, nil
}

func (db *Repo) SubtasksByTaskId(ctx context.Context, ownerId uuid.UUID, parentId uuid.UUID) ([]Task, error) {
	ctx, span := db.tracer.Start(ctx, "repo.SubtasksByTaskId")
	defer span.End()

	query := `
		SELECT id, parent_task_id, title, description, owner_id, start_at, end_at,
		completed_at IS NOT NULL as completed
		FROM todo_app.tasks
		WHERE owner_id = $1 AND parent_task_id = $2
		ORDER BY created_at DESC
	`

	rows, err := db.DB().QueryContext(ctx, query, ownerId, parentId)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("tasksrepos.SubtasksByTaskId: query: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.Id,
			&task.ParentId,
			&task.Title,
			&task.Description,
			&task.OwnerId,
			&task.StartAt,
			&task.EndAt,
			&task.CompletedAt,
		); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("tasksrepos.SubtasksByTaskId: scan: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("tasksrepos.SubtasksByTaskId: rows: %w", err)
	}

	return tasks, nil
}

func (db *Repo) TaskById(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID) (Task, error) {
	ctx, span := db.tracer.Start(ctx, "repo.TaskById")
	defer span.End()

	query := `
		SELECT id, title, description, owner_id, start_at, end_at, parent_task_id
		FROM todo_app.tasks
		WHERE id = $1 AND owner_id = $2
	`

	var task Task
	err := db.DB().QueryRowContext(ctx, query, taskId, ownerId).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.OwnerId,
		&task.StartAt,
		&task.EndAt,
		&task.ParentId,
	)
	if err != nil {
		span.RecordError(err)
		return Task{}, fmt.Errorf("tasksrepos.TaskById: query: %w", err)
	}

	return task, nil
}
