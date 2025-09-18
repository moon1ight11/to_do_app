package settings

import (
	"fmt"
	"time"
	"github.com/google/uuid"
)

// изменение временной зоны по id
func (db *Base) UpdateTZ(newTZ string, id uuid.UUID) error {
	query := `
				UPDATE settings
				SET default_tz = $1
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newTZ, id)
	if err != nil {
		return fmt.Errorf("Error in UpdateTZ query: %w", err)
	}

	return nil
}

// изменение продолжительности по id
func (db *Base) UpdateDuration(newDuration time.Duration, id uuid.UUID) error {
	query := `
				UPDATE settings
				SET default_duration = $1
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newDuration, id)
	if err != nil {
		return fmt.Errorf("Error in UpdateDuration query: %w", err)
	}

	return nil
}
