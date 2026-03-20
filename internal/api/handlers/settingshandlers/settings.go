package settingshandlers

import (
	"todoapp/internal/services"
	"todoapp/pkg/logger"
)

type SettingsHandler struct {
	settingsService services.SettingsServiceInterface
	logger          logger.LoggerInterface
}

func NewSettingsHandler(settingsService services.SettingsServiceInterface, logger logger.LoggerInterface) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
	}
}
