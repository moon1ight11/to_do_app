package tasks

import "fmt"

// создание задачи
func (db *Base) CreateTask(NewTask Task) error {
	query := `
				INSERT INTO todo_app.tasks (title, description, start_at, end_at, owner_id, parent_task_id)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
	_, err := db.DB.Exec(
		query,
		NewTask.Title,
		NewTask.Description,
		NewTask.Start_at,
		NewTask.End_at,
		NewTask.Owner_id,
		NewTask.Parent_id,
	)
	if err != nil {
		return fmt.Errorf("Error in CreateTask query: %w", err)
	}

	return nil
}
