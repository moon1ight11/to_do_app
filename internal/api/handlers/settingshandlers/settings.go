package settingshandlers

import (
	"todoapp/internal/metrics"
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
	metrics         *metrics.Metrics
}

func NewSettingsHandler(
	settingsService services.SettingsServiceInterface,
	logger logger.LoggerInterface,
	cacheService cache.CacheInterface,
	tracer trace.Tracer,
	metrics *metrics.Metrics,
) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
		cacheService:    cacheService,
		tracer:          tracer,
		metrics:         metrics,
	}
}
