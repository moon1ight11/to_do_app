package taskshandlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"
)

// создание задачи
func (t *TasksHandler) CreateTask(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in CreateTask:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var task models.Task

	// получаем задачу с фронта
	if err := c.ShouldBindJSON(&task); err != nil {
		t.logger.Error("Error in CreateTask ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// устанавливаем ownerId
	task.OwnerId = userId

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// добавляем задачу
	err = t.taskService.CreateTask(ctx, task)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in CreateTask:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in CreateTask:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Task created successfully", "user", userId)

	c.JSON(http.StatusCreated, gin.H{"message": "successful"})
}

// получение списка задач пользователя
func (t *TasksHandler) GetTasks(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in GetTasks:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем задачи
	tasks, err := t.taskService.GetAllTasks(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in GetTasks:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in GetTasks:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Tasks getted successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// получение одной задачи по id
func (t *TasksHandler) GetTaskById(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in GetTaskById:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in GetTaskById:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем задачу
	task, err := t.taskService.GetOneTask(ctx, taskId, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in GetTaskById:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in GetTaskById:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("One task getted successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// изменение полей задач
func (t *TasksHandler) UpdateTask(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in UpdateTask:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var updatedTask models.Task
	// получаем измененную задачу с фронта
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		t.logger.Error("Error in UpdateTask ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id задачи которую нужно изменить
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in UpdateTask:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// указываем, какая задача должна быть изменена
	updatedTask.Id = &taskId

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	task, err := t.taskService.GetOneTask(ctx, taskId, userId)
	if err != nil {
		t.logger.Error("Error in UpdateTask:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Task updated successfully", "task", taskId, "user", userId)

	c.JSON(http.StatusOK, gin.H{"updated_task": task})
}

// удаление задачи
func (t *TasksHandler) DeleteTask(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		t.logger.Error("Error in DeleteTask:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in DeleteTask:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// процесс удаления задачи и ее подзадач
	err = t.taskService.DeleteTask(ctx, taskId, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in DeleteTask:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in DeleteTask:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Task deleted successfully", "task", taskId, "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
