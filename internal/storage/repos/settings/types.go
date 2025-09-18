package settings

import (
	"time"
	"todoapp/internal/storage"
)

type Base struct {
	storage.DataBase
}

type Setting struct {
	DefaultTz       time.Location
	DefaultDuration time.Duration
}
