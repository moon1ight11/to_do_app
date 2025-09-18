package users

import (
	"todoapp/internal/storage"

	"github.com/google/uuid"
)

type Base struct {
	storage.DataBase
}

type User struct {
	Id    uuid.UUID
	Name  string
	Email string
	Pass  string
}
