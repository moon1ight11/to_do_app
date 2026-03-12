package tasksrepos

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id          *uuid.UUID `json:"task_id"`
	ParentId    *uuid.UUID `json:"parent_id"`
	OwnerId     uuid.UUID  `json:"owner_id"`
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	CompletedAt *bool      `json:"completed_at"`
	Subtasks    *[]Task    `json:"subtasks"`
}
