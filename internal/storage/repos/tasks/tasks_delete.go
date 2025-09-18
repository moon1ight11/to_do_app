package tasks

import "log"

// удаление задачи\подзадачи
func (db *Base) DeleteTask(id int) error {
	query := `
				DELETE FROM tasks
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("Error in DeleteTask query", err)
		return err
	}
	return nil
}

// удаление всех задач пользователя
func (db *Base) DeleteAllTasks(owner_id int) error {
	query := `
				DELETE FROM tasks
				WHERE owner_id = $1
			`
	_, err := db.DB.Exec(query, owner_id)
	if err != nil {
		log.Println("Error in DeleteAllTasks query", err)
		return err
	}
	return nil
}
