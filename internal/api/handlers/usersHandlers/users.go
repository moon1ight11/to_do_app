package usershandlers

import (
	"todoapp/internal/metrics"
	"todoapp/internal/services"
	"todoapp/internal/storage/cache"
	"todoapp/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type UserHandler struct {
	userService  services.UsersServiceInterface
	logger       logger.LoggerInterface
	cacheService cache.CacheInterface
	tracer       trace.Tracer
	metrics      *metrics.Metrics
}

func NewUserHandler(
	userService services.UsersServiceInterface,
	logger logger.LoggerInterface,
	cacheService cache.CacheInterface,
	tracer trace.Tracer,
	metrics *metrics.Metrics,
) *UserHandler {
	return &UserHandler{
		userService:  userService,
		logger:       logger,
		cacheService: cacheService,
		tracer:       tracer,
		metrics:      metrics,
	}
}
