package usersrepos

import (
	"fmt"
	"github.com/google/uuid"
)

// удаление пользователя по id
func (db *Repo) DeleteUser(user_id uuid.UUID) error {
	transaction, err := db.DB.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	querySettings := `
					DELETE FROM todo_app.settings
					WHERE user_id = $1
					`
	_, err = transaction.Exec(querySettings, user_id)
	if err != nil {
		return fmt.Errorf("Error in DeleteUser query: %w", err)
	}

	query := `
				DELETE FROM todo_app.users
				WHERE id = $1
			`
	_, err = transaction.Exec(query, user_id)
	if err != nil {
		return fmt.Errorf("Error in DeleteUser query: %w", err)
	}

	transaction.Commit()
	
	return nil
}
