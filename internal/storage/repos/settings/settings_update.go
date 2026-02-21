package settings

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// изменение временной зоны по id
func (db *Repo) UpdateTZ(new_tz string, user_id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.settings
				SET default_tz = $1
				WHERE user_id = $2
			`
			
	_, err := tx.Exec(query, new_tz, user_id)
	if err != nil {
		return fmt.Errorf("Error in UpdateTZ query: %w", err)
	}

	return nil
}

// изменение продолжительности по id
func (db *Repo) UpdateDuration(new_duration float64, user_id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.settings
				SET default_duration = $1
				WHERE user_id = $2
			`
	_, err := tx.Exec(query, new_duration, user_id)
	if err != nil {
		return fmt.Errorf("Error in UpdateDuration query: %w", err)
	}

	return nil
}
