package tasksrepos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// удаление задачи
func (db *Repo) DeleteTask(ctx context.Context, taskId uuid.UUID, tx *sql.Tx) error {
	query := `
				DELETE FROM todo_app.tasks
				WHERE id = $1
			`
	_, err := tx.ExecContext(ctx, query, taskId)
	if err != nil {
		return fmt.Errorf("error in DeleteTask query: %w", err)
	}

	return nil
}

// удаление всех задач пользователя
func (db *Repo) DeleteAllTasks(ctx context.Context, ownerId uuid.UUID) error {
	query := `
				DELETE FROM todo_app.tasks
				WHERE owner_id = $1
			`
	_, err := db.DB.ExecContext(ctx, query, ownerId)
	if err != nil {
		return fmt.Errorf("error in DeleteAllTasks query: %w", err)
	}
	return nil
}
