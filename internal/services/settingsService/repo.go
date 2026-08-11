package settingsservice

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"todoapp/internal/storage/repos/settingsrepos"
)

type settingsRepo interface {
	DB() *sql.DB
	SettingsById(ctx context.Context, userId uuid.UUID) (settingsrepos.Setting, error)
	UpdateTZ(ctx context.Context, tz string, userId uuid.UUID, tx *sql.Tx) error
	UpdateDuration(ctx context.Context, duration float64, userId uuid.UUID, tx *sql.Tx) error
}
