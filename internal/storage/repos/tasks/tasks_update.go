package tasks

import (
	"log"
	"time"
)

// изменение названия задачи
func (db *Base) UpdateTitle(id int, newTitle string) error {
	query := `
				UPDATE tasks
				SET title = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newTitle)
	if err != nil {
		log.Println("Error in UpdateTitle query", err)
		return err
	}

	return nil
}

// изменение описания задачи
func (db *Base) UpdateDescription(id int, newDescription string) error {
	query := `
				UPDATE tasks
				SET description = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newDescription)
	if err != nil {
		log.Println("Error in UpdateDescription query", err)
		return err
	}

	return nil
}

// изменение времени начала задачи
func (db *Base) UpdateStartAt(id int, newStart time.Time) error {
	query := `
				UPDATE tasks
				SET start_at = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newStart)
	if err != nil {
		log.Println("Error in UpdateStartAt query", err)
		return err
	}

	return nil
}

// изменение времени окончания задачи
func (db *Base) UpdateEndAt(id int, newEnd time.Time) error {
	query := `
				UPDATE tasks
				SET end_at = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, id, newEnd)
	if err != nil {
		log.Println("Error in UpdateEndAt query", err)
		return err
	}

	return nil
}

// изменение статуса задачи (+ выполнено)
func (db *Base) TaskCompleted(id int) error {
	query := `
				UPDATE tasks
				SET completed_at = NOW(), updated_at = NOW()
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("Error in TaskCompleted query", err)
		return err
	}

	return nil
}
