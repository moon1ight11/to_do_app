package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"todoapp/internal/services"
	"todoapp/internal/storage/repos/tasks"
)

type TasksRouter struct {
	taskService *services.TasksService
}

// создание задачи
func (t *TasksRouter) AddTask(c *gin.Context) {
	var NewTask tasks.Task

	// получаем задачу с фронта
	if err := c.ShouldBindJSON(&NewTask); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
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
func (t *TasksRouter) GetTasksList(c *gin.Context) {
	// получаем id пользователя (пока заглушка, далее - из jwt)
	user_id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	Tasks, err := t.taskService.GetTasks(user_id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": Tasks})
}

// получение одной задачи по id
func (t *TasksRouter) GetOneTask(c *gin.Context) {
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
func (t *TasksRouter) UpdateTasks(c *gin.Context) {
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
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "Task updated"})
}

// удаление задачи
func (t *TasksRouter) DeleteTask(c *gin.Context) {
	// получаем id задачи
	idStr := c.Param("id")
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
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "sucessful"})
}
