package usershandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services/usersservice"
)

type UserHandler struct {
	userService *usersservice.UserService
	jwtService  jwt.TokenService
}

func NewUserHandler(userService *usersservice.UserService, jwtService jwt.TokenService) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}