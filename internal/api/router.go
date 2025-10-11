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

func UpRouterSet(settingsHandler *handlers.SettingsHandler) {

	defaultGroup := gin.Default()
	privateGroup := defaultGroup.Group("/v1/private")

	// получение настроек пользователя
	privateGroup.GET("/settings", settingsHandler.GetSettings)

	// изменение настроек
	privateGroup.PATCH("/settings", settingsHandler.UpdateSettings)

	defaultGroup.Run(":8080")

}

func UpRouter(userHandler *handlers.UserHandler) {

	defaultGroup := gin.Default()
	privateGroup := defaultGroup.Group("/v1/private")

	// регистрация
	defaultGroup.POST("/v1/register", userHandler.AddUser)

	// авторизация
	defaultGroup.GET("/v1/auth", userHandler.CheckUser)

	// обновление пользователя
	privateGroup.PATCH("/users", userHandler.UpdateUser)

	// удаление пользователя
	privateGroup.DELETE("/users", userHandler.DeleteUser)

	defaultGroup.Run(":8080")
}
