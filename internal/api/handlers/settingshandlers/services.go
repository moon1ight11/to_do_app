package settingshandlers

import (
	"context"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type settingsService interface {
	GetSettings(ctx context.Context, userId uuid.UUID) (models.Setting, error)
	UpdateSettings(ctx context.Context, userId uuid.UUID, duration *float64, tz *string) error
}
