package worker

import (
	"github.com/romanchechyotkin/orchestrator/internal/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type TasksStorage struct {
	tasks      map[uuid.UUID]*task.Task
	tasksCount uint
	queue      queue.Queue // todo implement own Queue
}
