package tasksrepos

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id          *uuid.UUID
	ParentId    *uuid.UUID
	OwnerId     uuid.UUID
	StartAt     *time.Time
	EndAt       *time.Time
	Title       *string
	Description *string
	CompletedAt *bool
	Subtasks    *[]Task
}
