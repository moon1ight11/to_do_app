package users

import (
	"fmt"

	"github.com/google/uuid"
)

// обновление имени
func (db *Base) UpdateName (newName string, id uuid.UUID) error {
	query := `
				UPDATE users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newName, id)
	if err != nil {
		return fmt.Errorf("Error in UpdateName query: %w", err)
	}

	return nil
}

// обновление пароля
func (db *Base) UpdatePass (newPass string, id uuid.UUID) error {
	query := `
				UPDATE users
				SET pass = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newPass, id)
	if err != nil {
		return fmt.Errorf("Error in UpdatePass query: %w", err)
	}
	return nil
}

// обновление почты
func (db *Base) UpdateEmail (newEmail string, id uuid.UUID) error {
	query := `
				UPDATE users
				SET email = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newEmail, id)
	if err != nil {
		return fmt.Errorf("Error in UpdateEmail query: %w", err)
	}
	return nil
}