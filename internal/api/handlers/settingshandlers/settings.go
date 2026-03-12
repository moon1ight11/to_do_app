package settingshandlers

import (
	"todoapp/internal/services/settingsservice"
)

type SettingsHandler struct {
	settingsService *settingsservice.SettingsService
}

func NewSettingsHandler(settingsService *settingsservice.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
	}
}
