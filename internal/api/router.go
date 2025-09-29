package api

import (
	"todoapp/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func Router() {
	settingsHandler := &routes.SettingsHandler{}

	c := gin.Default()

	// получение настроек пользователя
	c.GET("/v1/private/settings", settingsHandler.GetSettings)

	// изменение настроек
	c.PATCH("/v1/private/settings", settingsHandler.UpdateSettings)

	c.Run(":8080")

}
