package settingshandlers

import (
	"todoapp/internal/services"
	"todoapp/pkg/logger"
)

type SettingsHandler struct {
	settingsService services.SettingsServiceInterface
	logger          *logger.Logger
}

func NewSettingsHandler(settingsService services.SettingsServiceInterface, logger *logger.Logger) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
	}
}
