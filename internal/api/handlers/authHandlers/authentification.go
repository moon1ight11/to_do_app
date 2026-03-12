package authhandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services/usersservice"
)

type AuthHandler struct {
	userService *usersservice.UserService
	jwtService  jwt.TokenService
}

func NewAuthHandler(userService *usersservice.UserService, jwtService jwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}
