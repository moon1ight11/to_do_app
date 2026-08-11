package tasksrepos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (db *Repo) UpdateTitle(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, title string, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateTitle")
	defer span.End()

	query := `
		UPDATE todo_app.tasks
		SET title = $2, updated_at = NOW()
		WHERE id = $1 AND owner_id = $3
	`

	_, err := tx.ExecContext(ctx, query, taskId, title, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.UpdateTitle: exec: %w", err)
	}

	return nil
}

func (db *Repo) UpdateDescription(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, description string, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateDescription")
	defer span.End()

	query := `
		UPDATE todo_app.tasks
		SET description = $2, updated_at = NOW()
		WHERE id = $1 AND owner_id = $3
	`

	_, err := tx.ExecContext(ctx, query, taskId, description, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.UpdateDescription: exec: %w", err)
	}

	return nil
}

func (db *Repo) UpdateStartAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, start time.Time, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateStartAt")
	defer span.End()

	query := `
		UPDATE todo_app.tasks
		SET start_at = $2, updated_at = NOW()
		WHERE id = $1 AND owner_id = $3
	`

	_, err := tx.ExecContext(ctx, query, taskId, start, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.UpdateStartAt: exec: %w", err)
	}

	return nil
}

func (db *Repo) UpdateEndAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, end time.Time, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateEndAt")
	defer span.End()

	query := `
		UPDATE todo_app.tasks
		SET end_at = $2, updated_at = NOW()
		WHERE id = $1 AND owner_id = $3
	`

	_, err := tx.ExecContext(ctx, query, taskId, end, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.UpdateEndAt: exec: %w", err)
	}

	return nil
}

func (db *Repo) TaskCompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.TaskCompleted")
	defer span.End()

	query := `
		UPDATE todo_app.tasks
		SET completed_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND owner_id = $2
	`

	_, err := tx.ExecContext(ctx, query, taskId, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.TaskCompleted: exec: %w", err)
	}

	return nil
}

func (db *Repo) TaskUncompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.TaskUncompleted")
	defer span.End()

	query := `
		UPDATE todo_app.tasks
		SET completed_at = null, updated_at = NOW()
		WHERE id = $1 AND owner_id = $2
	`

	_, err := tx.ExecContext(ctx, query, taskId, ownerId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksrepos.TaskUncompleted: exec: %w", err)
	}

	return nil
}
