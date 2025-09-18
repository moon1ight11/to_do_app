package tasks

import "log"

// создание задачи
func (db *Base) CreateTask(NewTask Task) error {
	query := `
				INSERT INTO tasks (title, description, start_at, end_at, owner_id)
				VALUES ($1, $2, $3, $4, $5)
			`
	_, err := db.DB.Exec(query, NewTask.Title, NewTask.Description, NewTask.Start_at, NewTask.End_at, NewTask.Owner_id)
	if err != nil {
		log.Println("Error in CreateTask query", err)
		return err
	}

	return nil
}

// создание подзадачи
func (db *Base) CreateSubtask(NewSubtask Task) error {
	query := `
				INSERT INTO tasks (title, description, start_at, end_at, owner_id, parent_id)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
	_, err := db.DB.Exec(
		query,
		NewSubtask.Title,
		NewSubtask.Description,
		NewSubtask.Start_at,
		NewSubtask.End_at,
		NewSubtask.Owner_id,
		NewSubtask.Parent_id,
	)
	if err != nil {
		log.Println("Error in CreateSubtask query", err)
		return err
	}

	return nil
}
