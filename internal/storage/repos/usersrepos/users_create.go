package usersrepos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// Добавление пользователя
func (db *Repo) CreateUser(ctx context.Context, name string, hashPass string, email string) (uuid.UUID, error) {
	transaction, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
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
		return uuid.Nil, fmt.Errorf("error in AddUser query: %w", err)
	}

	querySettings := `
    					INSERT INTO todo_app.settings (user_id)
    					VALUES ($1)
					`

	_, err = transaction.ExecContext(ctx, querySettings, userId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in AddSettings query: %w", err)
	}

	err = transaction.Commit()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in create user commit: %w", err)
	}

	return userId, nil
}
