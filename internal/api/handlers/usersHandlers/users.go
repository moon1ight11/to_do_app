package usershandlers

import (
	"todoapp/internal/services/usersservice"
	"todoapp/pkg/logger"
)

type UserHandler struct {
	userService *usersservice.UserService
	logger      *logger.Logger
}

func NewUserHandler(userService *usersservice.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}
