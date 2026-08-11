package authhandlers

import (
	"context"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type userService interface {
	AddUser(ctx context.Context, user models.UserAuth) (uuid.UUID, error)
	CheckAndGetUser(ctx context.Context, user models.UserAuth) (models.UserRequest, error)
}