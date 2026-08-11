package usersservice

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"todoapp/internal/storage/repos/usersrepos"
)

type userRepo interface {
	DB() *sql.DB
	CreateUser(ctx context.Context, name string, hashPass string, email string) (uuid.UUID, error)
	DeleteUser(ctx context.Context, userId uuid.UUID) error
	UserById(ctx context.Context, userId uuid.UUID) (usersrepos.User, error)
	UserByEmail(ctx context.Context, userEmail string) (usersrepos.User, error)
	CheckUserName(ctx context.Context, name string) (bool, error)
	CheckUserEmail(ctx context.Context, email string) (bool, error)
	UpdateName(ctx context.Context, name string, userId uuid.UUID, tx *sql.Tx) error
	UpdatePass(ctx context.Context, pass string, userId uuid.UUID, tx *sql.Tx) error
	UpdateEmail(ctx context.Context, email string, userId uuid.UUID, tx *sql.Tx) error
}
