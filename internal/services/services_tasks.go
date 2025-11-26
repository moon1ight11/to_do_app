package services

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"todoapp/internal/storage/repos/tasks"
)

type TasksService struct {
	tasksRepo *tasks.Repo
}

func NewTasksService(tasksRepo *tasks.Repo) *TasksService {
	return &TasksService{tasksRepo: tasksRepo}
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
func (t *TasksService) GetAllTasks(user_id uuid.UUID) ([]tasks.Task, error) {
	Tasks, err := t.tasksRepo.TasksByOwnerId(user_id)
	if err != nil {
		return nil, err
	}

	for i := range Tasks {
		Subtasks, err := t.tasksRepo.SubtasksByTaskId(user_id, *Tasks[i].Id)
		if err != nil {
			return nil, err
		}

		Tasks[i].Subtasks = &Subtasks
	}

	return Tasks, err
}

// получение одной задачи
func (t *TasksService) GetOneTask(task_id uuid.UUID) (tasks.Task, error) {
	task, err := t.tasksRepo.TaskById(task_id)
	if err != nil {
		return tasks.Task{}, err
	}

	return task, nil
}

// изменение полей задачи
func (t *TasksService) ChangeTask(task_id uuid.UUID, new_title *string, new_description *string, new_start_at *time.Time, new_end_at *time.Time, completed_at *bool) error {
	// запускаем транзакцию
	transaction, err := t.tasksRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// находим задачу, которую нужно изменить
	OldTask, err := t.tasksRepo.TaskById(task_id)
	if err != nil {
		return err
	}

	// проверяем, чтобы время старта было раньше времени конца
	// если меняются оба времени
	if new_end_at != nil && new_start_at != nil {
		if new_end_at.Before(*new_start_at) {
			return fmt.Errorf("end_at cannot be before start_at")
		}
	}

	// если меняется толко время начала
	if new_end_at == nil && new_start_at != nil && OldTask.End_at != nil {
		if OldTask.End_at.Before(*new_start_at) {
			return fmt.Errorf("end_at cannot be before start_at")
		}
	}

	// если меняется только конец
	if new_start_at == nil && new_end_at != nil && OldTask.Start_at != nil {
		if new_end_at.Before(*OldTask.Start_at) {
			return fmt.Errorf("end_at cannot be before start_at")
		}
	}

	// меняем название
	if new_title != nil {
		err := t.tasksRepo.UpdateTitle(task_id, *new_title, transaction)
		if err != nil {
			return err
		}
	}

	// меняем описание
	if new_description != nil {
		err := t.tasksRepo.UpdateDescription(task_id, *new_description, transaction)
		if err != nil {
			return err
		}
	}

	// меняем время начала
	if new_start_at != nil {
		err := t.tasksRepo.UpdateStartAt(task_id, *new_start_at, transaction)
		if err != nil {
			return err
		}
	}

	// меняем время завершения
	if new_end_at != nil {
		err := t.tasksRepo.UpdateEndAt(task_id, *new_end_at, transaction)
		if err != nil {
			return err
		}
	}

	// меняем статус задачи
	if completed_at != nil {
		if *completed_at == false {
			err := t.tasksRepo.TaskUncompleted(task_id, transaction)
			if err != nil {
				return err
			}
		} else if *completed_at == true {
			err := t.tasksRepo.TaskCompleted(task_id, transaction)
			if err != nil {
				return err
			}
		}
	}

	// если ошибок нет - подтверждаем транзакцию
	transaction.Commit()
	return nil
}

// удаление задачи
func (t *TasksService) DeleteTask(task_id uuid.UUID) error {
	// запускаем транзакцию
	transaction, err := t.tasksRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// ищем задачу по id
	DeletedTask, err := t.tasksRepo.TaskById(task_id)
	if err != nil {
		return err
	}

	// удаляем задачу с подзадачами
	err = t.tasksRepo.DeleteTask(*DeletedTask.Id, transaction)
	if err != nil {
		return err
	}

	// если ошибок нет - подтверждаем транзакцию
	transaction.Commit()
	return nil
}
