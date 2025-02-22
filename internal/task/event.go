package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestemp time.Time
	Task      *Task
}
