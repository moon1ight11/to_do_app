package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"todoapp/internal/api/jwt"
	"todoapp/internal/services"
	"todoapp/internal/storage/repos/tasks"
)

type TasksHandler struct {
	taskService *services.TasksService
	jwtService  jwt.TokenService
}

func NewTasksHandler(taskService *services.TasksService, jwtService jwt.TokenService) *TasksHandler {
	return &TasksHandler{
		taskService: taskService,
		jwtService:  jwtService}
}

// создание задачи
func (t *TasksHandler) CreateTask(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var NewTask tasks.Task

	// получаем задачу с фронта
	if err := c.ShouldBindJSON(&NewTask); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// устанавливаем owner_id
	NewTask.Owner_id = UserId

	// если есть временные рамки - время на выполнение не должно быть отрицательным
	if NewTask.Start_at != nil && NewTask.End_at != nil {
		if NewTask.End_at.Before(*NewTask.Start_at) {
			log.Println("End_at cannot be before start_at")
			c.JSON((http.StatusBadRequest), gin.H{"error": "End_at cannot be before start_at"})
			return
		}
	}

	// добавляем задачу
	err := t.taskService.CreateTask(NewTask)
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
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	Tasks, err := t.taskService.GetAllTasks(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"tasks": Tasks})
}

// получение одной задачи по id
func (t *TasksHandler) GetTaskById(c *gin.Context) {
	// получаем id задачи
	idStr := c.Param("task_id")
	task_id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	task, err := t.taskService.GetOneTask(task_id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"task": task})
}

// изменение полей задач
func (t *TasksHandler) UpdateTask(c *gin.Context) {
	var UpdatedTask tasks.Task
	// получаем измененную задачу с фронта
	if err := c.ShouldBindJSON(&UpdatedTask); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// получаем id задачи которую нужно изменить
	idStr := c.Param("task_id")
	task_id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	// указываем, какая задача должна быть изменена
	UpdatedTask.Id = &task_id

	// меняем необходимые поля
	err = t.taskService.ChangeTask(*UpdatedTask.Id, UpdatedTask.Title, UpdatedTask.Description, UpdatedTask.Start_at, UpdatedTask.End_at, UpdatedTask.Completed_at)
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
	task_id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	// процесс удаления задачи и ее подзадач
	err = t.taskService.DeleteTask(task_id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "sucessful"})
}
