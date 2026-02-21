package tasks

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id           *uuid.UUID `json:"task_id"`
	Parent_id    *uuid.UUID `json:"parent_id"`
	Owner_id     uuid.UUID  `json:"owner_id"`
	Start_at     *time.Time `json:"start_at"`
	End_at       *time.Time `json:"end_at"`
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	Completed_at *bool      `json:"completed_at"`
	Subtasks     *[]Task    `json:"subtasks"`
}
