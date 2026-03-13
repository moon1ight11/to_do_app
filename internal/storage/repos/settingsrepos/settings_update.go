package settingsrepos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// изменение временной зоны по id
func (db *Repo) UpdateTZ(ctx context.Context, tz string, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.settings
				SET default_tz = $1
				WHERE user_id = $2
			`

	_, err := tx.ExecContext(ctx, query, tz, userId)
	if err != nil {
		return fmt.Errorf("error in UpdateTZ query: %w", err)
	}

	return nil
}

// изменение продолжительности по id
func (db *Repo) UpdateDuration(ctx context.Context, duration float64, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.settings
				SET default_duration = $1
				WHERE user_id = $2
			`
	_, err := tx.ExecContext(ctx, query, duration, userId)
	if err != nil {
		return fmt.Errorf("error in UpdateDuration query: %w", err)
	}

	return nil
}
