package usersservice

import (
	"todoapp/internal/storage/repos/usersrepos"

	"go.opentelemetry.io/otel/trace"
)

type UserService struct {
	userRepo *usersrepos.Repo
	tracer   trace.Tracer
}

func NewUserService(userRepo *usersrepos.Repo, tracer trace.Tracer) *UserService {
	return &UserService{
		userRepo: userRepo,
		tracer:   tracer,
	}
}
