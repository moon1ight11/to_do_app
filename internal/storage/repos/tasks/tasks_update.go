package tasks

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/google/uuid"
)

// изменение названия задачи
func (db *Base) UpdateTitle(id uuid.UUID, newTitle string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET title = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, id, newTitle)
	if err != nil {
		return fmt.Errorf("Error in UpdateTitle query: %w", err)
	}

	return nil
}

// изменение описания задачи
func (db *Base) UpdateDescription(id uuid.UUID, newDescription string, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET description = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, id, newDescription)
	if err != nil {
		return fmt.Errorf("Error in UpdateDescription query: %w", err)
	}

	return nil
}

// изменение времени начала задачи
func (db *Base) UpdateStartAt(id uuid.UUID, newStart time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET start_at = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, id, newStart)
	if err != nil {
		return fmt.Errorf("Error in UpdateStartAt query: %w", err)
	}

	return nil
}

// изменение времени окончания задачи
func (db *Base) UpdateEndAt(id uuid.UUID, newEnd time.Time, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET end_at = $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, id, newEnd)
	if err != nil {
		return fmt.Errorf("Error in UpdateEndAt query: %w", err)
	}

	return nil
}

// изменение статуса задачи (+ выполнено)
func (db *Base) TaskCompleted(id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.tasks
				SET completed_at = NOW(), updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error in TaskCompleted query: %w", err)
	}

	return nil
}
