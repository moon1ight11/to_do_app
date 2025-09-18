package tasks

import (
	"time"
	"todoapp/internal/storage"
)

type Base struct {
	storage.DataBase
}

type Task struct {
	Id          int
	Parent_id   int
	Owner_id	int
	Start_at    time.Time
	End_at      time.Time
	Title       string
	Description string
}
