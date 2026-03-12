package settingsrepos

import "todoapp/internal/storage"

type Repo struct {
	storage.DataBase
}

func NewSettingsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}