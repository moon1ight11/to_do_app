package api

import (
	"todoapp/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func Router() {
	tasksRouter := &routes.TasksRouter{}

	c := gin.Default()

	c.POST("/create_task", tasksRouter.AddTask)
	c.GET("/tasks/:id", tasksRouter.GetTasksList)
	c.PATCH("/tasks/change", tasksRouter.UpdateTasks)
	c.DELETE("tasks/delete/:id", tasksRouter.DeleteTask)


	c.Run(":8080")
}
