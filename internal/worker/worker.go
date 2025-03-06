package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/romanchechyotkin/orchestrator/internal/task"
	"github.com/romanchechyotkin/orchestrator/pkg/docker"
)

type Worker struct {
	name string

	tasksStorage *TasksStorage
	dockerClient *docker.Client
}

func New(name string, dc *docker.Client) *Worker {
	return &Worker{
		name:         name,
		tasksStorage: NewTasksStorage(),
		dockerClient: dc,
	}
}

func (w *Worker) CollectStats() {
	fmt.Println("will collect stats")
}

func (w *Worker) RunTask() *docker.Result {
	queuedTask := w.tasksStorage.Pop()
	if queuedTask == nil {
		log.Println("No tasks in the queue")
		return &docker.Result{}
	}

	persistedTask, ok := w.tasksStorage.Get(queuedTask.ID)
	if !ok {
		persistedTask = queuedTask
		w.tasksStorage.Set(queuedTask.ID, queuedTask)
	}

	var result docker.Result
	if task.ValidStateTransition(persistedTask.State, queuedTask.State) {
		switch queuedTask.State {
		case task.Scheduled:
			result = *w.StartTask(queuedTask)
		case task.Completed:
			result = *w.StopTask(queuedTask)
		default:
			result.Error = errors.New("We should not get here")
		}
	} else {
		err := fmt.Errorf("Invalid transition from %v to %v", persistedTask.State, persistedTask.State)
		result.Error = err
	}

	return &result
}

func (w *Worker) StartTask(t *task.Task) *docker.Result {
	log.Printf("worker %s starting task %s\n", w.name, t.ID)

	cfg := task.NewConfig(t)
	t.StartTime = time.Now()

	res := w.dockerClient.Run(context.TODO(), cfg)

	if res.Error != nil {
		log.Printf("failed to start container %s for task %s: %v\n", t.ContainerID, t.ID, res.Error)
		t.State = task.Failed
		t.FinishTime = time.Now()
		w.tasksStorage.Set(t.ID, t)
		return res
	}

	t.ContainerID = res.ContainerID
	t.State = task.Running
	w.tasksStorage.Set(t.ID, t)

	log.Printf("started container %s for task %s\n", t.ContainerID, t.ID)

	return res
}

func (w *Worker) StopTask(t *task.Task) *docker.Result {
	log.Printf("worker %s stopping task %s\n", w.name, t.ContainerID)

	res := w.dockerClient.Stop(t.ContainerID)
	if res.Error != nil {
		log.Printf("failed to stop container %s for task %s: %v\n", t.ContainerID, t.ID, res.Error)
		t.State = task.Failed
		t.FinishTime = time.Now()
		w.tasksStorage.Set(t.ID, t)

		return res
	}

	t.FinishTime = time.Now()
	t.State = task.Completed
	w.tasksStorage.Set(t.ID, t)

	log.Printf("stopped and removed container %s for task %s\n", t.ContainerID, t.ID)

	return res
}

func (w *Worker) AddTask(t *task.Task) {
	w.tasksStorage.Push(t)
}
