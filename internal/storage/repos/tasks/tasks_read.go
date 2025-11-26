package tasks

import (
	"fmt"
	"github.com/google/uuid"
)

// отображение всех родительских задач пользователя
func (db *Repo) TasksByOwnerId(owner_id uuid.UUID) ([]Task, error) {
	query := `
				SELECT id, title, description, start_at, end_at,
				completed_at IS NOT NULL as completed
				FROM todo_app.tasks
				WHERE owner_id = $1 AND parent_task_id IS NULL
				ORDER BY created_at DESC
			`
	var tasks []Task
	rows, err := db.DB.Query(query, owner_id)
	if err != nil {
		return nil, fmt.Errorf("Error in ParentTasks query: %w", err)
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
			&task.Completed_at,
		)
		if err != nil {
			return nil, fmt.Errorf("Error in ParentTasks scan: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// отображение всех подзадач одной родительской задачи пользователя
func (db *Repo) SubtasksByTaskId(owner_id uuid.UUID, parent_id uuid.UUID) ([]Task, error) {
	query := `
				SELECT id, parent_task_id, title, description, start_at, end_at,
				completed_at IS NOT NULL as completed
				FROM todo_app.tasks
				WHERE owner_id = $1 AND parent_task_id = $2
				ORDER BY created_at DESC
			`
	var tasks []Task
	rows, err := db.DB.Query(query, owner_id, parent_id)
	if err != nil {
		return nil, fmt.Errorf("Error in Subtasks query: %w", err)
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
			&task.Completed_at,
		)
		if err != nil {
			return nil, fmt.Errorf("Error in Subtasks scan: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// поиск задачи по id
func (db *Repo) TaskById(task_id uuid.UUID) (Task, error) {
	query := `
				SELECT id, title, description, start_at, end_at, parent_task_id
				FROM todo_app.tasks
				WHERE id = $1
			`

	var task Task
	err := db.DB.QueryRow(query, task_id).Scan(&task.Id, &task.Title, &task.Description, &task.Start_at, &task.End_at, &task.Parent_id)

	if err != nil {
		return Task{}, fmt.Errorf("Task not found")
	}

	return task, nil
}
