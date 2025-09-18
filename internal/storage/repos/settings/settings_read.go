package settings

import "log"

// получение настроек по id
func (db *Base) SettingsById(id int) (Setting, error) {
	query := `
				SELECT default_tz, default_duration
				FROM settings
				WHERE id = $1
			`
	var setting Setting 
	err := db.DB.QueryRow(query, id).Scan(&setting.DefaultTz, &setting.DefaultDuration)
	if err != nil {
		log.Println("Error in SettingsById query", err)
		return Setting{}, err
	}
	return setting, nil
}
