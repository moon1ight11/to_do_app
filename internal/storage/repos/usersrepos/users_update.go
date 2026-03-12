package usersrepos

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// обновление имени
func (db *Repo) UpdateName (newName string, user_id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newName, user_id)
	if err != nil {
		return fmt.Errorf("Error in UpdateName query: %w", err)
	}

	return nil
}

// обновление пароля
func (db *Repo) UpdatePass (newPass string, user_id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET pass = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newPass, user_id)
	if err != nil {
		return fmt.Errorf("Error in UpdatePass query: %w", err)
	}
	return nil
}

// обновление почты
func (db *Repo) UpdateEmail (newEmail string, user_id uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE todo_app.users
				SET email = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newEmail, user_id)
	if err != nil {
		return fmt.Errorf("Error in UpdateEmail query: %w", err)
	}
	return nil
}