package worker

import (
	"fmt"
	"log"
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
	result := worker.RunTask()
	require.NoError(t, result.Error)

	testTask.ContainerID = result.ContainerID
	log.Printf("task %s is running in container %s\n", testTask.ID, testTask.ContainerID)

	log.Println("Sleepy time")
	time.Sleep(time.Second * 5)
	fmt.Printf("stopping task %s\n", testTask.ID)

	testTask.State = task.Completed
	worker.AddTask(testTask)
	result = worker.RunTask()
	require.NoError(t, result.Error)
	require.Equal(t, result.ContainerID, testTask.ContainerID)
}
