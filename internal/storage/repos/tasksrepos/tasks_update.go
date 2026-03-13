package tasksrepos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// изменение названия задачи
func (db *Repo) UpdateTitle(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, title string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET title = $2, updated_at = NOW()
				WHERE id = $1 AND owner_id = $3
			`
	_, err := tx.ExecContext(ctx, query, taskId, title, ownerId)
	if err != nil {
		return fmt.Errorf("error in UpdateTitle query: %w", err)
	}

	return nil
}

// изменение описания задачи
func (db *Repo) UpdateDescription(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, description string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET description = $2, updated_at = NOW()
				WHERE id = $1 AND owner_id = $3
			`
	_, err := tx.ExecContext(ctx, query, taskId, description, ownerId)
	if err != nil {
		return fmt.Errorf("error in UpdateDescription query: %w", err)
	}

	return nil
}

// изменение времени начала задачи
func (db *Repo) UpdateStartAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, start time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET start_at = $2, updated_at = NOW()
				WHERE id = $1 AND owner_id = $3
			`
	_, err := tx.ExecContext(ctx, query, taskId, start, ownerId)
	if err != nil {
		return fmt.Errorf("error in UpdateStartAt query: %w", err)
	}

	return nil
}

// изменение времени окончания задачи
func (db *Repo) UpdateEndAt(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, end time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET end_at = $2, updated_at = NOW()
				WHERE id = $1 AND owner_id = $3
			`
	_, err := tx.ExecContext(ctx, query, taskId, end, ownerId)
	if err != nil {
		return fmt.Errorf("error in UpdateEndAt query: %w", err)
	}

	return nil
}

// изменение статуса задачи (+ выполнено)
func (db *Repo) TaskCompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET completed_at = NOW(), updated_at = NOW()
				WHERE id = $1 AND owner_id = $2
			`
	_, err := tx.ExecContext(ctx, query, taskId, ownerId)
	if err != nil {
		return fmt.Errorf("error in TaskCompleted query: %w", err)
	}

	return nil
}

// изменение статуса задачи (- невыполнено)
func (db *Repo) TaskUncompleted(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET completed_at = null, updated_at = NOW()
				WHERE id = $1 AND owner_id = $2
			`
	_, err := tx.ExecContext(ctx, query, taskId, ownerId)
	if err != nil {
		return fmt.Errorf("error in TaskUncompleted query: %w", err)
	}

	return nil
}
