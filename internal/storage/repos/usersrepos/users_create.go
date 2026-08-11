package usersrepos

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (db *Repo) CreateUser(ctx context.Context, name string, hashPass string, email string) (uuid.UUID, error) {
	ctx, span := db.tracer.Start(ctx, "repo.CreateUser")
	defer span.End()

	transaction, err := db.DB().BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersrepos.CreateUser: begin tx: %w", err)
	}
	defer transaction.Rollback()

	query := `
		INSERT INTO todo_app.users (name, pass, email)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var userId uuid.UUID
	err = transaction.QueryRowContext(ctx, query, name, hashPass, email).Scan(&userId)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersrepos.CreateUser: insert user: %w", err)
	}

	querySettings := `
		INSERT INTO todo_app.settings (user_id)
		VALUES ($1)
	`

	_, err = transaction.ExecContext(ctx, querySettings, userId)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersrepos.CreateUser: insert settings: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersrepos.CreateUser: commit: %w", err)
	}

	return userId, nil
}
