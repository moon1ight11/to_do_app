package usersrepos

import (
	"fmt"
	"github.com/google/uuid"
)

// удаление пользователя по id
func (db *Repo) DeleteUser(userId uuid.UUID) error {
	transaction, err := db.DB.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	querySettings := `
					DELETE FROM todo_app.settings
					WHERE user_id = $1
					`
	_, err = transaction.Exec(querySettings, userId)
	if err != nil {
		return fmt.Errorf("Error in DeleteUser query: %w", err)
	}

	query := `
				DELETE FROM todo_app.users
				WHERE id = $1
			`
	_, err = transaction.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("Error in DeleteUser query: %w", err)
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("error in delete user commit: %w", err)
	}
	
	return nil
}
