package settingshandlers

import (
	"todoapp/internal/services"
	"todoapp/internal/storage/cache"
	"todoapp/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type SettingsHandler struct {
	settingsService services.SettingsServiceInterface
	logger          logger.LoggerInterface
	cacheService    cache.CacheInterface
	tracer          trace.Tracer
}

func NewSettingsHandler(
	settingsService services.SettingsServiceInterface,
	logger logger.LoggerInterface,
	cacheService cache.CacheInterface,
	tracer trace.Tracer,
) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
		cacheService:    cacheService,
		tracer:          tracer,
	}
}
