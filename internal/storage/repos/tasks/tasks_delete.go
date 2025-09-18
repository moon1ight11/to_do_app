package tasks

import (
	"fmt"
	"github.com/google/uuid"
)

// удаление задачи\подзадачи
func (db *Base) DeleteTask(id uuid.UUID) error {
	query := `
				DELETE FROM tasks
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error in DeleteTask query: %w", err)
	}

	return nil
}

// удаление всех задач пользователя
func (db *Base) DeleteAllTasks(owner_id uuid.UUID) error {
	query := `
				DELETE FROM tasks
				WHERE owner_id = $1
			`
	_, err := db.DB.Exec(query, owner_id)
	if err != nil {
		return fmt.Errorf("Error in DeleteAllTasks query: %w", err)
	}
	return nil
}
