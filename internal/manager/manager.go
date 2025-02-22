package manager

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/romanchechyotkin/orchestrator/internal/task"
)

type Manager struct {
	PengingTasks queue.Queue
	Workers      []string

	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string

	EventsStorage map[string]*task.TaskEvent
	TasksStorage  map[string]*task.Task
}

func (m *Manager) SelectWorker() {
	fmt.Println("will select worker")
}

func (m *Manager) SendTask() {
	fmt.Println("will send task")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("will update tasks")
}
