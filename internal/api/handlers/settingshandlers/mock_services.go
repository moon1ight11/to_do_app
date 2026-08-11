package settingshandlers

import (
	"context"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

type MockSettingsService struct {
	GetSettingsFunc    func(ctx context.Context, userId uuid.UUID) (models.Setting, error)
	UpdateSettingsFunc func(ctx context.Context, userId uuid.UUID, duration *float64, tz *string) error
}

func (m *MockSettingsService) GetSettings(ctx context.Context, userId uuid.UUID) (models.Setting, error) {
	return m.GetSettingsFunc(ctx, userId)
}

func (m *MockSettingsService) UpdateSettings(ctx context.Context, userId uuid.UUID, duration *float64, tz *string) error {
	return m.UpdateSettingsFunc(ctx, userId, duration, tz)
}
