package users

import (
	"fmt"

	"github.com/google/uuid"
)

// получение списка всех пользователей
func (db *Base) AllUsers() ([]User, error) {
	query := `
				SELECT id, name, email
				FROM users
			`
	var users []User

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error in AllUsers query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("Error in AllUsers scan: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// получение количества пользователей
func (db *Base) AmountUsers() (int, error) {
	query := `
				SELECT COUNT(*)
				FROM users
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
				FROM users
				WHERE id = $1
			`
	var user User
	err := db.DB.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		return User{}, fmt.Errorf("Error in UserById query: %w", err)
	}

	return user, nil
}

// получение пользователя по имени
func (db *Base) UserByName(name string) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM users
				WHERE name = $1
			`
	var user User
	err := db.DB.QueryRow(query, name).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		return User{}, fmt.Errorf("Error in UserByName query: %w", err)
	}

	return user, nil
}

// получение пользователя по почте
func (db *Base) UserByEmail(email string) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM users
				WHERE email = $1
			`
	var user User
	err := db.DB.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		return User{}, fmt.Errorf("Error in UserByEmail query: %w", err)
	}

	return user, nil
}
