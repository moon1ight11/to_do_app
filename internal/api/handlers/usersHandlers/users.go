package usershandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services"
)

type UserHandler struct {
	userService *services.UserService
	jwtService  jwt.TokenService
}

func NewUserHandler(userService *services.UserService, jwtService jwt.TokenService) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}