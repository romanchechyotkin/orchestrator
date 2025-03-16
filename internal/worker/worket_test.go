package worker

import (
	"log"
	"net/rpc"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romanchechyotkin/orchestrator/internal/task"
	"github.com/romanchechyotkin/orchestrator/pkg/docker"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

func TestRunTask(t *testing.T) {
	client, err := docker.NewClient(
		client.WithHost("unix:///Users/rchechetkin/.colima/ozon/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	worker := New("test", client)

	testTask := &task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	worker.AddTask(testTask)
	result := worker.runTask()
	require.NoError(t, result.Error)

	testTask.ContainerID = result.ContainerID
	log.Printf("task %s is running in container %s\n", testTask.ID, testTask.ContainerID)

	log.Println("Sleepy time")
	time.Sleep(time.Second * 5)
	log.Printf("stopping task %+v\n", testTask)

	completeTask := &task.Task{
		ID:          testTask.ID,
		State:       task.Completed,
		ContainerID: testTask.ContainerID,
	}
	worker.AddTask(completeTask)
	result = worker.runTask()
	require.NoError(t, result.Error)
	require.Equal(t, result.ContainerID, completeTask.ContainerID)
}

func TestGetAllTasks(t *testing.T) {
	worker := New("test", nil)
	go worker.serve()

	require.Eventually(t, func() bool {
		c, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
		if err == nil {
			_ = c.Close()
			return true
		}
		return false
	}, 5*time.Second, 100*time.Millisecond)

	task1 := &task.Task{
		ID: uuid.New(),
	}

	task2 := &task.Task{
		ID: uuid.New(),
	}

	worker.tasksStorage.Set(task1.ID, task1)
	worker.tasksStorage.Set(task2.ID, task2)
	require.Len(t, worker.tasksStorage.tasks, 2)

	c, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
	require.NoError(t, err)
	defer c.Close()

	res := &GetAllTasksResponse{}

	err = c.Call("Worker.GetAllTasks", &GetAllTasksRequest{}, res)
	require.NoError(t, err)
	require.NoError(t, res.Err)
	require.Len(t, res.Tasks, 2)
}
