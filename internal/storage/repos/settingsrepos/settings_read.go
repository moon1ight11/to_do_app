package settingsrepos

import (
	"fmt"
	"github.com/google/uuid"
)

// получение настроек по id
func (db *Repo) SettingsById(userId uuid.UUID) (Setting, error) {
	query := `
				SELECT default_tz, default_duration
				FROM todo_app.settings
				WHERE user_id = $1
			`
	var setting Setting

	err := db.DB.QueryRow(query, userId).Scan(&setting.UserTz, &setting.TimeDuration)
	if err != nil {
		return Setting{}, fmt.Errorf("Error in SettingsById query: %w", err)
	}

	return setting, nil
}
