package settingsrepos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// получение настроек по id
func (db *Repo) SettingsById(ctx context.Context, userId uuid.UUID) (Setting, error) {
	query := `
				SELECT default_tz, default_duration
				FROM todo_app.settings
				WHERE user_id = $1
			`
	var setting Setting

	err := db.DB.QueryRowContext(ctx, query, userId).Scan(&setting.UserTz, &setting.TimeDuration)
	if err != nil {
		return Setting{}, fmt.Errorf("error in SettingsById query: %w", err)
	}

	return setting, nil
}
