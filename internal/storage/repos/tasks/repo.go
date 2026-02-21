package tasks

import (
	"todoapp/internal/storage"
)

type Repo struct {
	storage.DataBase
}

func NewTasksRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}