package taskshandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"
	"todoapp/internal/metrics"
)

func (t *TasksHandler) CreateTask(c *gin.Context) {
	ctx, span := t.tracer.Start(c.Request.Context(), "handler.CreateTask")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrForbidden), "CreateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.CreateTask: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		t.metrics.RecordError(string(metrics.ErrBadRequest), "CreateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.CreateTask: bind json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task.OwnerId = userId

	if err := t.taskService.CreateTask(ctx, task); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.metrics.RecordError(string(metrics.ErrInternal), "CreateTask")
			span.RecordError(err)
			t.logger.Error("taskshandlers.CreateTask: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.metrics.RecordError(string(metrics.ErrInternal), "CreateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.CreateTask: create", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if t.cacheService != nil {
		cacheKey := fmt.Sprintf("user_tasks:%s", userId.String())
		if err := t.cacheService.Delete(ctx, cacheKey); err != nil {
			t.logger.Error("taskshandlers.CreateTask: cache delete", "key", cacheKey, "error", err)
		}
	}

	t.logger.Info("taskshandlers.CreateTask: success", "user", userId)

	c.JSON(http.StatusCreated, gin.H{"message": "successful"})
}

func (t *TasksHandler) GetTasks(c *gin.Context) {
	ctx, span := t.tracer.Start(c.Request.Context(), "handler.GetTasks")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrForbidden), "GetTasks")
		span.RecordError(err)
		t.logger.Error("taskshandlers.GetTasks: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	cacheKey := fmt.Sprintf("user_tasks:%s", userId.String())

	if t.cacheService != nil {
		var cachedTasks []models.Task
		if err := t.cacheService.Get(ctx, cacheKey, &cachedTasks); err == nil {
			t.logger.Info("taskshandlers.GetTasks: cache hit", "user", userId)
			c.JSON(http.StatusOK, gin.H{"tasks": cachedTasks})
			return
		}
		t.logger.Info("taskshandlers.GetTasks: cache miss", "key", cacheKey)
	}

	tasks, err := t.taskService.GetAllTasks(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.metrics.RecordError(string(metrics.ErrInternal), "GetTasks")
			span.RecordError(err)
			t.logger.Error("taskshandlers.GetTasks: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.metrics.RecordError(string(metrics.ErrInternal), "GetTasks")
		span.RecordError(err)
		t.logger.Error("taskshandlers.GetTasks: get all", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if t.cacheService != nil {
		if err := t.cacheService.Set(ctx, cacheKey, tasks, 10*time.Minute); err != nil {
			t.logger.Error("taskshandlers.GetTasks: cache set", "key", cacheKey, "error", err)
		}
	}

	t.logger.Info("taskshandlers.GetTasks: success", "user", userId)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (t *TasksHandler) GetTaskById(c *gin.Context) {
	ctx, span := t.tracer.Start(c.Request.Context(), "handler.GetTaskById")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrForbidden), "GetTaskById")
		span.RecordError(err)
		t.logger.Error("taskshandlers.GetTaskById: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrBadRequest), "GetTaskById")
		span.RecordError(err)
		t.logger.Error("taskshandlers.GetTaskById: parse uuid", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	cacheKey := fmt.Sprintf("task:%s", taskId.String())

	if t.cacheService != nil {
		var cachedTask models.Task
		if err := t.cacheService.Get(ctx, cacheKey, &cachedTask); err == nil {
			t.logger.Info("taskshandlers.GetTaskById: cache hit", "user", userId)
			c.JSON(http.StatusOK, gin.H{"task": cachedTask})
			return
		}
		t.logger.Info("taskshandlers.GetTaskById: cache miss", "key", cacheKey)
	}

	task, err := t.taskService.GetOneTask(ctx, taskId, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.metrics.RecordError(string(metrics.ErrInternal), "GetTaskById")
			span.RecordError(err)
			t.logger.Error("taskshandlers.GetTaskById: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.metrics.RecordError(string(metrics.ErrInternal), "GetTaskById")
		span.RecordError(err)
		t.logger.Error("taskshandlers.GetTaskById: get one", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if t.cacheService != nil {
		if err := t.cacheService.Set(ctx, cacheKey, task, 10*time.Minute); err != nil {
			t.logger.Error("taskshandlers.GetTaskById: cache set", "key", cacheKey, "error", err)
		}
	}

	t.logger.Info("taskshandlers.GetTaskById: success", "user", userId)

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (t *TasksHandler) UpdateTask(c *gin.Context) {
	ctx, span := t.tracer.Start(c.Request.Context(), "handler.UpdateTask")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrForbidden), "UpdateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.UpdateTask: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		t.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.UpdateTask: bind json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.UpdateTask: parse uuid", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	updatedTask.Id = &taskId

	if err := t.taskService.ChangeTask(
		ctx,
		*updatedTask.Id,
		userId,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.StartAt,
		updatedTask.EndAt,
		updatedTask.CompletedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.metrics.RecordError(string(metrics.ErrInternal), "UpdateTask")
			span.RecordError(err)
			t.logger.Error("taskshandlers.UpdateTask: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.metrics.RecordError(string(metrics.ErrInternal), "UpdateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.UpdateTask: change", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	task, err := t.taskService.GetOneTask(ctx, taskId, userId)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrInternal), "UpdateTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.UpdateTask: get one after update", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if t.cacheService != nil {
		taskCacheKey := fmt.Sprintf("task:%s", taskId.String())
		if err := t.cacheService.Delete(ctx, taskCacheKey); err != nil {
			t.logger.Error("taskshandlers.UpdateTask: cache delete task", "key", taskCacheKey, "error", err)
		}

		userTasksCacheKey := fmt.Sprintf("user_tasks:%s", userId.String())
		if err := t.cacheService.Delete(ctx, userTasksCacheKey); err != nil {
			t.logger.Error("taskshandlers.UpdateTask: cache delete user tasks", "key", userTasksCacheKey, "error", err)
		}
	}

	t.logger.Info("taskshandlers.UpdateTask: success", "task", taskId, "user", userId)

	c.JSON(http.StatusOK, gin.H{"updated_task": task})
}

func (t *TasksHandler) DeleteTask(c *gin.Context) {
	ctx, span := t.tracer.Start(c.Request.Context(), "handler.DeleteTask")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrForbidden), "DeleteTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.DeleteTask: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.metrics.RecordError(string(metrics.ErrBadRequest), "DeleteTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.DeleteTask: parse uuid", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := t.taskService.DeleteTask(ctx, taskId, userId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.metrics.RecordError(string(metrics.ErrInternal), "DeleteTask")
			span.RecordError(err)
			t.logger.Error("taskshandlers.DeleteTask: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.metrics.RecordError(string(metrics.ErrInternal), "DeleteTask")
		span.RecordError(err)
		t.logger.Error("taskshandlers.DeleteTask: delete", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if t.cacheService != nil {
		taskCacheKey := fmt.Sprintf("task:%s", taskId.String())
		if err := t.cacheService.Delete(ctx, taskCacheKey); err != nil {
			t.logger.Error("taskshandlers.DeleteTask: cache delete task", "key", taskCacheKey, "error", err)
		}

		userTasksCacheKey := fmt.Sprintf("user_tasks:%s", userId.String())
		if err := t.cacheService.Delete(ctx, userTasksCacheKey); err != nil {
			t.logger.Error("taskshandlers.DeleteTask: cache delete user tasks", "key", userTasksCacheKey, "error", err)
		}
	}

	t.logger.Info("taskshandlers.DeleteTask: success", "task", taskId, "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
