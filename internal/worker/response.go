package worker

import (
	"github.com/romanchechyotkin/orchestrator/internal/task"

	"github.com/google/uuid"
)

type GetAllTasksResponse struct {
	Tasks []*task.Task
	Err   error
}

type GetTaskResponse struct {
	Task *task.Task
	Err  error
}

type CreateTaskResponse struct {
	ID  uuid.UUID
	Err error
}
