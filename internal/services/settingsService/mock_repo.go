package settingsservice

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"todoapp/internal/storage/repos/settingsrepos"
)

type MockSettingsRepo struct {
	DBField            *sql.DB
	SettingsByIdFunc   func(ctx context.Context, userId uuid.UUID) (settingsrepos.Setting, error)
	UpdateTZFunc       func(ctx context.Context, tz string, userId uuid.UUID, tx *sql.Tx) error
	UpdateDurationFunc func(ctx context.Context, duration float64, userId uuid.UUID, tx *sql.Tx) error
}

func (m *MockSettingsRepo) DB() *sql.DB { return m.DBField }

func (m *MockSettingsRepo) SettingsById(ctx context.Context, userId uuid.UUID) (settingsrepos.Setting, error) {
	return m.SettingsByIdFunc(ctx, userId)
}

func (m *MockSettingsRepo) UpdateTZ(ctx context.Context, tz string, userId uuid.UUID, tx *sql.Tx) error {
	return m.UpdateTZFunc(ctx, tz, userId, tx)
}

func (m *MockSettingsRepo) UpdateDuration(ctx context.Context, duration float64, userId uuid.UUID, tx *sql.Tx) error {
	return m.UpdateDurationFunc(ctx, duration, userId, tx)
}