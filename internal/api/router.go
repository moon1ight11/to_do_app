package api

import (
	"github.com/gin-gonic/gin"
	"todoapp/internal/api/handlers"
)

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
