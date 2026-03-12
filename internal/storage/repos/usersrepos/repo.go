package usersrepos

import (
	"todoapp/internal/storage"
)

type Repo struct {
	storage.DataBase
}

func NewUserRepo(db *storage.DataBase) *Repo {
    return &Repo{DataBase: *db}
}