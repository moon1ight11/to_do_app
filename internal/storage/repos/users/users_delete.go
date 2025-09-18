package users

import (
	"fmt"

	"github.com/google/uuid"
)

// удаление пользователя по id
func (db *Base) DeleteUser(id uuid.UUID) error {
	query := `
				DELETE FROM users
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error in AddUser query: %w", err)
	}
	return nil
}