package usersrepos

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (db *Repo) UserById(ctx context.Context, userId uuid.UUID) (User, error) {
	ctx, span := db.tracer.Start(ctx, "repo.UserById")
	defer span.End()

	query := `
		SELECT id, name, email, pass
		FROM todo_app.users
		WHERE id = $1
	`

	var user User
	err := db.DB().QueryRowContext(ctx, query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		span.RecordError(err)
		return User{}, fmt.Errorf("usersrepos.UserById: query: %w", err)
	}

	return user, nil
}

func (db *Repo) UserByEmail(ctx context.Context, userEmail string) (User, error) {
	ctx, span := db.tracer.Start(ctx, "repo.UserByEmail")
	defer span.End()

	query := `
		SELECT id, name, email, pass
		FROM todo_app.users
		WHERE email = $1
	`

	var user User
	err := db.DB().QueryRowContext(ctx, query, userEmail).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		span.RecordError(err)
		return User{}, fmt.Errorf("usersrepos.UserByEmail: query: %w", err)
	}

	return user, nil
}

func (db *Repo) CheckUserName(ctx context.Context, name string) (bool, error) {
	ctx, span := db.tracer.Start(ctx, "repo.CheckUserName")
	defer span.End()

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM todo_app.users
			WHERE name = $1
		)
	`

	var exist bool
	err := db.DB().QueryRowContext(ctx, query, name).Scan(&exist)
	if err != nil {
		span.RecordError(err)
		return false, fmt.Errorf("usersrepos.CheckUserName: query: %w", err)
	}

	return exist, nil
}

func (db *Repo) CheckUserEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := db.tracer.Start(ctx, "repo.CheckUserEmail")
	defer span.End()

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM todo_app.users
			WHERE email = $1
		)
	`

	var exist bool
	err := db.DB().QueryRowContext(ctx, query, email).Scan(&exist)
	if err != nil {
		span.RecordError(err)
		return false, fmt.Errorf("usersrepos.CheckUserEmail: query: %w", err)
	}

	return exist, nil
}
