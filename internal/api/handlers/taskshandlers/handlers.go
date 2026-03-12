package taskshandlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
	"todoapp/internal/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// создание задачи
func (t *TasksHandler) CreateTask(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var task models.Task

	// получаем задачу с фронта
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Println("Error in ShouldBindJSON", err)
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
            c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
            return
        }
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sucessful"})
}

// получение списка задач пользователя
func (t *TasksHandler) GetTasks(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
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
            c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
            return
        }
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// получение одной задачи по id
func (t *TasksHandler) GetTaskById(c *gin.Context) {
	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем задачу
	task, err := t.taskService.GetOneTask(ctx, taskId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
            c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
            return
        }
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// изменение полей задач
func (t *TasksHandler) UpdateTask(c *gin.Context) {
	var updatedTask models.Task
	// получаем измененную задачу с фронта
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id задачи которую нужно изменить
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
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
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.StartAt,
		updatedTask.EndAt,
		updatedTask.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
            c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
            return
        }
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated succesfull"})
}

// удаление задачи
func (t *TasksHandler) DeleteTask(c *gin.Context) {
	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// процесс удаления задачи и ее подзадач
	err = t.taskService.DeleteTask(ctx, taskId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
            c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
            return
        }
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sucessful"})
}
