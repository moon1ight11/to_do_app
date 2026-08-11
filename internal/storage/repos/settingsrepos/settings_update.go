package settingsrepos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (db *Repo) UpdateTZ(ctx context.Context, tz string, userId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateTZ")
	defer span.End()

	query := `
		UPDATE todo_app.settings
		SET default_tz = $1
		WHERE user_id = $2
	`

	_, err := tx.ExecContext(ctx, query, tz, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("settingsrepos.UpdateTZ: exec: %w", err)
	}

	return nil
}

func (db *Repo) UpdateDuration(ctx context.Context, duration float64, userId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateDuration")
	defer span.End()

	query := `
		UPDATE todo_app.settings
		SET default_duration = $1
		WHERE user_id = $2
	`

	_, err := tx.ExecContext(ctx, query, duration, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("settingsrepos.UpdateDuration: exec: %w", err)
	}

	return nil
}
