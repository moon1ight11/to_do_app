package settings

import (
	"time"
	"todoapp/internal/storage"
)

type Base struct {
	storage.DataBase
}

type Setting struct {
	DefaultTz       string        `json:"default_tz" binding:"required"`
	DefaultDuration time.Duration `json:"default_duration" binding:"required"`
}
