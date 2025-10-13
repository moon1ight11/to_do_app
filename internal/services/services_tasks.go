package services

import (
	"fmt"
	"time"
	"todoapp/internal/storage/repos/tasks"
	"github.com/google/uuid"
)

type TasksService struct {
	tasksRepo *tasks.Base
}

func NewTasksService(tasksRepo *tasks.Base) *TasksService {
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

	// находим задачу, которую нужно изменить
	OldTask, err := t.tasksRepo.GetTaskById(id)
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
		err := t.tasksRepo.UpdateTitle(id, *new_title, transaction)
		if err != nil {
			return err
		}
	}

	// меняем описание
	if new_description != nil {
		err := t.tasksRepo.UpdateDescription(id, *new_description, transaction)
		if err != nil {
			return err
		}
	}

	// меняем время начала
	if new_start_at != nil {
		err := t.tasksRepo.UpdateStartAt(id, *new_start_at, transaction)
		if err != nil {
			return err
		}
	}

	// меняем время завершения
	if new_end_at != nil {
		err := t.tasksRepo.UpdateEndAt(id, *new_end_at, transaction)
		if err != nil {
			return err
		}
	}

	// меняем статус задачи
	if completed_at != nil {
		err := t.tasksRepo.TaskCompleted(id, transaction)
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

	// // ищем подзадачи для задачи
	// DeletedSubtasks, err := t.tasksRepo.Subtasks(DeletedTask.Owner_id, *DeletedTask.Id)
	// if err != nil {
	// 	return err
	// }

	// // если подзадачи есть - удаляем
	// if len(DeletedSubtasks) > 0 {
	// 	for i := range DeletedSubtasks {
	// 		St := DeletedSubtasks[i]
	// 		err := t.tasksRepo.DeleteTask(*St.Id, transaction)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// удаляем задачу с подзадачами
	err = t.tasksRepo.DeleteTask(*DeletedTask.Id, transaction)
	if err != nil {
		return err
	}

	// если ошибок нет - подтверждаем транзакцию
	transaction.Commit()
	return nil
}
