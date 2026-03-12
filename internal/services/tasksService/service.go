package tasksservice

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"todoapp/internal/api/models"
	"todoapp/internal/storage/repos/tasksrepos"
)

// создание задачи
func (t *TasksService) CreateTask(task models.Task) error {
	// если есть временные рамки - время на выполнение не должно быть отрицательным
	if task.StartAt != nil && task.EndAt != nil {
		if task.EndAt.Before(*task.StartAt) {
			return fmt.Errorf("End_at cannot be before start_at")
		}
	}

	// приводим тип
	taskrepo := taskFromApiToRepo(&task)

	// добавляем задачу в репозиторий
	err := t.tasksRepo.CreateTask(*taskrepo)
	if err != nil {
		return err
	}

	return nil
}

// получение списка всех задач пользователя
func (t *TasksService) GetAllTasks(userId uuid.UUID) ([]tasksrepos.Task, error) {
	// получаем список всех задач пользователя
	tasks, err := t.tasksRepo.TasksByOwnerId(userId)
	if err != nil {
		return nil, err
	}

	// среди них ищем подзадачи
	for i := range tasks {
		subtasks, err := t.tasksRepo.SubtasksByTaskId(userId, *tasks[i].Id)
		if err != nil {
			return nil, err
		}

		tasks[i].Subtasks = &subtasks
	}

	return tasks, err
}

// получение одной задачи
func (t *TasksService) GetOneTask(taskId uuid.UUID) (tasksrepos.Task, error) {
	task, err := t.tasksRepo.TaskById(taskId)
	if err != nil {
		return tasksrepos.Task{}, err
	}

	return task, nil
}

// изменение полей задачи
func (t *TasksService) ChangeTask(
	taskId uuid.UUID,
	title *string,
	description *string,
	startAt *time.Time,
	endAt *time.Time,
	completedAt *bool,
) error {
	// запускаем транзакцию
	transaction, err := t.tasksRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// находим задачу, которую нужно изменить
	foundTask, err := t.tasksRepo.TaskById(taskId)
	if err != nil {
		return err
	}

	// проверяем, чтобы время старта было раньше времени конца
	// если меняются оба времени
	if endAt != nil && startAt != nil {
		if endAt.Before(*startAt) {
			return fmt.Errorf("end_at cannot be before start_at")
		}
	}

	// если меняется толко время начала
	if endAt == nil && startAt != nil && foundTask.EndAt != nil {
		if foundTask.EndAt.Before(*startAt) {
			return fmt.Errorf("end_at cannot be before start_at")
		}
	}

	// если меняется только конец
	if startAt == nil && endAt != nil && foundTask.StartAt != nil {
		if endAt.Before(*foundTask.StartAt) {
			return fmt.Errorf("end_at cannot be before start_at")
		}
	}

	// меняем название
	if title != nil {
		err := t.tasksRepo.UpdateTitle(taskId, *title, transaction)
		if err != nil {
			return err
		}
	}

	// меняем описание
	if description != nil {
		err := t.tasksRepo.UpdateDescription(taskId, *description, transaction)
		if err != nil {
			return err
		}
	}

	// меняем время начала
	if startAt != nil {
		err := t.tasksRepo.UpdateStartAt(taskId, *startAt, transaction)
		if err != nil {
			return err
		}
	}

	// меняем время завершения
	if endAt != nil {
		err := t.tasksRepo.UpdateEndAt(taskId, *endAt, transaction)
		if err != nil {
			return err
		}
	}

	// меняем статус задачи
	if completedAt != nil {
		if *completedAt == false {
			err := t.tasksRepo.TaskUncompleted(taskId, transaction)
			if err != nil {
				return err
			}
		} else if *completedAt == true {
			err := t.tasksRepo.TaskCompleted(taskId, transaction)
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

// перевод типа из апи в репо
func taskFromApiToRepo(modelTask *models.Task) *tasksrepos.Task {
	if modelTask == nil {
		return nil
	}

	var subtasks *[]tasksrepos.Task
	if modelTask.Subtasks != nil {
		mappedSubtasks := make([]tasksrepos.Task, len(*modelTask.Subtasks))
		for i, t := range *modelTask.Subtasks {
			mappedSubtasks[i] = *taskFromApiToRepo(&t)
		}
		subtasks = &mappedSubtasks
	}

	return &tasksrepos.Task{
		Id:          modelTask.Id,
		ParentId:    modelTask.ParentId,
		OwnerId:     modelTask.OwnerId,
		StartAt:     modelTask.StartAt,
		EndAt:       modelTask.EndAt,
		Title:       modelTask.Title,
		Description: modelTask.Description,
		CompletedAt: modelTask.CompletedAt,
		Subtasks:    subtasks,
	}
}
