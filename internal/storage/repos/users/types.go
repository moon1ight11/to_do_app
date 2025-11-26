package users

import (
	"todoapp/internal/storage"
	"github.com/google/uuid"
)

type Repo struct {
	storage.DataBase
}

func NewUserRepo(db *storage.DataBase) *Repo {
    return &Repo{DataBase: *db}
}

type User struct {
	Id    uuid.UUID `json:"user_id"`
	Name  string    `json:"user_name"`
	Email string    `json:"user_email" binding:"required,email"`
	Pass  string    `json:"user_pass"`
}
