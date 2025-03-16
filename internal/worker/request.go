package worker

import (
	"github.com/romanchechyotkin/orchestrator/internal/task"

	"github.com/google/uuid"
)

type GetAllTasksRequest struct{}

type GetTaskRequest struct {
	ID uuid.UUID
}

type CreateTaskRequest struct {
	TaskEvent *task.TaskEvent
}
