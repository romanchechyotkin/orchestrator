package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/romanchechyotkin/orchestrator/internal/task"
)

func main() {
	t := &task.Task{
		ID: uuid.New(),
	}

	ev := task.TaskEvent{
		ID:   uuid.New(),
		Task: t,
	}
	fmt.Printf("task %+v\n", t)
	fmt.Printf("task event %+v\n", ev)
}
