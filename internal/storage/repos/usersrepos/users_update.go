package usersrepos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// обновление имени
func (db *Repo) UpdateName(ctx context.Context, name string, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, name, userId)
	if err != nil {
		return fmt.Errorf("error in UpdateName query: %w", err)
	}

	return nil
}

// обновление пароля
func (db *Repo) UpdatePass(ctx context.Context, pass string, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET pass = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, pass, userId)
	if err != nil {
		return fmt.Errorf("error in UpdatePass query: %w", err)
	}
	return nil
}

// обновление почты
func (db *Repo) UpdateEmail(ctx context.Context, email string, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET email = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, email, userId)
	if err != nil {
		return fmt.Errorf("error in UpdateEmail query: %w", err)
	}
	return nil
}
