package settingshandlers

import (
	"todoapp/internal/services/settingsservice"
	"todoapp/pkg/logger"
)

type SettingsHandler struct {
	settingsService *settingsservice.SettingsService
	logger          *logger.Logger
}

func NewSettingsHandler(settingsService *settingsservice.SettingsService, logger *logger.Logger) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
	}
}
