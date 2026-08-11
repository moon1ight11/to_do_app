package usersrepos

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (db *Repo) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	ctx, span := db.tracer.Start(ctx, "repo.DeleteUser")
	defer span.End()

	transaction, err := db.DB().BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersrepos.DeleteUser: begin tx: %w", err)
	}
	defer transaction.Rollback()

	_, err = transaction.ExecContext(ctx, `DELETE FROM todo_app.settings WHERE user_id = $1`, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersrepos.DeleteUser: delete settings: %w", err)
	}

	_, err = transaction.ExecContext(ctx, `DELETE FROM todo_app.users WHERE id = $1`, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersrepos.DeleteUser: delete user: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersrepos.DeleteUser: commit: %w", err)
	}

	return nil
}
