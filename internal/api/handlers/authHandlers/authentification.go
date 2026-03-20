package authhandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services"
	"todoapp/pkg/logger"
)

type AuthHandler struct {
	userService services.UsersServiceInterface
	jwtService  jwt.TokenService
	logger      logger.LoggerInterface
}

func NewAuthHandler(userService services.UsersServiceInterface, jwtService jwt.TokenService, logger logger.LoggerInterface) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
	}
}
