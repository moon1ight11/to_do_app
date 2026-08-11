package usershandlers

import (
	"context"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type MockUserService struct {
	GetUserFunc    func(ctx context.Context, userId uuid.UUID) (models.UserRequest, error)
	UpdateUserFunc func(ctx context.Context, name *string, pass *string, email *string, userId uuid.UUID) error
	DeleteUserFunc func(ctx context.Context, userId uuid.UUID) error
}

func (m *MockUserService) GetUser(ctx context.Context, userId uuid.UUID) (models.UserRequest, error) {
	return m.GetUserFunc(ctx, userId)
}

func (m *MockUserService) UpdateUser(ctx context.Context, name *string, pass *string, email *string, userId uuid.UUID) error {
	return m.UpdateUserFunc(ctx, name, pass, email, userId)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	return m.DeleteUserFunc(ctx, userId)
}