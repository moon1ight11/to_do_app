package taskshandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// создание задачи
func (t *TasksHandler) CreateTask(c *gin.Context) {
	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in CreateTask:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var task models.Task

	// получаем задачу с фронта
	if err := c.ShouldBindJSON(&task); err != nil {
		t.logger.Error("Error in CreateTask ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// устанавливаем ownerId
	task.OwnerId = userId

	// добавляем задачу
	err = t.taskService.CreateTask(ctx, task)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in CreateTask:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in CreateTask:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// инвалидируем список задач пользователя
	if t.cacheService != nil {
		cacheKey := fmt.Sprintf("user_tasks:%s", userId.String())
		if err := t.cacheService.Delete(ctx, cacheKey); err != nil {
			t.logger.Error("Failed to invalidate tasks cache", "key", cacheKey, "error", err)
		}
	}

	t.logger.Info("Task created successfully", "user", userId)

	c.JSON(http.StatusCreated, gin.H{"message": "successful"})
}

// получение списка задач пользователя
func (t *TasksHandler) GetTasks(c *gin.Context) {
	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in GetTasks:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	cacheKey := fmt.Sprintf("user_tasks:%s", userId.String())

	// пробуем получить из кэша
	if t.cacheService != nil {
		var cachedTasks []models.Task
		err := t.cacheService.Get(ctx, cacheKey, &cachedTasks)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"tasks": cachedTasks})
			t.logger.Info("Tasks retrieved from cache", "user", userId)
			return
		}

		// если нет - идем в БД
		t.logger.Info("Cache miss for tasks", "key", cacheKey, "error", err)
	}

	// получаем задачи
	tasks, err := t.taskService.GetAllTasks(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in GetTasks:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in GetTasks:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// сохраняем список задач в кэш
	if t.cacheService != nil {
		if err := t.cacheService.Set(ctx, cacheKey, tasks, 10*time.Minute); err != nil {
			t.logger.Error("Failed to set cache", "key", cacheKey, "error", err)
		}
	}

	t.logger.Info("Tasks retrieved successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// получение одной задачи по id
func (t *TasksHandler) GetTaskById(c *gin.Context) {
	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in GetTaskById:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in GetTaskById:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// ключ для кэша
	cacheKey := fmt.Sprintf("task:%s", taskId.String())

	// пробуем найти настройки в кэше
	if t.cacheService != nil {
		var cachedTask models.Task

		err := t.cacheService.Get(ctx, cacheKey, &cachedTask)
		if err == nil {
			t.logger.Info("Task retrieved from cache", "user", userId)
			c.JSON(http.StatusOK, gin.H{"task": cachedTask})
			return
		}

		// если нет - идем в БД
		t.logger.Info("Cache miss for task", "key", cacheKey, "error", err)
	}

	// получаем задачу
	task, err := t.taskService.GetOneTask(ctx, taskId, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in GetTaskById:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in GetTaskById:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// сохраняем в кэш
	if t.cacheService != nil {
		if err := t.cacheService.Set(ctx, cacheKey, task, 10*time.Minute); err != nil {
			t.logger.Error("Failed to set cache", "key", cacheKey, "error", err)
		}
	}

	t.logger.Info("One task retrieved successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// изменение полей задач
func (t *TasksHandler) UpdateTask(c *gin.Context) {
	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in UpdateTask:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var updatedTask models.Task
	// получаем измененную задачу с фронта
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		t.logger.Error("Error in UpdateTask ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// получаем id задачи которую нужно изменить
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in UpdateTask:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// указываем, какая задача должна быть изменена
	updatedTask.Id = &taskId

	// меняем необходимые поля
	err = t.taskService.ChangeTask(
		ctx,
		*updatedTask.Id,
		userId,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.StartAt,
		updatedTask.EndAt,
		updatedTask.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in UpdateTask:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in UpdateTask:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	task, err := t.taskService.GetOneTask(ctx, taskId, userId)
	if err != nil {
		t.logger.Error("Error in UpdateTask:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if t.cacheService != nil {
		// удаляем кэш конкретно этой задачи
		taskCacheKey := fmt.Sprintf("task:%s", taskId.String())
		if err := t.cacheService.Delete(ctx, taskCacheKey); err != nil {
			t.logger.Error("Failed to invalidate task cache", "key", taskCacheKey, "error", err)
		}

		// список задач также удаляем
		userTasksCacheKey := fmt.Sprintf("user_tasks:%s", userId.String())
		if err := t.cacheService.Delete(ctx, userTasksCacheKey); err != nil {
			t.logger.Error("Failed to invalidate user tasks cache", "key", userTasksCacheKey, "error", err)
		}
	}

	t.logger.Info("Task updated successfully", "task", taskId, "user", userId)

	c.JSON(http.StatusOK, gin.H{"updated_task": task})
}

// удаление задачи
func (t *TasksHandler) DeleteTask(c *gin.Context) {
	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in DeleteTask:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in DeleteTask:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// процесс удаления задачи и ее подзадач
	err = t.taskService.DeleteTask(ctx, taskId, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in DeleteTask:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in DeleteTask:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// удаляем из кэша конкретную задачу и список задач
	if t.cacheService != nil {
		taskCacheKey := fmt.Sprintf("task:%s", taskId.String())
		if err := t.cacheService.Delete(ctx, taskCacheKey); err != nil {
			t.logger.Error("Failed to invalidate task cache", "key", taskCacheKey, "error", err)
		}

		userTasksCacheKey := fmt.Sprintf("user_tasks:%s", userId.String())
		if err := t.cacheService.Delete(ctx, userTasksCacheKey); err != nil {
			t.logger.Error("Failed to invalidate user tasks cache", "key", userTasksCacheKey, "error", err)
		}
	}

	t.logger.Info("Task deleted successfully", "task", taskId, "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
