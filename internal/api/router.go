package api

import (
	"todoapp/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func Router() {
	settingsHandler := &routes.SettingsHandler{}

	c := gin.Default()

	c.GET("/settings/:id", settingsHandler.GetSettings)
	c.POST("/update_settings", settingsHandler.UpdateSettings)

	c.Run(":8080")

}
