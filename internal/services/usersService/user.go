package usersservice

import "go.opentelemetry.io/otel/trace"

type UserService struct {
	userRepo userRepo
	tracer   trace.Tracer
}

func NewUserService(userRepo userRepo, tracer trace.Tracer) *UserService {
	return &UserService{
		userRepo: userRepo,
		tracer:   tracer,
	}
}
