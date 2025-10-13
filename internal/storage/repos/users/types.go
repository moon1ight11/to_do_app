package users

import (
	"todoapp/internal/storage"
	"github.com/google/uuid"
)

type Base struct {
	storage.DataBase
}

func NewUserBase(db *storage.DataBase) *Base {
    return &Base{DataBase: *db}
}

type User struct {
	Id    uuid.UUID `json:"user_id"`
	Name  string    `json:"user_name"`
	Email string    `json:"user_email" binding:"required,email"`
	Pass  string    `json:"user_pass"`
}
