package tasks

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// изменение названия задачи
func (db *Repo) UpdateTitle(task_id uuid.UUID, newTitle string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET title = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, task_id, newTitle)
	if err != nil {
		return fmt.Errorf("Error in UpdateTitle query: %w", err)
	}

	return nil
}

// изменение описания задачи
func (db *Repo) UpdateDescription(task_id uuid.UUID, newDescription string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET description = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, task_id, newDescription)
	if err != nil {
		return fmt.Errorf("Error in UpdateDescription query: %w", err)
	}

	return nil
}

// изменение времени начала задачи
func (db *Repo) UpdateStartAt(task_id uuid.UUID, newStart time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET start_at = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, task_id, newStart)
	if err != nil {
		return fmt.Errorf("Error in UpdateStartAt query: %w", err)
	}

	return nil
}

// изменение времени окончания задачи
func (db *Repo) UpdateEndAt(task_id uuid.UUID, newEnd time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET end_at = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, task_id, newEnd)
	if err != nil {
		return fmt.Errorf("Error in UpdateEndAt query: %w", err)
	}

	return nil
}

// изменение статуса задачи (+ выполнено)
func (db *Repo) TaskCompleted(task_id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET completed_at = NOW(), updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, task_id)
	if err != nil {
		return fmt.Errorf("Error in TaskCompleted query: %w", err)
	}

	return nil
}

// изменение статуса задачи (- невыполнено)
func (db *Repo) TaskUncompleted(task_id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET completed_at = null, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, task_id)
	if err != nil {
		return fmt.Errorf("Error in TaskUncompleted query: %w", err)
	}

	return nil
}
