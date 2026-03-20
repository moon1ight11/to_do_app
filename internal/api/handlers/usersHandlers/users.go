package usershandlers

import (
	"todoapp/internal/services"
	"todoapp/pkg/logger"
)

type UserHandler struct {
	userService services.UsersServiceInterface
	logger      logger.LoggerInterface
}

func NewUserHandler(userService services.UsersServiceInterface, logger logger.LoggerInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}
