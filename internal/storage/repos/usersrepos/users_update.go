package usersrepos

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// обновление имени
func (db *Repo) UpdateName (name string, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, name, userId)
	if err != nil {
		return fmt.Errorf("Error in UpdateName query: %w", err)
	}

	return nil
}

// обновление пароля
func (db *Repo) UpdatePass (pass string, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET pass = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, pass, userId)
	if err != nil {
		return fmt.Errorf("Error in UpdatePass query: %w", err)
	}
	return nil
}

// обновление почты
func (db *Repo) UpdateEmail (email string, userId uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET email = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, email, userId)
	if err != nil {
		return fmt.Errorf("Error in UpdateEmail query: %w", err)
	}
	return nil
}