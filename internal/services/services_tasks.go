package services

import (
	"github.com/google/uuid"
	"time"
	"todoapp/internal/storage/repos/tasks"
)

type TasksService struct {
	tasksRepo *tasks.Base
}

// создание задачи(подзадачи)
func (t *TasksService) CreateTask(newTask tasks.Task) error {
	err := t.tasksRepo.CreateTask(newTask)
	if err != nil {
		return err
	}
	return nil
}

// получение списка всех задач пользователя
func (t *TasksService) GetTasks(user_id uuid.UUID) ([]tasks.Task, error) {
	Tasks, err := t.tasksRepo.ParentTasks(user_id)
	if err != nil {
		return nil, err
	}

	for i := range Tasks {
		Subtasks, err := t.tasksRepo.Subtasks(user_id, *Tasks[i].Id)
		if err != nil {
			return nil, err
		}

		Tasks[i].Subtasks = &Subtasks
	}

	return Tasks, err
}

// получение одной задачи
func (t *TasksService) GetOneTask(task_id uuid.UUID) (tasks.Task, error) {
	task, err := t.tasksRepo.GetTaskById(task_id)
	if err != nil {
		return tasks.Task{}, err
	}

	return task, nil
}

// изменение полей задачи
func (t *TasksService) ChangeTask(id uuid.UUID, new_title *string, new_description *string, new_start_at *time.Time, new_end_at *time.Time, completed_at *time.Time) error {
	// запускаем транзакцию
	transaction, err := t.tasksRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// меняем название
	if new_title != nil {
		err := t.tasksRepo.UpdateTitle(id, *new_title)
		if err != nil {
			return err
		}
	}

	// меняем описание
	if new_description != nil {
		err := t.tasksRepo.UpdateDescription(id, *new_description)
		if err != nil {
			return err
		}
	}

	// меняем время начала
	if new_start_at != nil {
		err := t.tasksRepo.UpdateStartAt(id, *new_start_at)
		if err != nil {
			return err
		}
	}

	// меняем время завершения
	if new_end_at != nil {
		err := t.tasksRepo.UpdateEndAt(id, *new_end_at)
		if err != nil {
			return err
		}
	}

	// меняем статус задачи
	if completed_at != nil {
		err := t.tasksRepo.TaskCompleted(id)
		if err != nil {
			return err
		}
	}

	// если ошибок нет - подтверждаем транзакцию
	transaction.Commit()
	return nil
}

// удаление задачи
func (t *TasksService) DeleteTask(id uuid.UUID) error {
	// запускаем транзакцию
	transaction, err := t.tasksRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// ищем задачу по id
	DeletedTask, err := t.tasksRepo.GetTaskById(id)
	if err != nil {
		return err
	}

	// ищем подзадачи для задачи
	DeletedSubtasks, err := t.tasksRepo.Subtasks(DeletedTask.Owner_id, *DeletedTask.Id)
	if err != nil {
		return err
	}

	// если подзадачи есть - удаляем
	if len(DeletedSubtasks) > 0 {
		for i := range DeletedSubtasks {
			St := DeletedSubtasks[i]
			err := t.tasksRepo.DeleteTask(*St.Id)
			if err != nil {
				return err
			}
		}
	}

	// удаляем саму задачу
	err = t.tasksRepo.DeleteTask(*DeletedTask.Id)
	if err != nil {
		return err
	}

	// если ошибок нет - подтверждаем транзакцию
	transaction.Commit()
	return nil
}
