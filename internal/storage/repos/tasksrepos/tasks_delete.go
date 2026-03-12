package tasksrepos

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// удаление задачи
func (db *Repo) DeleteTask(taskId uuid.UUID, tx *sql.Tx) error {
	query := `
				DELETE FROM todo_app.tasks
				WHERE id = $1
			`
	_, err := tx.Exec(query, taskId)
	if err != nil {
		return fmt.Errorf("Error in DeleteTask query: %w", err)
	}

	return nil
}

// удаление всех задач пользователя
func (db *Repo) DeleteAllTasks(ownerId uuid.UUID) error {
	query := `
				DELETE FROM todo_app.tasks
				WHERE owner_id = $1
			`
	_, err := db.DB.Exec(query, ownerId)
	if err != nil {
		return fmt.Errorf("Error in DeleteAllTasks query: %w", err)
	}
	return nil
}
