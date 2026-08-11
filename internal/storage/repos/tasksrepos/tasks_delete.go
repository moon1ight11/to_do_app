package tasksrepos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (db *Repo) DeleteTask(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.DeleteTask")
	defer span.End()

	query := `
		DELETE FROM todo_app.tasks
		WHERE id = $1 AND owner_id = $2
	`

	_, err := tx.ExecContext(ctx, query, taskId, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.DeleteTask: exec: %w", err)
	}

	return nil
}

func (db *Repo) DeleteAllTasks(ctx context.Context, ownerId uuid.UUID) error {
	ctx, span := db.tracer.Start(ctx, "repo.DeleteAllTasks")
	defer span.End()

	query := `
		DELETE FROM todo_app.tasks
		WHERE owner_id = $1
	`

	_, err := db.DB().ExecContext(ctx, query, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.DeleteAllTasks: exec: %w", err)
	}

	return nil
}
