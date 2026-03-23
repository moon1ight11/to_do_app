package usershandlers

import (
	"todoapp/internal/services"
	"todoapp/internal/storage/cache"
	"todoapp/pkg/logger"
)

type UserHandler struct {
	userService  services.UsersServiceInterface
	logger       logger.LoggerInterface
	cacheService cache.CacheInterface
}

func NewUserHandler(userService services.UsersServiceInterface, logger logger.LoggerInterface, cacheService cache.CacheInterface) *UserHandler {
	return &UserHandler{
		userService:  userService,
		logger:       logger,
		cacheService: cacheService,
	}
}
