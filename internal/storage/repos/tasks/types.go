package tasks

import (
	"time"
	"todoapp/internal/storage"
	"github.com/google/uuid"
)

type Base struct {
	storage.DataBase
}

type Task struct {
	Id          uuid.UUID
	Parent_id   *uuid.UUID
	Owner_id    uuid.UUID
	Start_at    time.Time
	End_at      time.Time
	Title       string
	Description string
}
