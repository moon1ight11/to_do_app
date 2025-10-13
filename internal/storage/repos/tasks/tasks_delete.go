package tasks

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// удаление задачи\подзадачи
func (db *Base) DeleteTask(id uuid.UUID, tx *sql.Tx) error {
	query := `
				DELETE FROM todo_app.tasks
				WHERE id = $1
			`
	_, err := tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error in DeleteTask query: %w", err)
	}

	return nil
}

// удаление всех задач пользователя
func (db *Base) DeleteAllTasks(owner_id uuid.UUID) error {
	query := `
				DELETE FROM todo_app.tasks
				WHERE owner_id = $1
			`
	_, err := db.DB.Exec(query, owner_id)
	if err != nil {
		return fmt.Errorf("Error in DeleteAllTasks query: %w", err)
	}
	return nil
}
