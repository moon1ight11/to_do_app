package settingshandlers

import (
	"todoapp/internal/services"
	"todoapp/internal/storage/cache"
	"todoapp/pkg/logger"
)

type SettingsHandler struct {
	settingsService services.SettingsServiceInterface
	logger          logger.LoggerInterface
	cacheService    cache.CacheInterface
}

func NewSettingsHandler(settingsService services.SettingsServiceInterface, logger logger.LoggerInterface, cacheService cache.CacheInterface) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
		cacheService:    cacheService,
	}
}
