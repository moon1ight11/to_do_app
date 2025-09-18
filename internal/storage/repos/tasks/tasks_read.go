package tasks

import "log"

// отображение всех родительских задач пользователя
func (db *Base) ParentTasks(owner_id int) ([]Task, error) {
	query := `
				SELECT id, title, description, start_at, end_at
				FROM tasks
				WHERE owner_id = $1 AND parent_task_id IS NULL
				ORDER BY created_at DESC
			`
	var tasks []Task
	rows, err := db.DB.Query(query, owner_id)
	if err != nil {
		log.Println("Error in ParentTasks query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Start_at,
			&task.End_at,
		)
		if err != nil {
			log.Println("Error in ParentTasks scan", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// отображение всех подзадач одной родительской задачи пользователя
func (db *Base) Subtasks(owner_id int, parent_id int) ([]Task, error) {
	query := `
				SELECT id, parent_task_id, title, description, start_at, end_at
				FROM tasks
				WHERE owner_id = $1, AND parent_task_id = $2
				ORDER BY created_at DESC
			`
	var tasks []Task
	rows, err := db.DB.Query(query, owner_id, parent_id)
	if err != nil {
		log.Println("Error in Subtasks query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.Id,
			&task.Parent_id,
			&task.Title,
			&task.Description,
			&task.Start_at,
			&task.End_at,
		)
		if err != nil {
			log.Println("Error in Subtasks scan", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// отображение только невыполненных родительских задач
func (db *Base) OpenParentTasks(owner_id int) ([]Task, error) {
	query := `
				SELECT id, title, description, start_at, end_at
				FROM tasks
				WHERE owner_id = $1, AND parent_task_id IS NULL AND completed_at IS NULL
				ORDER BY created_at DESC
			`
	var tasks []Task
	rows, err := db.DB.Query(query, owner_id)
	if err != nil {
		log.Println("Error in OpenParentTasks query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Start_at,
			&task.End_at,
		)
		if err != nil {
			log.Println("Error in OpenParentTasks scan", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// отображение только невыполненных подзадач
func (db *Base) OpenSubtasks(owner_id int, parent_id int) ([]Task, error) {
	query := `
				SELECT id, parent_task_id, title, description, start_at, end_at
				FROM tasks
				WHERE owner_id = $1, AND parent_task_id = $2, AND completed_at IS NULL
				ORDER BY created_at DESC
			`
	var tasks []Task
	rows, err := db.DB.Query(query, owner_id, parent_id)
	if err != nil {
		log.Println("Error in OpenSubtasks query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.Id,
			&task.Parent_id,
			&task.Title,
			&task.Description,
			&task.Start_at,
			&task.End_at,
		)
		if err != nil {
			log.Println("Error in OpenSubtasks scan", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
