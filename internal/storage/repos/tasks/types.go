package tasks

import "time"

type Task struct {
	Id          int
	Parent_id   int
	Start_at    time.Time
	End_at      time.Time
	Title       string
	Description string
}
