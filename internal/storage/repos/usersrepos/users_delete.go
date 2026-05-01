package usersrepos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// удаление пользователя по id
func (db *Repo) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	ctx, span := db.tracer.Start(ctx, "repo.DeleteUser")
	defer span.End()

	transaction, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteUser BeginTx: %w", err)
	}

	defer transaction.Rollback()

	querySettings := `
					DELETE FROM todo_app.settings
					WHERE user_id = $1
					`
	_, err = transaction.ExecContext(ctx, querySettings, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteSettings query: %w", err)
	}

	query := `
				DELETE FROM todo_app.users
				WHERE id = $1
			`
	_, err = transaction.ExecContext(ctx, query, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteUser query: %w", err)
	}

	err = transaction.Commit()
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteUser commit: %w", err)
	}

	return nil
}
