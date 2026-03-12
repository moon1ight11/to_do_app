package usersrepos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// получение пользователя по id
func (db *Repo) UserById(ctx context.Context, userId uuid.UUID) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM todo_app.users
				WHERE id = $1
			`
	var user User
	err := db.DB.QueryRowContext(ctx, query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		return User{}, fmt.Errorf("error in UserById query: %w", err)
	}

	return user, nil
}

// получение пользователя по почте
func (db *Repo) UserByEmail(ctx context.Context, userEmail string) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM todo_app.users
				WHERE email = $1
			`
	var user User
	err := db.DB.QueryRowContext(ctx, query, userEmail).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		return User{}, fmt.Errorf("error in UserByEmail query: %w", err)
	}

	return user, nil
}

// проверка занятости имени пользователя
func (db *Repo) CheckUserName(ctx context.Context, name string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM todo_app.users
					WHERE name = $1)
			`
	var exist bool

	err := db.DB.QueryRowContext(ctx, query, name).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("error in CheckUserName query: %w", err)
	}

	return exist, nil
}

// проверка занятости почты
func (db *Repo) CheckUserEmail(ctx context.Context, email string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM todo_app.users
					WHERE email = $1)
			`
	var exist bool

	err := db.DB.QueryRowContext(ctx, query, email).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("error in CheckUserEmail query: %w", err)
	}

	return exist, nil
}
