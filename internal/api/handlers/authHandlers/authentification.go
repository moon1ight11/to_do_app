package authhandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services/usersservice"
	"todoapp/pkg/logger"
)

type AuthHandler struct {
	userService *usersservice.UserService
	jwtService  jwt.TokenService
	logger      *logger.Logger
}

func NewAuthHandler(userService *usersservice.UserService, jwtService jwt.TokenService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
	}
}
