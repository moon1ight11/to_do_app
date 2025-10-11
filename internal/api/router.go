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
