package users

import (
	"fmt"
	"github.com/google/uuid"
)

// Добавление пользователя
func (db *Repo) CreateUser(name string, hashPass string, email string) (uuid.UUID, error) {
	transaction, err := db.DB.Begin()
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

	err = transaction.QueryRow(query, name, hashPass, email).Scan(&userId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in AddUser query: %w", err)
	}

	querySettings := `
    					INSERT INTO todo_app.settings (user_id)
    					VALUES ($1)
					`

	_, err = transaction.Exec(querySettings, userId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in AddSettings query: %w", err)
	}

	transaction.Commit()

	return userId, nil
}
