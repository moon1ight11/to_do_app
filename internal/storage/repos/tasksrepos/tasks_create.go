package tasksrepos

import "fmt"

// создание задачи
func (db *Repo) CreateTask(task Task) error {
	query := `
				INSERT INTO todo_app.tasks (title, description, start_at, end_at, owner_id, parent_task_id)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
	_, err := db.DB.Exec(
		query,
		task.Title,
		task.Description,
		task.StartAt,
		task.EndAt,
		task.OwnerId,
		task.ParentId,
	)
	if err != nil {
		return fmt.Errorf("Error in CreateTask query: %w", err)
	}

	return nil
}
