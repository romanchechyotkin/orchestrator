package worker

import (
	"github.com/romanchechyotkin/orchestrator/internal/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type TasksStorage struct {
	tasks      map[uuid.UUID]*task.Task
	tasksCount uint
	queue      *queue.Queue // todo implement own Queue
}

func NewTasksStorage() *TasksStorage {
	return &TasksStorage{
		tasks:      make(map[uuid.UUID]*task.Task),
		tasksCount: 0,
		queue:      queue.New(),
	}
}

func (ts *TasksStorage) Set(id uuid.UUID, task *task.Task) {
	ts.tasks[id] = task
}

func (ts *TasksStorage) Get(id uuid.UUID) (*task.Task, bool) {
	value, ok := ts.tasks[id]
	return value, ok
}

func (ts *TasksStorage) GetAll() []*task.Task {
	tasks := make([]*task.Task, 0, len(ts.tasks))

	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (ts *TasksStorage) Push(t *task.Task) {
	ts.queue.Enqueue(t)
}

func (ts *TasksStorage) Pop() *task.Task {
	if ts.queue.Len() == 0 {
		return nil
	}

	return ts.queue.Dequeue().(*task.Task)
}
