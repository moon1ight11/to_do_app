package users

import (
	"todoapp/internal/storage"
	"github.com/google/uuid"
)

type Base struct {
	storage.DataBase
}

func NewBase(db *storage.DataBase) *Base {
    if db == nil || db.DB == nil {
        panic("database connection cannot be nil")
    }
    return &Base{DataBase: *db}
}

type User struct {
	Id    uuid.UUID `json:"user_id"`
	Name  string    `json:"user_name"`
	Email string    `json:"user_email" binding:"required,email"`
	Pass  string    `json:"user_pass"`
}
