package usershandlers

import (
	"context"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type userService interface {
	GetUser(ctx context.Context, userId uuid.UUID) (models.UserRequest, error)
	UpdateUser(ctx context.Context, name *string, pass *string, email *string, userId uuid.UUID) error
	DeleteUser(ctx context.Context, userId uuid.UUID) error
}