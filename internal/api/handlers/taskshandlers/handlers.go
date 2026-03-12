package taskshandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"todoapp/internal/api/models"
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
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// устанавливаем owner_id
	task.OwnerId = userId

	// добавляем задачу
	err := t.taskService.CreateTask(task)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "sucessful"})
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

	tasks, err := t.taskService.GetAllTasks(userId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"tasks": tasks})
}

// получение одной задачи по id
func (t *TasksHandler) GetTaskById(c *gin.Context) {
	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	task, err := t.taskService.GetOneTask(taskId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"task": task})
}

// изменение полей задач
func (t *TasksHandler) UpdateTask(c *gin.Context) {
	var updatedTask models.Task
	// получаем измененную задачу с фронта
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// получаем id задачи которую нужно изменить
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	// указываем, какая задача должна быть изменена
	updatedTask.Id = &taskId

	// меняем необходимые поля
	err = t.taskService.ChangeTask(
		*updatedTask.Id,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.StartAt,
		updatedTask.EndAt,
		updatedTask.CompletedAt,
	)

	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "Task updated succesfull"})
}

// удаление задачи
func (t *TasksHandler) DeleteTask(c *gin.Context) {
	// получаем id задачи
	idStr := c.Param("task_id")
	taskId, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	// процесс удаления задачи и ее подзадач
	err = t.taskService.DeleteTask(taskId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "sucessful"})
}
