package tasksrepos

import (
	"context"
	"fmt"
)

func (db *Repo) CreateTask(ctx context.Context, task Task) error {
	ctx, span := db.tracer.Start(ctx, "repo.CreateTask")
	defer span.End()

	query := `
		INSERT INTO todo_app.tasks (title, description, start_at, end_at, owner_id, parent_task_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.DB().ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.StartAt,
		task.EndAt,
		task.OwnerId,
		task.ParentId,
	)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.CreateTask: exec: %w", err)
	}

	return nil
}
