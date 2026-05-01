package usershandlers

import (
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
}

func NewUserHandler(
	userService services.UsersServiceInterface,
	logger logger.LoggerInterface,
	cacheService cache.CacheInterface,
	tracer trace.Tracer,
) *UserHandler {
	return &UserHandler{
		userService:  userService,
		logger:       logger,
		cacheService: cacheService,
		tracer:       tracer,
	}
}
