package usersservice

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"todoapp/internal/storage/repos/usersrepos"
)

type MockUserRepo struct {
	DBField             *sql.DB
	CreateUserFunc      func(ctx context.Context, name string, hashPass string, email string) (uuid.UUID, error)
	DeleteUserFunc      func(ctx context.Context, userId uuid.UUID) error
	UserByIdFunc        func(ctx context.Context, userId uuid.UUID) (usersrepos.User, error)
	UserByEmailFunc     func(ctx context.Context, userEmail string) (usersrepos.User, error)
	CheckUserNameFunc   func(ctx context.Context, name string) (bool, error)
	CheckUserEmailFunc  func(ctx context.Context, email string) (bool, error)
	UpdateNameFunc      func(ctx context.Context, name string, userId uuid.UUID, tx *sql.Tx) error
	UpdatePassFunc      func(ctx context.Context, pass string, userId uuid.UUID, tx *sql.Tx) error
	UpdateEmailFunc     func(ctx context.Context, email string, userId uuid.UUID, tx *sql.Tx) error
	BeginTxFunc         func(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func (m *MockUserRepo) DB() *sql.DB { return m.DBField }

func (m *MockUserRepo) CreateUser(ctx context.Context, name string, hashPass string, email string) (uuid.UUID, error) {
	return m.CreateUserFunc(ctx, name, hashPass, email)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	return m.DeleteUserFunc(ctx, userId)
}

func (m *MockUserRepo) UserById(ctx context.Context, userId uuid.UUID) (usersrepos.User, error) {
	return m.UserByIdFunc(ctx, userId)
}

func (m *MockUserRepo) UserByEmail(ctx context.Context, userEmail string) (usersrepos.User, error) {
	return m.UserByEmailFunc(ctx, userEmail)
}

func (m *MockUserRepo) CheckUserName(ctx context.Context, name string) (bool, error) {
	return m.CheckUserNameFunc(ctx, name)
}

func (m *MockUserRepo) CheckUserEmail(ctx context.Context, email string) (bool, error) {
	return m.CheckUserEmailFunc(ctx, email)
}

func (m *MockUserRepo) UpdateName(ctx context.Context, name string, userId uuid.UUID, tx *sql.Tx) error {
	return m.UpdateNameFunc(ctx, name, userId, tx)
}

func (m *MockUserRepo) UpdatePass(ctx context.Context, pass string, userId uuid.UUID, tx *sql.Tx) error {
	return m.UpdatePassFunc(ctx, pass, userId, tx)
}

func (m *MockUserRepo) UpdateEmail(ctx context.Context, email string, userId uuid.UUID, tx *sql.Tx) error {
	return m.UpdateEmailFunc(ctx, email, userId, tx)
}

func (m *MockUserRepo) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return m.BeginTxFunc(ctx, opts)
}