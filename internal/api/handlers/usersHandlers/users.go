package usershandlers

import (
	"todoapp/internal/services/usersservice"
)

type UserHandler struct {
	userService *usersservice.UserService
}

func NewUserHandler(userService *usersservice.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
