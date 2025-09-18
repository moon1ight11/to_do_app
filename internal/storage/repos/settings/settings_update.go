package settings

import (
	"log"
	"time"
)

// изменение временной зоны по id
func (db *Base) UpdateTZ(newTZ string, id int) error {
	query := `
				UPDATE settings
				SET default_tz = $1
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newTZ, id)
	if err != nil {
		log.Println("Error in UpdateTZ query", err)
		return err
	}

	return nil
}

// изменение продолжительности по id
func (db *Base) UpdateDuration(newDuration time.Duration, id int) error {
	query := `
				UPDATE settings
				SET default_duration = $1
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newDuration, id)
	if err != nil {
		log.Println("Error in UpdateDuration query", err)
		return err
	}

	return nil
}
