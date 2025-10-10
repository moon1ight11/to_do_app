package users

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// получение количества пользователей
func (db *Base) AmountUsers() (int, error) {
	query := `
				SELECT COUNT(*)
				FROM todo_app.users
			`
	var amount int
	err := db.DB.QueryRow(query).Scan(&amount)
	if err != nil {
		return 0, fmt.Errorf("Error in AmountUsers query: %w", err)
	}

	return amount, nil
}

// получение пользователя по id
func (db *Base) UserById(id uuid.UUID) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM todo_app.users
				WHERE id = $1
			`
	var user User
	err := db.DB.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		return User{}, fmt.Errorf("Error in UserById query: %w", err)
	}

	return user, nil
}

// получение пользователя по почте
func (db *Base) UserByEmail(email string) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM todo_app.users
				WHERE email = $1
			`
	var user User
	err := db.DB.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
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
func (db *Base) CheckUserName(userName string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM todo_app.users
					WHERE name = $1)
			`
	var exist bool

	err := db.DB.QueryRow(query, userName).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("Error in CheckUserName query: %w", err)
	}

	return exist, nil
}

// проверка свободности почты
func (db *Base) CheckUserEmail(userEmail string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM todo_app.users
					WHERE email = $1)
			`
	var exist bool

	err := db.DB.QueryRow(query, userEmail).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("Error in CheckUserEmail query: %w", err)
	}

	return exist, nil
}
