package settings

import (
	"fmt"
	"github.com/google/uuid"
)

// получение настроек по id
func (db *Base) SettingsById(id uuid.UUID) (Setting, error) {
	query := `
				SELECT default_tz, default_duration
				FROM settings
				WHERE id = $1
			`
	var setting Setting
	err := db.DB.QueryRow(query, id).Scan(&setting.DefaultTz, &setting.DefaultDuration)
	if err != nil {
		return Setting{}, fmt.Errorf("Error in SettingsById query: %w", err)
	}
	return setting, nil
}
