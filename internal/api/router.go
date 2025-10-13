package api

import (
	"github.com/gin-gonic/gin"
	"todoapp/internal/api/handlers"
)

func UpRouter(userHandler *handlers.UserHandler, settingsHandler *handlers.SettingsHandler, tasksHandler *handlers.TasksHandler) {
	// группировка роутов
	defaultGroup := gin.Default()
	privateGroup := defaultGroup.Group("/v1/private")

	// ЮЗЕРЫ //

	// регистрация
	defaultGroup.POST("/v1/register", userHandler.AddUser)

	// авторизация
	defaultGroup.GET("/v1/auth", userHandler.CheckUser)

	// обновление пользователя
	privateGroup.PATCH("/users", userHandler.UpdateUser)

	// удаление пользователя
	privateGroup.DELETE("/users", userHandler.DeleteUser)

	// НАСТРОЙКИ //

	// получение настроек пользователя
	privateGroup.GET("/settings", settingsHandler.GetSettings)

	// изменение настроек
	privateGroup.PATCH("/settings", settingsHandler.UpdateSettings)

	// ЗАДАЧИ //

	// создание задачи
	privateGroup.POST("/tasks", tasksHandler.AddTask)

	// получение всех задач пользователя
	privateGroup.GET("/tasks", tasksHandler.GetTasksList)

	// получение одной задачи по id
	privateGroup.GET("/tasks/:task_id", tasksHandler.GetOneTask)

	// изменение задачи
	privateGroup.PATCH("/tasks/:task_id", tasksHandler.UpdateTasks)

	// удаление задачи
	privateGroup.DELETE("/tasks/:task_id", tasksHandler.DeleteTask)

	defaultGroup.Run(":8080")
}
