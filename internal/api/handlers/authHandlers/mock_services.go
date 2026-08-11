package authhandlers

import (
	"context"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type MockUserService struct {
	AddUserFunc         func(ctx context.Context, user models.UserAuth) (uuid.UUID, error)
	CheckAndGetUserFunc func(ctx context.Context, user models.UserAuth) (models.UserRequest, error)
}

func (m *MockUserService) AddUser(ctx context.Context, user models.UserAuth) (uuid.UUID, error) {
	return m.AddUserFunc(ctx, user)
}

func (m *MockUserService) CheckAndGetUser(ctx context.Context, user models.UserAuth) (models.UserRequest, error) {
	return m.CheckAndGetUserFunc(ctx, user)
}
