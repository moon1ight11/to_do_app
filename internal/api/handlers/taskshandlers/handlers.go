package taskshandlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"todoapp/internal/api/models"
)

// создание задачи
func (t *TasksHandler) CreateTask(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		t.logger.Error("Error in CreateTask: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		t.logger.Error("Error in CreateTask: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var task models.Task

	// получаем задачу с фронта
	if err := c.ShouldBindJSON(&task); err != nil {
		t.logger.Error("Error in CreateTask ShouldBindJSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// устанавливаем ownerId
	task.OwnerId = userId

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// добавляем задачу
	err := t.taskService.CreateTask(ctx, task)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in CreateTask:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in CreateTask:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Task created successfully by user %v", userId)

	c.JSON(http.StatusCreated, gin.H{"message": "successful"})
}

// получение списка задач пользователя
func (t *TasksHandler) GetTasks(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		t.logger.Error("Error in GetTasks: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		t.logger.Error("Error in GetTasks: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем задачи
	tasks, err := t.taskService.GetAllTasks(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.logger.Error("Error in GetTasks:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in GetTasks:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Tasks getted successfully by user %v", userId)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// получение одной задачи по id
func (t *TasksHandler) GetTaskById(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		t.logger.Error("Error in GetTaskById: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		t.logger.Error("Error in GetTaskById: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in GetTaskById:", err)
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
			t.logger.Error("Error in GetTaskById:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in GetTaskById:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("One task getted successfully by user %v", userId)

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// изменение полей задач
func (t *TasksHandler) UpdateTask(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		t.logger.Error("Error in UpdateTask: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		t.logger.Error("Error in UpdateTask: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var updatedTask models.Task
	// получаем измененную задачу с фронта
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		t.logger.Error("Error in UpdateTask ShouldBindJSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id задачи которую нужно изменить
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in UpdateTask:", err)
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
			t.logger.Error("Error in UpdateTask:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in UpdateTask:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	task, err := t.taskService.GetOneTask(ctx, taskId, userId)
	if err != nil {
		t.logger.Error("Error in UpdateTask:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Task %v updated successfully by user %v", taskId, userId)

	c.JSON(http.StatusOK, gin.H{"updated_task": task})
}

// удаление задачи
func (t *TasksHandler) DeleteTask(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		t.logger.Error("Error in DeleteTask: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		t.logger.Error("Error in DeleteTask: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		t.logger.Error("Error in parse uuid in DeleteTask:", err)
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
			t.logger.Error("Error in DeleteTask:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		t.logger.Error("Error in DeleteTask:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.logger.Info("Task %v deleted successfully by user %v", taskId, userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
