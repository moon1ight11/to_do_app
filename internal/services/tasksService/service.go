package tasksservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
	"todoapp/internal/storage/repos/tasksrepos"
)

func (t *TasksService) CreateTask(ctx context.Context, task models.Task) error {
	ctx, span := t.tracer.Start(ctx, "service.CreateTask")
	defer span.End()

	if task.StartAt != nil && task.EndAt != nil {
		if task.EndAt.Before(*task.StartAt) {
			err := fmt.Errorf("end_at before start_at")
			span.RecordError(err)
			return fmt.Errorf("tasksservice.CreateTask: %w", err)
		}
	}

	taskrepo := taskFromApiToRepo(&task)

	if err := t.tasksRepo.CreateTask(ctx, *taskrepo); err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.CreateTask: %w", err)
	}

	return nil
}

func (t *TasksService) GetAllTasks(ctx context.Context, userId uuid.UUID) ([]models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "service.GetAllTasks")
	defer span.End()

	tasks, err := t.tasksRepo.TasksByOwnerId(ctx, userId)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("tasksservice.GetAllTasks: %w", err)
	}

	var tasksApi []models.Task
	for i := range tasks {
		taskApi := t.buildTaskWithSubtasks(ctx, &tasks[i], userId)
		tasksApi = append(tasksApi, taskApi)
	}

	return tasksApi, nil
}

func (t *TasksService) GetOneTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "service.GetOneTask")
	defer span.End()

	task, err := t.tasksRepo.TaskById(ctx, taskId, userId)
	if err != nil {
		span.RecordError(err)
		return models.Task{}, fmt.Errorf("tasksservice.GetOneTask: %w", err)
	}

	if task.OwnerId == uuid.Nil {
		err := fmt.Errorf("ownerId is nil")
		span.RecordError(err)
		return models.Task{}, fmt.Errorf("tasksservice.GetOneTask: %w", err)
	}

	return t.buildTaskWithSubtasks(ctx, &task, task.OwnerId), nil
}

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

	transaction, err := t.tasksRepo.DB().BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.ChangeTask: begin tx: %w", err)
	}
	defer transaction.Rollback()

	foundTask, err := t.tasksRepo.TaskById(ctx, taskId, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.ChangeTask: find task: %w", err)
	}

	if err := validateTimeRange(startAt, endAt, foundTask.StartAt, foundTask.EndAt); err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.ChangeTask: %w", err)
	}

	if title != nil {
		if err := t.tasksRepo.UpdateTitle(ctx, taskId, userId, *title, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("tasksservice.ChangeTask: update title: %w", err)
		}
	}

	if description != nil {
		if err := t.tasksRepo.UpdateDescription(ctx, taskId, userId, *description, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("tasksservice.ChangeTask: update description: %w", err)
		}
	}

	if startAt != nil {
		if err := t.tasksRepo.UpdateStartAt(ctx, taskId, userId, *startAt, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("tasksservice.ChangeTask: update start_at: %w", err)
		}
	}

	if endAt != nil {
		if err := t.tasksRepo.UpdateEndAt(ctx, taskId, userId, *endAt, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("tasksservice.ChangeTask: update end_at: %w", err)
		}
	}

	if completedAt != nil {
		if *completedAt {
			if err := t.tasksRepo.TaskCompleted(ctx, taskId, userId, transaction); err != nil {
				span.RecordError(err)
				return fmt.Errorf("tasksservice.ChangeTask: complete: %w", err)
			}
		} else {
			if err := t.tasksRepo.TaskUncompleted(ctx, taskId, userId, transaction); err != nil {
				span.RecordError(err)
				return fmt.Errorf("tasksservice.ChangeTask: uncomplete: %w", err)
			}
		}
	}

	if err := transaction.Commit(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.ChangeTask: commit: %w", err)
	}

	return nil
}

func (t *TasksService) DeleteTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
	ctx, span := t.tracer.Start(ctx, "service.DeleteTask")
	defer span.End()

	transaction, err := t.tasksRepo.DB().BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.DeleteTask: begin tx: %w", err)
	}
	defer transaction.Rollback()

	deletedTask, err := t.tasksRepo.TaskById(ctx, taskId, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.DeleteTask: find task: %w", err)
	}

	if err := t.tasksRepo.DeleteTask(ctx, *deletedTask.Id, userId, transaction); err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.DeleteTask: delete: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("tasksservice.DeleteTask: commit: %w", err)
	}

	return nil
}

func validateTimeRange(startAt, endAt, existingStart, existingEnd *time.Time) error {
	if endAt != nil && startAt != nil {
		if endAt.Before(*startAt) {
			return fmt.Errorf("end_at before start_at")
		}
	}

	if endAt == nil && startAt != nil && existingEnd != nil {
		if existingEnd.Before(*startAt) {
			return fmt.Errorf("end_at before new start_at")
		}
	}

	if startAt == nil && endAt != nil && existingStart != nil {
		if endAt.Before(*existingStart) {
			return fmt.Errorf("new end_at before start_at")
		}
	}

	return nil
}

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

func (t *TasksService) buildTaskWithSubtasks(ctx context.Context, task *tasksrepos.Task, userId uuid.UUID) models.Task {
	subtasks, err := t.tasksRepo.SubtasksByTaskId(ctx, userId, *task.Id)
	if err != nil {
		return taskFromRepoToApi(task)
	}

	var subtasksApi []models.Task
	for i := range subtasks {
		subtaskApi := t.buildTaskWithSubtasks(ctx, &subtasks[i], userId)
		subtasksApi = append(subtasksApi, subtaskApi)
	}

	taskApi := taskFromRepoToApi(task)
	if len(subtasksApi) > 0 {
		taskApi.Subtasks = &subtasksApi
	}

	return taskApi
}
