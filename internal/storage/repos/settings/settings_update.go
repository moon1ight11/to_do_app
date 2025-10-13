package settings

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// изменение временной зоны по id
func (db *Base) UpdateTZ(newTZ string, id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.settings
				SET default_tz = $1
				WHERE user_id = $2
			`
	_, err := tx.Exec(query, newTZ, id)
	if err != nil {
		return fmt.Errorf("Error in UpdateTZ query: %w", err)
	}

	return nil
}

// изменение продолжительности по id
func (db *Base) UpdateDuration(newDuration float64, id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.settings
				SET default_duration = $1
				WHERE user_id = $2
			`
	_, err := tx.Exec(query, newDuration, id)
	if err != nil {
		return fmt.Errorf("Error in UpdateDuration query: %w", err)
	}

	return nil
}
