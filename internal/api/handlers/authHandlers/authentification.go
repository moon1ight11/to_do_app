package authhandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services"
	"todoapp/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type AuthHandler struct {
	userService services.UsersServiceInterface
	jwtService  jwt.TokenService
	logger      logger.LoggerInterface
	tracer      trace.Tracer
}

func NewAuthHandler(
	userService services.UsersServiceInterface,
	jwtService jwt.TokenService,
	logger logger.LoggerInterface,
	tracer trace.Tracer,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
		tracer:      tracer,
	}
}
