package usershandlers

import (
	"todoapp/internal/services"
	"todoapp/pkg/logger"
)

type UserHandler struct {
	userService services.UsersServiceInterface
	logger      *logger.Logger
}

func NewUserHandler(userService services.UsersServiceInterface, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}
