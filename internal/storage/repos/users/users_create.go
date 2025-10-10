package users

import (
	"fmt"

	"github.com/google/uuid"
)

// Добавление пользователя в DB и возврат его id
func (db *Base) AddUser(NewUser User) (uuid.UUID, error) {
	query := `
				INSERT INTO todo_app.users (name, pass, email)
				VALUES ($1, $2, $3)
				RETURNING id
			`
	var user_id uuid.UUID
	
	err := db.DB.QueryRow(query, NewUser.Name, NewUser.Pass, NewUser.Email).Scan(&user_id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in AddUser query: %w", err)
	}

	return user_id, nil
}