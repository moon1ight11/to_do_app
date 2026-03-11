package authhandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
	jwtService  jwt.TokenService
}

func NewAuthHandler(userService *services.UserService, jwtService jwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}