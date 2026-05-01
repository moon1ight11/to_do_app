package tasksservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"todoapp/internal/api/models"
	"todoapp/internal/storage/repos/tasksrepos"
)

// создание задачи
func (t *TasksService) CreateTask(ctx context.Context, task models.Task) error {
	ctx, span := t.tracer.Start(ctx, "service.CreateTask")
	defer span.End()

	// если есть временные рамки - время на выполнение не должно быть отрицательным
	if task.StartAt != nil && task.EndAt != nil {
		if task.EndAt.Before(*task.StartAt) {
			span.RecordError(fmt.Errorf("error in CreateTask: end_at cannot be before start_at"))
			return fmt.Errorf("error in CreateTask: end_at cannot be before start_at")
		}
	}

	// приводим тип
	taskrepo := taskFromApiToRepo(&task)

	// добавляем задачу в репозиторий
	err := t.tasksRepo.CreateTask(ctx, *taskrepo)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in CreateTask: %w", err)
	}

	return nil
}

// получение списка всех задач пользователя
func (t *TasksService) GetAllTasks(ctx context.Context, userId uuid.UUID) ([]models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "service.GetAllTasks")
	defer span.End()

	// получаем список всех задач пользователя
	tasks, err := t.tasksRepo.TasksByOwnerId(ctx, userId)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("error in GetAllTasks: %w", err)
	}

	// рекурсивно загружаем подзадачи для каждой родительской задачи
	var tasksApi []models.Task
	for i := range tasks {
		taskApi := t.buildTaskWithSubtasks(ctx, &tasks[i], userId)
		tasksApi = append(tasksApi, taskApi)
	}

	return tasksApi, nil
}

// получение одной задачи
func (t *TasksService) GetOneTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "service.GetOneTask")
	defer span.End()

	// получаем задачу
	task, err := t.tasksRepo.TaskById(ctx, taskId, userId)
	if err != nil {
		span.RecordError(err)
		return models.Task{}, fmt.Errorf("error in GetOneTask: %w", err)
	}

	ownerId := task.OwnerId
	if ownerId == uuid.Nil {
		span.RecordError(fmt.Errorf("error in GetOneTask: ownerId is nil"))
		return models.Task{}, fmt.Errorf("error in GetOneTask: ownerId is nil")
	}

	taskApi := t.buildTaskWithSubtasks(ctx, &task, ownerId)

	return taskApi, nil
}

// изменение полей задачи
func (t *TasksService) ChangeTask(
	ctx context.Context,
	taskId uuid.UUID,
	userId uuid.UUID,
	title *string,
	description *string,
	startAt *time.Time,
	endAt *time.Time,
	completedAt *bool,
) error {
	ctx, span := t.tracer.Start(ctx, "service.ChangeTask")
	defer span.End()

	// запускаем транзакцию
	transaction, err := t.tasksRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in ChangeTask BeginTx: %w", err)
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// находим задачу, которую нужно изменить
	foundTask, err := t.tasksRepo.TaskById(ctx, taskId, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in ChangeTask: %w", err)
	}

	// проверяем, чтобы время старта было раньше времени конца
	// если меняются оба времени
	if endAt != nil && startAt != nil {
		if endAt.Before(*startAt) {
			span.RecordError(fmt.Errorf("error in ChangeTask: end_at cannot be before start_at"))
			return fmt.Errorf("error in ChangeTask: end_at cannot be before start_at")
		}
	}

	// если меняется толко время начала
	if endAt == nil && startAt != nil && foundTask.EndAt != nil {
		if foundTask.EndAt.Before(*startAt) {
			span.RecordError(fmt.Errorf("error in ChangeTask: end_at cannot be before start_at"))
			return fmt.Errorf("error in ChangeTask: end_at cannot be before start_at")
		}
	}

	// если меняется только конец
	if startAt == nil && endAt != nil && foundTask.StartAt != nil {
		if endAt.Before(*foundTask.StartAt) {
			span.RecordError(fmt.Errorf("error in ChangeTask: end_at cannot be before start_at"))
			return fmt.Errorf("error in ChangeTask: end_at cannot be before start_at")
		}
	}

	// меняем название
	if title != nil {
		err := t.tasksRepo.UpdateTitle(ctx, taskId, userId, *title, transaction)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in ChangeTask: %w", err)
		}
	}

	// меняем описание
	if description != nil {
		err := t.tasksRepo.UpdateDescription(ctx, taskId, userId, *description, transaction)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in ChangeTask: %w", err)
		}
	}

	// меняем время начала
	if startAt != nil {
		err := t.tasksRepo.UpdateStartAt(ctx, taskId, userId, *startAt, transaction)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in ChangeTask: %w", err)
		}
	}

	// меняем время завершения
	if endAt != nil {
		err := t.tasksRepo.UpdateEndAt(ctx, taskId, userId, *endAt, transaction)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in ChangeTask: %w", err)
		}
	}

	// меняем статус задачи
	if completedAt != nil {
		if *completedAt == false {
			err := t.tasksRepo.TaskUncompleted(ctx, taskId, userId, transaction)
			if err != nil {
				span.RecordError(err)
				return fmt.Errorf("error in ChangeTask: %w", err)
			}
		} else if *completedAt == true {
			err := t.tasksRepo.TaskCompleted(ctx, taskId, userId, transaction)
			if err != nil {
				span.RecordError(err)
				return fmt.Errorf("error in ChangeTask: %w", err)
			}
		}
	}

	// если ошибок нет - подтверждаем транзакцию
	err = transaction.Commit()
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in ChangeTask commit: %w", err)
	}

	return nil
}

// удаление задачи
func (t *TasksService) DeleteTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
	ctx, span := t.tracer.Start(ctx, "service.DeleteTask")
	defer span.End()

	// запускаем транзакцию
	transaction, err := t.tasksRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteTask BeginTx: %w", err)
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// ищем задачу по id
	deletedTask, err := t.tasksRepo.TaskById(ctx, taskId, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteTask: %w", err)
	}

	// удаляем задачу с подзадачами
	err = t.tasksRepo.DeleteTask(ctx, *deletedTask.Id, userId, transaction)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteTask: %w", err)
	}

	// если ошибок нет - подтверждаем транзакцию
	err = transaction.Commit()
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteTask commit: %w", err)
	}
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

// перевод типа из репо в апи
func taskFromRepoToApi(repoTask *tasksrepos.Task) models.Task {
	if repoTask == nil {
		return models.Task{}
	}

	var subtasks *[]models.Task
	if repoTask.Subtasks != nil {
		mappedSubtasks := make([]models.Task, len(*repoTask.Subtasks))
		for i, t := range *repoTask.Subtasks {
			mappedSubtasks[i] = taskFromRepoToApi(&t)
		}
		subtasks = &mappedSubtasks
	}

	return models.Task{
		Id:          repoTask.Id,
		ParentId:    repoTask.ParentId,
		OwnerId:     repoTask.OwnerId,
		StartAt:     repoTask.StartAt,
		EndAt:       repoTask.EndAt,
		Title:       repoTask.Title,
		Description: repoTask.Description,
		CompletedAt: repoTask.CompletedAt,
		Subtasks:    subtasks,
	}
}

// конструктор для рекурсии подзадач
func (t *TasksService) buildTaskWithSubtasks(ctx context.Context, task *tasksrepos.Task, userId uuid.UUID) models.Task {
	// получаем прямые подзадачи текущей задачи
	subtasks, err := t.tasksRepo.SubtasksByTaskId(ctx, userId, *task.Id)
	if err != nil {
		return taskFromRepoToApi(task)
	}

	// рекурсивно загружаем подзадачи для каждой подзадачи
	var subtasksApi []models.Task
	for i := range subtasks {
		subtaskApi := t.buildTaskWithSubtasks(ctx, &subtasks[i], userId)
		subtasksApi = append(subtasksApi, subtaskApi)
	}

	// конвертируем основную задачу
	taskApi := taskFromRepoToApi(task)

	// добавляем подзадачи
	if len(subtasksApi) > 0 {
		taskApi.Subtasks = &subtasksApi
	}

	return taskApi
}
