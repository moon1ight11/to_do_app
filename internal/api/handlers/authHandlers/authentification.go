package authhandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/metrics"
	"todoapp/internal/services"
	"todoapp/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type AuthHandler struct {
	userService services.UsersServiceInterface
	jwtService  jwt.TokenService
	logger      logger.LoggerInterface
	tracer      trace.Tracer
	metrics     *metrics.Metrics
}

func NewAuthHandler(
	userService services.UsersServiceInterface,
	jwtService jwt.TokenService,
	logger logger.LoggerInterface,
	tracer trace.Tracer,
	metrics *metrics.Metrics,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
		tracer:      tracer,
		metrics:     metrics,
	}
}
