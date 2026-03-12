package tasksrepos

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// изменение названия задачи
func (db *Repo) UpdateTitle(taskId uuid.UUID, title string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET title = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, taskId, title)
	if err != nil {
		return fmt.Errorf("Error in UpdateTitle query: %w", err)
	}

	return nil
}

// изменение описания задачи
func (db *Repo) UpdateDescription(taskId uuid.UUID, description string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET description = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, taskId, description)
	if err != nil {
		return fmt.Errorf("Error in UpdateDescription query: %w", err)
	}

	return nil
}

// изменение времени начала задачи
func (db *Repo) UpdateStartAt(taskId uuid.UUID, start time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET start_at = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, taskId, start)
	if err != nil {
		return fmt.Errorf("Error in UpdateStartAt query: %w", err)
	}

	return nil
}

// изменение времени окончания задачи
func (db *Repo) UpdateEndAt(taskId uuid.UUID, end time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET end_at = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, taskId, end)
	if err != nil {
		return fmt.Errorf("Error in UpdateEndAt query: %w", err)
	}

	return nil
}

// изменение статуса задачи (+ выполнено)
func (db *Repo) TaskCompleted(taskId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET completed_at = NOW(), updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, taskId)
	if err != nil {
		return fmt.Errorf("Error in TaskCompleted query: %w", err)
	}

	return nil
}

// изменение статуса задачи (- невыполнено)
func (db *Repo) TaskUncompleted(taskId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET completed_at = null, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, taskId)
	if err != nil {
		return fmt.Errorf("Error in TaskUncompleted query: %w", err)
	}

	return nil
}
