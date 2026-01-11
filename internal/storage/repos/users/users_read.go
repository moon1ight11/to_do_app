package users

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// получение пользователя по id
func (db *Repo) UserById(user_id uuid.UUID) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM todo_app.users
				WHERE id = $1
			`
	var user User
	err := db.DB.QueryRow(query, user_id).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		return User{}, fmt.Errorf("Error in UserById query: %w", err)
	}

	return user, nil
}

// получение пользователя по почте
func (db *Repo) UserByEmail(user_email string) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM todo_app.users
				WHERE email = $1
			`
	var user User
	err := db.DB.QueryRow(query, user_email).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
        // Проверяем, это "запись не найдена" или реальная ошибка БД
        if err == sql.ErrNoRows {
            return User{}, fmt.Errorf("No users with this email")
        }
        return User{}, fmt.Errorf("Error in UserByEmail query: %w", err)
    }

	return user, nil
}

// проверка занятости имени пользователя 
func (db *Repo) CheckUserName(user_name string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM todo_app.users
					WHERE name = $1)
			`
	var exist bool

	err := db.DB.QueryRow(query, user_name).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("Error in CheckUserName query: %w", err)
	}

	return exist, nil
}

// проверка занятости почты
func (db *Repo) CheckUserEmail(user_email string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM todo_app.users
					WHERE email = $1)
			`
	var exist bool

	err := db.DB.QueryRow(query, user_email).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("Error in CheckUserEmail query: %w", err)
	}

	return exist, nil
}
