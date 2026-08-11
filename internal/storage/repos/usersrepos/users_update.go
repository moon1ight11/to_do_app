package usersrepos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (db *Repo) UpdateName(ctx context.Context, name string, userId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateName")
	defer span.End()

	query := `
		UPDATE todo_app.users
		SET name = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, name, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersrepos.UpdateName: exec: %w", err)
	}

	return nil
}

func (db *Repo) UpdatePass(ctx context.Context, pass string, userId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdatePass")
	defer span.End()

	query := `
		UPDATE todo_app.users
		SET pass = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, pass, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersrepos.UpdatePass: exec: %w", err)
	}

	return nil
}

func (db *Repo) UpdateEmail(ctx context.Context, email string, userId uuid.UUID, tx *sql.Tx) error {
	ctx, span := db.tracer.Start(ctx, "repo.UpdateEmail")
	defer span.End()

	query := `
		UPDATE todo_app.users
		SET email = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, email, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersrepos.UpdateEmail: exec: %w", err)
	}

	return nil
}
