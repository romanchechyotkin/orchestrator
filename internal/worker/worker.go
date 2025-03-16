package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
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
	w := &Worker{
		name:         name,
		tasksStorage: NewTasksStorage(),
		dockerClient: dc,
	}

	return w
}

func (w *Worker) serve() {
	if err := rpc.Register(w); err != nil {
		log.Printf("failed to register rpc server for worker %s: %v\n", w.name, err)
		os.Exit(1)
	}
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("failed to listen socket for worker %s: %v", w.name, err)
		os.Exit(1)
	}

	log.Printf("worker %s listening on port 8080\n", w.name)

	if err = http.Serve(l, nil); err != nil {
		log.Printf("failed to serve http server for worker %s: %v", w.name, err)
		os.Exit(1)
	}
}

func (w *Worker) GetAllTasks(_ *GetAllTasksRequest, reply *GetAllTasksResponse) error {
	reply.Tasks = w.tasksStorage.GetAll()
	reply.Err = nil

	return nil
}

func (w *Worker) GetTask(args *GetTaskRequest, reply *GetTaskResponse) error {
	task, ok := w.tasksStorage.Get(args.ID)
	if !ok {
		reply.Err = errors.New("task not found")
	}

	reply.Task = task
	return nil
}

func (w *Worker) CreateTask(args *CreateTaskRequest, reply *CreateTaskResponse) error {
	return nil
}

func (w *Worker) CollectStats() {
	fmt.Println("will collect stats")
}

func (w *Worker) runTask() *docker.Result {
	queuedTask := w.tasksStorage.Pop()
	if queuedTask == nil {
		log.Println("No tasks in the queue")
		return &docker.Result{}
	}

	log.Printf("got task %+v from queue\n", queuedTask)

	persistedTask, ok := w.tasksStorage.Get(queuedTask.ID)
	if !ok {
		persistedTask = queuedTask
		w.tasksStorage.Set(queuedTask.ID, queuedTask)
	}

	var result docker.Result
	if task.ValidStateTransition(persistedTask.State, queuedTask.State) {
		switch queuedTask.State {
		case task.Scheduled:
			result = *w.startTask(queuedTask)
		case task.Completed:
			result = *w.stopTask(queuedTask)
		default:
			result.Error = errors.New("We should not get here")
		}
	} else {
		err := fmt.Errorf("Invalid transition from %v to %v", persistedTask.State, persistedTask.State)
		result.Error = err
	}

	return &result
}

func (w *Worker) startTask(t *task.Task) *docker.Result {
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

func (w *Worker) stopTask(t *task.Task) *docker.Result {
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
