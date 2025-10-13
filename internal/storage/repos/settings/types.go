package settings

import (
	"todoapp/internal/storage"
)

type Base struct {
	storage.DataBase
}

func NewSettingsBase(db *storage.DataBase) *Base {
	return &Base{DataBase: *db}
}

type Setting struct {
	DefaultTz       string  `json:"default_tz" binding:"required"`
	DefaultDuration float64 `json:"default_duration" binding:"required"`
}
