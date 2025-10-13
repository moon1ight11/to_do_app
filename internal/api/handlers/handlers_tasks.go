package handlers

import (
	"log"
	"net/http"
	"todoapp/internal/services"
	"todoapp/internal/storage/repos/tasks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TasksHandler struct {
	taskService *services.TasksService
}

func NewTasksHandler(taskService *services.TasksService) *TasksHandler {
	return &TasksHandler{taskService: taskService}
}

// создание задачи
func (t *TasksHandler) AddTask(c *gin.Context) {
	var NewTask tasks.Task

	// получаем задачу с фронта
	if err := c.ShouldBindJSON(&NewTask); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// owner_id получаем из куков
	NewTask.Owner_id, _ = uuid.Parse("b6609ddd-95f4-42f0-993e-7f07f3fa1d5b")

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
func (t *TasksHandler) GetTasksList(c *gin.Context) {
	// получаем id пользователя (пока заглушка, далее - из jwt)
	user_id, _ := uuid.Parse("8b1bbae9-6e4d-41dc-984f-3a4a6c0abb17")

	Tasks, err := t.taskService.GetTasks(user_id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"tasks": Tasks})
}

// получение одной задачи по id
func (t *TasksHandler) GetOneTask(c *gin.Context) {
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
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"task": task})
}

// изменение полей задач
func (t *TasksHandler) UpdateTasks(c *gin.Context) {
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

	// из куков получаем id пользователя
	UpdatedTask.Owner_id, _ = uuid.Parse("b6609ddd-95f4-42f0-993e-7f07f3fa1d5b")

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
