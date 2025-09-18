package tasks

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// изменение названия задачи
func (db *Base) UpdateTitle(id uuid.UUID, newTitle string) error {
	query := `
				UPDATE tasks
				SET title = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newTitle)
	if err != nil {
		return fmt.Errorf("Error in UpdateTitle query: %w", err)
	}

	return nil
}

// изменение описания задачи
func (db *Base) UpdateDescription(id uuid.UUID, newDescription string) error {
	query := `
				UPDATE tasks
				SET description = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newDescription)
	if err != nil {
		return fmt.Errorf("Error in UpdateDescription query: %w", err)
	}

	return nil
}

// изменение времени начала задачи
func (db *Base) UpdateStartAt(id uuid.UUID, newStart time.Time) error {
	query := `
				UPDATE tasks
				SET start_at = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newStart)
	if err != nil {
		return fmt.Errorf("Error in UpdateStartAt query: %w", err)
	}

	return nil
}

// изменение времени окончания задачи
func (db *Base) UpdateEndAt(id uuid.UUID, newEnd time.Time) error {
	query := `
				UPDATE tasks
				SET end_at = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newEnd)
	if err != nil {
		return fmt.Errorf("Error in UpdateEndAt query: %w", err)
	}

	return nil
}

// изменение статуса задачи (+ выполнено)
func (db *Base) TaskCompleted(id uuid.UUID) error {
	query := `
				UPDATE tasks
				SET completed_at = NOW(), updated_at = NOW()
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error in TaskCompleted query: %w", err)
	}

	return nil
}
