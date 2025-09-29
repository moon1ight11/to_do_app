package api

import (
	"todoapp/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func Router() {
	tasksRouter := &routes.TasksRouter{}

	c := gin.Default()

	// создание задачи
	c.POST("/v1/private/tasks", tasksRouter.AddTask)

	// получение всех задач пользователя
	c.GET("/v1/private/tasks", tasksRouter.GetTasksList)

	// получение одной задачи по id
	c.GET("/v1/private/tasks/:task_id", tasksRouter.GetOneTask)

	// изменение задачи
	c.PATCH("/v1/private/tasks/:task_id", tasksRouter.UpdateTasks)

	// удаление задачи
	c.DELETE("/v1/private/tasks/:task_id", tasksRouter.DeleteTask)


	c.Run(":8080")
} 
