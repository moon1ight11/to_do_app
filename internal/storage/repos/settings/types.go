package settings

import (
	"todoapp/internal/storage"
)

type Repo struct {
	storage.DataBase
}

func NewSettingsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}

type Setting struct {
	DefaultTz       string  `json:"default_tz" binding:"required"`
	DefaultDuration float64 `json:"default_duration" binding:"required"`
}
